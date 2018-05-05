// Package goredis is another redis client with full features which writter in golang
//
// Protocol Specification: http://redis.io/topics/protocol.
//
// Redis reply has five types: status, error, integer, bulk, multi bulk.
// A Status Reply is in the form of a single line string starting with "+" terminated by "\r\n".
// Error Replies are very similar to Status Replies. The only difference is that the first byte is "-".
// Integer reply is just a CRLF terminated string representing an integer, prefixed by a ":" byte.
// Bulk replies are used by the server in order to return a single binary safe string up to 512 MB in length.
// A Multi bulk reply is used to return an array of other replies.
// Every element of a Multi Bulk Reply can be of any kind, including a nested Multi Bulk Reply.
// So five reply type is defined:
//  const (
//  	ErrorReply = iota
//  	StatusReply
//  	IntegerReply
//  	BulkReply
//  	MultiReply
//  )
// And then a Reply struct which represent the redis response data is defined:
//  type Reply struct {
//  	Type    int
//  	Error   string
//  	Status  string
//  	Integer int64  // Support Redis 64bit integer
//  	Bulk    []byte // Support Redis Null Bulk Reply
//  	Multi   []*Reply
//  }
// Reply struct has many useful methods:
//  func (rp *Reply) IntegerValue() (int64, error)
//  func (rp *Reply) BoolValue() (bool, error)
//  func (rp *Reply) StatusValue() (string, error)
//  func (rp *Reply) OKValue() error
//  func (rp *Reply) BytesValue() ([]byte, error)
//  func (rp *Reply) StringValue() (string, error)
//  func (rp *Reply) MultiValue() ([]*Reply, error)
//  func (rp *Reply) HashValue() (map[string]string, error)
//  func (rp *Reply) ListValue() ([]string, error)
//  func (rp *Reply) BytesArrayValue() ([][]byte, error)
//  func (rp *Reply) BoolArrayValue() ([]bool, error)
//
// Connect redis has two function: Dial and DialURL, for example:
//  client, err := Dial()
//  client, err := Dial(&DialConfig{Address: "127.0.0.1:6379"})
//  client, err := Dial(&DialConfig{"tcp", "127.0.0.1:6379", 0, "", 10*time.Second, 10})
//  client, err := DialURL("tcp://auth:password@127.0.0.1:6379/0?timeout=10s&maxidle=1")
//
// DialConfig can also take named options for connection config:
//   config := &DialConfig {
//     Network:  "tcp",
//     Address:  "127.0.0.1:6379",
//     Database: 0,
//     Password: "yourpasswordhere"
//     Timeout:  10*time.Second,
//     MaxIdle:  10
//   }
//
// Try a redis command is simple too, let's do GET/SET:
//  err := client.Set("key", "value", 0, 0, false, false)
//  value, err := client.Get("key")
//
// Or you can execute customer command with Redis.ExecuteCommand method:
//  reply, err := client.ExecuteCommand("SET", "key", "value")
//  err := reply.OKValue()
//
// Redis Pipelining is defined as:
//  type Pipelined struct {
//  	redis *Redis
//  	conn  *Connection
//  	times int
//  }
//  func (p *Pipelined) Close()
//  func (p *Pipelined) Command(args ...interface{})
//  func (p *Pipelined) Receive() (*Reply, error)
//  func (p *Pipelined) ReceiveAll() ([]*Reply, error)
//
// Transaction, Lua Eval, Publish/Subscribe, Monitor, Scan, Sort are also supported.
//
package goredis

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func packArgs(items ...interface{}) (args []interface{}) {
	for _, item := range items {
		v := reflect.ValueOf(item)
		switch v.Kind() {
		case reflect.Slice:
			if v.IsNil() {
				continue
			}
			for i := 0; i < v.Len(); i++ {
				args = append(args, v.Index(i).Interface())
			}
		case reflect.Map:
			if v.IsNil() {
				continue
			}
			for _, key := range v.MapKeys() {
				args = append(args, key.Interface(), v.MapIndex(key).Interface())
			}
		default:
			args = append(args, v.Interface())
		}
	}
	return args
}

func numLen(i int64) int64 {
	n, pos10 := int64(1), int64(10)
	if i < 0 {
		i = -i
		n++
	}
	for i >= pos10 {
		n++
		pos10 *= 10
	}
	return n
}

func packCommand(args ...interface{}) ([]byte, error) {
	n := len(args)
	res := make([]byte, 0, 16*n)
	res = append(res, byte('*'))
	res = strconv.AppendInt(res, int64(n), 10)
	res = append(res, byte('\r'), byte('\n'))
	for _, arg := range args {
		res = append(res, byte('$'))
		switch v := arg.(type) {
		case []byte:
			res = strconv.AppendInt(res, int64(len(v)), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = append(res, v...)
		case string:
			res = strconv.AppendInt(res, int64(len(v)), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = append(res, []byte(v)...)
		case int:
			res = strconv.AppendInt(res, numLen(int64(v)), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = strconv.AppendInt(res, int64(v), 10)
		case int64:
			res = strconv.AppendInt(res, numLen(v), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = strconv.AppendInt(res, int64(v), 10)
		case uint64:
			res = strconv.AppendInt(res, numLen(int64(v)), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = strconv.AppendUint(res, uint64(v), 10)
		case float64:
			var buf []byte
			buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
			res = strconv.AppendInt(res, int64(len(buf)), 10)
			res = append(res, byte('\r'), byte('\n'))
			res = append(res, buf...)
		default:
			return nil, errors.New("invalid argument type when pack command")
		}
		res = append(res, byte('\r'), byte('\n'))
	}
	return res, nil
}

type connection struct {
	Conn   net.Conn
	Reader *bufio.Reader
}

func (c *connection) SendCommand(args ...interface{}) error {
	request, err := packCommand(args...)
	if err != nil {
		return err
	}
	if _, err := c.Conn.Write(request); err != nil {
		return err
	}
	return nil
}

func (c *connection) RecvReply() (*Reply, error) {
	line, err := c.Reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-2]
	switch line[0] {
	case '-':
		return &Reply{
			Type:  ErrorReply,
			Error: string(line[1:]),
		}, nil
	case '+':
		return &Reply{
			Type:   StatusReply,
			Status: string(line[1:]),
		}, nil
	case ':':
		i, err := strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &Reply{
			Type:    IntegerReply,
			Integer: i,
		}, nil
	case '$':
		size, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		bulk, err := c.ReadBulk(size)
		if err != nil {
			return nil, err
		}
		return &Reply{
			Type: BulkReply,
			Bulk: bulk,
		}, nil
	case '*':
		i, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		rp := &Reply{Type: MultiReply}
		if i >= 0 {
			multi := make([]*Reply, i)
			for j := 0; j < i; j++ {
				rp, err := c.RecvReply()
				if err != nil {
					return nil, err
				}
				multi[j] = rp
			}
			rp.Multi = multi
		}
		return rp, nil
	}
	return nil, errors.New("redis protocol error")
}

func (c *connection) ReadBulk(size int) ([]byte, error) {
	// If the requested value does not exist the bulk reply will use the special value -1 as data length
	if size < 0 {
		return nil, nil
	}
	buf := make([]byte, size+2)
	if _, err := io.ReadFull(c.Reader, buf); err != nil {
		return nil, err
	}
	return buf[:size], nil
}

type connPool struct {
	MaxIdle int
	Dial    func() (*connection, error)

	idle   *list.List
	closed bool
	mutex  sync.Mutex
}

func (p *connPool) Close() {
	p.mutex.Lock()
	p.closed = true
	for e := p.idle.Front(); e != nil; e = e.Next() {
		e.Value.(*connection).Conn.Close()
	}
	p.mutex.Unlock()
}

func (p *connPool) Get() (*connection, error) {
	p.mutex.Lock()
	if p.closed {
		p.mutex.Unlock()
		return nil, errors.New("connection pool closed")
	}
	if p.idle.Len() > 0 {
		back := p.idle.Back()
		p.idle.Remove(back)
		p.mutex.Unlock()
		return back.Value.(*connection), nil
	}
	p.mutex.Unlock()
	return p.Dial()
}

func (p *connPool) Put(c *connection) {
	p.mutex.Lock()
	if c == nil {
		p.mutex.Unlock()
		return
	}
	if p.closed {
		c.Conn.Close()
		p.mutex.Unlock()
		return
	}
	if p.idle.Len() >= p.MaxIdle {
		p.idle.Remove(p.idle.Front())
	}
	p.idle.PushBack(c)
	p.mutex.Unlock()
}

// Redis client struct
// Containers connection parameters and connection pool
type Redis struct {
	network  string
	address  string
	db       int
	password string
	timeout  time.Duration
	pool     *connPool
}

// ExecuteCommand send any raw redis command and receive reply from redis server
func (r *Redis) ExecuteCommand(args ...interface{}) (*Reply, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	if err := c.SendCommand(args...); err != nil {
		if err != io.EOF {
			return nil, err
		}
		c, err = r.pool.Get()
		if err != nil {
			return nil, err
		}
		if err = c.SendCommand(args...); err != nil {
			return nil, err
		}
	}
	rp, err := c.RecvReply()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		c, err = r.pool.Get()
		if err != nil {
			return nil, err
		}
		if err = c.SendCommand(args...); err != nil {
			return nil, err
		}
		rp, err = c.RecvReply()
	}
	if err == nil {
		r.pool.Put(c)
	}
	return rp, err
}

func (r *Redis) dialConnection() (*connection, error) {
	conn, err := net.DialTimeout(r.network, r.address, r.timeout)
	if err != nil {
		return nil, err
	}
	c := &connection{conn, bufio.NewReader(conn)}
	if r.password != "" {
		if err := c.SendCommand("AUTH", r.password); err != nil {
			return nil, err
		}
		rp, err := c.RecvReply()
		if err != nil {
			return nil, err
		}
		if rp.Type == ErrorReply {
			return nil, errors.New(rp.Error)
		}
	}
	if r.db > 0 {
		if err := c.SendCommand("SELECT", r.db); err != nil {
			return nil, err
		}
		rp, err := c.RecvReply()
		if err != nil {
			return nil, err
		}
		if rp.Type == ErrorReply {
			return nil, errors.New(rp.Error)
		}
	}
	return c, nil
}

// ClosePool close the redis client under connection pool
// this will close all the connections which in the pool
func (r *Redis) ClosePool() {
	r.pool.Close()
}

const (
	// DefaultNetwork is the default value of network
	DefaultNetwork = "tcp"

	// DefaultAddress is the default value of address(host:port)
	DefaultAddress = ":6379"

	// DefaultTimeout is the default value of connect timeout
	DefaultTimeout = 15 * time.Second

	// DefaultMaxIdle is the default value of connection pool size
	DefaultMaxIdle = 1
)

// DialConfig is redis client connect to server parameters
type DialConfig struct {
	Network  string
	Address  string
	Database int
	Password string
	Timeout  time.Duration
	MaxIdle  int
}

func newDialConfigFromURLString(rawurl string) (*DialConfig, error) {
	ul, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	scheme := DefaultNetwork
	if ul.Scheme != "" {
		scheme = ul.Scheme
	}
	host := DefaultAddress
	if ul.Host != "" {
		host = ul.Host
	}
	password := ""
	if ul.User != nil {
		if pw, set := ul.User.Password(); set {
			password = pw
		}
	}
	db := 0
	path := strings.Trim(ul.Path, "/")
	if path != "" {
		db, err = strconv.Atoi(path)
		if err != nil {
			return nil, err
		}
	}
	timeout := DefaultTimeout
	if ul.Query().Get("timeout") != "" {
		timeout, err = time.ParseDuration(ul.Query().Get("timeout"))
		if err != nil {
			return nil, err
		}
	}
	maxidle := DefaultMaxIdle
	if ul.Query().Get("maxidle") != "" {
		maxidle, err = strconv.Atoi(ul.Query().Get("maxidle"))
		if err != nil {
			return nil, err
		}
	}
	return &DialConfig{scheme, host, db, password, timeout, maxidle}, nil
}

// Dial new a redis client with DialConfig
func Dial(cfg *DialConfig) (*Redis, error) {
	if cfg == nil {
		cfg = &DialConfig{}
	}
	if cfg.Network == "" {
		cfg.Network = DefaultNetwork
	}
	if cfg.Address == "" {
		cfg.Address = DefaultAddress
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
	if cfg.MaxIdle == 0 {
		cfg.MaxIdle = DefaultMaxIdle
	}
	r := &Redis{
		network:  cfg.Network,
		address:  cfg.Address,
		db:       cfg.Database,
		password: cfg.Password,
		timeout:  cfg.Timeout,
	}
	r.pool = &connPool{
		MaxIdle: cfg.MaxIdle,
		Dial:    r.dialConnection,
		idle:    list.New(),
	}
	conn, err := r.dialConnection()
	if err != nil {
		return nil, err
	}
	r.pool.Put(conn)
	return r, nil
}

// DialURL new a redis client with URL-like argument
func DialURL(rawurl string) (*Redis, error) {
	dialConfig, err := newDialConfigFromURLString(rawurl)
	if err != nil {
		return nil, err
	}
	return Dial(dialConfig)
}

// Reply Type: Status, Integer, Bulk, Multi Bulk
// Error Reply Type return error directly
const (
	ErrorReply = iota
	StatusReply
	IntegerReply
	BulkReply
	MultiReply
)

// Reply struct Represent Redis Reply
type Reply struct {
	Type    int
	Error   string
	Status  string
	Integer int64  // Support Redis 64bit integer
	Bulk    []byte // Support Redis Null Bulk Reply
	Multi   []*Reply
}

// IntegerValue returns redis reply number value
func (rp *Reply) IntegerValue() (int64, error) {
	if rp.Type == ErrorReply {
		return 0, errors.New(rp.Error)
	}
	if rp.Type != IntegerReply {
		return 0, errors.New("invalid reply type, not integer")
	}
	return rp.Integer, nil
}

// BoolValue returns redis reply integer
// which are also extensively used in order to return true or false.
// For instance commands like EXISTS or SISMEMBER will return 1 for true and 0 for false.
func (rp *Reply) BoolValue() (bool, error) {
	if rp.Type == ErrorReply {
		return false, errors.New(rp.Error)
	}
	if rp.Type != IntegerReply {
		return false, errors.New("invalid reply type, not integer")
	}
	return rp.Integer != 0, nil
}

// StatusValue indicates redis reply a status string
func (rp *Reply) StatusValue() (string, error) {
	if rp.Type == ErrorReply {
		return "", errors.New(rp.Error)
	}
	if rp.Type != StatusReply {
		return "", errors.New("invalid reply type, not status")
	}
	return rp.Status, nil
}

// OKValue indicates redis reply a OK status string
func (rp *Reply) OKValue() error {
	if rp.Type == ErrorReply {
		return errors.New(rp.Error)
	}
	if rp.Type != StatusReply {
		return errors.New("invalid reply type, not status")
	}
	if rp.Status == "OK" {
		return nil
	}
	return errors.New(rp.Status)
}

// BytesValue indicates redis reply a bulk which maybe nil
func (rp *Reply) BytesValue() ([]byte, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != BulkReply {
		return nil, errors.New("invalid reply type, not bulk")
	}
	return rp.Bulk, nil
}

// StringValue indicates redis reply a bulk which should not be nil
func (rp *Reply) StringValue() (string, error) {
	if rp.Type == ErrorReply {
		return "", errors.New(rp.Error)
	}
	if rp.Type != BulkReply {
		return "", errors.New("invalid reply type, not bulk")
	}
	if rp.Bulk == nil {
		return "", nil
	}
	return string(rp.Bulk), nil
}

// MultiValue indicates redis reply a multi bulk
func (rp *Reply) MultiValue() ([]*Reply, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("invalid reply type, not multi bulk")
	}
	return rp.Multi, nil
}

// HashValue indicates redis reply a multi value which represent hash map
func (rp *Reply) HashValue() (map[string]string, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("invalid reply type, not multi bulk")
	}
	result := make(map[string]string)
	if rp.Multi != nil {
		length := len(rp.Multi)
		for i := 0; i < length/2; i++ {
			key, err := rp.Multi[i*2].StringValue()
			if err != nil {
				return nil, err
			}
			value, err := rp.Multi[i*2+1].StringValue()
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
	}
	return result, nil
}

// ListValue indicates redis reply a multi value which represent list
func (rp *Reply) ListValue() ([]string, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("invalid reply type, not multi bulk")
	}
	var result []string
	if rp.Multi != nil {
		for _, subrp := range rp.Multi {
			item, err := subrp.StringValue()
			if err != nil {
				return nil, err
			}
			result = append(result, item)
		}
	}
	return result, nil
}

// BytesArrayValue indicates redis reply a multi value
// which represent list, but item in the list maybe nil
func (rp *Reply) BytesArrayValue() ([][]byte, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("invalid reply type, not multi bulk")
	}
	var result [][]byte
	if rp.Multi != nil {
		for _, subrp := range rp.Multi {
			b, err := subrp.BytesValue()
			if err != nil {
				return nil, err
			}
			result = append(result, b)
		}
	}
	return result, nil
}

// BoolArrayValue indicates redis reply a multi value
// each bulk is an integer(bool)
func (rp *Reply) BoolArrayValue() ([]bool, error) {
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("invalid reply type, not multi bulk")
	}
	var result []bool
	if rp.Multi != nil {
		for _, subrp := range rp.Multi {
			b, err := subrp.BoolValue()
			if err != nil {
				return nil, err
			}
			result = append(result, b)
		}
	}
	return result, nil
}

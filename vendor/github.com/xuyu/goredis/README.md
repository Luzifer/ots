goredis
=======

[![GoDoc](https://godoc.org/github.com/xuyu/goredis?status.png)](https://godoc.org/github.com/xuyu/goredis)

redis client in golang

[Go or Golang](http://golang.org) is an open source programming language that makes it easy to build simple, reliable, and efficient software.

[Redis](http://redis.io) is an open source, BSD licensed, advanced key-value store. It is often referred to as a data structure server since keys can contain strings, hashes, lists, sets and sorted sets.

- Pure golang, and doesn't depend on any 3rd party libraries;
- Hight test coverage and will continue to improve;
- Tested under Go 1.2 and Redis 2.8.3;
- Tested under Go 1.2.1 and Redis 2.8.4;


Features
--------

* Python Redis Client Like API
* Support [Pipeling](http://godoc.org/github.com/xuyu/goredis#Pipelined)
* Support [Transaction](http://godoc.org/github.com/xuyu/goredis#Transaction)
* Support [Publish Subscribe](http://godoc.org/github.com/xuyu/goredis#PubSub)
* Support [Lua Eval](http://godoc.org/github.com/xuyu/goredis#Redis.Eval)
* Support [Connection Pool](http://godoc.org/github.com/xuyu/goredis#ConnPool)
* Support [Dial URL-Like](http://godoc.org/github.com/xuyu/goredis#DialURL)
* Support [monitor](http://godoc.org/github.com/xuyu/goredis#MonitorCommand), [sort](http://godoc.org/github.com/xuyu/goredis#SortCommand), [scan](http://godoc.org/github.com/xuyu/goredis#Redis.Scan), [slowlog](http://godoc.org/github.com/xuyu/goredis#SlowLog) .etc


Document
--------

- [Redis Commands](http://redis.io/commands)
- [Redis Protocol](http://redis.io/topics/protocol)
- [GoDoc](http://godoc.org/github.com/xuyu/goredis)


Simple Example
--------------

Connect:

	client, err := Dial()
	client, err := Dial(&DialConfig{Address: "127.0.0.1:6379"})
	client, err := DialURL("tcp://auth:password@127.0.0.1:6379/0?timeout=10s&maxidle=1")

Try a redis command is simple too, let's do GET/SET:

	err := client.Set("key", "value", 0, 0, false, false)
	value, err := client.Get("key")

Or you can execute a custom command with Redis.ExecuteCommand method:

	reply, err := client.ExecuteCommand("SET", "key", "value")
	err := reply.OKValue()

And then a Reply struct which represent the redis response data is defined:
	
	type Reply struct {
		Type    int
		Error   string
		Status  string
		Integer int64  // Support Redis 64bit integer
		Bulk    []byte // Support Redis Null Bulk Reply
		Multi   []*Reply
	}

Reply.Type is defined as:

	const (
		ErrorReply = iota
		StatusReply
		IntegerReply
		BulkReply
		MultiReply
	)

Reply struct has many useful methods:

	func (rp *Reply) IntegerValue() (int64, error)
	func (rp *Reply) BoolValue() (bool, error)
	func (rp *Reply) StatusValue() (string, error)
	func (rp *Reply) OKValue() error
	func (rp *Reply) BytesValue() ([]byte, error)
	func (rp *Reply) StringValue() (string, error)
	func (rp *Reply) MultiValue() ([]*Reply, error)
	func (rp *Reply) HashValue() (map[string]string, error)
	func (rp *Reply) ListValue() ([]string, error)
	func (rp *Reply) BytesArrayValue() ([][]byte, error)
	func (rp *Reply) BoolArrayValue() ([]bool, error)

You can find more examples in test files.


Run Test
--------

normal test:

	go test

coverage test:

	go test -cover

coverage test with html result:

	go test -coverprofile=cover.out
	go tool cover -html=cover.out

Welcome to report issues :)


Run Benchmark
-------------

	go test -test.run=none -test.bench="Benchmark.*"

At my virtualbox Ubuntu 13.04 with single CPU: Intel(R) Core(TM) i5-3450 CPU @ 3.10GHz, get result:

	BenchmarkPing	   50000	     40100 ns/op
	BenchmarkLPush	   50000	     34939 ns/op
	BenchmarkLRange	   50000	     41420 ns/op
	BenchmarkGet	   50000	     37948 ns/op
	BenchmarkIncr	   50000	     44460 ns/op
	BenchmarkSet	   50000	     41300 ns/op

Welcome to show your benchmark result :)


License
-------

[The MIT License (MIT) Copyright (c) 2013 xuyu](http://opensource.org/licenses/MIT)

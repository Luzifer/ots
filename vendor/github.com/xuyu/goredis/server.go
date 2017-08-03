package goredis

import (
	"errors"
	"io"
	"net"
	"strconv"
)

// BgRewriteAof Instruct Redis to start an Append Only File rewrite process.
// The rewrite will create a small optimized version of the current Append Only File.
func (r *Redis) BgRewriteAof() error {
	_, err := r.ExecuteCommand("BGREWRITEAOF")
	return err
}

// BgSave save the DB in background.
// The OK code is immediately returned.
// Redis forks, the parent continues to serve the clients, the child saves the DB on disk then exits.
// A client my be able to check if the operation succeeded using the LASTSAVE command.
func (r *Redis) BgSave() error {
	_, err := r.ExecuteCommand("BGSAVE")
	return err
}

// ClientKill closes a given client connection identified by ip:port.
// Due to the single-treaded nature of Redis,
// it is not possible to kill a client connection while it is executing a command.
// However, the client will notice the connection has been closed
// only when the next command is sent (and results in network error).
// Status code reply: OK if the connection exists and has been closed
func (r *Redis) ClientKill(ip string, port int) error {
	rp, err := r.ExecuteCommand("CLIENT", "KILL", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ClientList returns information and statistics
// about the client connections server in a mostly human readable format.
// Bulk reply: a unique string, formatted as follows:
// One client connection per line (separated by LF)
// Each line is composed of a succession of property=value fields separated by a space character.
func (r *Redis) ClientList() (string, error) {
	rp, err := r.ExecuteCommand("CLIENT", "LIST")
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// ClientGetName returns the name of the current connection as set by CLIENT SETNAME.
// Since every new connection starts without an associated name,
// if no name was assigned a null bulk reply is returned.
func (r *Redis) ClientGetName() ([]byte, error) {
	rp, err := r.ExecuteCommand("CLIENT", "GETNAME")
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// ClientPause stops the server processing commands from clients for some time.
func (r *Redis) ClientPause(timeout uint64) error {
	rp, err := r.ExecuteCommand("CLIENT", "PAUSE", timeout)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ClientSetName assigns a name to the current connection.
func (r *Redis) ClientSetName(name string) error {
	rp, err := r.ExecuteCommand("CLIENT", "SETNAME", name)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigGet is used to read the configuration parameters of a running Redis server.
// Not all the configuration parameters are supported in Redis 2.4,
// while Redis 2.6 can read the whole configuration of a server using this command.
// CONFIG GET takes a single argument, which is a glob-style pattern.
func (r *Redis) ConfigGet(parameter string) (map[string]string, error) {
	rp, err := r.ExecuteCommand("CONFIG", "GET", parameter)
	if err != nil {
		return nil, err
	}
	return rp.HashValue()
}

// ConfigRewrite rewrites the redis.conf file the server was started with,
// applying the minimal changes needed to make it reflecting the configuration currently used by the server,
// that may be different compared to the original one because of the use of the CONFIG SET command.
// Available since 2.8.0.
func (r *Redis) ConfigRewrite() error {
	rp, err := r.ExecuteCommand("CONFIG", "REWRITE")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigSet is used in order to reconfigure the server at run time without the need to restart Redis.
// You can change both trivial parameters or switch from one to another persistence option using this command.
func (r *Redis) ConfigSet(parameter, value string) error {
	rp, err := r.ExecuteCommand("CONFIG", "SET", parameter, value)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigResetStat resets the statistics reported by Redis using the INFO command.
// These are the counters that are reset:
// Keyspace hits
// Keyspace misses
// Number of commands processed
// Number of connections received
// Number of expired keys
// Number of rejected connections
// Latest fork(2) time
// The aof_delayed_fsync counter
func (r *Redis) ConfigResetStat() error {
	_, err := r.ExecuteCommand("CONFIG", "RESETSTAT")
	return err
}

// DBSize return the number of keys in the currently-selected database.
func (r *Redis) DBSize() (int64, error) {
	rp, err := r.ExecuteCommand("DBSIZE")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// DebugObject is a debugging command that should not be used by clients.
func (r *Redis) DebugObject(key string) (string, error) {
	rp, err := r.ExecuteCommand("DEBUG", "OBJECT", key)
	if err != nil {
		return "", err
	}
	return rp.StatusValue()
}

// FlushAll delete all the keys of all the existing databases,
// not just the currently selected one.
// This command never fails.
func (r *Redis) FlushAll() error {
	_, err := r.ExecuteCommand("FLUSHALL")
	return err
}

// FlushDB delete all the keys of the currently selected DB.
// This command never fails.
func (r *Redis) FlushDB() error {
	_, err := r.ExecuteCommand("FLUSHDB")
	return err
}

// Info returns information and statistics about the server
// in a format that is simple to parse by computers and easy to read by humans.
// format document at http://redis.io/commands/info
func (r *Redis) Info(section string) (string, error) {
	args := packArgs("INFO", section)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// LastSave return the UNIX TIME of the last DB save executed with success.
// A client may check if a BGSAVE command succeeded reading the LASTSAVE value,
// then issuing a BGSAVE command and checking at regular intervals every N seconds if LASTSAVE changed.
// Integer reply: an UNIX time stamp.
func (r *Redis) LastSave() (int64, error) {
	rp, err := r.ExecuteCommand("LASTSAVE")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// MonitorCommand is a debugging command that streams back every command processed by the Redis server.
type MonitorCommand struct {
	redis *Redis
	conn  *connection
}

// Monitor sned MONITOR command to redis server.
func (r *Redis) Monitor() (*MonitorCommand, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	if err := c.SendCommand("MONITOR"); err != nil {
		return nil, err
	}
	rp, err := c.RecvReply()
	if err != nil {
		return nil, err
	}
	if err := rp.OKValue(); err != nil {
		return nil, err
	}
	return &MonitorCommand{r, c}, nil
}

// Receive read from redis server and return the reply.
func (m *MonitorCommand) Receive() (string, error) {
	rp, err := m.conn.RecvReply()
	if err != nil {
		return "", err
	}
	return rp.StatusValue()
}

// Close closes current monitor command.
func (m *MonitorCommand) Close() error {
	return m.conn.SendCommand("QUIT")
}

// Save performs a synchronous save of the dataset
// producing a point in time snapshot of all the data inside the Redis instance,
// in the form of an RDB file.
// You almost never want to call SAVE in production environments
// where it will block all the other clients. Instead usually BGSAVE is used.
func (r *Redis) Save() error {
	rp, err := r.ExecuteCommand("SAVE")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Shutdown behavior is the following:
// Stop all the clients.
// Perform a blocking SAVE if at least one save point is configured.
// Flush the Append Only File if AOF is enabled.
// Quit the server.
func (r *Redis) Shutdown(save, noSave bool) error {
	args := packArgs("SHUTDOWN")
	if save {
		args = append(args, "SAVE")
	} else if noSave {
		args = append(args, "NOSAVE")
	}
	rp, err := r.ExecuteCommand(args...)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New(rp.Status)
}

// SlaveOf can change the replication settings of a slave on the fly.
// If a Redis server is already acting as slave, the command SLAVEOF NO ONE will turn off the replication,
// turning the Redis server into a MASTER.
// In the proper form SLAVEOF hostname port will make the server a slave of
// another server listening at the specified hostname and port.
//
// If a server is already a slave of some master,
// SLAVEOF hostname port will stop the replication against the old server
// and start the synchronization against the new one, discarding the old dataset.
// The form SLAVEOF NO ONE will stop replication, turning the server into a MASTER,
// but will not discard the replication.
// So, if the old master stops working,
// it is possible to turn the slave into a master and set the application to use this new master in read/write.
// Later when the other Redis server is fixed, it can be reconfigured to work as a slave.
func (r *Redis) SlaveOf(host, port string) error {
	rp, err := r.ExecuteCommand("SLAVEOF", host, port)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// SlowLog is used in order to read and reset the Redis slow queries log.
type SlowLog struct {
	ID           int64
	Timestamp    int64
	Microseconds int64
	Command      []string
}

// SlowLogGet returns slow logs.
func (r *Redis) SlowLogGet(n int) ([]*SlowLog, error) {
	rp, err := r.ExecuteCommand("SLOWLOG", "GET", n)
	if err != nil {
		return nil, err
	}
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("slowlog get protocol error")
	}
	var slow []*SlowLog
	for _, subrp := range rp.Multi {
		if subrp.Multi == nil || len(subrp.Multi) != 4 {
			return nil, errors.New("slowlog get protocol error")
		}
		id, err := subrp.Multi[0].IntegerValue()
		if err != nil {
			return nil, err
		}
		timestamp, err := subrp.Multi[1].IntegerValue()
		if err != nil {
			return nil, err
		}
		microseconds, err := subrp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		command, err := subrp.Multi[3].ListValue()
		if err != nil {
			return nil, err
		}
		slow = append(slow, &SlowLog{id, timestamp, microseconds, command})
	}
	return slow, nil
}

// SlowLogLen Obtaining the current length of the slow log
func (r *Redis) SlowLogLen() (int64, error) {
	rp, err := r.ExecuteCommand("SLOWLOG", "LEN")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SlowLogReset resetting the slow log.
// Once deleted the information is lost forever.
func (r *Redis) SlowLogReset() error {
	rp, err := r.ExecuteCommand("SLOWLOG", "RESET")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Time returns a multi bulk reply containing two elements:
// unix time in seconds,
// microseconds.
func (r *Redis) Time() ([]string, error) {
	rp, err := r.ExecuteCommand("TIME")
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

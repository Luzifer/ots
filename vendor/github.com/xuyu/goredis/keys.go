package goredis

import (
	"errors"
	"strconv"
)

// Del removes the specified keys.
// A key is ignored if it does not exist.
// Integer reply: The number of keys that were removed.
func (r *Redis) Del(keys ...string) (int64, error) {
	args := packArgs("DEL", keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Dump serialize the value stored at key in a Redis-specific format and return it to the user.
// The returned value can be synthesized back into a Redis key using the RESTORE command.
// Return []byte for maybe big data
func (r *Redis) Dump(key string) ([]byte, error) {
	rp, err := r.ExecuteCommand("DUMP", key)
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// Exists returns true if key exists.
func (r *Redis) Exists(key string) (bool, error) {
	rp, err := r.ExecuteCommand("EXISTS", key)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// Expire set a second timeout on key.
// After the timeout has expired, the key will automatically be deleted.
// A key with an associated timeout is often said to be volatile in Redis terminology.
func (r *Redis) Expire(key string, seconds int) (bool, error) {
	rp, err := r.ExecuteCommand("EXPIRE", key, seconds)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// ExpireAt has the same effect and semantic as expire,
// but instead of specifying the number of seconds representing the TTL (time to live),
// it takes an absolute Unix timestamp (seconds since January 1, 1970).
func (r *Redis) ExpireAt(key string, timestamp int64) (bool, error) {
	rp, err := r.ExecuteCommand("EXPIREAT", key, timestamp)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// Keys returns all keys matching pattern.
func (r *Redis) Keys(pattern string) ([]string, error) {
	rp, err := r.ExecuteCommand("KEYS", pattern)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// Atomically transfer a key from a source Redis instance to a destination Redis instance.
// On success the key is deleted from the original instance and is guaranteed to exist in the target instance.
//
// The command is atomic and blocks the two instances for the time required to transfer the key,
// at any given time the key will appear to exist in a given instance or in the other instance,
// unless a timeout error occurs.
//
// The timeout specifies the maximum idle time in any moment of the communication
// with the destination instance in milliseconds.
//
// COPY -- Do not remove the key from the local instance.
// REPLACE -- Replace existing key on the remote instance.
//
// Status code reply: The command returns OK on success.
// MIGRATE host port key destination-db timeout [COPY] [REPLACE]
//
// func (r *Redis) Migrate(host, port, key string, db, timeout int, cp, replace bool) error {
// 	args := packArgs("MIGRATE", host, port, key, db, timeout)
// 	if cp {
// 		args = append(args, "COPY")
// 	}
// 	if replace {
// 		args = append(args, "REPLACE")
// 	}
// 	rp, err := r.ExecuteCommand(args...)
// 	if err != nil {
// 		return err
// 	}
// 	return rp.OKValue()
// }

// Move moves key from the currently selected database (see SELECT)
// to the specified destination database.
// When key already exists in the destination database,
// or it does not exist in the source database, it does nothing.
func (r *Redis) Move(key string, db int) (bool, error) {
	rp, err := r.ExecuteCommand("MOVE", key, db)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// Object inspects the internals of Redis Objects associated with keys.
// It is useful for debugging or to understand if your keys are using the specially encoded data types to save space.
// Your application may also use the information reported by the OBJECT command
// to implement application level key eviction policies
// when using Redis as a Cache.
func (r *Redis) Object(subcommand string, arguments ...string) (*Reply, error) {
	args := packArgs("OBJECT", subcommand, arguments)
	return r.ExecuteCommand(args...)
}

// Persist removes the existing timeout on key,
// turning the key from volatile (a key with an expire set) to persistent
// (a key that will never expire as no timeout is associated).
// True if the timeout was removed.
// False if key does not exist or does not have an associated timeout.
func (r *Redis) Persist(key string) (bool, error) {
	rp, err := r.ExecuteCommand("PERSIST", key)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// PExpire works exactly like EXPIRE
// but the time to live of the key is specified in milliseconds instead of seconds.
func (r *Redis) PExpire(key string, milliseconds int) (bool, error) {
	rp, err := r.ExecuteCommand("PEXPIRE", key, milliseconds)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// PExpireAt has the same effect and semantic as EXPIREAT,
// but the Unix time at which the key will expire is specified in milliseconds instead of seconds.
func (r *Redis) PExpireAt(key string, timestamp int64) (bool, error) {
	rp, err := r.ExecuteCommand("PEXPIREAT", key, timestamp)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// PTTL returns the remaining time to live of a key that has an expire set,
// with the sole difference that TTL returns the amount of remaining time in seconds
// while PTTL returns it in milliseconds.
func (r *Redis) PTTL(key string) (int64, error) {
	rp, err := r.ExecuteCommand("PTTL", key)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// RandomKey returns a random key from the currently selected database.
// Bulk reply: the random key, or nil when the database is empty.
func (r *Redis) RandomKey() ([]byte, error) {
	rp, err := r.ExecuteCommand("RANDOMKEY")
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// Rename renames key to newkey.
// It returns an error when the source and destination names are the same, or when key does not exist.
// If newkey already exists it is overwritten, when this happens RENAME executes an implicit DEL operation,
// so if the deleted key contains a very big value it may cause high latency
// even if RENAME itself is usually a constant-time operation.
func (r *Redis) Rename(key, newkey string) error {
	rp, err := r.ExecuteCommand("RENAME", key, newkey)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Renamenx renames key to newkey if newkey does not yet exist.
// It returns an error under the same conditions as RENAME.
func (r *Redis) Renamenx(key, newkey string) (bool, error) {
	rp, err := r.ExecuteCommand("RENAMENX", key, newkey)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// Restore creates a key associated with a value that is obtained by deserializing
// the provided serialized value (obtained via DUMP).
// If ttl is 0 the key is created without any expire, otherwise the specified expire time (in milliseconds) is set.
// RESTORE checks the RDB version and data checksum. If they don't match an error is returned.
func (r *Redis) Restore(key string, ttl int, serialized string) error {
	rp, err := r.ExecuteCommand("RESTORE", key, ttl, serialized)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// TTL returns the remaining time to live of a key that has a timeout.
// Integer reply: TTL in seconds, or a negative value in order to signal an error (see the description above).
func (r *Redis) TTL(key string) (int64, error) {
	rp, err := r.ExecuteCommand("TTL", key)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Type returns the string representation of the type of the value stored at key.
// The different types that can be returned are: string, list, set, zset and hash.
// Status code reply: type of key, or none when key does not exist.
func (r *Redis) Type(key string) (string, error) {
	rp, err := r.ExecuteCommand("TYPE", key)
	if err != nil {
		return "", err
	}
	return rp.StatusValue()
}

// Scan command:
// SCAN cursor [MATCH pattern] [COUNT count]
func (r *Redis) Scan(cursor uint64, pattern string, count int) (uint64, []string, error) {
	args := packArgs("SCAN", cursor)
	if pattern != "" {
		args = append(args, "MATCH", pattern)
	}
	if count > 0 {
		args = append(args, "COUNT", count)
	}
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, nil, err
	}
	if rp.Type == ErrorReply {
		return 0, nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return 0, nil, errors.New("scan protocol error")
	}
	first, err := rp.Multi[0].StringValue()
	if err != nil {
		return 0, nil, err
	}
	next, err := strconv.ParseUint(first, 10, 64)
	if err != nil {
		return 0, nil, err
	}
	list, err := rp.Multi[1].ListValue()
	return next, list, err
}

package goredis

import (
	"strconv"
)

// SAdd add the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored.
// If key does not exist, a new set is created before adding the specified members.
// An error is returned when the value stored at key is not a set.
//
// Integer reply: the number of elements that were added to the set,
// not including all the elements already present into the set.
func (r *Redis) SAdd(key string, members ...string) (int64, error) {
	args := packArgs("SADD", key, members)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SCard returns the set cardinality (number of elements) of the set stored at key.
func (r *Redis) SCard(key string) (int64, error) {
	rp, err := r.ExecuteCommand("SCARD", key)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SDiff returns the members of the set resulting from the difference
// between the first set and all the successive sets.
// Keys that do not exist are considered to be empty sets.
// Multi-bulk reply: list with members of the resulting set.
func (r *Redis) SDiff(keys ...string) ([]string, error) {
	args := packArgs("SDIFF", keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// SDiffStore is equal to SDIFF, but instead of returning the resulting set,
// it is stored in destination.
// If destination already exists, it is overwritten.
// Integer reply: the number of elements in the resulting set.
func (r *Redis) SDiffStore(destination string, keys ...string) (int64, error) {
	args := packArgs("SDIFFSTORE", destination, keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SInter returns the members of the set resulting from the intersection of all the given sets.
// Multi-bulk reply: list with members of the resulting set.
func (r *Redis) SInter(keys ...string) ([]string, error) {
	args := packArgs("SINTER", keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// SInterStore is equal to SINTER, but instead of returning the resulting set,
// it is stored in destination.
// If destination already exists, it is overwritten.
// Integer reply: the number of elements in the resulting set.
func (r *Redis) SInterStore(destination string, keys ...string) (int64, error) {
	args := packArgs("SINTERSTORE", destination, keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SIsMember returns if member is a member of the set stored at key.
func (r *Redis) SIsMember(key, member string) (bool, error) {
	rp, err := r.ExecuteCommand("SISMEMBER", key, member)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// SMembers returns all the members of the set value stored at key.
func (r *Redis) SMembers(key string) ([]string, error) {
	rp, err := r.ExecuteCommand("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// SMove moves member from the set at source to the set at destination.
// This operation is atomic.
// In every given moment the element will appear to be a member of source or destination for other clients.
func (r *Redis) SMove(source, destination, member string) (bool, error) {
	rp, err := r.ExecuteCommand("SMOVE", source, destination, member)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// SPop removes and returns a random element from the set value stored at key.
// Bulk reply: the removed element, or nil when key does not exist.
func (r *Redis) SPop(key string) ([]byte, error) {
	rp, err := r.ExecuteCommand("SPOP", key)
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// SRandMember returns a random element from the set value stored at key.
// Bulk reply: the command returns a Bulk Reply with the randomly selected element,
// or nil when key does not exist.
func (r *Redis) SRandMember(key string) ([]byte, error) {
	rp, err := r.ExecuteCommand("SRANDMEMBER", key)
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// SRandMemberCount returns an array of count distinct elements if count is positive.
// If called with a negative count the behavior changes and the command
// is allowed to return the same element multiple times.
// In this case the numer of returned elements is the absolute value of the specified count.
// returns an array of elements, or an empty array when key does not exist.
func (r *Redis) SRandMemberCount(key string, count int) ([]string, error) {
	rp, err := r.ExecuteCommand("SRANDMEMBER", key, count)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// SRem remove the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
// An error is returned when the value stored at key is not a set.
// Integer reply: the number of members that were removed from the set,
// not including non existing members.
func (r *Redis) SRem(key string, members ...string) (int64, error) {
	args := packArgs("SREM", key, members)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SUnion returns the members of the set resulting from the union of all the given sets.
// Multi-bulk reply: list with members of the resulting set.
func (r *Redis) SUnion(keys ...string) ([]string, error) {
	args := packArgs("SUNION", keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// SUnionStore is equal to SUnion.
// If destination already exists, it is overwritten.
// Integer reply: the number of elements in the resulting set.
func (r *Redis) SUnionStore(destination string, keys ...string) (int64, error) {
	args := packArgs("SUNIONSTORE", destination, keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SScan key cursor [MATCH pattern] [COUNT count]
func (r *Redis) SScan(key string, cursor uint64, pattern string, count int) (uint64, []string, error) {
	args := packArgs("SSCAN", key, cursor)
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

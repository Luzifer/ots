package goredis

// PFAdd adds all the element arguments to the HyperLogLog data structure
// stored at the variable name specified as first argument.
func (r *Redis) PFAdd(key string, elements ...string) (int64, error) {
	args := packArgs("PFADD", key, elements)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// PFCount returns the approximated cardinality computed by the HyperLogLog
// data structure stored at the specified variable,
// which is 0 if the variable does not exist.
// When called with multiple keys, returns the approximated cardinality of
// the union of the HyperLogLogs passed, by internally merging the HyperLogLogs
// stored at the provided keys into a temporary hyperLogLog.
func (r *Redis) PFCount(keys ...string) (int64, error) {
	args := packArgs("PFCOUNT", keys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// PFMerge merges multiple HyperLogLog values into an unique value
// that will approximate the cardinality of the union of the observed
// Sets of the source HyperLogLog structures.
// The computed merged HyperLogLog is set to the destination variable,
// which is created if does not exist (defauling to an empty HyperLogLog).
func (r *Redis) PFMerge(destkey string, sourcekeys ...string) error {
	args := packArgs("PFMERGE", destkey, sourcekeys)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

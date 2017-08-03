package goredis

// ScriptExists returns information about the existence of the scripts in the script cache.
// Multi-bulk reply The command returns an array of integers
// that correspond to the specified SHA1 digest arguments.
// For every corresponding SHA1 digest of a script that actually exists in the script cache.
func (r *Redis) ScriptExists(scripts ...string) ([]bool, error) {
	args := packArgs("SCRIPT", "EXISTS", scripts)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.BoolArrayValue()
}

// ScriptFlush flush the Lua scripts cache.
// Please refer to the EVAL documentation for detailed information about Redis Lua scripting.
func (r *Redis) ScriptFlush() error {
	rp, err := r.ExecuteCommand("SCRIPT", "FLUSH")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ScriptKill kills the currently executing Lua script,
// assuming no write operation was yet performed by the script.
func (r *Redis) ScriptKill() error {
	rp, err := r.ExecuteCommand("SCRIPT", "KILL")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ScriptLoad Load a script into the scripts cache, without executing it.
// After the specified command is loaded into the script cache
// it will be callable using EVALSHA with the correct SHA1 digest of the script,
// exactly like after the first successful invocation of EVAL.
// Bulk reply This command returns the SHA1 digest of the script added into the script cache.
func (r *Redis) ScriptLoad(script string) (string, error) {
	rp, err := r.ExecuteCommand("SCRIPT", "LOAD", script)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// Eval first argument of EVAL is a Lua 5.1 script.
// The script does not need to define a Lua function (and should not).
// It is just a Lua program that will run in the context of the Redis server.
// The second argument of EVAL is the number of arguments that follows the script
// (starting from the third argument) that represent Redis key names.
// This arguments can be accessed by Lua using the KEYS global variable
// in the form of a one-based array (so KEYS[1], KEYS[2], ...).
func (r *Redis) Eval(script string, keys []string, args []string) (*Reply, error) {
	cmds := packArgs("EVAL", script, len(keys), keys, args)
	return r.ExecuteCommand(cmds...)
}

// EvalSha evaluates a script cached on the server side by its SHA1 digest.
// Scripts are cached on the server side using the SCRIPT LOAD command.
func (r *Redis) EvalSha(sha1 string, keys []string, args []string) (*Reply, error) {
	cmds := packArgs("EVALSHA", sha1, len(keys), keys, args)
	return r.ExecuteCommand(cmds...)
}

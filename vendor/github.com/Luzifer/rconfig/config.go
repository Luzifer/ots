// Package rconfig implements a CLI configuration reader with struct-embedded
// defaults, environment variables and posix compatible flag parsing using
// the pflag library.
package rconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	validator "gopkg.in/validator.v2"
)

type afterFunc func() error

var (
	autoEnv          bool
	fs               *pflag.FlagSet
	variableDefaults map[string]string

	timeParserFormats = []string{
		// Default constants
		time.RFC3339Nano, time.RFC3339,
		time.RFC1123Z, time.RFC1123,
		time.RFC822Z, time.RFC822,
		time.RFC850, time.RubyDate, time.UnixDate, time.ANSIC,
		"2006-01-02 15:04:05.999999999 -0700 MST",
		// More uncommon time formats
		"2006-01-02 15:04:05", "2006-01-02 15:04:05Z07:00", // Simplified ISO time format
		"01/02/2006 15:04:05", "01/02/2006 15:04:05Z07:00", // US time format
		"02.01.2006 15:04:05", "02.01.2006 15:04:05Z07:00", // DE time format
	}
)

func init() {
	variableDefaults = make(map[string]string)
}

// Parse takes the pointer to a struct filled with variables which should be read
// from ENV, default or flag. The precedence in this is flag > ENV > default. So
// if a flag is specified on the CLI it will overwrite the ENV and otherwise ENV
// overwrites the default specified.
//
// For your configuration struct you can use the following struct-tags to control
// the behavior of rconfig:
//
//     default: Set a default value
//     vardefault: Read the default value from the variable defaults
//     env: Read the value from this environment variable
//     flag: Flag to read in format "long,short" (for example "listen,l")
//     description: A help text for Usage output to guide your users
//
// The format you need to specify those values you can see in the example to this
// function.
//
func Parse(config interface{}) error {
	return parse(config, nil)
}

// ParseAndValidate works exactly like Parse but implements an additional run of
// the go-validator package on the configuration struct. Therefore additonal struct
// tags are supported like described in the readme file of the go-validator package:
//
// https://github.com/go-validator/validator/tree/v2#usage
func ParseAndValidate(config interface{}) error {
	return parseAndValidate(config, nil)
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return fs.Args()
}

// AddTimeParserFormats adds custom formats to parse time.Time fields
func AddTimeParserFormats(f ...string) {
	timeParserFormats = append(timeParserFormats, f...)
}

// AutoEnv enables or disables automated env variable guessing. If no `env` struct
// tag was set and AutoEnv is enabled the env variable name is derived from the
// name of the field: `MyFieldName` will get `MY_FIELD_NAME`
func AutoEnv(enable bool) {
	autoEnv = enable
}

// Usage prints a basic usage with the corresponding defaults for the flags to
// os.Stdout. The defaults are derived from the `default` struct-tag and the ENV.
func Usage() {
	if fs != nil && fs.Parsed() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
	}
}

// SetVariableDefaults presets the parser with a map of default values to be used
// when specifying the vardefault tag
func SetVariableDefaults(defaults map[string]string) {
	variableDefaults = defaults
}

func parseAndValidate(in interface{}, args []string) error {
	if err := parse(in, args); err != nil {
		return err
	}

	return validator.Validate(in)
}

func parse(in interface{}, args []string) error {
	if args == nil {
		args = os.Args
	}

	fs = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	afterFuncs, err := execTags(in, fs)
	if err != nil {
		return err
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if afterFuncs != nil {
		for _, f := range afterFuncs {
			if err := f(); err != nil {
				return err
			}
		}
	}

	return nil
}

func execTags(in interface{}, fs *pflag.FlagSet) ([]afterFunc, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("Calling parser with non-pointer")
	}

	if reflect.ValueOf(in).Elem().Kind() != reflect.Struct {
		return nil, errors.New("Calling parser with pointer to non-struct")
	}

	afterFuncs := []afterFunc{}

	st := reflect.ValueOf(in).Elem()
	for i := 0; i < st.NumField(); i++ {
		valField := st.Field(i)
		typeField := st.Type().Field(i)

		if typeField.Tag.Get("default") == "" && typeField.Tag.Get("env") == "" && typeField.Tag.Get("flag") == "" && typeField.Type.Kind() != reflect.Struct {
			// None of our supported tags is present and it's not a sub-struct
			continue
		}

		value := varDefault(typeField.Tag.Get("vardefault"), typeField.Tag.Get("default"))
		value = envDefault(typeField, value)
		parts := strings.Split(typeField.Tag.Get("flag"), ",")

		switch typeField.Type {
		case reflect.TypeOf(time.Duration(0)):
			v, err := time.ParseDuration(value)
			if err != nil {
				if value == "" {
					v = time.Duration(0)
				} else {
					return nil, err
				}
			}

			if typeField.Tag.Get("flag") != "" {
				if len(parts) == 1 {
					fs.DurationVar(valField.Addr().Interface().(*time.Duration), parts[0], v, typeField.Tag.Get("description"))
				} else {
					fs.DurationVarP(valField.Addr().Interface().(*time.Duration), parts[0], parts[1], v, typeField.Tag.Get("description"))
				}
			} else {
				valField.Set(reflect.ValueOf(v))
			}
			continue

		case reflect.TypeOf(time.Time{}):
			var sVar string

			if typeField.Tag.Get("flag") != "" {
				if len(parts) == 1 {
					fs.StringVar(&sVar, parts[0], value, typeField.Tag.Get("description"))
				} else {
					fs.StringVarP(&sVar, parts[0], parts[1], value, typeField.Tag.Get("description"))
				}
			} else {
				sVar = value
			}

			afterFuncs = append(afterFuncs, func(valField reflect.Value, sVar *string) func() error {
				return func() error {
					if *sVar == "" {
						// No time, no problem
						return nil
					}

					// Check whether we could have a timestamp
					if ts, err := strconv.ParseInt(*sVar, 10, 64); err == nil {
						t := time.Unix(ts, 0)
						valField.Set(reflect.ValueOf(t))
						return nil
					}

					// We haven't so lets walk through possible time formats
					matched := false
					for _, tf := range timeParserFormats {
						if t, err := time.Parse(tf, *sVar); err == nil {
							matched = true
							valField.Set(reflect.ValueOf(t))
							return nil
						}
					}

					if !matched {
						return fmt.Errorf("Value %q did not match expected time formats", *sVar)
					}

					return nil
				}
			}(valField, &sVar))

			continue
		}

		switch typeField.Type.Kind() {
		case reflect.String:
			if typeField.Tag.Get("flag") != "" {
				if len(parts) == 1 {
					fs.StringVar(valField.Addr().Interface().(*string), parts[0], value, typeField.Tag.Get("description"))
				} else {
					fs.StringVarP(valField.Addr().Interface().(*string), parts[0], parts[1], value, typeField.Tag.Get("description"))
				}
			} else {
				valField.SetString(value)
			}

		case reflect.Bool:
			v := value == "true"
			if typeField.Tag.Get("flag") != "" {
				if len(parts) == 1 {
					fs.BoolVar(valField.Addr().Interface().(*bool), parts[0], v, typeField.Tag.Get("description"))
				} else {
					fs.BoolVarP(valField.Addr().Interface().(*bool), parts[0], parts[1], v, typeField.Tag.Get("description"))
				}
			} else {
				valField.SetBool(v)
			}

		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			vt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				if value == "" {
					vt = 0
				} else {
					return nil, err
				}
			}
			if typeField.Tag.Get("flag") != "" {
				registerFlagInt(typeField.Type.Kind(), fs, valField.Addr().Interface(), parts, vt, typeField.Tag.Get("description"))
			} else {
				valField.SetInt(vt)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vt, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				if value == "" {
					vt = 0
				} else {
					return nil, err
				}
			}
			if typeField.Tag.Get("flag") != "" {
				registerFlagUint(typeField.Type.Kind(), fs, valField.Addr().Interface(), parts, vt, typeField.Tag.Get("description"))
			} else {
				valField.SetUint(vt)
			}

		case reflect.Float32, reflect.Float64:
			vt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				if value == "" {
					vt = 0.0
				} else {
					return nil, err
				}
			}
			if typeField.Tag.Get("flag") != "" {
				registerFlagFloat(typeField.Type.Kind(), fs, valField.Addr().Interface(), parts, vt, typeField.Tag.Get("description"))
			} else {
				valField.SetFloat(vt)
			}

		case reflect.Struct:
			afs, err := execTags(valField.Addr().Interface(), fs)
			if err != nil {
				return nil, err
			}
			afterFuncs = append(afterFuncs, afs...)

		case reflect.Slice:
			switch typeField.Type.Elem().Kind() {
			case reflect.Int:
				def := []int{}
				for _, v := range strings.Split(value, ",") {
					it, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
					if err != nil {
						return nil, err
					}
					def = append(def, int(it))
				}
				if len(parts) == 1 {
					fs.IntSliceVar(valField.Addr().Interface().(*[]int), parts[0], def, typeField.Tag.Get("description"))
				} else {
					fs.IntSliceVarP(valField.Addr().Interface().(*[]int), parts[0], parts[1], def, typeField.Tag.Get("description"))
				}
			case reflect.String:
				del := typeField.Tag.Get("delimiter")
				if len(del) == 0 {
					del = ","
				}
				var def = []string{}
				if value != "" {
					def = strings.Split(value, del)
				}
				if len(parts) == 1 {
					fs.StringSliceVar(valField.Addr().Interface().(*[]string), parts[0], def, typeField.Tag.Get("description"))
				} else {
					fs.StringSliceVarP(valField.Addr().Interface().(*[]string), parts[0], parts[1], def, typeField.Tag.Get("description"))
				}
			}
		}
	}

	return afterFuncs, nil
}

func registerFlagFloat(t reflect.Kind, fs *pflag.FlagSet, field interface{}, parts []string, vt float64, desc string) {
	switch t {
	case reflect.Float32:
		if len(parts) == 1 {
			fs.Float32Var(field.(*float32), parts[0], float32(vt), desc)
		} else {
			fs.Float32VarP(field.(*float32), parts[0], parts[1], float32(vt), desc)
		}
	case reflect.Float64:
		if len(parts) == 1 {
			fs.Float64Var(field.(*float64), parts[0], float64(vt), desc)
		} else {
			fs.Float64VarP(field.(*float64), parts[0], parts[1], float64(vt), desc)
		}
	}
}

func registerFlagInt(t reflect.Kind, fs *pflag.FlagSet, field interface{}, parts []string, vt int64, desc string) {
	switch t {
	case reflect.Int:
		if len(parts) == 1 {
			fs.IntVar(field.(*int), parts[0], int(vt), desc)
		} else {
			fs.IntVarP(field.(*int), parts[0], parts[1], int(vt), desc)
		}
	case reflect.Int8:
		if len(parts) == 1 {
			fs.Int8Var(field.(*int8), parts[0], int8(vt), desc)
		} else {
			fs.Int8VarP(field.(*int8), parts[0], parts[1], int8(vt), desc)
		}
	case reflect.Int32:
		if len(parts) == 1 {
			fs.Int32Var(field.(*int32), parts[0], int32(vt), desc)
		} else {
			fs.Int32VarP(field.(*int32), parts[0], parts[1], int32(vt), desc)
		}
	case reflect.Int64:
		if len(parts) == 1 {
			fs.Int64Var(field.(*int64), parts[0], int64(vt), desc)
		} else {
			fs.Int64VarP(field.(*int64), parts[0], parts[1], int64(vt), desc)
		}
	}
}

func registerFlagUint(t reflect.Kind, fs *pflag.FlagSet, field interface{}, parts []string, vt uint64, desc string) {
	switch t {
	case reflect.Uint:
		if len(parts) == 1 {
			fs.UintVar(field.(*uint), parts[0], uint(vt), desc)
		} else {
			fs.UintVarP(field.(*uint), parts[0], parts[1], uint(vt), desc)
		}
	case reflect.Uint8:
		if len(parts) == 1 {
			fs.Uint8Var(field.(*uint8), parts[0], uint8(vt), desc)
		} else {
			fs.Uint8VarP(field.(*uint8), parts[0], parts[1], uint8(vt), desc)
		}
	case reflect.Uint16:
		if len(parts) == 1 {
			fs.Uint16Var(field.(*uint16), parts[0], uint16(vt), desc)
		} else {
			fs.Uint16VarP(field.(*uint16), parts[0], parts[1], uint16(vt), desc)
		}
	case reflect.Uint32:
		if len(parts) == 1 {
			fs.Uint32Var(field.(*uint32), parts[0], uint32(vt), desc)
		} else {
			fs.Uint32VarP(field.(*uint32), parts[0], parts[1], uint32(vt), desc)
		}
	case reflect.Uint64:
		if len(parts) == 1 {
			fs.Uint64Var(field.(*uint64), parts[0], uint64(vt), desc)
		} else {
			fs.Uint64VarP(field.(*uint64), parts[0], parts[1], uint64(vt), desc)
		}
	}
}

func envDefault(field reflect.StructField, def string) string {
	value := def

	env := field.Tag.Get("env")
	if env == "" && autoEnv {
		env = deriveEnvVarName(field.Name)
	}

	if env != "" {
		if e := os.Getenv(env); e != "" {
			value = e
		}
	}

	return value
}

func varDefault(name, def string) string {
	value := def

	if name != "" {
		if v, ok := variableDefaults[name]; ok {
			value = v
		}
	}

	return value
}

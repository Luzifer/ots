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

var (
	fs               *pflag.FlagSet
	variableDefaults map[string]string
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
	if err := execTags(in, fs); err != nil {
		return err
	}

	return fs.Parse(args)
}

func execTags(in interface{}, fs *pflag.FlagSet) error {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return errors.New("Calling parser with non-pointer")
	}

	if reflect.ValueOf(in).Elem().Kind() != reflect.Struct {
		return errors.New("Calling parser with pointer to non-struct")
	}

	st := reflect.ValueOf(in).Elem()
	for i := 0; i < st.NumField(); i++ {
		valField := st.Field(i)
		typeField := st.Type().Field(i)

		if typeField.Tag.Get("default") == "" && typeField.Tag.Get("env") == "" && typeField.Tag.Get("flag") == "" && typeField.Type.Kind() != reflect.Struct {
			// None of our supported tags is present and it's not a sub-struct
			continue
		}

		value := varDefault(typeField.Tag.Get("vardefault"), typeField.Tag.Get("default"))
		value = envDefault(typeField.Tag.Get("env"), value)
		parts := strings.Split(typeField.Tag.Get("flag"), ",")

		switch typeField.Type {
		case reflect.TypeOf(time.Duration(0)):
			v, err := time.ParseDuration(value)
			if err != nil {
				if value == "" {
					v = time.Duration(0)
				} else {
					return err
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
					return err
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
					return err
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
					return err
				}
			}
			if typeField.Tag.Get("flag") != "" {
				registerFlagFloat(typeField.Type.Kind(), fs, valField.Addr().Interface(), parts, vt, typeField.Tag.Get("description"))
			} else {
				valField.SetFloat(vt)
			}

		case reflect.Struct:
			if err := execTags(valField.Addr().Interface(), fs); err != nil {
				return err
			}

		case reflect.Slice:
			switch typeField.Type.Elem().Kind() {
			case reflect.Int:
				def := []int{}
				for _, v := range strings.Split(value, ",") {
					it, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
					if err != nil {
						return err
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
				def := strings.Split(value, del)
				if len(parts) == 1 {
					fs.StringSliceVar(valField.Addr().Interface().(*[]string), parts[0], def, typeField.Tag.Get("description"))
				} else {
					fs.StringSliceVarP(valField.Addr().Interface().(*[]string), parts[0], parts[1], def, typeField.Tag.Get("description"))
				}
			}
		}
	}

	return nil
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

func envDefault(env, def string) string {
	value := def

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

package cliconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Opts are options used for mapping viper config keys and env vars to cli flags.
type Opts struct {
	// ViperPrefix adds a prefix to viper config key and env var.
	// Seperated is a - for config value and _ for env var.
	// No prefix is added when not set.
	ViperPrefix string
	// EnvPrefix sets a prefix just for the env lookup,
	// 	without modifying other lookup sources.
	// No prefix is added when empty.
	// When used with ViperPrefix, the EnvPrefix will be prefixed to ViperPrefix.
	EnvPrefix string
	// Vipers is a list of viper instances used for config lookup.
	// The global viper instance is used when not set.
	Vipers []*viper.Viper
}

// withDefaults populates Opts with default values
func (p *Opts) withDefaults() {
	if len(p.Vipers) == 0 {
		p.Vipers = []*viper.Viper{viper.GetViper()}
	}
}

// SetFlags registers flags from struct tags using `arg:"name"`.
// The str arg can be either a struct or a pointer to a struct.
// Must be called before the Populate function is called.
func SetFlags(flags *flag.FlagSet, str interface{}) (err error) {
	// incase str is a pointer to struct, get indirect
	indirect := reflect.Indirect(reflect.ValueOf(str))
	if indirect.Kind() != reflect.Struct {
		return newKindError(indirect.Kind(), []reflect.Kind{reflect.Struct},
			"str arg invalid")
	}
	t := indirect.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("arg")
		// skip fields with missing tag
		if name == "" {
			continue
		}
		desc := f.Tag.Get("desc")
		short := f.Tag.Get("short")
		def := f.Tag.Get("default")
		// default to all fields being required
		required := f.Tag.Get("required") != "false"

		switch f.Type.Kind() {
		case reflect.String:
			flags.StringP(name, short, def, desc)
		case reflect.Bool:
			flags.BoolP(name, short, def == "true", desc)
		case reflect.Slice:
			sliceKind := f.Type.Elem().Kind()
			switch sliceKind {
			case reflect.String:
				defArr := strings.Split(def, ",")
				flags.StringArrayP(name, short, defArr, desc)
			default:
				return newKindError(sliceKind, []reflect.Kind{reflect.String},
					"struct's type is invalid")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			defInt := 0
			if def != "" {
				defInt, err = strconv.Atoi(def)
				if err != nil {
					return err
				}
			}
			flags.IntP(name, short, defInt, desc)
		case reflect.Struct:
			SetFlags(flags, f)
		default:
			return newKindError(f.Type.Kind(),
				[]reflect.Kind{reflect.String, reflect.Bool, reflect.Int, reflect.Slice},
				"struct field not supported")
		}
		if required {
			// Set required annoations for bash and cobra
			flags.SetAnnotation(name, cobra.BashCompOneRequiredFlag, []string{"true"})
		}
	}
	return nil
}

// Populate populates structs values by using the arg tag to pull values
// either from the cli, viper or an env var (in that order).
func Populate(fs *flag.FlagSet, ptr interface{}, opts Opts) error {
	opts.withDefaults()
	// Set flags to viper lookup value when not set by on the cli.
	// This allows validations on the flags to work as expected (e.x. Cobra's MarkFlagRequired).
	viperSetFlags(fs, opts)
	valueOf := reflect.ValueOf(ptr)
	// ensure ptr is a pointer before getting the Elem()
	if valueOf.Kind() != reflect.Ptr {
		return newKindError(valueOf.Kind(), []reflect.Kind{reflect.Ptr},
			"ptr must be pointer to struct")
	}
	elem := valueOf.Elem()
	if !elem.CanSet() || elem.Kind() != reflect.Struct {
		return newKindError(valueOf.Kind(), []reflect.Kind{reflect.Struct},
			"ptr must be pointer to struct")
	}
	t := elem.Type()
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		elemField := elem.Field(i)
		arg := typeField.Tag.Get("arg")
		if arg == "" {
			continue
		}
		switch typeField.Type.Kind() {
		case reflect.String:
			val, err := fs.GetString(arg)
			if err != nil {
				return err
			}
			elemField.SetString(val)
		case reflect.Bool:
			val, err := fs.GetBool(arg)
			if err != nil {
				return err
			}
			elemField.SetBool(val)
		case reflect.Slice:
			sliceType := typeField.Type.Elem().Kind()
			switch sliceType {
			case reflect.String:
				vals, err := fs.GetStringArray(arg)
				if err != nil {
					return err
				}
				elemField.Set(reflect.ValueOf(vals))
			default:
				return errors.Errorf("slice type of %v is not supported", sliceType)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := fs.GetInt(arg)
			if err != nil {
				return err
			}
			elemField.SetInt(int64(val))
		case reflect.Struct:
			Populate(fs, &elemField, opts)
		default:
			return errors.Errorf("type of %v is not supported", typeField.Type.Kind())
		}
	}
	return nil
}

func viperSetFlags(flags *flag.FlagSet, opts Opts) error {
	var result error
	for _, v := range opts.Vipers {
		flags.VisitAll(func(f *pflag.Flag) {
			key := f.Name
			if opts.ViperPrefix != "" {
				key = fmt.Sprintf("%s-%s", opts.ViperPrefix, key)
			}

			// bind to env var
			envVar := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
			if opts.EnvPrefix != "" {
				envVar = fmt.Sprintf("%s_%s", opts.EnvPrefix, envVar)
			}
			if err := v.BindEnv(key, envVar); err != nil {
				result = err
				return
			}
			// If flag is not set and viper has config value, set flag to viper's value.
			if !f.Changed && v.IsSet(key) {
				val := v.Get(key)
				if err := flags.Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
					result = err
					return
				}
			}
		})
	}
	return result
}

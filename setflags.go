package cliconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetFlags registers flags from struct tags using `arg:"name"`
// The str arg can be either a struct or a pointer to a struct
func SetFlags(flags *flag.FlagSet, str interface{}) (err error) {
	// incase str is a pointer to struct, get indirect
	val := reflect.Indirect(reflect.ValueOf(str))
	if val.Kind() != reflect.Struct {
		return errors.New("str must be a struct")
	}
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("arg")
		// skip fields with missing tags
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
			defArr := strings.Split(def, ",")
			flags.StringArrayP(name, short, defArr, desc)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			defInt := 0
			if def != "" {
				defInt, err = strconv.Atoi(def)
				if err != nil {
					return err
				}
			}
			flags.IntP(name, short, defInt, desc)
		default:
			fmt.Println("Skipping SetFlags for", f.Type.Kind().String())
			continue
		}
		if required {
			bashCompOneRequiredFlag := "cobra_annotation_bash_completion_one_required_flag"
			flags.SetAnnotation(name, bashCompOneRequiredFlag, []string{"true"})
		}
	}
	return nil
}

// BindViperDefaults binds each cobra flag to its associated viper configuration (config file and environment variable)
func BindViperDefaults(flags flag.FlagSet, prefix string) error {
	return BindViperDefaultsIntance(flags, viper.GetViper(), prefix)
}

// BindViperDefaultsIntance uses a viper instance as apposed to global viper from BindViperDefaults
func BindViperDefaultsIntance(flags flag.FlagSet, v *viper.Viper, prefix string) error {
	var result error
	flags.VisitAll(func(f *pflag.Flag) {
		key := f.Name
		if prefix != "" {
			key = fmt.Sprintf("%s-%s", prefix, key)
		}

		// bind to env var
		envVar := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		fmt.Printf("key %s binding to %s\n", key, envVar)
		if err := v.BindEnv(key, envVar); err != nil {
			result = err
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(key) {
			val := v.Get(key)
			if err := flags.Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				result = err
			}
		}
	})
	return result
}

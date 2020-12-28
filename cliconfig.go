package cliconfig

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	NotSettableErr error = errors.New("ptr value is not a settable struct pointer")
)

// SetFlags registers flags from struct tags using `arg:"name"`
// The str arg can be either a struct or a pointer to a struct
func SetFlags(flags *flag.FlagSet, str interface{}) (err error) {
	// incase str is a pointer to struct, get indirect
	val := reflect.Indirect(reflect.ValueOf(str))
	if val.Kind() != reflect.Struct {
		return newKindError(val.Kind(), []reflect.Kind{reflect.Struct},
			"str arg invalid")
	}
	t := val.Type()
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

// PopulateWithViper first sets any unset flags to their Viper lookup value before populating fields
func PopulateWithViper(fs *flag.FlagSet, ptr interface{}, prefix string, vs ...*viper.Viper) error {
	ViperSetFlags(fs, prefix, vs...)
	return Populate(fs, ptr)
}

// Populate populates structs values matching the flags name to the arg tag
func Populate(fs *flag.FlagSet, ptr interface{}) error {
	valueOf := reflect.ValueOf(ptr)
	if valueOf.Kind() != reflect.Ptr {
		return NotSettableErr
	}
	elem := valueOf.Elem()
	if !elem.CanSet() || elem.Kind() != reflect.Struct {
		return NotSettableErr
	}
	t := elem.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		arg := f.Tag.Get("arg")
		if arg == "" {
			continue
		}
		switch f.Type.Kind() {
		case reflect.String:
			val, err := fs.GetString(arg)
			if err != nil {
				return err
			}
			elem.Field(i).SetString(val)
		case reflect.Bool:
			val, err := fs.GetBool(arg)
			if err != nil {
				return err
			}
			elem.Field(i).SetBool(val)
		case reflect.Slice:
			sliceType := f.Type.Elem().Kind()
			switch sliceType {
			case reflect.String:
				vals, err := fs.GetStringArray(arg)
				if err != nil {
					return err
				}
				elem.Field(i).Set(reflect.ValueOf(vals))
			default:
				return errors.Errorf("slice type of %v is not supported", sliceType)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := fs.GetInt(arg)
			if err != nil {
				return err
			}
			valueOf.Elem().Field(i).SetInt(int64(val))
		default:
			return errors.Errorf("type of %v is not supported", f.Type.Kind())
		}
	}
	return nil
}

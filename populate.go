package cliconfig

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	NotSettableErr error = errors.New("ptr value is not a settable struct pointer")
)

// Populate populates structs values
func Populate(cmd *cobra.Command, ptr interface{}) error {
	valueOf := reflect.ValueOf(ptr)
	indirect := reflect.Indirect(valueOf)
	t := indirect.Type()
	if !indirect.CanSet() || valueOf.Kind() != reflect.Ptr || indirect.Kind() != reflect.Struct {
		return NotSettableErr
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		arg := f.Tag.Get("arg")
		if arg == "" {
			continue
		}
		fmt.Printf("arg: %s\n", arg)

		switch f.Type.Kind() {
		case reflect.String:
			val, err := cmd.Flags().GetString(arg)
			if err != nil {
				return err
			}
			reflect.ValueOf(ptr).Elem().Field(i).SetString(val)
		case reflect.Bool:
			val, err := cmd.Flags().GetBool(arg)
			if err != nil {
				return err
			}
			reflect.ValueOf(ptr).Elem().Field(i).SetBool(val)
		case reflect.Slice:
			vals, err := cmd.Flags().GetStringArray(arg)
			if err != nil {
				return err
			}
			reflect.ValueOf(ptr).Elem().Field(i).Set(reflect.ValueOf(vals))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := cmd.Flags().GetInt(arg)
			if err != nil {
				return err
			}
			reflect.ValueOf(ptr).Elem().Field(i).SetInt(int64(val))
		default:
			continue
		}
	}
	return nil
}

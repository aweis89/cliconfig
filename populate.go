package cliconfig

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

var (
	StructPtrErr error = errors.New("value is not a struct pointer")
)

// Populate populates structs values
func Populate(cmd *cobra.Command, ptr interface{}) error {
	valueOf := reflect.ValueOf(ptr)
	if valueOf.Kind() != reflect.Ptr {
		return StructPtrErr
	}
	indirect := reflect.Indirect(valueOf)
	t := indirect.Type()
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

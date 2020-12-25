package cliconfig

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
)

// Populate populates structs values
func Populate(cmd *cobra.Command, ptr interface{}) error {
	newVal := reflect.Indirect(reflect.ValueOf(ptr))
	t := newVal.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		arg := f.Tag.Get("arg")
		if arg != "" {
			fmt.Printf("arg: %s\n", arg)
			val, err := cmd.Flags().GetString(arg)
			if err != nil {
				return err
			}
			reflect.ValueOf(ptr).Elem().Field(i).SetString(val)
		}
	}
	return nil
}

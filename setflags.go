package cliconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetFlags registers flags from struct tags.
// Example: field `arg:"flag-name" required:"false" desc:"description" short:"fm"`
func SetFlags(cmd *cobra.Command, str interface{}) {
	t := reflect.TypeOf(str)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("arg")
		desc := f.Tag.Get("desc")
		short := f.Tag.Get("short")
		def := f.Tag.Get("default")
		// default to all fields being required
		required := f.Tag.Get("required") != "false"

		switch f.Type.Kind() {
		case reflect.String:
			cmd.Flags().StringP(name, short, def, desc)
			// case reflect.SliceOf(reflect:
			// 	cmd.Flags().StringArrayP(name, short, []string{}, usage)
		}
		if required {
			cmd.MarkFlagRequired(name)
		}
	}
}

// BindViperDefaults each cobra flag to its associated viper configuration (config file and environment variable)
func BindViperDefaults(cmd *cobra.Command, prefix string) error {
	var result error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := f.Name
		if prefix != "" {
			key = fmt.Sprintf("%s-%s", prefix, key)
		}

		// bind to env var
		envVar := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		fmt.Printf("key %s binding to %s\n", key, envVar)
		if err := viper.BindEnv(key, envVar); err != nil {
			result = err
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(key) {
			val := viper.Get(key)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				result = err
			}
		}
	})
	return result
}

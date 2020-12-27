package cliconfig

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BindViperDefaults binds each cobra flag to its associated viper configuration (config file and environment variable)
func BindViperDefaults(flags flag.FlagSet, prefix string) error {
	return BindViperIntanceDefaults(flags, viper.GetViper(), prefix)
}

// BindViperIntanceDefaults uses a viper instance as apposed to global viper from BindViperDefaults
func BindViperIntanceDefaults(flags flag.FlagSet, v *viper.Viper, prefix string) error {
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

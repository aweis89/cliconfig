package cliconfig

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ViperSetFlags(flags *flag.FlagSet, prefix string, vs ...*viper.Viper) error {
	var result error
	// default to global viper instance
	if len(vs) == 0 {
		vs = []*viper.Viper{viper.GetViper()}
	}
	for _, v := range vs {
		flags.VisitAll(func(f *pflag.Flag) {
			key := f.Name
			if prefix != "" {
				key = fmt.Sprintf("%s-%s", prefix, key)
			}

			// bind to env var
			envVar := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
			if err := v.BindEnv(key, envVar); err != nil {
				result = err
				return
			}

			// If flag is not set and viper has config value, set flag to viper's value.
			if !f.Changed && v.IsSet(key) {
				fmt.Println("setting viper", key)
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

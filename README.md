[![Actions Status](https://github.com/aweis89/cliconfig/workflows/build/badge.svg)](https://github.com/aweis89/cliconfig/actions)
[![codecov](https://codecov.io/gh/aweis98/cliconfig/branch/master/graph/badge.svg)](https://codecov.io/gh/aweis98/cliconfig)

# cliconfig
Combines Viper and Cobra libraries for flexible configurable values.

```go
import(
	"github.com/aweis89/viper/cliconfig"
	"github.com/spf13/cobra"
)

type myStruct struct {
	// The arg tag is used as the CLI name and Viper lookup key when binding to viper, see below.
	SomeArg  string   `arg:"foo-arg" short:"an" desc:"does fooing stuff"`
	// By default all args are required to be set, either by the CLI or viper config when binding to viper
	Optional string   `arg:"some-optional-arg" required:"false"`
	Slice    []string `arg:"my-slice"`
	Bool     bool     `arg:"my-bool"`
	Int      int      `arg:"my-int"`

}

var mycmd = &cobra.Command{
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// When an arg is not set on the CLI, the arg will get set to the viper lookup value (using the global viper instance).
		// Bind all args to viper keys using prefix-<arg> and env vars PREFIX_<upcased arg>.
		// For example, in this case a viper registered config with `prefix-foo-arg` or an env variable of `PREFIX_FOO_ARG` will be used 
		// assuming `--foo-arg` is not specified on the CLI.
		return cliconfig.BindViperDefaults(cmd, "prefix")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ms := myStruct{}
		if err := cliconfig.Populate(cmd.Flags(), &ms); err != nil {
			return err
		}
		fmt.Printf("%+v", ms)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mycmd)
	cliconfig.SetFlags(gmc, myStruct{})
}
```

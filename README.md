# cliconfig
Combines Viper and Cobra libraries for flexible configurable values.

```go
import(
	"github.com/aweis89/viper/cliconfig"
	"github.com/spf13/cobra"
)

type myStruct struct {
	SomeArg string `arg:"foo-arg" short:"an" required:"false" desc:"does fooing stuff"`
}

var mycmd = &cobra.Command{
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Bind all args to viper keys using prefix-<arg> and env vars PREFIX_<upcased arg>.
		// When an arg is not set on CLI, the arg will get set to the viper lookup value (using the global viper instance).
		return cliconfig.BindViperDefaults(cmd, "prefix")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ms := myStruct{}
		if err := cliconfig.Populate(cmd, &ms); err != nil {
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

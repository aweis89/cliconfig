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

&cobra.Command{
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cliconfig.BindViperDefaults(cmd, "git")
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
```

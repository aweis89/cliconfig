[![Actions Status](https://github.com/aweis89/cliconfig/workflows/build/badge.svg)](https://github.com/aweis89/cliconfig/actions)
[![codecov](https://codecov.io/gh/aweis89/cliconfig/branch/master/graph/badge.svg)](https://codecov.io/gh/aweis89/cliconfig)

### cliconfig
Uses struct field tags to set flags using [pflags](https://github.com/spf13/pflags).
Can also be set to use Viper as default fallback when cli arg is missing.

<details>
<summary>Example using https://github.com/spf13/cobra:</summary>

```console
$ bat main.go
```	
```go
package main

import (
	"fmt"

	"github.com/aweis89/cliconfig"
	"github.com/spf13/cobra"
)

type myStruct struct {
	// The arg tag is used as the CLI name and Viper lookup key when binding to viper, see below.
	SomeArg string `arg:"foo-arg" short:"f" desc:"does fooing stuff"`
	// By default all args are required to be set, either by the CLI or viper config when binding to viper
	Optional string   `arg:"some-optional-arg" required:"false"`
	Slice    []string `arg:"my-slice"`
	Bool     bool     `arg:"my-bool"`
	Int      int      `arg:"my-int"`
}

func main() {
	cmd := &cobra.Command{
		Use: "testcmd",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// When an arg is not set on the CLI, the arg will get set to the viper lookup value (using the global viper instance).
			// Bind all args to viper keys using prefix-<arg> and env vars PREFIX_<upcased arg>.
			// For example, in this case a viper registered config with `prefix-foo-arg` or an env variable of `PREFIX_FOO_ARG` will be used
			// assuming `--foo-arg` is not specified on the CLI.
			return cliconfig.ViperSetFlags(cmd.Flags(), "prefix")

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ms := myStruct{}
			if err := cliconfig.Populate(cmd.Flags(), &ms); err != nil {
				return err
			}
			fmt.Printf("Populated struct: %+v\n", ms)
			return nil
		},
	}
	panicIfErr(cliconfig.SetFlags(cmd.Flags(), myStruct{}))
	panicIfErr(cmd.Execute())
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
```
```console
$ go run ./ --help
Usage:
  testcmd [flags]

Flags:
  -f, --foo-arg string             does fooing stuff
  -h, --help                       help for testcmd
      --my-bool
      --my-int int
      --my-slice stringArray
      --some-optional-arg string

$ go run ./
Error: required flag(s) "foo-arg", "my-bool", "my-int", "my-slice" not set
Usage:
  testcmd [flags]

Flags:
  -f, --foo-arg string             does fooing stuff
  -h, --help                       help for testcmd
      --my-bool
      --my-int int
      --my-slice stringArray
      --some-optional-arg string

$ go run ./ --foo-arg cli --my-bool --my-int 10 --my-slice one --my-slice two
Populated struct: {SomeArg:cli Optional: Slice:[one two] Bool:true Int:10}

# For viper usage, first setup viper
$ cat <<EOF >> ./main.go
func init() {
	viper.SetConfig("config")
	panicIfErr(viper.ReadInConfig())
}
EOF

# Add cli args to config.yaml with prefix specified in PreRunE `ViperSetFlags`
$ cat <<EOF > config.yaml
prefix-my-int: 88
prefix-my-bool: true
prefix-my-slice:
- one
- two
EOF

$ go run ./
Populated struct: {SomeArg:cli Optional: Slice:[one two] Bool:true Int:10}
```
</details>

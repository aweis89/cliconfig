package cliconfig

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setArgs(str ...string) (cleanup func()) {
	old := os.Args
	os.Args = append(os.Args[:1], str...)
	cleanup = func() {
		os.Args = old
	}
	return
}

func errTest(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		if testing.Verbose() {
			debug.PrintStack()
		}
	}
}

func ensureBool(t *testing.T, test bool, subject string) {
	if !test {
		t.Errorf("failed equality test regarding %s\n", subject)
		if testing.Verbose() {
			debug.PrintStack()
		}
	}
}

func ensureEq(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		t.Errorf("expected `%+v` to eq `%+v` while testing: `%s`", a, b, msg)
		if testing.Verbose() {
			debug.PrintStack()
		}
	}
}

// Test basic types are getting set
type typesStruct struct {
	String string   `arg:"string" desc:"github auth token"`
	Slice  []string `arg:"slice"`
	Bool   bool     `arg:"bool"`
	Int    int      `arg:"integer"`
}

func TestPopulate(t *testing.T) {
	cmd := cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Executing command..")
			ts := typesStruct{}
			errTest(t, Populate(cmd.Flags(), &ts))
			fmt.Printf("got struct: %+v\n", ts)
			ensureEq(t, ts.Bool, true, "struct bool field")
			ensureEq(t, ts.String, "string-cli", "string field not set correctly")
			ensureEq(t, len(ts.Slice), 2, "struct slice incorrect size")
			ensureEq(t, ts.Slice[0], "slice-cli-a", "slice args in struct")
			ensureEq(t, ts.Slice[1], "slice-cli-b", "slice args in struct")
			return nil
		},
	}
	defer setArgs("--slice=slice-cli-a",
		"--slice=slice-cli-b",
		"--string=string-cli",
		"--integer", "10",
		"--bool")()
	ts := typesStruct{}
	flags := cmd.Flags()
	errTest(t, SetFlags(flags, ts))
	errTest(t, cmd.Execute())
	strErr := SetFlags(flags, "not struct pointer")
	ensureBool(t, strErr != nil, "expecting error from SetFlags args")
	strErr = Populate(flags, "not struct pointer")
	ensureBool(t, strErr != nil, "expecting error from Populate")
}

func TestPopulateWithViper(t *testing.T) {
	v := viper.New()
	v.SetConfigType("yaml")
	config := `string: file-arg`
	v.ReadConfig(bytes.NewBuffer([]byte(config)))
	t.Run("use config value", func(t *testing.T) {
		str := typesStruct{}
		emptyFlagSet := pflag.FlagSet{}
		SetFlags(&emptyFlagSet, str)
		PopulateWithViper(&emptyFlagSet, &str, "", v)
		ensureEq(t, str.String, "file-arg", "flag value is set from viper config file")
	})
	t.Run("use cli override", func(t *testing.T) {
		str := typesStruct{}
		flagSet := pflag.FlagSet{}
		SetFlags(&flagSet, str)
		flagSet.Parse([]string{"--string", "cli-arg"})
		PopulateWithViper(&flagSet, &str, "", v)
		ensureEq(t, str.String, "cli-arg", "flag value is set from from cli override")
	})
}

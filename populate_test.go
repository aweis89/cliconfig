package cliconfig

import (
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/spf13/cobra"
)

func TestPopulate(t *testing.T) {
	type typesStruct struct {
		String string   `arg:"string" desc:"github auth token"`
		Slice  []string `arg:"slice"`
		Bool   bool     `arg:"bool"`
		Int    int      `arg:"integer"`
	}

	cmd := cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Executing command..")
			ts := typesStruct{}
			errTest(t, Populate(cmd, &ts))
			fmt.Printf("got struct: %+v\n", ts)
			ensureBool(t, ts.Bool == true, "struct bool field")
			ensureBool(t, ts.String == "string-cli", "string field not set correctly")
			ensureBool(t, len(ts.Slice) == 2, "struct slice incorrect size")
			ensureBool(t, ts.Slice[0] == "slice-cli-a", "incorrect slice args in struct")
			ensureBool(t, ts.Slice[1] == "slice-cli-b", "incorrect slice args in struct")
			return nil
		},
	}

	cliArgs := []string{
		"--slice", "slice-cli-a",
		"--slice", "slice-cli-b",
		"--string", "string-cli",
		"--bool", "true",
		"--integer", "10"}

	cmd.SetArgs(cliArgs)

	ts := typesStruct{}
	// errTest(t, SetFlags(&cmd, &ts))
	errTest(t, SetFlags(&cmd, ts))
	errTest(t, cmd.Execute())
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

package cliconfig

import (
	"fmt"
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
			reportErr(t, Populate(cmd, &ts))
			fmt.Printf("%+v", ts)
			// setBool, err := cmd.Flags().GetBool("bool")
			// reportErr(t, err)
			ensureBool(t, ts.Bool == true, "struct bool field")

			//var strarg string
			//strarg, err = cmd.Flags().GetString("string")
			// reportErr(t, err)
			ensureBool(t, ts.String == "string-cli", "string field not set correctly")

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
	}
}

func ensureBool(t *testing.T, test bool, subject string) {
	if !test {
		t.Errorf("failed equality test regarding %s\n", subject)
	}
}

func reportErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

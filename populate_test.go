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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Executing command..")
			ts := typesStruct{}
			Populate(cmd, &ts)
			fmt.Printf("%+v", ts)
			if setBool, err := cmd.Flags().GetBool("bool"); setBool != ts.Bool || err != nil {
				reportErr(err, t)
				t.Errorf("expecting %v got %v\n", setBool, ts.Bool)
			}

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
	errTest(t, SetFlags(&cmd, ts))
	errTest(t, cmd.Execute())
}

func errTest(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func boolTest(t *testing.T, pass bool, msg string) {
	if !pass {
		t.Error(msg)
	}
}

func reportErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestFullName(t *testing.T) {
	cmdWithName := func(cmdName string) *cobra.Command {
		return &cobra.Command{
			Use: cmdName,
		}
	}
	cmdWithParent := func(cmd *cobra.Command, parent *cobra.Command) *cobra.Command {
		parent.AddCommand(cmd)
		return cmd
	}

	tests := []struct {
		cmd            *cobra.Command
		expectedResult string
	}{
		{
			cmd:            nil,
			expectedResult: "",
		},
		{
			cmd:            cmdWithName("cmd"),
			expectedResult: "cmd",
		},
		{
			cmd:            cmdWithParent(cmdWithName("sub"), cmdWithName("top")),
			expectedResult: "top sub",
		},
		{
			cmd:            cmdWithParent(cmdWithName("sub2"), cmdWithParent(cmdWithName("sub"), cmdWithName("top"))),
			expectedResult: "top sub sub2",
		},
	}

	for count, test := range tests {
		result := FullName(test.cmd)
		if result != test.expectedResult {
			t.Errorf("Test[%d] Failed: Expected '%s' but got '%s'", count, test.expectedResult, result)
		}
	}
}

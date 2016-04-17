package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// RunHelp is a Command "Run" compatible function
// that just prints the command's help to stdout
func RunHelp(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		fmt.Printf("Failed to print help message! %v", err)
	}
}

// FullName returns the full name of the given command,
// which is the concatenation of all the names of all the parents commands
// example: "top-cmd sub-cmd sub-sub-cmd"
func FullName(cmd *cobra.Command) string {
	names := []string{}
	for parent := cmd; parent != nil; parent = parent.Parent() {
		names = append(names, parent.Name())
	}

	fullName := ""
	for _, name := range names {
		fullName = fmt.Sprintf("%s %s", name, fullName)
	}

	return strings.TrimSpace(fullName)
}

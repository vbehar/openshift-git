package main

import (
	"fmt"

	"github.com/vbehar/openshift-git/pkg/cmd"

	// init all the commands
	_ "github.com/vbehar/openshift-git/pkg/cmd/export"
	_ "github.com/vbehar/openshift-git/pkg/cmd/importer"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

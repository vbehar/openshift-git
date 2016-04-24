package cmd

import (
	"flag"

	"github.com/spf13/cobra"

	// init glog to get its flags
	_ "github.com/golang/glog"
)

var (
	// RootCmd is the main command
	RootCmd = &cobra.Command{
		Use:   "openshift-git",
		Short: "Import/Export OpenShift resources from/to a Git repository",
		Long: `openshift-git is a standalone OpenShift Client which helps you export/import OpenShift resources to/from a Git repository.

Run either of the following commands to see the usage:
$ openshift-git export --help
$ openshift-git import --help

More informations at https://github.com/vbehar/openshift-git`,
		Run: RunHelp,
	}
)

func init() {
	// add glog flags
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

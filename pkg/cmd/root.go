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
		Short: "Nagios plugin for interacting with an OpenShift cluster",
		Long:  `openshift-git is a standalone OpenShift Client which helps you setup Nagios checks`,
		Run:   RunHelp,
	}
)

func init() {
	// add glog flags
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

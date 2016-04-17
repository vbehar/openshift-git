package importer

import (
	"github.com/vbehar/openshift-git/pkg/cmd"
	"github.com/vbehar/openshift-git/pkg/openshift"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import OpenShift resources from a Git repository",
		Long: `
Imports OpenShift resources from a Git repository.

Note that it behaves like the standard OpenShift Client (oc) to connect to the OpenShift Cluster.
By default, if a ~/.kube/config file exists, it will be used.
Otherwise, you can use the same option as the OpenShift Client (oc):
--config to use a custom kube config file
--server and --token to specify the master URL and (service account) token directly`,
		Run: func(cmd *cobra.Command, args []string) {
			glog.Fatal("The import command is not implemented yet!")
		},
	}
)

func init() {
	cmd.RootCmd.AddCommand(importCmd)
	importCmd.Flags().AddFlagSet(openshift.Flags)
}

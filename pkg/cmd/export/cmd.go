package export

import (
	"fmt"
	"time"

	"github.com/vbehar/openshift-git/pkg/cmd"
	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export OpenShift resources from a Git repository",
		Long: `
Exports OpenShift resources from a Git repository.

Note that it behaves like the standard OpenShift Client (oc) to connect to the OpenShift Cluster.
By default, if a ~/.kube/config file exists, it will be used.
Otherwise, you can use the same option as the OpenShift Client (oc):
--config to use a custom kube config file
--server and --token to specify the master URL and (service account) token directly`,
		ValidArgs: []string{"all"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Missing export type.")
			}
			if args[0] != "all" {
				return fmt.Errorf("Invalid export type '%s'. The only supported type for the moment is 'all'.", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			repo, err := git.NewRepository(exportOptions.RepositoryPath,
				exportOptions.RepositoryBranch,
				exportOptions.RepositoryRemote,
				exportOptions.RepositoryContextDir)
			if err != nil {
				glog.Fatalf("Failed to init git repo: %v", err)
			}

			if len(repo.RemoteURL) > 0 {
				glog.Infof("Exporting %s to %s ...", args[0], repo.RemoteURL)
			} else {
				glog.Infof("Exporting %s to %s ...", args[0], repo.Path)
			}

			if exportOptions.Watch {
				runWatch(repo)
			} else {
				runList(repo)
			}
		},
	}

	exportOptions = &ExportOptions{}
)

func init() {
	cmd.RootCmd.AddCommand(exportCmd)
	exportCmd.Flags().AddFlagSet(openshift.Flags)
	exportCmd.Flags().StringVar(&exportOptions.RepositoryPath, "repository-path", "/tmp/repo", "Path of the git repository")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryBranch, "repository-branch", "master", "Branch of the git repository")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryRemote, "repository-remote", "", "Remote of the git repository")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryContextDir, "repository-context-dir", "", "Dir in the git repository")
	exportCmd.Flags().StringVar(&exportOptions.Format, "format", "yaml", "Format (json, yaml)")
	exportCmd.Flags().StringVarP(&exportOptions.LabelSelector, "selector", "l", "", "Selector (label query) to filter on")
	exportCmd.Flags().BoolVar(&exportOptions.AllNamespaces, "all-namespaces", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	exportCmd.Flags().BoolVarP(&exportOptions.Watch, "watch", "w", false, "Keep watching any changes")
	exportCmd.Flags().DurationVar(&exportOptions.ResyncPeriod, "resync-period", 1*time.Hour, "If not zero, define the interval of time to perform a full resync")
	exportCmd.Flags().DurationVar(&exportOptions.RepositoryPullPeriod, "repository-pull-period", 2*time.Minute, "If not zero, define the interval of time to perform a pull of the git repository")
	exportCmd.Flags().DurationVar(&exportOptions.RepositoryPushPeriod, "repository-push-period", 2*time.Minute, "If not zero, define the interval of time to perform a push of the git repository")
}

// ExportOptions represents the options of the export command
type ExportOptions struct {
	AllNamespaces        bool
	Namespace            string
	Format               string
	Watch                bool
	LabelSelector        string
	ResyncPeriod         time.Duration
	RepositoryPath       string
	RepositoryBranch     string
	RepositoryRemote     string
	RepositoryContextDir string
	RepositoryPullPeriod time.Duration
	RepositoryPushPeriod time.Duration
}

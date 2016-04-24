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
	exportCmdLongDescription = `
Exports OpenShift resources to a Git repository - optionally pushing to a configured remote.

It can either be run:
- as a one-time operation, just exporting all resources
- as a daemon, first exporting all resources and then watching for changes (with the '--watch' option)

It expects a comma-separated list of types to export, like buildconfig, pods, routes and so on.
You can use the special 'all' alias (expanded by OpenShift to [bc builds is dc rc routes svc pods]),
or the recommended 'everything' alias (expanded by openshift-git to %[1]s).

The '--repository-path' flag is mandatory: it defines where the files will be saved on the filesystem.
If there is no existing repository at this path, a new repository will be created.
If a remote repository is provided with the '--repository-remote' flag, it will be cloned to the local repository.

By default, resources will be exported in the YAML format, but the '--format' flag can be used to export as JSON.

Note that it behaves like the standard OpenShift Client (oc) to connect to the OpenShift Cluster.
By default, if a ~/.kube/config file exists, it will be used.
Otherwise, you can use the same option as the OpenShift Client (oc):
--config to use a custom kube config file
--server and --token to specify the master URL and (service account) token directly`
	exportCmdExample = `
	# Basic usage: export everything from the current namespace to a new Git repository at /tmp/export
	$ %[1]s everything --repository-path=/tmp/export

	# Export everything from the "my-namespace" namespace, to a new Git repository
	$ %[1]s everything -n my-namespace --repository-path=/tmp/export
	
	# Export specific types from the "my-namespace" namespace, to a new Git repository
	$ %[1]s bc,dc,is,svc,route,ns -n my-namespace --repository-path=/tmp/export

	# Export everything from all namespaces, to a new Git repository
	# Note that it requires at least the cluster-reader role
	$ %[1]s everything --all-namespaces --repository-path=/tmp/export

	# Export everything from all namespaces, and keep watching for changes
	$ %[1]s everything --all-namespaces --repository-path=/tmp/export -w`

	exportCmd = &cobra.Command{
		Use:   "export TYPE",
		Short: "Export OpenShift resources to a Git repository",
		PreRunE: func(command *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Missing export type.")
			}
			if len(exportOptions.RepositoryPath) == 0 {
				return fmt.Errorf("Missing repository path.")
			}
			return nil
		},
		Run: func(command *cobra.Command, args []string) {
			repo, err := git.NewRepository(exportOptions.RepositoryPath,
				exportOptions.RepositoryBranch,
				exportOptions.RepositoryRemote,
				exportOptions.RepositoryContextDir)
			if err != nil {
				glog.Fatalf("Failed to init git repo: %v", err)
			}

			if exportOptions.Watch {
				err = runWatch(args[0], repo)
			} else {
				err = runList(args[0], repo)
			}

			if err != nil {
				glog.Fatalf("Failed: %v", err)
			}
		},
	}

	exportOptions = &ExportOptions{}
)

func init() {
	cmd.RootCmd.AddCommand(exportCmd)
	exportCmd.Long = fmt.Sprintf(exportCmdLongDescription, openshift.AllKinds)
	exportCmd.Example = fmt.Sprintf(exportCmdExample, cmd.FullName(exportCmd))
	exportCmd.Flags().AddFlagSet(openshift.Flags)
	exportCmd.Flags().StringVar(&exportOptions.RepositoryPath, "repository-path", "", "Mandatory. Path of the git repository on the filesystem. A new repository will be created if the path does not exists.")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryBranch, "repository-branch", "master", "Branch of the git repository to use for commits.")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryRemote, "repository-remote", "", "Optional URL of a remote git repository. If present, periodic push/pull operations will be scheduled, to keep the local and remote repositories in sync.")
	exportCmd.Flags().StringVar(&exportOptions.RepositoryContextDir, "repository-context-dir", "", "Optional relative directory (in the repository) that will be used to store data.")
	exportCmd.Flags().StringVar(&exportOptions.Format, "format", "yaml", "Format of the exported resources ('json' or 'yaml')")
	exportCmd.Flags().StringVarP(&exportOptions.LabelSelector, "selector", "l", "", "Selector (label query) to filter on")
	exportCmd.Flags().BoolVar(&exportOptions.AllNamespaces, "all-namespaces", false, "If present, export the requested resources across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	exportCmd.Flags().BoolVar(&exportOptions.UseDefaultSelector, "default-selector", true, "If present, some default label selectors will be applied (for example, ignore build and deploy pods, ignore pods managed by RC or DC, or ignore RC managed by DC)")
	exportCmd.Flags().BoolVarP(&exportOptions.Watch, "watch", "w", false, "After exporting the requested types, watch for changes.")
	exportCmd.Flags().DurationVar(&exportOptions.ResyncPeriod, "resync-period", 1*time.Hour, "If not zero, defines the interval of time to perform a full resync of the OpenShift resources to export.")
	exportCmd.Flags().DurationVar(&exportOptions.RepositoryPullPeriod, "repository-pull-period", 2*time.Minute, "If not zero, defines the interval of time to perform a pull of the remote git repository.")
	exportCmd.Flags().DurationVar(&exportOptions.RepositoryPushPeriod, "repository-push-period", 2*time.Minute, "If not zero, defines the interval of time to perform a push to the remote git repository.")
}

// ExportOptions represents the options of the export command
type ExportOptions struct {
	AllNamespaces        bool
	Namespace            string
	Format               string
	Watch                bool
	UseDefaultSelector   bool
	LabelSelector        string
	ResyncPeriod         time.Duration
	RepositoryPath       string
	RepositoryBranch     string
	RepositoryRemote     string
	RepositoryContextDir string
	RepositoryPullPeriod time.Duration
	RepositoryPushPeriod time.Duration
}

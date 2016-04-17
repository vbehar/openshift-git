package export

import (
	"time"

	"github.com/vbehar/openshift-git/pkg/git"
	"github.com/vbehar/openshift-git/pkg/openshift"

	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl"

	"github.com/golang/glog"
)

// saveResources saves all the resources coming from the given channel to the given git repository.
// it pulls/pushes from/to the remote repository at configured interval if the git repository has a remote.
// should be run in a single goroutine (the git-related operations are not thread-safe)
func saveResources(repo *git.Repository, resourcesChan <-chan openshift.Resource, mapper meta.RESTMapper, printer kubectl.ResourcePrinter) {
	var saved, deleted int64
	pullTicker := time.NewTicker(exportOptions.RepositoryPullPeriod)
	pushTicker := time.NewTicker(exportOptions.RepositoryPushPeriod)

	for {
		select {

		case <-pullTicker.C:
			if err := repo.Pull(); err != nil {
				glog.Errorf("Failed to pull from %s: %v", repo.RemoteURL, err)
			}

		case <-pushTicker.C:
			if err := repo.Push(); err != nil {
				glog.Errorf("Failed to push to %s: %v", repo.RemoteURL, err)
			}

		case resource, open := <-resourcesChan:
			if !open {
				glog.Infof("Closing ! Stats: %d resources saved, and %d resources deleted.", saved, deleted)
				return
			}

			if resource.Exists {
				if err := saveResource(repo, &resource, mapper, printer); err != nil {
					glog.Errorf("Failed to save %s: %v", resource.String(), err)
				} else {
					saved++
				}
			} else {
				if err := deleteResource(repo, &resource); err != nil {
					glog.Errorf("Failed to delete %s: %v", resource.String(), err)
				} else {
					deleted++
				}
			}
		}
	}
}

// saveResource saves (and commit) the single given resource to the given git repository
func saveResource(repo *git.Repository, resource *openshift.Resource, mapper meta.RESTMapper, printer kubectl.ResourcePrinter) error {
	glog.V(2).Infof("Saving %s", resource)

	printer, err := upgradePrinterForObject(printer, resource.Object, mapper)
	if err != nil {
		return err
	}

	gitResource := git.NewGitResource(repo, resource, exportOptions.Format)

	if err := gitResource.Open(); err != nil {
		return err
	}

	if err := printer.PrintObj(resource.Object, gitResource); err != nil {
		gitResource.Close()
		return err
	}
	gitResource.Close()

	if err = gitResource.Commit(); err != nil {
		return err
	}

	return nil
}

// deleteResource deletes (and commit) the single given resource from the given git repository
func deleteResource(repo *git.Repository, resource *openshift.Resource) error {
	glog.V(3).Infof("Deleting %s", resource.String())

	gitResource := git.NewGitResource(repo, resource, exportOptions.Format)

	if err := gitResource.Delete(); err != nil {
		return err
	}

	if err := gitResource.Commit(); err != nil {
		return err
	}

	return nil
}

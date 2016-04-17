package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vbehar/openshift-git/pkg/openshift"

	git "github.com/gogits/git-module"
)

// GitResource represents an OpenShift resource in a Git repository
type GitResource struct {
	// repository is the repository in which the resource is stored
	repository *Repository

	// resource is the underlying resource
	resource *openshift.Resource

	// format is the storage format of the resource (like YAML or JSON)
	format string

	// path is the full absolute path on the filesystem where the resource is stored
	path string

	// file is the resource's file on the filesystem
	file *os.File
}

// NewGitResource instantiates a new GitResource in the given repository, for the given resource, in the given format
func NewGitResource(repository *Repository, resource *openshift.Resource, format string) *GitResource {
	path := repository.PathForResource(resource, format)

	gitResource := &GitResource{
		repository: repository,
		resource:   resource,
		format:     format,
		path:       path,
	}
	return gitResource
}

// Open opens the resource so that it could then be used as an io.Writer
// It then needs to be closed at the end.
func (gr *GitResource) Open() error {
	dir := filepath.Dir(gr.path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var err error
	gr.file, err = os.OpenFile(gr.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	return err
}

// Write writes some data to the resource
// implements the io.Writer interface
func (gr *GitResource) Write(p []byte) (n int, err error) {
	return gr.file.Write(p)
}

// Close closes the underlying file
// Needs to be called after writing
// implements the io.Closer interface
func (gr *GitResource) Close() error {
	return gr.file.Close()
}

// Delete deletes the resource from the filesystem
// does not complains if the file does not exists
func (gr *GitResource) Delete() error {
	err := os.Remove(gr.path)
	if os.IsNotExist(err) {
		// already deleted
		return nil
	}
	return err
}

// Commit commits the resource to the git repository
func (gr *GitResource) Commit() error {
	needCommit, err := IsFileNewOrModified(gr.repository.Path, gr.path)
	if err != nil {
		return err
	}
	if !needCommit {
		return nil
	}

	if err := git.AddChanges(gr.repository.Path, false, gr.path); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("%s %s", gr.resource.Status, gr.resource)
	if err := git.CommitChanges(gr.repository.Path, commitMsg, nil); err != nil {
		git.ResetHEAD(gr.repository.Path, false, "HEAD")
		return err
	}

	return nil
}

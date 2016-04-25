package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vbehar/openshift-git/pkg/openshift"

	git "github.com/gogits/git-module"
	"github.com/golang/glog"
)

// Repository represents a Git repository
type Repository struct {
	// The underlying git repository
	*git.Repository

	// Path is the path of the repository on the FileSystem
	Path string

	// Branch is the branch that we will use to commit our changes
	Branch string

	// RemoteURL is the URL of the remote repository (optional).
	// If specified, we will pull/push from/to this remote
	RemoteURL string

	// ContextDir is the (optional) path (relative to `Path`)
	// which will be used inside the repository.
	ContextDir string
}

// NewRepository instantiates a new Git repository at the given path.
// If there is nothing at the given path, a new repository will be initialized:
// - If a remoteURL is provided, we will clone from this configured remote
// - Otherwise, we will just create an new empty repository
func NewRepository(path, branch, remoteURL, contextDir, userName, userEmail string) (*Repository, error) {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		if len(remoteURL) > 0 {
			glog.Infof("Cloning from %s to %s ...", remoteURL, path)
			if err := git.Clone(remoteURL, path, git.CloneRepoOptions{}); err != nil {
				return nil, err
			}
		} else {
			glog.Infof("Initializing a new empty repository at %s ...", path)
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return nil, err
			}
			if err := git.InitRepository(path, false); err != nil {
				return nil, err
			}
		}

		if len(userName) > 0 {
			if err := SetUserName(path, userName); err != nil {
				return nil, err
			}
		}
		if len(userEmail) > 0 {
			if err := SetUserEmail(path, userEmail); err != nil {
				return nil, err
			}
		}

		err = nil
	}

	if err != nil {
		return nil, err
	}

	if fi != nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	if len(remoteURL) > 0 {
		glog.Infof("Scheduling push/pull to/from remote repository %s ...", remoteURL)
		// TODO do a first pull if using an existing local repo
	} else {
		glog.Infof("No remote repository configured.")
	}

	repository := &Repository{
		Repository: repo,
		Path:       path,
		Branch:     branch,
		RemoteURL:  remoteURL,
		ContextDir: contextDir,
	}

	if err := repo.SetDefaultBranch(branch); err != nil {
		return nil, err
	}

	fi, err = os.Stat(repository.PathWithContextDir())

	if os.IsNotExist(err) {
		if err := os.MkdirAll(repository.PathWithContextDir(), os.ModePerm); err != nil {
			return nil, err
		}

		err = nil
	}

	if err != nil {
		return nil, err
	}

	if fi != nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", repository.PathWithContextDir())
	}

	commits, err := git.CommitsCount(path, "HEAD")
	if commits == 0 || err != nil {
		glog.V(1).Infof("Creating initial commit...")
		if err := ioutil.WriteFile(
			filepath.Join(repository.PathWithContextDir(), "README.md"),
			[]byte("# Export of OpenShift resources\n\nAutomatically managed by `openshift-git`."),
			0644); err != nil {
			return nil, err
		}
		if err := git.AddChanges(path, true); err != nil {
			return nil, err
		}
		if err := git.CommitChanges(path, "Initial commit", nil); err != nil {
			return nil, err
		}
		if err := repository.Push(); err != nil {
			return nil, err
		}
	}

	return repository, nil
}

// Pull pulls from the configured remote
// (if a remote as been configured)
func (r *Repository) Pull() error {
	if len(r.RemoteURL) > 0 {
		if err := Pull(r.Path, "origin", r.Branch); err != nil {
			return err
		}
	}
	return nil
}

// Push pushes to the configured remote
// (if a remote as been configured)
func (r *Repository) Push() error {
	if len(r.RemoteURL) > 0 {
		if err := git.Push(r.Path, "origin", r.Branch); err != nil {
			return err
		}
	}
	return nil
}

// PathForResource returns the full absolute path of the given resource, for the given format
func (r *Repository) PathForResource(resource *openshift.Resource, format string) string {
	filename := fmt.Sprintf("%s.%s", resource.Name, format)
	path := r.PathWithContextDir()
	if resource.IsNamespaced() {
		path = filepath.Join(path, "Namespace", resource.Namespace)
	}
	path = filepath.Join(path, resource.Kind, filename)

	return path
}

// ResourceFromPath returns a (minimalist) representation of the resource
// stored at the given path.
// Returns nil if no resource could be found at that path.
// Note that the returned resource contains only a reference
// (with kind, namespace and name), not the resource (content) itself.
func (r *Repository) ResourceFromPath(path string) *openshift.Resource {
	if strings.HasPrefix(path, r.PathWithContextDir()+"/") {
		path = strings.TrimPrefix(path, r.PathWithContextDir()+"/")
		elems := strings.Split(path, "/")
		switch len(elems) {
		case 4:
			namespace := elems[1]
			kind := elems[2]
			nameWithExtension := elems[3]
			extension := filepath.Ext(nameWithExtension)
			name := strings.TrimSuffix(nameWithExtension, extension)
			return openshift.NewResource(kind, fmt.Sprintf("%s/%s", namespace, name))
		case 2:
			kind := elems[0]
			nameWithExtension := elems[1]
			extension := filepath.Ext(nameWithExtension)
			name := strings.TrimSuffix(nameWithExtension, extension)
			return openshift.NewResource(kind, name)
		}
	}

	return nil
}

// PathWithContextDir returns the full path of the directory
// where we will write the exported resources
func (r *Repository) PathWithContextDir() string {
	if len(r.ContextDir) > 0 {
		return filepath.Join(r.Path, r.ContextDir)
	}
	return r.Path
}

// KeyListFuncForKind returns a ListKeys function, that implements the cache.KeyLister interface
// It is a function that returns the list of keys ("namespace/name" format)
// that we "know about" (to get a 2-way sync) for the given kind of resources.
// It simply walks the FS to list all the resources matching the given kind.
func (r *Repository) KeyListFuncForKind(kind string) func() []string {
	return func() []string {
		keys := []string{}

		err := filepath.Walk(r.PathWithContextDir(), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			resource := r.ResourceFromPath(path)
			if resource == nil {
				return nil
			}
			if resource.Kind == kind {
				key := resource.NamespacedName()
				glog.V(4).Infof("Found %s for %s at %s for %s", key, resource, path, kind)
				keys = append(keys, key)
			}
			return nil
		})
		if err != nil {
			glog.Errorf("Failed to walk FS %s for kind %s: %v", r.PathWithContextDir(), kind, err)
			return []string{}
		}

		glog.V(2).Infof("Found %d local keys for %s", len(keys), kind)
		return keys
	}
}

// KeyGetFuncForKindAndFormat returns a GetByKey function, implements the cache.KeyGetter interface
// It is a function that returns the object that we "know about"
// for the given key ("namespace/name" format) - and a boolean if it exists
// for the given kind and format.
func (r *Repository) KeyGetFuncForKindAndFormat(kind, format string) func(key string) (interface{}, bool, error) {
	return func(key string) (interface{}, bool, error) {
		resource := openshift.NewResource(kind, key)
		path := r.PathForResource(resource, format)

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				glog.V(3).Infof("key %s for kind %s does not exists at %s ! %v", key, kind, path, err)
				return "", false, nil
			}
			return "", false, err
		}

		glog.V(4).Infof("Found %v for %s %s at %s", resource, kind, key, path)
		return *resource, true, nil
	}
}

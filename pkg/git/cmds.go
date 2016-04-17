package git

import (
	git "github.com/gogits/git-module"
)

// IsFileNewOrModified checks if the given file is either new or modified
// in the given git repository.
func IsFileNewOrModified(repoPath, file string) (bool, error) {
	output, err := git.NewCommand("ls-files", "-m", "-o", file).RunInDir(repoPath)
	return len(output) > 0, err
}

// Pull pulls changes from given remote branch.
func Pull(repoPath, remote, branch string) error {
	_, err := git.NewCommand("pull", remote, branch).RunInDir(repoPath)
	return err
}

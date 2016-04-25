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

// SetUserName sets the user's name for the given repository
func SetUserName(repoPath, userName string) error {
	_, err := git.NewCommand("config", "user.name", userName).RunInDir(repoPath)
	return err
}

// SetUserEmail sets the user's email for the given repository
func SetUserEmail(repoPath, userEmail string) error {
	_, err := git.NewCommand("config", "user.email", userEmail).RunInDir(repoPath)
	return err
}

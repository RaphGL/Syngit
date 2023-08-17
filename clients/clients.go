package clients

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type GitRepo interface {
    // Returns the name of the repository
	GetName() string
    // Returns the name of the repository in the format "Username/Repo"
    GetFullName() string
    // Returns the Git URL for the repository
	GetURL() string
	IsPrivate() bool
	IsFork() bool
}

func CloneRepo(r GitRepo) {
	repoURL := r.GetURL()
    // TODO: use path from config or fallback to a default
    const test = "/home/raph/Documents/Test/"
	_, err := git.PlainClone(test + r.GetName(), false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

package client

import "github.com/google/go-github/v33/github"

// GitRepoInterface is for testing
type GitRepoInterface interface {
	GetContributors() ([]*github.Contributor, error)
}

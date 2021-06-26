package client

import (
	ctx "context"
	"errors"
	"github-contributors-action/internal/pkg/configs"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// GHClient is the custom handler for all requests
type GHClient struct {
	Client  *github.Client
	Context ctx.Context
	Config  configs.Config
}

// NewGHClient creates a new instance of GitHub client
func NewGHClient(config configs.Config) GitRepoInterface {
	context := ctx.Background()
	if config.GitHubToken == "" {
		return GHClient{
			Client:  github.NewClient(nil),
			Context: context,
			Config:  config,
		}
	}
	oauth2Client := oauth2.NewClient(context, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	))
	return GHClient{
		Client:  github.NewClient(oauth2Client),
		Context: context,
		Config:  config,
	}
}

func (client GHClient) GetContributors() ([]*github.Contributor, error) {
	repoParts := strings.Split(client.Config.SourceRepo, "/")
	// do not check for input, because it is already checked
	owner := repoParts[0]
	repository := repoParts[1]
	listContributorsOptions := &github.ListContributorsOptions{
		Anon: "1",
		ListOptions: github.ListOptions{
			PerPage: 20,
		},
	}
	var listOfContributors []*github.Contributor
	// Loop over until all contributors are listed, queries
	// PerPage number of entries
	for {
		contributors, response, err :=
			client.Client.Repositories.ListContributors(
				client.Context, owner, repository, listContributorsOptions)
		if err != nil {
			return nil, err
		}
		log.Printf("Response: %v", response)
		if response.StatusCode != http.StatusOK {
			return nil, errors.New("could not get the response")
		}
		listOfContributors = append(listOfContributors, contributors...)
		if response.NextPage == 0 {
			log.Println("Breaking from the loop of repositories")
			break
		}
		// assign next page
		listContributorsOptions.Page = response.NextPage
	}
	return listOfContributors, nil
}

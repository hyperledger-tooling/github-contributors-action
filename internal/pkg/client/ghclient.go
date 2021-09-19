package client

import (
	ctx "context"
	"errors"
	"fmt"
	"github-contributors-action/internal/pkg/configs"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
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

// GetContributors returns the list of all contributors.
// Refer to client configuration for the repositories information.
func (client GHClient) GetContributors() ([]*github.Contributor, error) {
	owner := client.Config.SourceRepos.Owner
	var listOfContributors []*github.Contributor

	// Get contributors from all of the repositories
	for _, repository := range client.Config.SourceRepos.Repositories {
		log.Printf("Requested to get contributors of %s/%s \n",
			client.Config.SourceRepos.Owner,
			*repository.Name,
		)
		listContributorsOptions := &github.ListContributorsOptions{
			Anon: "1",
			ListOptions: github.ListOptions{
				PerPage: 20,
			},
		}
		// Loop over until all contributors are listed, queries
		// PerPage number of entries
		for {
			contributors, response, err :=
				client.Client.Repositories.ListContributors(
					client.Context, owner, *repository.Name, listContributorsOptions)
			if err != nil {
				return nil, err
			}
			log.Printf("Response: %v", response)
			err, isSkip := filterRepoResponse(response)
			if isSkip {
				// continue to next repository, outer loop takes care of it
				break
			}
			if err != nil {
				return nil, err
			}
			listOfContributors = append(listOfContributors, contributors...)
			if response.NextPage == 0 {
				log.Println("Breaking from the loop of repositories")
				break
			}
			// assign next page
			listContributorsOptions.Page = response.NextPage
		}
	}

	var finalList []*github.Contributor
	// Loses association of repository in the Contributor list if any
	contributorsMap := make(map[int64]bool)
	for _, contributor := range listOfContributors {
		if _, present := contributorsMap[contributor.GetID()]; !present {
			contributorsMap[contributor.GetID()] = true
			finalList = append(finalList, contributor)
		}
	}

	// return final list of unique contributors
	return finalList, nil
}

// filterRepoResponse informs whether to skip the result of this repository.
// For example, skip the empty repository because there are no contributors.
func filterRepoResponse(response *github.Response) (error, bool) {
	// not interested in an empty repository
	if response.StatusCode == http.StatusNoContent {
		// nothing to do
		return nil, true
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response code: %d", response.StatusCode), false
	}
	// success case
	return nil, false
}

// GetRepos queries and fetches all the repositories that an organization
// has listed against the GITHUB_TOKEN used.
func (client GHClient) GetRepos() ([]*github.Repository, error) {
	var listOfRepositories []*github.Repository
	listOption := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 20,
		},
	}
	for {
		repositories, response, err :=
			client.Client.Repositories.ListByOrg(client.Context, client.Config.SourceRepos.Owner, listOption)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != http.StatusOK {
			return nil, errors.New("could not get the response")
		}
		for _, repository := range repositories {
			listOfRepositories = append(listOfRepositories, repository)
		}
		if response.NextPage == 0 {
			break
		}
		// assign next page
		listOption.Page = response.NextPage
	}
	log.Printf("Obtained list of repositories from %s organization\n",
		client.Config.SourceRepos.Owner,
	)
	return listOfRepositories, nil
}

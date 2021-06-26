package main

import (
	"errors"
	client2 "github-contributors-action/internal/pkg/client"
	"github-contributors-action/internal/pkg/configs"
	"github-contributors-action/internal/pkg/templates"
	"github.com/google/go-github/v33/github"
	"log"
	"os"
	"strings"
)

// AppVersion is to maintain the versioning of the application
var AppVersion = ""

const AppName = "GitHub Contributors Action"

func init() {
	if AppVersion == "" {
		AppVersion = "Unknown" // expect to set the version at build time
	}
}

func main() {

	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
	}()

	log.Printf("%s version: %s\n", AppName, AppVersion)

	// Read inputs
	// Need GitHub token
	// Need a source repository to fetch contributors information from
	// - input can be a ssh git URL
	// - input can be a https git URL
	// - input can also be owner/repository name
	// Need input templates file to generate the output in required format
	// Need a output file where the result is written
	// Input templates is optional, default templates is available in the
	// repository
	config, err := readConfiguration()
	if err != nil {
		return
	}

	// Create a client (with token if input, public repositories may not
	// require a token)
	client := client2.NewGHClient(config)

	// Update the input repositories into the config which required
	// client object to be constructed
	err, isUpdated := updateRepositories(client, &config)
	if err != nil {
		return
	}
	if isUpdated {
		log.Println("New client created because of updated config")
		client = client2.NewGHClient(config)
	}

	// Program is intelligent to fetch required information from the input
	// Call GetContributorsStats
	contributors, err := client.GetContributors()
	if err != nil {
		return
	}

	// Generate the resulting file from the input templates
	err = templates.ApplyTemplate(contributors, config)
	if err != nil {
		return
	}
}

func readConfiguration() (configs.Config, error) {
	// Read inputs from ENV
	token := getEnvOrDefault("GITHUB_AUTH_TOKEN", "")
	repo :=
		getEnvOrDefault("SOURCE_GITHUB_REPOSITORY",
			"hyperledger-tooling/github-contributors-action")
	pattern :=
		getEnvOrDefault("CONTRIBUTORS_SECTION_PATTERN", "## Contributors")
	endPattern := getEnvOrDefault("CONTRIBUTORS_SECTION_END_PATTERN", "## Contributions")
	inputTemplate :=
		getEnvOrDefault("INPUT_TEMPLATE_FILE", "assets/minimal.md")
	outputFile := getEnvOrDefault("FILE_WITH_PATTERN", "README.md")

	// Convert input repo to owner/repo form
	repos, err := parseRepoField(repo)
	if err != nil {
		return configs.Config{}, err
	}

	return configs.Config{
		GitHubToken:     token,
		SourceRepos:     repos,
		Pattern:         pattern,
		EndPattern:      endPattern,
		TemplateFile:    inputTemplate,
		FileWithPattern: outputFile,
	}, nil
}

func parseRepoField(inputRepo string) (configs.Repos, error) {
	// Check if the input is in the form owner/repository
	repoParts := strings.Split(inputRepo, "/")

	// if the input is present, error out
	if len(repoParts) < 1 || len(repoParts) > 2 {
		return configs.Repos{},
			errors.New(
				"input should have owner and repository at the least. Expected format owner/repo or owner")
	}

	repos := configs.Repos{
		Owner: repoParts[0],
	}
	if len(repoParts) == 2 {
		// single repository
		repository := github.Repository{Name: &repoParts[1]}
		repos.Repositories = []*github.Repository{&repository}
	}
	// if the input is just the owner name, get the repositories
	// need to update the repositories from GitHub
	return repos, nil
}

// updateRepositories returns error status and the status. If the config
// updated then it is expected that a new client is created with the updated
// config.
func updateRepositories(
	client client2.GitRepoInterface,
	config *configs.Config,
) (error, bool) {
	// Update by querying all repositories
	if len(config.SourceRepos.Repositories) == 0 {
		repos, err := client.GetRepos()
		if err != nil {
			return err, false
		}
		config.SourceRepos.Repositories = repos
		return nil, true
	}
	return nil, false
}

func getEnvOrDefault(env, defaultValue string) string {
	value, isPresent := os.LookupEnv(env)
	if !isPresent {
		value = defaultValue
	}
	return value
}

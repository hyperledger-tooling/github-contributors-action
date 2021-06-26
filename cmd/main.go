package main

import (
	"errors"
	client2 "github-contributors-action/internal/pkg/client"
	"github-contributors-action/internal/pkg/configs"
	"github-contributors-action/internal/pkg/templates"
	"log"
	"os"
	"strings"
)

// AppVersion is to maintain the versioning of the application
var AppVersion string = ""

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
	repoParts := strings.Split(repo, "/")
	if len(repoParts) < 2 {
		return configs.Config{},
			errors.New(
				"input should have owner and repository at the least. Expected format owner/repo")
	}
	repo = strings.Join(repoParts[len(repoParts)-2:], "/")

	return configs.Config{
		GitHubToken:     token,
		SourceRepo:      repo,
		Pattern:         pattern,
		EndPattern:      endPattern,
		TemplateFile:    inputTemplate,
		FileWithPattern: outputFile,
	}, nil
}

func getEnvOrDefault(env, defaultValue string) string {
	value, isPresent := os.LookupEnv(env)
	if !isPresent {
		value = defaultValue
	}
	return value
}

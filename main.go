package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

func main() {

	var err error
	defer func() {
		if err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
	}()

	// Read inputs
	// Need GitHub token
	// Need a source repository to fetch contributors information from
	// - input can be a ssh git URL
	// - input can be a https git URL
	// - input can also be owner/repository name
	// Need input template file to generate the output in required format
	// Need a output file where the result is written
	// Input template is optional, default template is aailable in the
	// repository
	config, err := readConfiguration()
	if err != nil {
		return
	}

	// Create a client (with token if input, public repositories may not
	// require a token)
	client := NewClient(config)
	if err != nil {
		return
	}

	// Program is intelligent to fetch required information from the input
	// Call GetContributorsStats
	contributors, err := client.GetContributors()
	if err != nil {
		return
	}

	// Generate the resulting file from the input template
	err = ApplyTemplate(contributors, config)
	if err != nil {
		return
	}
}

func readConfiguration() (Config, error) {
	// Read inputs from ENV
	token := getEnvOrDefault("GITHUB_AUTH_TOKEN", "")
	repo :=
		getEnvOrDefault("GITHUB_REPOSITORY",
			"arsulegai/github-contributors-action")
	pattern :=
		getEnvOrDefault("CONTRIBUTORS_SECTION_PATTERN", "{{PlaceHolder}}")
	endPattern := getEnvOrDefault("CONTRIBUTORS_SECTION_END_PATTERN", "")
	inputTemplate :=
		getEnvOrDefault("INPUT_TEMPLATE_FILE", "templates/minimal.tpl")
	outputFile := getEnvOrDefault("FILE_WITH_PATTERN", "README.md")

	// Convert input repo to owner/repo form
	repoParts := strings.Split(repo, "/")
	if len(repoParts) < 2 {
		return Config{},
			errors.New(
				"Input should have owner and repository at the least. Expected format owner/repo")
	}
	repo = strings.Join(repoParts[len(repoParts)-2:], "/")

	return Config{
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

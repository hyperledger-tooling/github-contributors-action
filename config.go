package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	GitHubToken     string
	SourceRepo      string
	Pattern         string
	TemplateFile    string
	FileWithPattern string
}

func readConfiguration() (Config, error) {
	// Read inputs from ENV
	token := getEnvOrDefault("GITHUB_AUTH_TOKEN", "")
	repo :=
		getEnvOrDefault("GITHUB_REPOSITORY_URL",
			"arsulegai/github-contributors-action")
	pattern :=
		getEnvOrDefault("CONTRIBUTORS_SECTION_PATTERN", "{{PlaceHolder}}")
	inputTemplate :=
		getEnvOrDefault("INPUT_TEMPLATE_FILE", "templates/minimal.tpl")
	outputFile := getEnvOrDefault("FILE_WITH_PATTERN", "README.md")

	// Convert input repo to owner/repo form
	repoParts := strings.Split(repo, "/")
	if len(repoParts) < 2 {
		return Config{},
			errors.New(
				fmt.Sprintf("%s", "Input should have owner and repository at the least. Expected format owner/repo"))
	}

	return Config{
		GitHubToken:     token,
		SourceRepo:      repo,
		Pattern:         pattern,
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

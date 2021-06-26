package configs

import "github.com/google/go-github/v33/github"

// Config reads the input that the program reads while it is run.
type Config struct {
	GitHubToken     string
	SourceRepos     Repos
	Pattern         string
	EndPattern      string
	TemplateFile    string
	FileWithPattern string
}

// Repos is to maintain structure of the repository.
// Parse the input and store the owner and repository list
// information separately in this case.
type Repos struct {
	Owner        string
	Repositories []*github.Repository
}

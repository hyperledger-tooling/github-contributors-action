package main

type Config struct {
	GitHubToken     string
	SourceRepo      string
	Pattern         string
	EndPattern      string
	TemplateFile    string
	FileWithPattern string
}

package configs

// Config reads the input that the program reads while it is run.
type Config struct {
	GitHubToken     string
	SourceRepo      string
	Pattern         string
	EndPattern      string
	TemplateFile    string
	FileWithPattern string
}

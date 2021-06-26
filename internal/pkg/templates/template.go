package templates

import (
	"fmt"
	"github-contributors-action/internal/pkg/configs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-github/v33/github"
)

func ApplyTemplate(contributors []*github.Contributor, config configs.Config) error {

	// First read the templates file
	// Generate output from the templates
	templ, err := template.ParseFiles(config.TemplateFile)
	if err != nil {
		return err
	}
	templateFileBytes, err := ioutil.ReadFile(config.TemplateFile)
	log.Printf("Before applying templates: %v", string(templateFileBytes))
	fileHandler, err :=
		ioutil.TempFile(filepath.Dir(config.FileWithPattern), "generated")

	err = templ.Execute(fileHandler, contributors)
	if err != nil {
		return err
	}
	afterTemplate, err := ioutil.ReadFile(fileHandler.Name())
	if err != nil {
		return err
	}
	stringToReplace := string(afterTemplate)
	log.Printf("After applying templates: %v", stringToReplace)
	err = os.Remove(fileHandler.Name())
	if err != nil {
		return err
	}

	// Find the pattern from the output file
	// Replace pattern with the generated templates
	fileContents, err := ioutil.ReadFile(config.FileWithPattern)
	if err != nil {
		return err
	}
	fileString := string(fileContents)

	leftString := strings.Split(fileString, config.Pattern)[0]
	rightString := ""
	if config.EndPattern != "" {
		rightString = strings.Split(fileString, config.EndPattern)[1]
	}
	finalString :=
		fmt.Sprintf("%s%s\n%s\n%s%s",
			leftString,
			config.Pattern,
			stringToReplace,
			config.EndPattern,
			rightString,
		)
	log.Printf("Final: %s\n", finalString)

	info, err := os.Stat(config.FileWithPattern)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		config.FileWithPattern,
		[]byte(finalString),
		info.Mode(),
	)
}

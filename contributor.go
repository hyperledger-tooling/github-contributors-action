package main

import "log"

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
	client, err := getClient()
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
	err := GenerateTemplate(contributors, config)
	if err != nil {
		return
	}
}

/*	var client *github.Client
	ctx := context.Background()

	if token == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}



	stats, _, err := client.Repositories.ListContributorsStats(ctx, owner, repository_url)
	if err != nil {
		fmt.Printf("ListContributorsStats returned error: %v\n", err)
	}

	fmt.Printf("stats: %v", stats)*/

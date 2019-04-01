package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Service service
type Service struct {
	Name            string `yaml:"name"`
	GithubRepo      string `yaml:"github_repo,omitempty"`
	Dependent       string `yaml:"dependent,omitempty"`
	Site            string `yaml:"site,omitempty"`
	StargazersCount int
	LicenseName     string
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Requier token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	mapService := map[string]*Service{}
	mustYamlFileToGoObject("serverless-framework.yml", mapService)

	client := github.NewClient(tc)
	for k, v := range mapService {
		v.Name = k
		if v.GithubRepo != "" {
			userAndRepo := strings.Split(v.GithubRepo, "/")
			rep, _, err := client.Repositories.Get(context.Background(), userAndRepo[0], userAndRepo[1])
			if err != nil {
				log.Fatal(err)
			}
			if rep.StargazersCount != nil {
				v.StargazersCount = *rep.StargazersCount
			}
			if rep.License != nil {
				v.LicenseName = rep.License.GetName()
			}
		}
	}

	mustWrite("README.md", mustExecuteReadmemd(mapService))
}

func mustExecuteReadmemd(mapService map[string]*Service) string {
	return mustExecuteTemplateFile("./README.md.tmpl", mapService, nil)
}

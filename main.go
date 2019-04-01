package main

import (
	"context"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Category category
type Category struct {
	Title      string              `yaml:"title"`
	URL        string              `yaml:"url"`
	MapService map[string]*Service `yaml:"services"`
	Services   []*Service          `yaml:"-"`
}

// Service service
type Service struct {
	Name            string `yaml:"name"`
	GithubRepo      string `yaml:"github_repo,omitempty"`
	Dependent       string `yaml:"dependent,omitempty"`
	Site            string `yaml:"site,omitempty"`
	StargazersCount int
	LicenseName     string
}

// Checker checker
type Checker struct {
	Client     *github.Client
	Categories []*Category
}

// NewChecker newChecker
func NewChecker(token string) *Checker {
	if token == "" {
		log.Fatal("Requier token")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Checker{
		Client: github.NewClient(tc),
	}
}
func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Requier token")
	}

	checker := NewChecker(token)
	checker.addCategoies("serverless-framework.yml", "kubernetes-configuration.yml")

	readmeTmpl := mustRead("README.md.tmpl")
	output, err := executeTemplateSource(string(readmeTmpl), checker.Categories, nil)
	if err != nil {
		log.Fatal(err)
	}

	mustWrite("README.md", output)
}

func (c *Checker) addCategoies(categoryPaths ...string) {
	for _, categoryPath := range categoryPaths {
		category := c.createCategory(categoryPath)
		c.Categories = append(c.Categories, category)
	}
}

//func (c *Checker) serviceGroup(metaPath, templatePath string) string {
func (c *Checker) createCategory(metaPath string) *Category {
	category := &Category{}
	mustYamlFileToGoObject(metaPath, category)

	category.Services = []*Service{}
	for k, v := range category.MapService {
		v.Name = k
		if v.GithubRepo != "" {
			userAndRepo := strings.Split(v.GithubRepo, "/")
			if rep, _, err := c.Client.Repositories.Get(context.Background(), userAndRepo[0], userAndRepo[1]); err == nil {
				if rep.StargazersCount != nil {
					v.StargazersCount = *rep.StargazersCount
				}
				if rep.License != nil {
					v.LicenseName = rep.License.GetName()
				}
			} else {
				log.Fatal(err)
			}
		}
		category.Services = append(category.Services, v)
	}
	sortService(category.Services)
	return category
}

func sortService(services []*Service) {
	sort.Slice(services, func(i, j int) bool {
		return services[i].StargazersCount > services[j].StargazersCount
	})
}

package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type IssueRecap struct {
	Title       string
	URL         string
	LastUpdated string
	Labels      []string
}

func main() {
	githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	daysBack := os.Getenv("SINCE")
	var since int
	if daysBack == "" {
		since = -1
	} else {
		since, _ = strconv.Atoi(daysBack)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	orgs := make(map[string][]string)
	orgs["cloudfoundry"] = []string{
		"cf-networking-notes",
		"cf-networking-release",
		"cf-routing-test-helpers",
		"cf-tcp-router",
		"copilot",
		"gorouter",
		"istio-acceptance-tests",
		"istio-release",
		"istio-scaling",
		"istio-workspace",
		"multierror",
		"nats-release",
		"route-registrar",
		"routing-acceptance-tests",
		"routing-api",
		"routing-api-cli",
		"routing-ci",
		"routing-info",
		"routing-perf-release",
		"routing-release",
		"routing-sample-apps",
		"silk",
		"silk-release",
	}

	orgs["cloudfoundry-incubator"] = []string{
		"cfnetworking-cli-api",
		"routing-backup-and-restore-release",
		"uaa-go-client",
	}

	orgs["cloudfoundry-attic"] = []string{
		"tcp-emitter",
	}

	orgs["cloudfoundry-samples"] = []string{
		"logging-route-service",
	}

	summary := collectIssues(ctx, client, orgs, since)
	msg := createMessage(summary)
	fmt.Println(msg)
}

func collectIssues(ctx context.Context, client *github.Client, orgs map[string][]string, since int) []IssueRecap {
	var issueSummaries []IssueRecap

	for org, repos := range orgs {

		for _, repo := range repos {
			var options *github.IssueListByRepoOptions
			if since > 0 {
				options = &github.IssueListByRepoOptions{
					Since: time.Now().AddDate(0, 0, -1*since),
				}
			}

			issues, _, _ := client.Issues.ListByRepo(ctx, org, repo, options)
			for _, issue := range issues {
				var labels []string
				for _, l := range issue.Labels {
					labels = append(labels, *l.Name)
				}
				issueSummaries = append(issueSummaries, IssueRecap{
					Title:       *issue.Title,
					URL:         *issue.HTMLURL,
					LastUpdated: issue.UpdatedAt.Format("2006-01-02"),
					Labels:      labels,
				})
			}
		}
	}

	return issueSummaries
}

func createMessage(issues []IssueRecap) string {
	var msg string
	msg = "\nOpen issues sorted by most recently updated\n\n"

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].LastUpdated > issues[j].LastUpdated
	})

	for _, i := range issues {
		msg = msg + fmt.Sprintf("%s: %s %v\n%s\n\n", i.LastUpdated, i.Labels, i.Title, i.URL)
	}

	return msg
}

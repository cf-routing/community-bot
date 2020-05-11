package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type IssueRecap struct {
	Title       string
	URL         string
	LastUpdated string
	Labels      []string
}

type RepoSummary struct {
	Name   string
	Issues []IssueRecap
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

func collectIssues(ctx context.Context, client *github.Client, orgs map[string][]string, since int) []RepoSummary {
	summary := []RepoSummary{}

	for org, repos := range orgs {

		for _, repo := range repos {
			issueSummary := []IssueRecap{}

			var options *github.IssueListByRepoOptions
			if since > 0 {
				options = &github.IssueListByRepoOptions{
					Sort:      "updated",
					Direction: "asc",
					Since:     time.Now().AddDate(0, 0, -1*since),
				}
			} else {
				options = &github.IssueListByRepoOptions{
					Sort:      "updated",
					Direction: "asc",
				}
			}

			issues, _, _ := client.Issues.ListByRepo(ctx, org, repo, options)
			for _, issue := range issues {
				var labels []string
				for _, l := range issue.Labels {
					labels = append(labels, *l.Name)
				}
				issueSummary = append(issueSummary, IssueRecap{
					Title:       *issue.Title,
					URL:         *issue.HTMLURL,
					LastUpdated: issue.UpdatedAt.Format("2006-01-02"),
					Labels:      labels,
				})
			}
			summary = append(summary, RepoSummary{
				Name:   repo,
				Issues: issueSummary,
			})
		}
	}

	return summary
}

func createMessage(issues []RepoSummary) string {
	var msg, repo, issue string
	msg = "\nOpen issues sorted by least recently updated\n\n"
	for _, r := range issues {
		issue = ""
		if len(r.Issues) == 0 {
			continue
		}
		repo = fmt.Sprintf("%s (%d):\n\n", r.Name, len(r.Issues))
		for _, i := range r.Issues {
			issue = issue + fmt.Sprintf("%s\n%s\n%v -- Updated: %s\n\n", i.Title, i.URL, i.Labels, i.LastUpdated)
		}
		msg = msg + repo + issue + "\n"
	}

	return msg
}

func closeIssue(ctx context.Context, c *github.Client, issue github.Issue) {
	msg := "Closing due to lack of activity. Please re-open if the issue persists"
	_, _, err := c.Issues.CreateComment(ctx, "cloudfoundry", *issue.Repository.Name, *issue.Number, &github.IssueComment{
		Body: &msg,
	})
	if err != nil {
		panic(err)
	}
	newState := "closed"
	_, _, err = c.Issues.Edit(ctx, "cloudfoundry", *issue.Repository.Name, *issue.Number, &github.IssueRequest{
		State: &newState,
	})
	if err != nil {
		panic(err)
	}
}

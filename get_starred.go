package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "40723bd8f8b831eb5598b1199338b41254019982"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// get starred repos
	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var starredRepos []github.StarredRepository
	for {
		repos, resp, err := client.Activity.ListStarred("", opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		starredRepos = append(starredRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	fmt.Println(len(starredRepos))

	// TODO: get most starred repos from MAX..1000
}

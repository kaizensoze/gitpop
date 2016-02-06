package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	dat, err := ioutil.ReadFile("auth_token.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	authToken := string(dat)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
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

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
	starOpt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var starredRepos []github.StarredRepository
	for {
		repos, resp, err := client.Activity.ListStarred("", starOpt)
		if err != nil {
			fmt.Println(err)
			break
		}
		starredRepos = append(starredRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		starOpt.ListOptions.Page = resp.NextPage
	}
	fmt.Println(len(starredRepos))

	// get popular repos
	popOpt := &github.SearchOptions{
		Sort: "stars",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var popularRepos []github.Repository
	for {
		searchResults, resp, err := client.Search.Repositories("stars:>=1000",
																													 popOpt)
		if err != nil {
			fmt.Println(err)
			break
		}
		popularRepos = append(popularRepos, searchResults.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Println(resp.NextPage)
		popOpt.ListOptions.Page = resp.NextPage
	}

	for _, repo := range popularRepos {
		fmt.Println(*repo.Name)
	}
	fmt.Println(len(popularRepos))
}

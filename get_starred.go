package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	// "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

func handler(w http.ResponseWriter, r *http.Request) {
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

	var starredRepos []github.Repository
	for {
		repos, resp, err := client.Activity.ListStarred("", starOpt)
		if err != nil {
			fmt.Println(err)
			break
		}

		// add starred repos
		for _, repo := range repos {
			starredRepos = append(starredRepos, *repo.Repository)
		}

		// get next page
		if resp.NextPage == 0 {
			break
		}
		starOpt.ListOptions.Page = resp.NextPage
	}

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

		// add popular repos
		popularRepos = append(popularRepos, searchResults.Repositories...)

		// get next page
		if resp.NextPage == 0 {
			break
		}
		popOpt.ListOptions.Page = resp.NextPage
	}

	// print list
	fmt.Fprint(w, "<html><head></head><body>")
	fmt.Fprint(w, "<ul>")
	for _, repo := range popularRepos {
		fmt.Fprintf(w, "<li><a href=\"%s\">%s/%s (%d)</a></li>", *repo.HTMLURL, *repo.Owner.Login, *repo.Name, *repo.StargazersCount)
	}
	fmt.Fprint(w, "</ul></body></html>")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

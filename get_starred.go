package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil {
		start = math.MaxInt32
	}

	popularRepos := getPopularRepos(start)

	lastRepo := popularRepos[len(popularRepos)-1]
	lastStars := *lastRepo.StargazersCount

	// print list
	fmt.Fprint(w, "<html><head></head><body>")
	fmt.Fprint(w, "<ul>")
	for _, repo := range popularRepos {
		fmt.Fprintf(w, "<li><a href=\"%s\">%s/%s (%d)</a></li>",
								*repo.HTMLURL, *repo.Owner.Login, *repo.Name,
								*repo.StargazersCount)
	}
	fmt.Fprint(w, "</ul>")
	fmt.Fprintf(w, "<div style=\"padding-left: 40px;\"><a href=\"%s/?start=%d\">Next</a></div>", r.URL.Host, lastStars)
	fmt.Fprint(w, "</body></html>")
}

func getPopularRepos(start int) []github.Repository {
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
	q := fmt.Sprintf("stars:<=%d", start)
	for {
		searchResults, resp, err := client.Search.Repositories(q, popOpt)
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

	// remove already starred from popular repos list
	for i, repo := range popularRepos {
		if contains(starredRepos, repo) {
			popularRepos = append(popularRepos[:i], popularRepos[i+1:]...)
		}
	}

	return popularRepos
}

func contains(s []github.Repository, e github.Repository) bool {
    for _, a := range s {
        if *a.ID == *e.ID {
            return true
        }
    }
    return false
}

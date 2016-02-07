package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

var host string
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
	}
	defer db.Close()

	sqlStmt := `
		create table if not exists excludes (id integer not null primary key);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	host = r.URL.Host

	if r.Method == "GET" {
		start, err := strconv.Atoi(r.URL.Query().Get("start"))
		if err != nil {
			start = math.MaxInt32
		}
		printList(w, start)
	} else if r.Method == "POST" {
		// TODO: get id from URL
		addExclude()
	}
}

func printList(w http.ResponseWriter, start int) {
	popularRepos, lastStars := getPopularRepos(start)

	// print list
	fmt.Fprint(w, "<html><head></head><body>")
	fmt.Fprint(w, "<ul>")
	for _, repo := range popularRepos {
		fmt.Fprintf(w, `<li><a href="%s">%s/%s (%d)</a>
												<form action="/exclude?id=%d" method="POST" id="form1"
															style="display: inline;">
													<input type="submit" value="X">
												</form>
										</li>
										`,
			*repo.HTMLURL, *repo.Owner.Login, *repo.Name,
			*repo.StargazersCount, *repo.ID)
	}
	fmt.Fprint(w, "</ul>")
	fmt.Fprintf(w, `<div style="padding-left: 40px;">
										<a href="%s/?start=%d">Next</a>
									</div>`, host, lastStars)
	fmt.Fprint(w, "</body></html>")
}

func getPopularRepos(start int) ([]github.Repository, int) {
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
		Sort:        "stars",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var lastRepo github.Repository
	var lastStars int

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

		lastRepo = popularRepos[len(popularRepos)-1]
		lastStars = *lastRepo.StargazersCount

		// get next page
		if resp.NextPage == 0 {
			break
		}
		popOpt.ListOptions.Page = resp.NextPage
	}

	// get excludes
	rows, err := db.Query("select id from excludes")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var excludes []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		excludes = append(excludes, id)
	}

	// remove already starred from popular repos list
	for i, repo := range popularRepos {
		if contains(starredRepos, repo) || contains2(excludes, *repo.ID) {
			popularRepos = append(popularRepos[:i], popularRepos[i+1:]...)
		}
	}

	return popularRepos, lastStars
}

func addExclude() {
	fmt.Println("addExclude")
}

func contains(s []github.Repository, e github.Repository) bool {
	for _, a := range s {
		if *a.ID == *e.ID {
			return true
		}
	}
	return false
}

func contains2(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

# github-explore

Originally I had created a one-time list of popular github repos in the form of
a gist: https://gist.github.com/kaizensoze/00ccfb395ec8410daec2.

The issue is that the link states in the gist don't update and obviously the
star counts for the repos are going to change so I figured I'd write a simple
list generator that excludes already starred repos as an exercise in Go.

Turns out I'd really rather not program in Go and had to fight the urge to just
rewrite it in Python or Ruby, which is silly, because the best medium for this
would really be a Meteor or React app, with an additional x button next to each
link that removes it from the page and adds the given repo id to a list of
excludes to remember the repos you've gone through but didn't star, even if you
clear your browser history. It'd also use a github app for login instead of
providing an auth token.

...but I'll leave that as an exercise for the reader as this does the trick for
now for occasionally going through and finding new/popular repos to star [as
long as I don't clear my browser history...]

### usage

Create auth token file
```zsh
> auth_token.txt
<AUTH_TOKEN>
^C
```

Build and run the server
```zsh
go build get_starred.go && ./get_starred
```

Go to [http://localhost:4000](http://localhost:4000)

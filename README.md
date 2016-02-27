# github-explore

This flask app searches through the most popular github repos ordered by most
starred, excluding your starred repos. You can remove repos from the list,
which persist to a sqlite database so you don't need to worry about clearing
your browser cache and forgetting which ones you've already looked at.

### setup

Create auth token
```bash
> auth_token.txt
<AUTH_TOKEN>
```

Install Flask 1.0+ and run app
```
pip install https://github.com/mitsuhiko/flask/tarball/master
pip install PyGithub
flask --app=app initdb
flask --app=app run
```

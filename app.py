import os, sys
from sqlite3 import dbapi2 as sqlite3
from flask import Flask, request, session, g, redirect, url_for, abort, \
     render_template
from github import Github

f = open('auth_token.txt', 'r')
auth_token = f.read().strip()
github = Github(auth_token, per_page=100)

app = Flask(__name__)

app.config.update(dict(
    DATABASE=os.path.join(app.root_path, 'app.db'),
    DEBUG=True,
    SECRET_KEY='development key',
    USERNAME='admin',
    PASSWORD='default'
))
app.config.from_envvar('FLASKR_SETTINGS', silent=True)

def connect_db():
    """Connects to the specific database."""
    rv = sqlite3.connect(app.config['DATABASE'])
    rv.row_factory = sqlite3.Row
    return rv

def init_db():
    """Initializes the database."""
    db = get_db()
    with app.open_resource('schema.sql', mode='r') as f:
        db.cursor().executescript(f.read())
    db.commit()

@app.cli.command('initdb')
def initdb_command():
    """Creates the database tables."""
    init_db()
    print('Initialized the database.')

def get_db():
    """Opens a new database connection if there is none yet for the
    current application context.
    """
    if not hasattr(g, 'sqlite_db'):
        g.sqlite_db = connect_db()
    return g.sqlite_db

@app.teardown_appcontext
def close_db(error):
    """Closes the database again at the end of the request."""
    if hasattr(g, 'sqlite_db'):
        g.sqlite_db.close()

@app.route('/')
def index():
    start = 999999999
    if request.args.get('start'):
        start = request.args.get('start')

    starred_repos = github.get_user().get_starred()
    popular_repos = get_popular_repos(start)[:500]

    starred_map = {x.id: x for x in starred_repos}
    popular_map = {x.id: x for x in popular_repos}

    starred_ids = set([x.id for x in starred_repos])
    popular_ids = set([x.id for x in popular_repos])
    result_ids = popular_ids.difference(starred_ids)

    result_repos = [popular_map[id] for id in result_ids]
    result_repos.sort(key=lambda x: x.stargazers_count, reverse=True)

    # TODO: get ignored

    sample = [{
        "id": 1,
        "full_name": "blah",
        "html_url": "blah",
        "stargazers_count": 500
    }]

    return render_template('index.html', starred_repos=result_repos)

@app.route('/ignore', methods=['POST'])
def ignore():
    pass

def get_popular_repos(start):
    query = "stars:<{0}".format(start)
    popular_repos = github.search_repositories(
        query,
        sort="stars",
        order="desc"
    )
    return popular_repos

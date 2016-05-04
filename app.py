import json, os, sys
from sqlite3 import dbapi2 as sqlite3
from flask import Flask, request, session, g, redirect, url_for, abort, \
     render_template, jsonify
from github import Github
from flask_github import GitHub

# github
auth_token = open('auth_token.txt', 'r').read().strip()
github = Github(auth_token, per_page=100)

# flask app
app = Flask(__name__)

app.config.update(dict(
    DATABASE=os.path.join(app.root_path, 'app.db'),
    DEBUG=True,
    SECRET_KEY='development key',
    USERNAME='admin',
    PASSWORD='default'
))
app.config.from_envvar('FLASKR_SETTINGS', silent=True)

# flask-github
github_app_info = json.loads(open('github_app.json', 'r').read())
app.config['GITHUB_CLIENT_ID'] = str(github_app_info['Client_ID'])
app.config['GITHUB_CLIENT_SECRET'] = str(github_app_info['Client_Secret'])
github2 = GitHub(app)

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

@app.before_request
def before_request():
    g.user = None
    if 'user_id' in session:
        g.user = get_db().execute('select * from users where access_token = ?',
                    session['user_id']).fetchone()

@app.route('/')
def index():
    #start = 999999999
    #if request.args.get('start'):
    #    start = request.args.get('start')

    #starred_repos = github.get_user().get_starred()
    #popular_repos = get_popular_repos(start)[:500]

    #starred_map = {x.id: x for x in starred_repos}
    #popular_map = {x.id: x for x in popular_repos}

    #starred_ids = set([x.id for x in starred_repos])
    #popular_ids = set([x.id for x in popular_repos])

    ## get ignore ids
    #db = get_db()
    #cur = db.execute('select id from ignores')
    #ignore_ids = set([x[0] for x in cur.fetchall()])

    #result_ids = popular_ids - starred_ids - ignore_ids

    #result_repos = [popular_map[id] for id in result_ids]
    #result_repos.sort(key=lambda x: x.stargazers_count, reverse=True)

    #last_starred = min(x.stargazers_count for x in popular_repos)

    #return render_template('index.html',
    #    starred_repos=result_repos, last_starred=last_starred)
    return render_template('test.html')

@app.route('/ignore', methods=['POST'])
def ignore():
    repo_id = request.form['id']
    starred = request.form['starred']
    db = get_db()
    db.execute('insert or ignore into ignores (id, starred) \
                values (?, ?)', [repo_id, starred])
    db.commit()
    return 'blah'

def get_popular_repos(start):
    query = "stars:<{0}".format(start)
    popular_repos = github.search_repositories(
        query,
        sort="stars",
        order="desc"
    )
    return popular_repos


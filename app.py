import json, os, sys

from flask import Flask, request, session, g, redirect, url_for, abort, \
        render_template, render_template_string, jsonify
from flask_github import GitHub

from github import Github

from sqlalchemy import create_engine, Column, Integer, String, ForeignKey
from sqlalchemy.orm import scoped_session, sessionmaker
from sqlalchemy.ext.declarative import declarative_base

# flask app
app = Flask(__name__, instance_relative_config=True)

app.config.update(dict(
    DATABASE=os.path.join(app.root_path, 'app.db'),
    DEBUG=False,
))
app.config.from_envvar('FLASKR_SETTINGS', silent=True)
app.config.from_pyfile('config.py')

# flask-github
flask_github = GitHub(app)

# sqlalchemy
engine = create_engine("sqlite:///%s" % app.config['DATABASE'])
db_session = scoped_session(sessionmaker(autocommit=False,
                                         autoflush=False,
                                         bind=engine))
Base = declarative_base()
Base.query = db_session.query_property()

class User(Base):
    __tablename__ = 'users'

    id = Column(Integer, primary_key=True)
    username = Column(String(200))
    access_token = Column(String(200))

    def __init__(self, access_token):
        self.access_token = access_token

class Ignore(Base):
    __tablename__ = 'ignores'

    user_id = Column(Integer, ForeignKey("users.id"), primary_key=True)
    id = Column(Integer, primary_key=True)

    def __init__(self, id):
        self.id = id

def init_db():
    Base.metadata.create_all(bind=engine)

@app.cli.command('initdb')
def initdb_command():
    """Creates the database tables."""
    init_db()
    print('Initialized the database.')

@app.cli.command('cleardb')
def cleardb_command():
    for tbl in reversed(Base.metadata.sorted_tables):
        engine.execute(tbl.delete())

@app.teardown_appcontext
def close_db(error):
    """Closes the database again at the end of the request."""
    db_session.remove()

@app.before_request
def before_request():
    g.user = None
    if 'user_id' in session:
        g.user = User.query.get(session['user_id'])

@app.after_request
def after_request(response):
    db_session.remove()
    return response

@app.route('/')
def index():
    if g.user is None:
        t = '<a href="{{ url_for("login") }}">Login</a>'
        return render_template_string(t)

    start = 999999999
    if request.args.get('start'):
        start = request.args.get('start')

    starred_repos = Github(g.user.access_token, per_page=100).get_user().get_starred()
    starred_map = {x.id: x for x in starred_repos}
    starred_ids = set([x.id for x in starred_repos])

    popular_repos = get_popular_repos(start)[:500]
    popular_map = {x.id: x for x in popular_repos}
    popular_ids = set([x.id for x in popular_repos])

    # get ignore ids

    ignores = Ignore.query.filter_by(user_id=g.user.id).all()
    ignore_ids = set([x.id for x in ignores])

    result_ids = popular_ids - starred_ids - ignore_ids

    result_repos = [popular_map[id] for id in result_ids]
    result_repos.sort(key=lambda x: x.stargazers_count, reverse=True)

    last_starred = min(x.stargazers_count for x in popular_repos)

    return render_template('index.html',
        starred_repos=result_repos, last_starred=last_starred)

def get_popular_repos(start):
    query = "stars:<{0}".format(start)
    popular_repos = Github(g.user.access_token, per_page=100).search_repositories(
        query,
        sort="stars",
        order="desc"
    )
    return popular_repos

@app.route('/ignore', methods=['POST'])
def ignore():
    repo_id = request.form['id']
    if g.user is not None:
        ignore = Ignore(repo_id)
        ignore.user_id = g.user.id
        db_session.add(ignore)
        db_session.commit()
    return 'ignore' # doesn't matter

@flask_github.access_token_getter
def token_getter():
    user = g.user
    if user is not None:
        return user.access_token

@app.route('/github-callback')
@flask_github.authorized_handler
def authorized(access_token):
    next_url = request.args.get('next') or url_for('index')
    if access_token is None:
        return redirect(next_url)

    user = User.query.filter_by(access_token=access_token).first()
    if user is None:
        user = User(access_token)
        user_info = Github(access_token).get_user()
        user.username = user_info.login
        db_session.add(user)
    user.access_token = access_token
    db_session.commit()

    session['user_id'] = user.id
    return redirect(next_url)

@app.route('/login')
def login():
    if session.get('user_id', None) is None:
        return flask_github.authorize()
    else:
        return 'Already logged in'

@app.route('/logout')
def logout():
    session.pop('user_id', None)
    return redirect(url_for('index'))


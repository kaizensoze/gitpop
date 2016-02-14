Starred = new Mongo.Collection("starred");

if (Meteor.isClient) {
  // Meteor.call('getStarred', 'starred', function(error, result) {
  //   Session.set('starred', result);
  // });
  //
  // Template.starred.helpers({
  //   starred: function() {
  //     return Session.get('starred');
  //   }
  // });

  Template.starred.helpers({
    starred: function() {
      return Starred.find({});
    }
  });
}

if (Meteor.isServer) {
  var github;

  Meteor.startup(function() {
    var githubAuthToken = Assets.getText("auth_token.txt").trim();
    var GithubApi = Meteor.npmRequire('github4');
    github = new GithubApi({
    });
    github.authenticate({
      type: "oauth",
      token: githubAuthToken
    });

    var starred = getStarred();
    Starred.remove({});
    Starred.insert(starred);
  });

  function getStarred() {
    var starred = Async.runSync(function(done) {
      var starredRepos = [];
      github.activity.getStarredRepos({per_page: 100}, getStarredRepos);
      function getStarredRepos(err, res) {
        if (err) {
          return false;
        }
        starredRepos = starredRepos.concat(res);
        if (github.hasNextPage(res)) {
          github.getNextPage(res, getStarredRepos)
        } else {
          done(null, starredRepos);
        }
      }
    });
    return starred.result;
  }

  Meteor.methods({
  });
}

// github.search.repos({
//     q: "stars:>=20000",
//     sort: "stars",
//     order: "desc"
// }, function(err, res) {
//     for (var itemKey in res['items']) {
//         var item = res['items'][itemKey];
//         var url = item['html_url'];
//         var star_count = item['stargazers_count'];
//         console.log(url + " (" + star_count + ")");
//     }
// });

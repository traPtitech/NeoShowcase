package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	mock_eventbus "github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus/mock"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
	mock_repository "github.com/traPtitech/neoshowcase/pkg/interface/repository/mock"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func newWebhookReceiverHandlerExp(t *testing.T, eventbus eventbus.Bus, verifier usecase.GitPushWebhookService) *httpexpect.Expect {
	t.Helper()

	h := handler.NewWebhookReceiverHandler(eventbus, verifier)
	e := echo.New()
	e.Use(web.WrapContextMiddleware())
	e.POST("/_webhook", web.UnwrapHandler(h))
	httpserver := httptest.NewServer(e)
	t.Cleanup(func() { httpserver.Close() })

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  httpserver.URL,
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func TestWebhookReceiverHandler_HandleRequest(t *testing.T) {
	t.Parallel()
	t.Run("Gitea", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		bus := mock_eventbus.NewMockBus(mockCtrl)
		bus.EXPECT().
			Publish(event.WebhookRepositoryPush, eventbus.Fields{
				"repository_url": "https://git.trap.jp/xxpoxx/test_repo.git",
				"branch":         "heads/master",
			}).
			Times(1)
		repo := mock_repository.NewMockWebhookSecretRepository(mockCtrl)
		repo.EXPECT().
			GetWebhookSecretKeys(gomock.Any(), "https://git.trap.jp/xxpoxx/test_repo.git").
			Return([]string{"hogefugapopopo", "hoge", "fuga"}, nil).
			AnyTimes()
		verifier := usecase.NewGitPushWebhookService(repo)
		e := newWebhookReceiverHandlerExp(t, bus, verifier)

		headers := map[string]string{
			"Content-Type":      "application/json",
			"X-Gitea-Delivery":  "fae3fa5e-221d-4368-bdbe-08e14f6fb926",
			"X-GitHub-Delivery": "fae3fa5e-221d-4368-bdbe-08e14f6fb926",
			"X-GitHub-Event":    "push",
			"X-Gitea-Event":     "push",
			"X-Gitea-Signature": "ef2b2cbe1a2a9c32317df9b4f3de549dd997743e3e082755e214d5bd34c80dfc",
		}

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{
  "secret": "fuga",
  "ref": "refs/heads/master",
  "before": "0b4fb88e00b0d80b62abce08eb7034b70cfc9704",
  "after": "4b724983e8de8c00227d5f14aa1da5d3e3682d00",
  "compare_url": "https://git.trap.jp/xxpoxx/test_repo/compare/0b4fb88e00b0d80b62abce08eb7034b70cfc9704...4b724983e8de8c00227d5f14aa1da5d3e3682d00",
  "commits": [
    {
      "id": "4b724983e8de8c00227d5f14aa1da5d3e3682d00",
      "message": "'README.md' を更新\n",
      "url": "https://git.trap.jp/xxpoxx/test_repo/commit/4b724983e8de8c00227d5f14aa1da5d3e3682d00",
      "author": {
        "name": "Hiroki Sugiyama",
        "email": "xxpoxx@trap.jp",
        "username": "xxpoxx"
      },
      "committer": {
        "name": "Hiroki Sugiyama",
        "email": "xxpoxx@trap.jp",
        "username": "xxpoxx"
      },
      "verification": null,
      "timestamp": "2020-12-18T00:04:30+09:00",
      "added": [],
      "removed": [],
      "modified": [
        "README.md"
      ]
    }
  ],
  "head_commit": null,
  "repository": {
    "id": 2062,
    "owner": {
      "id": 391,
      "login": "xxpoxx",
      "full_name": "Hiroki Sugiyama",
      "email": "xxpoxx@trap.jp",
      "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
      "language": "ja-JP",
      "is_admin": false,
      "last_login": "2020-12-17T11:08:46+09:00",
      "created": "2019-05-07T17:28:05+09:00",
      "username": "xxpoxx"
    },
    "name": "test_repo",
    "full_name": "xxpoxx/test_repo",
    "description": "hoge",
    "empty": false,
    "private": false,
    "fork": false,
    "template": false,
    "parent": null,
    "mirror": false,
    "size": 22,
    "html_url": "https://git.trap.jp/xxpoxx/test_repo",
    "ssh_url": "ssh://git@git.trap.jp:2200/xxpoxx/test_repo.git",
    "clone_url": "https://git.trap.jp/xxpoxx/test_repo.git",
    "original_url": "",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "open_pr_counter": 0,
    "release_counter": 0,
    "default_branch": "master",
    "archived": false,
    "created_at": "2020-12-17T18:43:36+09:00",
    "updated_at": "2020-12-18T00:04:31+09:00",
    "permissions": {
      "admin": true,
      "push": true,
      "pull": true
    },
    "has_issues": true,
    "internal_tracker": {
      "enable_time_tracker": true,
      "allow_only_contributors_to_track_time": true,
      "enable_issue_dependencies": true
    },
    "has_wiki": false,
    "has_pull_requests": true,
    "ignore_whitespace_conflicts": false,
    "allow_merge_commits": true,
    "allow_rebase": true,
    "allow_rebase_explicit": true,
    "allow_squash_merge": true,
    "avatar_url": ""
  },
  "pusher": {
    "id": 391,
    "login": "xxpoxx",
    "full_name": "Hiroki Sugiyama",
    "email": "xxpoxx@trap.jp",
    "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
    "language": "ja-JP",
    "is_admin": false,
    "last_login": "2020-12-17T11:08:46+09:00",
    "created": "2019-05-07T17:28:05+09:00",
    "username": "xxpoxx"
  },
  "sender": {
    "id": 391,
    "login": "xxpoxx",
    "full_name": "Hiroki Sugiyama",
    "email": "xxpoxx@trap.jp",
    "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
    "language": "ja-JP",
    "is_admin": false,
    "last_login": "2020-12-17T11:08:46+09:00",
    "created": "2019-05-07T17:28:05+09:00",
    "username": "xxpoxx"
  }
}`)).
			Expect().
			Status(http.StatusNoContent)

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{
  "secret": "po",
  "ref": "refs/heads/master",
  "before": "0b4fb88e00b0d80b62abce08eb7034b70cfc9704",
  "after": "4b724983e8de8c00227d5f14aa1da5d3e3682d00",
  "compare_url": "https://git.trap.jp/xxpoxx/test_repo/compare/0b4fb88e00b0d80b62abce08eb7034b70cfc9704...4b724983e8de8c00227d5f14aa1da5d3e3682d00",
  "commits": [
    {
      "id": "4b724983e8de8c00227d5f14aa1da5d3e3682d00",
      "message": "'README.md' を更新\n",
      "url": "https://git.trap.jp/xxpoxx/test_repo/commit/4b724983e8de8c00227d5f14aa1da5d3e3682d00",
      "author": {
        "name": "Hiroki Sugiyama",
        "email": "xxpoxx@trap.jp",
        "username": "xxpoxx"
      },
      "committer": {
        "name": "Hiroki Sugiyama",
        "email": "xxpoxx@trap.jp",
        "username": "xxpoxx"
      },
      "verification": null,
      "timestamp": "2020-12-18T00:04:30+09:00",
      "added": [],
      "removed": [],
      "modified": [
        "README.md"
      ]
    }
  ],
  "head_commit": null,
  "repository": {
    "id": 2062,
    "owner": {
      "id": 391,
      "login": "xxpoxx",
      "full_name": "Hiroki Sugiyama",
      "email": "xxpoxx@trap.jp",
      "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
      "language": "ja-JP",
      "is_admin": false,
      "last_login": "2020-12-17T11:08:46+09:00",
      "created": "2019-05-07T17:28:05+09:00",
      "username": "xxpoxx"
    },
    "name": "test_repo",
    "full_name": "xxpoxx/test_repo",
    "description": "hoge",
    "empty": false,
    "private": false,
    "fork": false,
    "template": false,
    "parent": null,
    "mirror": false,
    "size": 22,
    "html_url": "https://git.trap.jp/xxpoxx/test_repo",
    "ssh_url": "ssh://git@git.trap.jp:2200/xxpoxx/test_repo.git",
    "clone_url": "https://git.trap.jp/xxpoxx/test_repo.git",
    "original_url": "",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "open_pr_counter": 0,
    "release_counter": 0,
    "default_branch": "master",
    "archived": false,
    "created_at": "2020-12-17T18:43:36+09:00",
    "updated_at": "2020-12-18T00:04:31+09:00",
    "permissions": {
      "admin": true,
      "push": true,
      "pull": true
    },
    "has_issues": true,
    "internal_tracker": {
      "enable_time_tracker": true,
      "allow_only_contributors_to_track_time": true,
      "enable_issue_dependencies": true
    },
    "has_wiki": false,
    "has_pull_requests": true,
    "ignore_whitespace_conflicts": false,
    "allow_merge_commits": true,
    "allow_rebase": true,
    "allow_rebase_explicit": true,
    "allow_squash_merge": true,
    "avatar_url": ""
  },
  "pusher": {
    "id": 391,
    "login": "xxpoxx",
    "full_name": "Hiroki Sugiyama",
    "email": "xxpoxx@trap.jp",
    "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
    "language": "ja-JP",
    "is_admin": false,
    "last_login": "2020-12-17T11:08:46+09:00",
    "created": "2019-05-07T17:28:05+09:00",
    "username": "xxpoxx"
  },
  "sender": {
    "id": 391,
    "login": "xxpoxx",
    "full_name": "Hiroki Sugiyama",
    "email": "xxpoxx@trap.jp",
    "avatar_url": "https://git.trap.jp/user/avatar/xxpoxx/-1",
    "language": "ja-JP",
    "is_admin": false,
    "last_login": "2020-12-17T11:08:46+09:00",
    "created": "2019-05-07T17:28:05+09:00",
    "username": "xxpoxx"
  }
}`)).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("GitHub", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		bus := mock_eventbus.NewMockBus(mockCtrl)
		bus.EXPECT().
			Publish(event.WebhookRepositoryPush, eventbus.Fields{
				"repository_url": "https://github.com/cskd8/test_repo.git",
				"branch":         "heads/main",
			}).
			Times(1)
		repo := mock_repository.NewMockWebhookSecretRepository(mockCtrl)
		repo.EXPECT().
			GetWebhookSecretKeys(gomock.Any(), "https://github.com/cskd8/test_repo.git").
			Return([]string{"hogefugapopopo", "hoge", "fuga"}, nil).
			AnyTimes()
		verifier := usecase.NewGitPushWebhookService(repo)
		e := newWebhookReceiverHandlerExp(t, bus, verifier)

		headers := map[string]string{
			"content-length":                         "6941",
			"user-agent":                             "GitHub-Hookshot/bb50ca3",
			"accept":                                 "*/*",
			"x-github-delivery":                      "50a1f424-41d0-11eb-95df-667359f4ebdd",
			"x-github-event":                         "push",
			"x-github-hook-id":                       "269355761",
			"x-github-hook-installation-target-id":   "322343963",
			"x-github-hook-installation-target-type": "repository",
			"x-hub-signature":                        "sha1=a530064097b106903eeb3d8b0bdbe80673053d33",
			"x-hub-signature-256":                    "sha256=3afdec1aa5864b246fce7c4204b994c6931c31a15b86eb8f0a606ac326aeacf5",
			"content-type":                           "application/json",
		}

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{"ref":"refs/heads/main","before":"d856906e3f0976156167856a31f338f7a3e0aaca","after":"707812d6ac52fcc296dce853692f172c54100825","repository":{"id":322343963,"node_id":"MDEwOlJlcG9zaXRvcnkzMjIzNDM5NjM=","name":"test_repo","full_name":"cskd8/test_repo","private":false,"owner":{"name":"cskd8","email":"57042565+cskd8@users.noreply.github.com","login":"cskd8","id":57042565,"node_id":"MDQ6VXNlcjU3MDQyNTY1","avatar_url":"https://avatars0.githubusercontent.com/u/57042565?v=4","gravatar_id":"","url":"https://api.github.com/users/cskd8","html_url":"https://github.com/cskd8","followers_url":"https://api.github.com/users/cskd8/followers","following_url":"https://api.github.com/users/cskd8/following{/other_user}","gists_url":"https://api.github.com/users/cskd8/gists{/gist_id}","starred_url":"https://api.github.com/users/cskd8/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/cskd8/subscriptions","organizations_url":"https://api.github.com/users/cskd8/orgs","repos_url":"https://api.github.com/users/cskd8/repos","events_url":"https://api.github.com/users/cskd8/events{/privacy}","received_events_url":"https://api.github.com/users/cskd8/received_events","type":"User","site_admin":false},"html_url":"https://github.com/cskd8/test_repo","description":"for NeoShowcase","fork":false,"url":"https://github.com/cskd8/test_repo","forks_url":"https://api.github.com/repos/cskd8/test_repo/forks","keys_url":"https://api.github.com/repos/cskd8/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/cskd8/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/cskd8/test_repo/teams","hooks_url":"https://api.github.com/repos/cskd8/test_repo/hooks","issue_events_url":"https://api.github.com/repos/cskd8/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/cskd8/test_repo/events","assignees_url":"https://api.github.com/repos/cskd8/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/cskd8/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/cskd8/test_repo/tags","blobs_url":"https://api.github.com/repos/cskd8/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/cskd8/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/cskd8/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/cskd8/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/cskd8/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/cskd8/test_repo/languages","stargazers_url":"https://api.github.com/repos/cskd8/test_repo/stargazers","contributors_url":"https://api.github.com/repos/cskd8/test_repo/contributors","subscribers_url":"https://api.github.com/repos/cskd8/test_repo/subscribers","subscription_url":"https://api.github.com/repos/cskd8/test_repo/subscription","commits_url":"https://api.github.com/repos/cskd8/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/cskd8/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/cskd8/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/cskd8/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/cskd8/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/cskd8/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/cskd8/test_repo/merges","archive_url":"https://api.github.com/repos/cskd8/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/cskd8/test_repo/downloads","issues_url":"https://api.github.com/repos/cskd8/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/cskd8/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/cskd8/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/cskd8/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/cskd8/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/cskd8/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/cskd8/test_repo/deployments","created_at":1608220778,"updated_at":"2020-12-18T12:52:55Z","pushed_at":1608364852,"git_url":"git://github.com/cskd8/test_repo.git","ssh_url":"git@github.com:cskd8/test_repo.git","clone_url":"https://github.com/cskd8/test_repo.git","svn_url":"https://github.com/cskd8/test_repo","homepage":null,"size":4,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"cskd8","email":"57042565+cskd8@users.noreply.github.com"},"sender":{"login":"cskd8","id":57042565,"node_id":"MDQ6VXNlcjU3MDQyNTY1","avatar_url":"https://avatars0.githubusercontent.com/u/57042565?v=4","gravatar_id":"","url":"https://api.github.com/users/cskd8","html_url":"https://github.com/cskd8","followers_url":"https://api.github.com/users/cskd8/followers","following_url":"https://api.github.com/users/cskd8/following{/other_user}","gists_url":"https://api.github.com/users/cskd8/gists{/gist_id}","starred_url":"https://api.github.com/users/cskd8/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/cskd8/subscriptions","organizations_url":"https://api.github.com/users/cskd8/orgs","repos_url":"https://api.github.com/users/cskd8/repos","events_url":"https://api.github.com/users/cskd8/events{/privacy}","received_events_url":"https://api.github.com/users/cskd8/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/cskd8/test_repo/compare/d856906e3f09...707812d6ac52","commits":[{"id":"707812d6ac52fcc296dce853692f172c54100825","tree_id":"825cddc4ebb7d572e5bc5b32f9a6516dddac6ef8","distinct":true,"message":"Update README.md","timestamp":"2020-12-19T17:00:51+09:00","url":"https://github.com/cskd8/test_repo/commit/707812d6ac52fcc296dce853692f172c54100825","author":{"name":"xxpoxx","email":"57042565+cskd8@users.noreply.github.com","username":"cskd8"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"707812d6ac52fcc296dce853692f172c54100825","tree_id":"825cddc4ebb7d572e5bc5b32f9a6516dddac6ef8","distinct":true,"message":"Update README.md","timestamp":"2020-12-19T17:00:51+09:00","url":"https://github.com/cskd8/test_repo/commit/707812d6ac52fcc296dce853692f172c54100825","author":{"name":"xxpoxx","email":"57042565+cskd8@users.noreply.github.com","username":"cskd8"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`)).
			Expect().
			Status(http.StatusNoContent)

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{"ref":"refs/heads/main_hoge","before":"d856906e3f0976156167856a31f338f7a3e0aaca","after":"707812d6ac52fcc296dce853692f172c54100825","repository":{"id":322343963,"node_id":"MDEwOlJlcG9zaXRvcnkzMjIzNDM5NjM=","name":"test_repo","full_name":"cskd8/test_repo","private":false,"owner":{"name":"cskd8","email":"57042565+cskd8@users.noreply.github.com","login":"cskd8","id":57042565,"node_id":"MDQ6VXNlcjU3MDQyNTY1","avatar_url":"https://avatars0.githubusercontent.com/u/57042565?v=4","gravatar_id":"","url":"https://api.github.com/users/cskd8","html_url":"https://github.com/cskd8","followers_url":"https://api.github.com/users/cskd8/followers","following_url":"https://api.github.com/users/cskd8/following{/other_user}","gists_url":"https://api.github.com/users/cskd8/gists{/gist_id}","starred_url":"https://api.github.com/users/cskd8/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/cskd8/subscriptions","organizations_url":"https://api.github.com/users/cskd8/orgs","repos_url":"https://api.github.com/users/cskd8/repos","events_url":"https://api.github.com/users/cskd8/events{/privacy}","received_events_url":"https://api.github.com/users/cskd8/received_events","type":"User","site_admin":false},"html_url":"https://github.com/cskd8/test_repo","description":"for NeoShowcase","fork":false,"url":"https://github.com/cskd8/test_repo","forks_url":"https://api.github.com/repos/cskd8/test_repo/forks","keys_url":"https://api.github.com/repos/cskd8/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/cskd8/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/cskd8/test_repo/teams","hooks_url":"https://api.github.com/repos/cskd8/test_repo/hooks","issue_events_url":"https://api.github.com/repos/cskd8/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/cskd8/test_repo/events","assignees_url":"https://api.github.com/repos/cskd8/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/cskd8/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/cskd8/test_repo/tags","blobs_url":"https://api.github.com/repos/cskd8/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/cskd8/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/cskd8/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/cskd8/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/cskd8/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/cskd8/test_repo/languages","stargazers_url":"https://api.github.com/repos/cskd8/test_repo/stargazers","contributors_url":"https://api.github.com/repos/cskd8/test_repo/contributors","subscribers_url":"https://api.github.com/repos/cskd8/test_repo/subscribers","subscription_url":"https://api.github.com/repos/cskd8/test_repo/subscription","commits_url":"https://api.github.com/repos/cskd8/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/cskd8/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/cskd8/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/cskd8/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/cskd8/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/cskd8/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/cskd8/test_repo/merges","archive_url":"https://api.github.com/repos/cskd8/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/cskd8/test_repo/downloads","issues_url":"https://api.github.com/repos/cskd8/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/cskd8/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/cskd8/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/cskd8/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/cskd8/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/cskd8/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/cskd8/test_repo/deployments","created_at":1608220778,"updated_at":"2020-12-18T12:52:55Z","pushed_at":1608364852,"git_url":"git://github.com/cskd8/test_repo.git","ssh_url":"git@github.com:cskd8/test_repo.git","clone_url":"https://github.com/cskd8/test_repo.git","svn_url":"https://github.com/cskd8/test_repo","homepage":null,"size":4,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"cskd8","email":"57042565+cskd8@users.noreply.github.com"},"sender":{"login":"cskd8","id":57042565,"node_id":"MDQ6VXNlcjU3MDQyNTY1","avatar_url":"https://avatars0.githubusercontent.com/u/57042565?v=4","gravatar_id":"","url":"https://api.github.com/users/cskd8","html_url":"https://github.com/cskd8","followers_url":"https://api.github.com/users/cskd8/followers","following_url":"https://api.github.com/users/cskd8/following{/other_user}","gists_url":"https://api.github.com/users/cskd8/gists{/gist_id}","starred_url":"https://api.github.com/users/cskd8/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/cskd8/subscriptions","organizations_url":"https://api.github.com/users/cskd8/orgs","repos_url":"https://api.github.com/users/cskd8/repos","events_url":"https://api.github.com/users/cskd8/events{/privacy}","received_events_url":"https://api.github.com/users/cskd8/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/cskd8/test_repo/compare/d856906e3f09...707812d6ac52","commits":[{"id":"707812d6ac52fcc296dce853692f172c54100825","tree_id":"825cddc4ebb7d572e5bc5b32f9a6516dddac6ef8","distinct":true,"message":"Update README.md","timestamp":"2020-12-19T17:00:51+09:00","url":"https://github.com/cskd8/test_repo/commit/707812d6ac52fcc296dce853692f172c54100825","author":{"name":"xxpoxx","email":"57042565+cskd8@users.noreply.github.com","username":"cskd8"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"707812d6ac52fcc296dce853692f172c54100825","tree_id":"825cddc4ebb7d572e5bc5b32f9a6516dddac6ef8","distinct":true,"message":"Update README.md","timestamp":"2020-12-19T17:00:51+09:00","url":"https://github.com/cskd8/test_repo/commit/707812d6ac52fcc296dce853692f172c54100825","author":{"name":"xxpoxx","email":"57042565+cskd8@users.noreply.github.com","username":"cskd8"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`)).
			Expect().
			Status(http.StatusBadRequest)
	})
}

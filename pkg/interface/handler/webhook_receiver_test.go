package handler_test

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	mock_eventbus "github.com/traPtitech/neoshowcase/pkg/domain/mock"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	mock_repository "github.com/traPtitech/neoshowcase/pkg/interface/repository/mock"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func newWebhookReceiverHandlerExp(t *testing.T, eventbus domain.Bus, verifier usecase.GitPushWebhookService) *httpexpect.Expect {
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
		rawurl := "https://git.trap.jp/hijiki51/git-test.git"
		bus.EXPECT().
			Publish(event.WebhookRepositoryPush, domain.Fields{
				"repository_url": rawurl,
				"branch":         "heads/master",
			}).
			Times(1)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		repo.EXPECT().
			GetProviderByHost(gomock.Any(), "git.trap.jp").
			Return(domain.Provider{
				ID:     "11ca352c-2556-4b8f-bcbf-1f873d3bb540",
				Secret: "ThisIsSecret",
			}, nil).
			AnyTimes()
		repo.EXPECT().
			GetRepository(gomock.Any(), rawurl).Return(domain.Repository{
			ID:        "9cf4d26d-0f35-474c-a4f2-18c3c7a9ffbf",
			RemoteURL: rawurl,
			Provider: domain.Provider{
				ID:     "11ca352c-2556-4b8f-bcbf-1f873d3bb540",
				Secret: "ThisIsSecret",
			},
		}, nil).AnyTimes()
		verifier := usecase.NewGitPushWebhookService(repo)
		e := newWebhookReceiverHandlerExp(t, bus, verifier)

		headers := map[string]string{
			"Content-Type":      "application/json",
			"X-Gitea-Delivery":  "e5e0b97b-740b-4c8b-8424-102333d9a977",
			"X-GitHub-Delivery": "e5e0b97b-740b-4c8b-8424-102333d9a977",
			"X-GitHub-Event":    "push",
			"X-Gitea-Event":     "push",
			"X-Gitea-Signature": "947e313360746f96ba332ea9d5c546bfbeb65e2472ea233f59477b13e03d380b",
		}

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{
  "ref": "refs/heads/master",
  "before": "aa9b56409f28fc6f8304ddd9313829284d22250b",
  "after": "aa9b56409f28fc6f8304ddd9313829284d22250b",
  "compare_url": "",
  "commits": [
    {
      "id": "aa9b56409f28fc6f8304ddd9313829284d22250b",
      "message": "Update 'README.md'\n",
      "url": "https://git.trap.jp/hijiki51/git-test/commit/aa9b56409f28fc6f8304ddd9313829284d22250b",
      "author": {
        "name": "Hibiki Seki",
        "email": "hibiki0719euph@gmail.com",
        "username": ""
      },
      "committer": {
        "name": "Hibiki Seki",
        "email": "hibiki0719euph@gmail.com",
        "username": ""
      },
      "verification": null,
      "timestamp": "0001-01-01T00:00:00Z",
      "added": null,
      "removed": null,
      "modified": null
    }
  ],
  "head_commit": {
    "id": "aa9b56409f28fc6f8304ddd9313829284d22250b",
    "message": "Update 'README.md'\n",
    "url": "https://git.trap.jp/hijiki51/git-test/commit/aa9b56409f28fc6f8304ddd9313829284d22250b",
    "author": {
      "name": "Hibiki Seki",
      "email": "hibiki0719euph@gmail.com",
      "username": ""
    },
    "committer": {
      "name": "Hibiki Seki",
      "email": "hibiki0719euph@gmail.com",
      "username": ""
    },
    "verification": null,
    "timestamp": "0001-01-01T00:00:00Z",
    "added": null,
    "removed": null,
    "modified": null
  },
  "repository": {
    "id": 2186,
    "owner": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
    "name": "git-test",
    "full_name": "hijiki51/git-test",
    "description": "",
    "empty": false,
    "private": false,
    "fork": false,
    "template": false,
    "parent": null,
    "mirror": false,
    "size": 26,
    "html_url": "https://git.trap.jp/hijiki51/git-test",
    "ssh_url": "ssh://git@git.trap.jp:2200/hijiki51/git-test.git",
    "clone_url": "https://git.trap.jp/hijiki51/git-test.git",
    "original_url": "",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "open_pr_counter": 1,
    "release_counter": 0,
    "default_branch": "master",
    "archived": false,
    "created_at": "2021-05-09T14:53:03+09:00",
    "updated_at": "2021-06-17T00:04:16+09:00",
    "permissions": {
      "admin": false,
      "push": false,
      "pull": false
    },
    "has_issues": true,
    "internal_tracker": {
      "enable_time_tracker": true,
      "allow_only_contributors_to_track_time": true,
      "enable_issue_dependencies": true
    },
    "has_wiki": false,
    "has_pull_requests": true,
    "has_projects": false,
    "ignore_whitespace_conflicts": false,
    "allow_merge_commits": true,
    "allow_rebase": true,
    "allow_rebase_explicit": true,
    "allow_squash_merge": true,
    "default_merge_style": "merge",
    "avatar_url": "",
    "internal": false,
    "mirror_interval": ""
  },
  "pusher": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
  "sender": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"}
}`)).
			Expect().
			Status(http.StatusNoContent)

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{
  "ref": "refs/heads/master_hoge",
  "before": "aa9b56409f28fc6f8304ddd9313829284d22250b",
  "after": "aa9b56409f28fc6f8304ddd9313829284d22250b",
  "compare_url": "",
  "commits": [
    {
      "id": "aa9b56409f28fc6f8304ddd9313829284d22250b",
      "message": "Update 'README.md'\n",
      "url": "https://git.trap.jp/hijiki51/git-test/commit/aa9b56409f28fc6f8304ddd9313829284d22250b",
      "author": {
        "name": "Hibiki Seki",
        "email": "hibiki0719euph@gmail.com",
        "username": ""
      },
      "committer": {
        "name": "Hibiki Seki",
        "email": "hibiki0719euph@gmail.com",
        "username": ""
      },
      "verification": null,
      "timestamp": "0001-01-01T00:00:00Z",
      "added": null,
      "removed": null,
      "modified": null
    }
  ],
  "head_commit": {
    "id": "aa9b56409f28fc6f8304ddd9313829284d22250b",
    "message": "Update 'README.md'\n",
    "url": "https://git.trap.jp/hijiki51/git-test/commit/aa9b56409f28fc6f8304ddd9313829284d22250b",
    "author": {
      "name": "Hibiki Seki",
      "email": "hibiki0719euph@gmail.com",
      "username": ""
    },
    "committer": {
      "name": "Hibiki Seki",
      "email": "hibiki0719euph@gmail.com",
      "username": ""
    },
    "verification": null,
    "timestamp": "0001-01-01T00:00:00Z",
    "added": null,
    "removed": null,
    "modified": null
  },
  "repository": {
    "id": 2186,
    "owner": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
    "name": "git-test",
    "full_name": "hijiki51/git-test",
    "description": "",
    "empty": false,
    "private": false,
    "fork": false,
    "template": false,
    "parent": null,
    "mirror": false,
    "size": 26,
    "html_url": "https://git.trap.jp/hijiki51/git-test",
    "ssh_url": "ssh://git@git.trap.jp:2200/hijiki51/git-test.git",
    "clone_url": "https://git.trap.jp/hijiki51/git-test.git",
    "original_url": "",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "open_pr_counter": 1,
    "release_counter": 0,
    "default_branch": "master",
    "archived": false,
    "created_at": "2021-05-09T14:53:03+09:00",
    "updated_at": "2021-06-17T00:04:16+09:00",
    "permissions": {
      "admin": false,
      "push": false,
      "pull": false
    },
    "has_issues": true,
    "internal_tracker": {
      "enable_time_tracker": true,
      "allow_only_contributors_to_track_time": true,
      "enable_issue_dependencies": true
    },
    "has_wiki": false,
    "has_pull_requests": true,
    "has_projects": false,
    "ignore_whitespace_conflicts": false,
    "allow_merge_commits": true,
    "allow_rebase": true,
    "allow_rebase_explicit": true,
    "allow_squash_merge": true,
    "default_merge_style": "merge",
    "avatar_url": "",
    "internal": false,
    "mirror_interval": ""
  },
  "pusher": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
  "sender": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"}
}`)).
			Expect().
			Status(http.StatusBadRequest)
		headers = map[string]string{
			"Content-Type":      "application/json",
			"X-Gitea-Delivery":  "44f4f939-5021-441c-b650-7be92a9ca6a1",
			"X-GitHub-Delivery": "44f4f939-5021-441c-b650-7be92a9ca6a1",
			"X-GitHub-Event":    "push",
			"X-Gitea-Event":     "push",
			"X-Gitea-Signature": "66abe18cdc2dc36f39ebf715c8792cec1c466a26d333d6c7c28fa29966debb02",
		}
		repo.EXPECT().
			GetRepository(gomock.Any(), rawurl).Return(domain.Repository{}, repository.ErrNotFound).AnyTimes()
		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{
  "ref": "refs/heads/main",
  "before": "42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
  "after": "42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
  "compare_url": "",
  "commits": [
    {
      "id": "42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
      "message": "Merge pull request 'タイトル' (#1) from develop into main\n",
      "url": "https://git.trap.jp/hijiki51/git-lecture/commit/42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
      "author": {
        "name": "Hibiki Seki",
        "email": "hijiki51@trap.jp",
        "username": ""
      },
      "committer": {
        "name": "Hibiki Seki",
        "email": "hijiki51@trap.jp",
        "username": ""
      },
      "verification": null,
      "timestamp": "0001-01-01T00:00:00Z",
      "added": null,
      "removed": null,
      "modified": null
    }
  ],
  "head_commit": {
    "id": "42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
    "message": "Merge pull request 'タイトル' (#1) from develop into main\n",
    "url": "https://git.trap.jp/hijiki51/git-lecture/commit/42b537aa4d424183de921fe1e5eab35ff1f9e3e6",
    "author": {
      "name": "Hibiki Seki",
      "email": "hijiki51@trap.jp",
      "username": ""
    },
    "committer": {
      "name": "Hibiki Seki",
      "email": "hijiki51@trap.jp",
      "username": ""
    },
    "verification": null,
    "timestamp": "0001-01-01T00:00:00Z",
    "added": null,
    "removed": null,
    "modified": null
  },
  "repository": {
    "id": 2193,
    "owner": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
    "name": "git-lecture",
    "full_name": "hijiki51/git-lecture",
    "description": "",
    "empty": false,
    "private": false,
    "fork": false,
    "template": false,
    "parent": null,
    "mirror": false,
    "size": 26,
    "html_url": "https://git.trap.jp/hijiki51/git-lecture",
    "ssh_url": "ssh://git@git.trap.jp:2200/hijiki51/git-lecture.git",
    "clone_url": "https://git.trap.jp/hijiki51/git-lecture.git",
    "original_url": "",
    "website": "",
    "stars_count": 0,
    "forks_count": 0,
    "watchers_count": 1,
    "open_issues_count": 0,
    "open_pr_counter": 0,
    "release_counter": 0,
    "default_branch": "main",
    "archived": false,
    "created_at": "2021-05-09T15:55:49+09:00",
    "updated_at": "2021-05-09T16:54:25+09:00",
    "permissions": {
      "admin": false,
      "push": false,
      "pull": false
    },
    "has_issues": true,
    "internal_tracker": {
      "enable_time_tracker": true,
      "allow_only_contributors_to_track_time": true,
      "enable_issue_dependencies": true
    },
    "has_wiki": false,
    "has_pull_requests": true,
    "has_projects": false,
    "ignore_whitespace_conflicts": false,
    "allow_merge_commits": true,
    "allow_rebase": true,
    "allow_rebase_explicit": true,
    "allow_squash_merge": true,
    "default_merge_style": "merge",
    "avatar_url": "",
    "internal": false,
    "mirror_interval": ""
  },
  "pusher": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"},
  "sender": {"id":538,"login":"hijiki51","full_name":"Hibiki Seki","email":"hibiki0719euph@gmail.com","avatar_url":"https://git.trap.jp/user/avatar/hijiki51/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2020-05-12T11:15:51+09:00","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":2,"username":"hijiki51"}
}`)).
			Expect().
			Status(http.StatusNotFound)
	})
	t.Run("GitHub", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		rawurl := "https://github.com/hijiki51/test_repo.git"
		bus := mock_eventbus.NewMockBus(mockCtrl)
		bus.EXPECT().
			Publish(event.WebhookRepositoryPush, domain.Fields{
				"repository_url": rawurl,
				"branch":         "heads/main",
			}).
			Times(1)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		repo.EXPECT().
			GetProviderByHost(gomock.Any(), "github.com").
			Return(domain.Provider{
				ID:     "6404c950-9bb8-4e5d-8151-5d053a724011",
				Secret: "ThisIsSecret",
			}, nil).
			AnyTimes()
		repo.EXPECT().
			GetRepository(gomock.Any(), rawurl).Return(domain.Repository{
			ID:        "9cf4d26d-0f35-474c-a4f2-18c3c7a9ffbf",
			RemoteURL: rawurl,
			Provider: domain.Provider{
				ID:     "11ca352c-2556-4b8f-bcbf-1f873d3bb540",
				Secret: "ThisIsSecret",
			},
		}, nil).
			AnyTimes()
		verifier := usecase.NewGitPushWebhookService(repo)
		e := newWebhookReceiverHandlerExp(t, bus, verifier)

		headers := map[string]string{
			"user-agent":                             "GitHub-Hookshot/e32936c",
			"accept":                                 "*/*",
			"x-github-delivery":                      "976b7fa0-5021-11ec-83bd-6b160fad0b02",
			"x-github-event":                         "push",
			"x-github-hook-id":                       "330625770",
			"x-github-hook-installation-target-id":   "432634274",
			"x-github-hook-installation-target-type": "repository",
			"x-hub-signature":                        "sha1=6942d447430487ce9962d4e4e82cef17cc7b2011",
			"x-hub-signature-256":                    "sha256=1eeba6b56e30b2b8ca9e23586e7293b2f44523b22dd84736096371069739d867",
			"content-type":                           "application/json",
		}

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{"ref":"refs/heads/main","before":"9670804420e775c7de34385b3305de417aa58fea","after":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","repository":{"id":432634274,"node_id":"R_kgDOGcl5og","name":"test_repo","full_name":"hijiki51/test_repo","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/test_repo","description":null,"fork":false,"url":"https://github.com/hijiki51/test_repo","forks_url":"https://api.github.com/repos/hijiki51/test_repo/forks","keys_url":"https://api.github.com/repos/hijiki51/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/test_repo/teams","hooks_url":"https://api.github.com/repos/hijiki51/test_repo/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/test_repo/events","assignees_url":"https://api.github.com/repos/hijiki51/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/test_repo/tags","blobs_url":"https://api.github.com/repos/hijiki51/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/test_repo/languages","stargazers_url":"https://api.github.com/repos/hijiki51/test_repo/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/test_repo/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/test_repo/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/test_repo/subscription","commits_url":"https://api.github.com/repos/hijiki51/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/test_repo/merges","archive_url":"https://api.github.com/repos/hijiki51/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/test_repo/downloads","issues_url":"https://api.github.com/repos/hijiki51/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/test_repo/deployments","created_at":1638080015,"updated_at":"2021-11-28T07:57:12Z","pushed_at":1638086574,"git_url":"git://github.com/hijiki51/test_repo.git","ssh_url":"git@github.com:hijiki51/test_repo.git","clone_url":"https://github.com/hijiki51/test_repo.git","svn_url":"https://github.com/hijiki51/test_repo","homepage":null,"size":1,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/test_repo/compare/9670804420e7...39f2505d83b5","commits":[{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`)).
			Expect().
			Status(http.StatusNoContent)

		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{"ref":"refs/heads/main_hoge","before":"9670804420e775c7de34385b3305de417aa58fea","after":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","repository":{"id":432634274,"node_id":"R_kgDOGcl5og","name":"test_repo","full_name":"hijiki51/test_repo","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/test_repo","description":null,"fork":false,"url":"https://github.com/hijiki51/test_repo","forks_url":"https://api.github.com/repos/hijiki51/test_repo/forks","keys_url":"https://api.github.com/repos/hijiki51/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/test_repo/teams","hooks_url":"https://api.github.com/repos/hijiki51/test_repo/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/test_repo/events","assignees_url":"https://api.github.com/repos/hijiki51/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/test_repo/tags","blobs_url":"https://api.github.com/repos/hijiki51/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/test_repo/languages","stargazers_url":"https://api.github.com/repos/hijiki51/test_repo/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/test_repo/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/test_repo/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/test_repo/subscription","commits_url":"https://api.github.com/repos/hijiki51/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/test_repo/merges","archive_url":"https://api.github.com/repos/hijiki51/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/test_repo/downloads","issues_url":"https://api.github.com/repos/hijiki51/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/test_repo/deployments","created_at":1638080015,"updated_at":"2021-11-28T07:57:12Z","pushed_at":1638086574,"git_url":"git://github.com/hijiki51/test_repo.git","ssh_url":"git@github.com:hijiki51/test_repo.git","clone_url":"https://github.com/hijiki51/test_repo.git","svn_url":"https://github.com/hijiki51/test_repo","homepage":null,"size":1,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/test_repo/compare/9670804420e7...39f2505d83b5","commits":[{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`)).
			Expect().
			Status(http.StatusBadRequest)

		repo.EXPECT().
			GetRepository(gomock.Any(), rawurl).
			Return(domain.Repository{}, repository.ErrNotFound).AnyTimes()

		headers = map[string]string{
			"user-agent":                             "GitHub-Hookshot/e32936c",
			"accept":                                 "*/*",
			"x-github-delivery":                      "2b4dbb46-5024-11ec-89cd-286f2bf3a569",
			"x-github-event":                         "push",
			"x-github-hook-id":                       "330627878",
			"x-github-hook-installation-target-id":   "432634274",
			"x-github-hook-installation-target-type": "repository",
			"x-hub-signature":                        "sha1=7209f11d2b353c757188b0ab598aee5d31c8b2aa",
			"x-hub-signature-256":                    "sha256=5614c2951296b61059c280005df8f7c76d3af63d07d202293cb7f4eed8610911",
			"content-type":                           "application/json",
		}
		e.POST("/_webhook").
			WithHeaders(headers).
			WithBytes([]byte(`{"ref":"refs/heads/master","before":"ffe6bae1e3c2318d49139eb59bf00c38885c2c9e","after":"6af3137590565acb17477be4774006ce46730d14","repository":{"id":362542470,"node_id":"MDEwOlJlcG9zaXRvcnkzNjI1NDI0NzA=","name":"git-lecture","full_name":"hijiki51/git-lecture","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/git-lecture","description":"2021年git講習会資料","fork":false,"url":"https://github.com/hijiki51/git-lecture","forks_url":"https://api.github.com/repos/hijiki51/git-lecture/forks","keys_url":"https://api.github.com/repos/hijiki51/git-lecture/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/git-lecture/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/git-lecture/teams","hooks_url":"https://api.github.com/repos/hijiki51/git-lecture/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/git-lecture/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/git-lecture/events","assignees_url":"https://api.github.com/repos/hijiki51/git-lecture/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/git-lecture/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/git-lecture/tags","blobs_url":"https://api.github.com/repos/hijiki51/git-lecture/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/git-lecture/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/git-lecture/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/git-lecture/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/git-lecture/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/git-lecture/languages","stargazers_url":"https://api.github.com/repos/hijiki51/git-lecture/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/git-lecture/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/git-lecture/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/git-lecture/subscription","commits_url":"https://api.github.com/repos/hijiki51/git-lecture/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/git-lecture/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/git-lecture/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/git-lecture/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/git-lecture/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/git-lecture/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/git-lecture/merges","archive_url":"https://api.github.com/repos/hijiki51/git-lecture/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/git-lecture/downloads","issues_url":"https://api.github.com/repos/hijiki51/git-lecture/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/git-lecture/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/git-lecture/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/git-lecture/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/git-lecture/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/git-lecture/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/git-lecture/deployments","created_at":1619628400,"updated_at":"2021-11-28T08:21:20Z","pushed_at":1638087681,"git_url":"git://github.com/hijiki51/git-lecture.git","ssh_url":"git@github.com:hijiki51/git-lecture.git","clone_url":"https://github.com/hijiki51/git-lecture.git","svn_url":"https://github.com/hijiki51/git-lecture","homepage":null,"size":2164,"stargazers_count":0,"watchers_count":0,"language":"HTML","has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":true,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"master","stargazers":0,"master_branch":"master"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/git-lecture/compare/ffe6bae1e3c2...6af313759056","commits":[{"id":"6af3137590565acb17477be4774006ce46730d14","tree_id":"3af9ffebfc1fda7f926f3f132db9153704222a85","distinct":true,"message":"Create README.md","timestamp":"2021-11-28T17:21:21+09:00","url":"https://github.com/hijiki51/git-lecture/commit/6af3137590565acb17477be4774006ce46730d14","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":[]}],"head_commit":{"id":"6af3137590565acb17477be4774006ce46730d14","tree_id":"3af9ffebfc1fda7f926f3f132db9153704222a85","distinct":true,"message":"Create README.md","timestamp":"2021-11-28T17:21:21+09:00","url":"https://github.com/hijiki51/git-lecture/commit/6af3137590565acb17477be4774006ce46730d14","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":[]}}`)).
			Expect().
			Status(http.StatusNotFound)
	})
}

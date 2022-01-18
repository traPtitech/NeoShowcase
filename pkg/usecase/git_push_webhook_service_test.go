package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	mock_repository "github.com/traPtitech/neoshowcase/pkg/interface/repository/mock"
)

func TestGitPushWebhookService_VerifySignature(t *testing.T) {
	t.Parallel()

	t.Run("VarifySignature(Success)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		s := NewGitPushWebhookService(repo)
		repo.EXPECT().
			GetProviderByHost(gomock.Any(), "github.com").
			Return(domain.Provider{
				ID:     "6404c950-9bb8-4e5d-8151-5d053a724011",
				Secret: "ThisIsSecret",
			}, nil).
			AnyTimes()

		valid, err := s.VerifySignature(
			context.Background(),
			"https://github.com/hijiki51/test_repo.git",
			"1eeba6b56e30b2b8ca9e23586e7293b2f44523b22dd84736096371069739d867",
			[]byte(`{"ref":"refs/heads/main","before":"9670804420e775c7de34385b3305de417aa58fea","after":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","repository":{"id":432634274,"node_id":"R_kgDOGcl5og","name":"test_repo","full_name":"hijiki51/test_repo","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/test_repo","description":null,"fork":false,"url":"https://github.com/hijiki51/test_repo","forks_url":"https://api.github.com/repos/hijiki51/test_repo/forks","keys_url":"https://api.github.com/repos/hijiki51/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/test_repo/teams","hooks_url":"https://api.github.com/repos/hijiki51/test_repo/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/test_repo/events","assignees_url":"https://api.github.com/repos/hijiki51/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/test_repo/tags","blobs_url":"https://api.github.com/repos/hijiki51/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/test_repo/languages","stargazers_url":"https://api.github.com/repos/hijiki51/test_repo/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/test_repo/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/test_repo/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/test_repo/subscription","commits_url":"https://api.github.com/repos/hijiki51/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/test_repo/merges","archive_url":"https://api.github.com/repos/hijiki51/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/test_repo/downloads","issues_url":"https://api.github.com/repos/hijiki51/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/test_repo/deployments","created_at":1638080015,"updated_at":"2021-11-28T07:57:12Z","pushed_at":1638086574,"git_url":"git://github.com/hijiki51/test_repo.git","ssh_url":"git@github.com:hijiki51/test_repo.git","clone_url":"https://github.com/hijiki51/test_repo.git","svn_url":"https://github.com/hijiki51/test_repo","homepage":null,"size":1,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/test_repo/compare/9670804420e7...39f2505d83b5","commits":[{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`),
		)
		assert.Nil(t, err)
		assert.True(t, valid)
	})

	t.Run("VarifySignature(Invalid)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		s := NewGitPushWebhookService(repo)
		repo.EXPECT().
			GetProviderByHost(gomock.Any(), "github.com").
			Return(domain.Provider{
				ID:     "6404c950-9bb8-4e5d-8151-5d053a724011",
				Secret: "ThisIsSecret",
			}, nil).
			AnyTimes()

		valid, err := s.VerifySignature(
			context.Background(),
			"https://github.com/hijiki51/test_repo.git",
			"1eeba6b56e30b2b8ca9e23586e7293b2f44523b22dd84736096371069739traP",
			[]byte(`{"ref":"refs/heads/main","before":"9670804420e775c7de34385b3305de417aa58fea","after":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","repository":{"id":432634274,"node_id":"R_kgDOGcl5og","name":"test_repo","full_name":"hijiki51/test_repo","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/test_repo","description":null,"fork":false,"url":"https://github.com/hijiki51/test_repo","forks_url":"https://api.github.com/repos/hijiki51/test_repo/forks","keys_url":"https://api.github.com/repos/hijiki51/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/test_repo/teams","hooks_url":"https://api.github.com/repos/hijiki51/test_repo/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/test_repo/events","assignees_url":"https://api.github.com/repos/hijiki51/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/test_repo/tags","blobs_url":"https://api.github.com/repos/hijiki51/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/test_repo/languages","stargazers_url":"https://api.github.com/repos/hijiki51/test_repo/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/test_repo/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/test_repo/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/test_repo/subscription","commits_url":"https://api.github.com/repos/hijiki51/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/test_repo/merges","archive_url":"https://api.github.com/repos/hijiki51/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/test_repo/downloads","issues_url":"https://api.github.com/repos/hijiki51/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/test_repo/deployments","created_at":1638080015,"updated_at":"2021-11-28T07:57:12Z","pushed_at":1638086574,"git_url":"git://github.com/hijiki51/test_repo.git","ssh_url":"git@github.com:hijiki51/test_repo.git","clone_url":"https://github.com/hijiki51/test_repo.git","svn_url":"https://github.com/hijiki51/test_repo","homepage":null,"size":1,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/test_repo/compare/9670804420e7...39f2505d83b5","commits":[{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`),
		)
		assert.Nil(t, err)
		assert.False(t, valid)
	})

	t.Run("VarifySignature(Provider Not found)", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		s := NewGitPushWebhookService(repo)
		repo.EXPECT().
			GetProviderByHost(gomock.Any(), "github.com").
			Return(domain.Provider{}, repository.ErrNotFound).
			AnyTimes()

		valid, err := s.VerifySignature(
			context.Background(),
			"https://github.com/hijiki51/test_repo.git",
			"1eeba6b56e30b2b8ca9e23586e7293b2f44523b22dd84736096371069739d867",
			[]byte(`{"ref":"refs/heads/main","before":"9670804420e775c7de34385b3305de417aa58fea","after":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","repository":{"id":432634274,"node_id":"R_kgDOGcl5og","name":"test_repo","full_name":"hijiki51/test_repo","private":false,"owner":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"html_url":"https://github.com/hijiki51/test_repo","description":null,"fork":false,"url":"https://github.com/hijiki51/test_repo","forks_url":"https://api.github.com/repos/hijiki51/test_repo/forks","keys_url":"https://api.github.com/repos/hijiki51/test_repo/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hijiki51/test_repo/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hijiki51/test_repo/teams","hooks_url":"https://api.github.com/repos/hijiki51/test_repo/hooks","issue_events_url":"https://api.github.com/repos/hijiki51/test_repo/issues/events{/number}","events_url":"https://api.github.com/repos/hijiki51/test_repo/events","assignees_url":"https://api.github.com/repos/hijiki51/test_repo/assignees{/user}","branches_url":"https://api.github.com/repos/hijiki51/test_repo/branches{/branch}","tags_url":"https://api.github.com/repos/hijiki51/test_repo/tags","blobs_url":"https://api.github.com/repos/hijiki51/test_repo/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hijiki51/test_repo/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hijiki51/test_repo/git/refs{/sha}","trees_url":"https://api.github.com/repos/hijiki51/test_repo/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hijiki51/test_repo/statuses/{sha}","languages_url":"https://api.github.com/repos/hijiki51/test_repo/languages","stargazers_url":"https://api.github.com/repos/hijiki51/test_repo/stargazers","contributors_url":"https://api.github.com/repos/hijiki51/test_repo/contributors","subscribers_url":"https://api.github.com/repos/hijiki51/test_repo/subscribers","subscription_url":"https://api.github.com/repos/hijiki51/test_repo/subscription","commits_url":"https://api.github.com/repos/hijiki51/test_repo/commits{/sha}","git_commits_url":"https://api.github.com/repos/hijiki51/test_repo/git/commits{/sha}","comments_url":"https://api.github.com/repos/hijiki51/test_repo/comments{/number}","issue_comment_url":"https://api.github.com/repos/hijiki51/test_repo/issues/comments{/number}","contents_url":"https://api.github.com/repos/hijiki51/test_repo/contents/{+path}","compare_url":"https://api.github.com/repos/hijiki51/test_repo/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hijiki51/test_repo/merges","archive_url":"https://api.github.com/repos/hijiki51/test_repo/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hijiki51/test_repo/downloads","issues_url":"https://api.github.com/repos/hijiki51/test_repo/issues{/number}","pulls_url":"https://api.github.com/repos/hijiki51/test_repo/pulls{/number}","milestones_url":"https://api.github.com/repos/hijiki51/test_repo/milestones{/number}","notifications_url":"https://api.github.com/repos/hijiki51/test_repo/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hijiki51/test_repo/labels{/name}","releases_url":"https://api.github.com/repos/hijiki51/test_repo/releases{/id}","deployments_url":"https://api.github.com/repos/hijiki51/test_repo/deployments","created_at":1638080015,"updated_at":"2021-11-28T07:57:12Z","pushed_at":1638086574,"git_url":"git://github.com/hijiki51/test_repo.git","ssh_url":"git@github.com:hijiki51/test_repo.git","clone_url":"https://github.com/hijiki51/test_repo.git","svn_url":"https://github.com/hijiki51/test_repo","homepage":null,"size":1,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":0,"license":null,"allow_forking":true,"is_template":false,"topics":[],"visibility":"public","forks":0,"open_issues":0,"watchers":0,"default_branch":"main","stargazers":0,"master_branch":"main"},"pusher":{"name":"hijiki51","email":"hibiki0719euph@gmail.com"},"sender":{"login":"hijiki51","id":19515624,"node_id":"MDQ6VXNlcjE5NTE1NjI0","avatar_url":"https://avatars.githubusercontent.com/u/19515624?v=4","gravatar_id":"","url":"https://api.github.com/users/hijiki51","html_url":"https://github.com/hijiki51","followers_url":"https://api.github.com/users/hijiki51/followers","following_url":"https://api.github.com/users/hijiki51/following{/other_user}","gists_url":"https://api.github.com/users/hijiki51/gists{/gist_id}","starred_url":"https://api.github.com/users/hijiki51/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hijiki51/subscriptions","organizations_url":"https://api.github.com/users/hijiki51/orgs","repos_url":"https://api.github.com/users/hijiki51/repos","events_url":"https://api.github.com/users/hijiki51/events{/privacy}","received_events_url":"https://api.github.com/users/hijiki51/received_events","type":"User","site_admin":false},"created":false,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hijiki51/test_repo/compare/9670804420e7...39f2505d83b5","commits":[{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}],"head_commit":{"id":"39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","tree_id":"a9f79354b4722c066b68f9fc0ed3b70f0a3f8c25","distinct":true,"message":"Update README.md","timestamp":"2021-11-28T17:02:54+09:00","url":"https://github.com/hijiki51/test_repo/commit/39f2505d83b5e9ed78bc2ca6db423c68ca29fe38","author":{"name":"hijiki51","email":"hibiki0719euph@gmail.com","username":"hijiki51"},"committer":{"name":"GitHub","email":"noreply@github.com","username":"web-flow"},"added":[],"removed":[],"modified":["README.md"]}}`),
		)
		assert.Equal(t, err, fmt.Errorf("provider not found"))
		assert.False(t, valid)
	})

	t.Run("CheckRepositoryExists(Success)", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		s := NewGitPushWebhookService(repo)
		repo.EXPECT().
			GetRepository(context.Background(), "https://github.com/hijiki51/test_repo.git").
			Return(domain.Repository{
				ID:        "9cf4d26d-0f35-474c-a4f2-18c3c7a9ffbf",
				RemoteURL: "https://github.com/hijiki51/test_repo.git",
				Provider: domain.Provider{
					ID:     "11ca352c-2556-4b8f-bcbf-1f873d3bb540",
					Secret: "ThisIsSecret",
				},
			}, nil).
			AnyTimes()
		exist, err := s.CheckRepositoryExists(
			context.Background(),
			"https://github.com/hijiki51/test_repo.git",
			"hijiki51",
			"test_repo",
		)
		assert.Nil(t, err)
		assert.True(t, exist)
	})

	t.Run("CheckRepositoryExists(Repository Not Found)", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		repo := mock_repository.NewMockGitRepositoryRepository(mockCtrl)
		s := NewGitPushWebhookService(repo)
		repo.EXPECT().
			GetRepository(context.Background(), "https://github.com/hijiki51/test_repo.git").
			Return(domain.Repository{}, repository.ErrNotFound).
			AnyTimes()
		exist, err := s.CheckRepositoryExists(
			context.Background(),
			"https://github.com/hijiki51/test_repo.git",
			"hijiki51",
			"test_repo",
		)
		assert.Nil(t, err)
		assert.False(t, exist)
	})

}

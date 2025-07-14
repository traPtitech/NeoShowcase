package apiserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const exampleCommitHash = "0123456789abcdef0123456789abcdef01234567"

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		useMariaDB      bool
		appName         string
		fqdn            string
		expectedDBCalls int
		checkEnvVars    bool
	}{
		{
			name:            "create application with MariaDB",
			useMariaDB:      true,
			appName:         "test-app-with-db",
			fqdn:            "create-application-test-db.example.com",
			expectedDBCalls: 1,
			checkEnvVars:    true,
		},
		{
			name:            "create application without MariaDB",
			useMariaDB:      false,
			appName:         "test-app-no-db",
			fqdn:            "create-application-test-no-db.example.com",
			expectedDBCalls: 0,
			checkEnvVars:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitmock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registrymock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			dbManagerMock := &mocks.MariaDBManagerMock{}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitmock),
				apiserver.WithRegistryMock(registrymock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := context.Background()

			user := apiserver.CreateUser(c, "test-user")
			web.SetUser(&ctx, user)

			repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
				Method:   domain.RepositoryAuthMethodBasic,
				Username: "test-user",
				Password: "test-password",
			}))
			if err != nil {
				t.Fatal(err)
			}

			app := &domain.Application{
				ID:           domain.NewID(),
				Name:         tt.appName,
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				DeployType:   domain.DeployTypeRuntime,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: tt.useMariaDB,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           tt.fqdn,
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				PortPublications: []*domain.PortPublication{},
				OwnerIDs:         []string{user.ID},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			// Act
			createdApp, err := svc.CreateApplication(ctx, app)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			diff := cmp.Diff(app, createdApp, cmpopts.EquateApproxTime(time.Second), cmpopts.IgnoreUnexported(domain.BuildConfigRuntimeBuildpack{}))
			if diff != "" {
				t.Errorf("created application should be equal to the input (-want +got):\n%s", diff)
			}
			assert.Len(t, dbManagerMock.CreateCalls(), tt.expectedDBCalls, "mariaDB create calls should match expected")

			if tt.checkEnvVars {
				envs, err := svc.GetEnvironmentVariables(ctx, app.ID)
				if err != nil {
					t.Fatal(err)
				}
				envMap := lo.SliceToMap(envs, func(e *domain.Environment) (string, bool) { return e.Key, true })
				assert.Contains(t, envMap, domain.EnvMariaDBHostnameKey, "mariaDB hostname should be set")
				assert.Contains(t, envMap, domain.EnvMariaDBPortKey, "mariaDB port should be set")
				assert.Contains(t, envMap, domain.EnvMariaDBUserKey, "mariaDB username should be set")
				assert.Contains(t, envMap, domain.EnvMariaDBPasswordKey, "mariaDB password should be set")
				assert.Contains(t, envMap, domain.EnvMariaDBDatabaseKey, "mariaDB database should be set")
			}
		})
	}
}

func TestCreateApplication_DuplicateURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		firstFQDN   string
		secondFQDN  string
		expectError bool
		errorDesc   string
	}{
		{
			name:        "exact same FQDN should fail",
			firstFQDN:   "duplicate-test.example.com",
			secondFQDN:  "duplicate-test.example.com",
			expectError: true,
			errorDesc:   "should fail when FQDN is exactly the same",
		},
		{
			name:        "different FQDN should succeed",
			firstFQDN:   "first-app.example.com",
			secondFQDN:  "second-app.example.com",
			expectError: false,
			errorDesc:   "should succeed when FQDN is different",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitmock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registrymock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			dbManagerMock := &mocks.MariaDBManagerMock{}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitmock),
				apiserver.WithRegistryMock(registrymock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := context.Background()

			user := apiserver.CreateUser(c, "test-user")
			web.SetUser(&ctx, user)

			repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
				Method:   domain.RepositoryAuthMethodBasic,
				Username: "test-user",
				Password: "test-password",
			}))
			if err != nil {
				t.Fatal(err)
			}

			// Create first application
			firstApp := &domain.Application{
				ID:           domain.NewID(),
				Name:         "first-app",
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				DeployType:   domain.DeployTypeRuntime,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: false,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           tt.firstFQDN,
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				PortPublications: []*domain.PortPublication{},
				OwnerIDs:         []string{user.ID},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			_, err = svc.CreateApplication(ctx, firstApp)
			if err != nil {
				t.Fatal(err)
			}

			// Attempt to create second application with potentially duplicate URL
			secondApp := &domain.Application{
				ID:           domain.NewID(),
				Name:         "second-app",
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				DeployType:   domain.DeployTypeRuntime,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: false,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           tt.secondFQDN,
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				PortPublications: []*domain.PortPublication{},
				OwnerIDs:         []string{user.ID},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			// Act
			_, err = svc.CreateApplication(ctx, secondApp)

			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.errorDesc)
			} else {
				assert.NoError(t, err, tt.errorDesc)
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		useMariaDB           bool
		appName              string
		fqdn                 string
		tags                 []string
		expectedDBDeletes    int
		expectedImageDeletes int
	}{
		{
			name:                 "delete application with MariaDB and multiple tags",
			useMariaDB:           true,
			appName:              "test-app-with-db",
			fqdn:                 "delete-application-test-db.example.com",
			tags:                 []string{"latest", "1.0.0", "1.1.0"},
			expectedDBDeletes:    1,
			expectedImageDeletes: 3,
		},
		{
			name:                 "delete application without MariaDB",
			useMariaDB:           false,
			appName:              "test-app-no-db",
			fqdn:                 "delete-application-test-no-db.example.com",
			tags:                 []string{"latest"},
			expectedDBDeletes:    0,
			expectedImageDeletes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitmock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registrymock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return tt.tags, nil
				},
			}
			dbManagerMock := &mocks.MariaDBManagerMock{
				DeleteFunc: func(ctx context.Context, args domain.DeleteArgs) error {
					return nil
				},
			}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitmock),
				apiserver.WithRegistryMock(registrymock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)

			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := context.Background()

			user := apiserver.CreateUser(c, "test-user")
			web.SetUser(&ctx, user)

			repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
				Method:   domain.RepositoryAuthMethodBasic,
				Username: "test-user",
				Password: "test-password",
			}))
			if err != nil {
				t.Fatal(err)
			}

			app := &domain.Application{
				ID:           domain.NewID(),
				Name:         tt.appName,
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				DeployType:   domain.DeployTypeRuntime,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: tt.useMariaDB,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           tt.fqdn,
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				OwnerIDs: []string{user.ID},
			}
			app, err = svc.CreateApplication(ctx, app)
			if err != nil {
				t.Fatal(err)
			}

			// Act
			err = svc.DeleteApplication(ctx, app.ID)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			_, err = svc.GetApplication(ctx, app.ID)
			assert.ErrorIs(t, err, repository.ErrNotFound, "application should be deleted")

			assert.Len(t, dbManagerMock.DeleteCalls(), tt.expectedDBDeletes, "mariaDB delete calls should match expected")
			assert.Len(t, registrymock.DeleteImageCalls(), tt.expectedImageDeletes, "image delete calls should match expected")
		})
	}
}

func TestUpdateApplication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		originalName string
		updatedName  string
		fqdn         string
	}{
		{
			name:         "update application name",
			originalName: "test-app",
			updatedName:  "updated-app-name",
			fqdn:         "update-application.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitmock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registrymock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitmock),
				apiserver.WithRegistryMock(registrymock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := context.Background()

			user := apiserver.CreateUser(c, "test-user")
			web.SetUser(&ctx, user)

			repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
				Method:   domain.RepositoryAuthMethodBasic,
				Username: "test-user",
				Password: "test-password",
			}))
			if err != nil {
				t.Fatal(err)
			}

			app := &domain.Application{
				ID:           domain.NewID(),
				Name:         tt.originalName,
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				DeployType:   domain.DeployTypeRuntime,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: false,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           tt.fqdn,
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				OwnerIDs: []string{user.ID},
			}
			app, err = svc.CreateApplication(ctx, app)
			if err != nil {
				t.Fatal(err)
			}

			// Act
			app.Name = tt.updatedName
			err = svc.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{
				Name: optional.From(app.Name),
			})
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			updatedApp, err := svc.GetApplication(ctx, app.ID)
			if err != nil {
				t.Fatal(err)
			}
			diff := cmp.Diff(app, updatedApp.App, cmpopts.EquateApproxTime(time.Second), cmpopts.IgnoreUnexported(domain.BuildConfigRuntimeBuildpack{}))
			if diff != "" {
				t.Errorf("updated application should be equal to the input (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateApplication_UpdateMariaDBConfigIsNotAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		initialUseMariaDB bool
		updatedUseMariaDB bool
		expectError       bool
		errorMessage      string
	}{
		{
			name:              "disable MariaDB should fail",
			initialUseMariaDB: true,
			updatedUseMariaDB: false,
			expectError:       true,
			errorMessage:      "should fail to disable mariadb",
		},
		{
			name:              "enable MariaDB should fail",
			initialUseMariaDB: false,
			updatedUseMariaDB: true,
			expectError:       true,
			errorMessage:      "should fail to enable mariadb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitmock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registrymock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitmock),
				apiserver.WithRegistryMock(registrymock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := context.Background()

			user := apiserver.CreateUser(c, "test-user")
			web.SetUser(&ctx, user)

			repo, err := svc.CreateRepository(ctx, "test-repo", "https://example.com/user/repo", optional.From(apiserver.CreateRepositoryAuth{
				Method:   domain.RepositoryAuthMethodBasic,
				Username: "test-user",
				Password: "test-password",
			}))
			if err != nil {
				t.Fatal(err)
			}

			app := &domain.Application{
				ID:           domain.NewID(),
				Name:         "test-app",
				RepositoryID: repo.ID,
				RefName:      "main",
				Commit:       domain.EmptyCommit,
				Config: domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: tt.initialUseMariaDB,
						},
						Context: ".",
					},
				},
				Websites: []*domain.Website{
					{
						ID:             domain.NewID(),
						FQDN:           "update-application-test.example.com",
						PathPrefix:     "/",
						HTTPPort:       80,
						Authentication: domain.AuthenticationTypeOff,
					},
				},
				PortPublications: []*domain.PortPublication{},
				OwnerIDs:         []string{user.ID},
			}
			app, err = svc.CreateApplication(ctx, app)
			if err != nil {
				t.Fatal(err)
			}

			// Act
			err = svc.UpdateApplication(ctx, app.ID, &domain.UpdateApplicationArgs{
				Config: optional.From(domain.ApplicationConfig{
					BuildConfig: &domain.BuildConfigRuntimeBuildpack{
						RuntimeConfig: domain.RuntimeConfig{
							UseMariaDB: tt.updatedUseMariaDB,
						},
					},
				}),
			})

			// Assert
			if tt.expectError {
				assert.Error(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err, "should not fail")
			}
		})
	}
}

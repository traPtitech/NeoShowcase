package apiserver_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/test/mocks"
	"github.com/traPtitech/neoshowcase/pkg/test/testhelper"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

const exampleCommitHash = "0123456789abcdef0123456789abcdef01234567"

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		useMariaDB bool
		appName    string
		fqdn       string
	}{
		{
			name:       "create application with MariaDB",
			useMariaDB: true,
			appName:    "test-app-with-db",
			fqdn:       "create-application-test-db.example.com",
		},
		{
			name:       "create application without MariaDB",
			useMariaDB: false,
			appName:    "test-app-no-db",
			fqdn:       "create-application-test-no-db.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			gitMock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registryMock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			dbManagerMock := &mocks.MariaDBManagerMock{}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitMock),
				apiserver.WithRegistryMock(registryMock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := t.Context()

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
						UseMariaDB: tt.useMariaDB,
						Context:    ".",
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
			expectedDBCalls := lo.Ternary(tt.useMariaDB, 1, 0)
			assert.Len(t, dbManagerMock.CreateCalls(), expectedDBCalls, "mariaDB create calls should match expected")

			if tt.useMariaDB {
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
			gitMock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registryMock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			dbManagerMock := &mocks.MariaDBManagerMock{}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitMock),
				apiserver.WithRegistryMock(registryMock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := t.Context()

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
						UseMariaDB: false,
						Context:    ".",
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
						UseMariaDB: false,
						Context:    ".",
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
			gitMock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registryMock := &mocks.RegistryClientMock{
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
				apiserver.WithGitMock(gitMock),
				apiserver.WithRegistryMock(registryMock),
				apiserver.WithMariaDBManagerMock(dbManagerMock),
			)

			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := t.Context()

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
						UseMariaDB: tt.useMariaDB,
						Context:    ".",
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
			_, typ, _ := domain.DecomposeError(err)
			assert.Equal(t, domain.ErrorTypeNotFound, typ, "application should be deleted")

			assert.Len(t, dbManagerMock.DeleteCalls(), tt.expectedDBDeletes, "mariaDB delete calls should match expected")
			assert.Len(t, registryMock.DeleteImageCalls(), tt.expectedImageDeletes, "image delete calls should match expected")
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
			gitMock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registryMock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitMock),
				apiserver.WithRegistryMock(registryMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := t.Context()

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
						UseMariaDB: false,
						Context:    ".",
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
			gitMock := &mocks.GitServiceMock{
				ResolveRefsFunc: func(ctx context.Context, repo *domain.Repository) (map[string]string, error) {
					return map[string]string{
						"main": exampleCommitHash,
					}, nil
				},
			}
			registryMock := &mocks.RegistryClientMock{
				GetTagsFunc: func(ctx context.Context, image string) ([]string, error) {
					return []string{"latest"}, nil
				},
			}
			c := testhelper.NewContainer(
				apiserver.DefaultOption(t),
				apiserver.WithGitMock(gitMock),
				apiserver.WithRegistryMock(registryMock),
			)
			svc := testhelper.Resolve[*apiserver.Service](c)
			ctx := t.Context()

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
						UseMariaDB: tt.initialUseMariaDB,
						Context:    ".",
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
						UseMariaDB: tt.updatedUseMariaDB,
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

func TestDeleteApplication_DeletesDatabasesNamedBySystemVariables(t *testing.T) {
	t.Parallel()
	f := newAppFixture(t, "delete-app-system-db", appDatabases{MariaDB: true, MongoDB: true})
	mariaDBName, ok := f.env(t, domain.EnvMariaDBDatabaseKey)
	require.True(t, ok)
	// NeoShowcase gives both databases the same name; a distinct MongoDB name
	// tells apart which variable each manager is given.
	const mongoDBName = "nsapp_mongo_only"
	f.setSystemEnv(t, domain.EnvMongoDBDatabaseKey, mongoDBName)

	err := f.svc.DeleteApplication(f.ctx, f.app.ID)

	require.NoError(t, err)
	mariaDBCalls := f.mariaDB.DeleteCalls()
	require.Len(t, mariaDBCalls, 1)
	assert.Equal(t, mariaDBName.Value, mariaDBCalls[0].Args.Database)
	mongoDBCalls := f.mongoDB.DeleteCalls()
	require.Len(t, mongoDBCalls, 1)
	assert.Equal(t, mongoDBName, mongoDBCalls[0].Args.Database)
}

// Applications whose system variables were overwritten before overwriting was
// rejected hold a user-chosen DB name, which may point to another application's DB.
// Deletion must stop before dropping any of the application's databases.
func TestDeleteApplication_RefusesToDeleteDatabaseNamedByUserVariable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		tamperedKey string
	}{
		{name: "mariadb", tamperedKey: domain.EnvMariaDBDatabaseKey},
		{name: "mongodb", tamperedKey: domain.EnvMongoDBDatabaseKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newAppFixture(t, "delete-app-tampered-"+tt.name, appDatabases{MariaDB: true, MongoDB: true})
			f.setUserEnv(t, tt.tamperedKey, "nsapp_someone_else")

			err := f.svc.DeleteApplication(f.ctx, f.app.ID)

			_, typ, _ := domain.DecomposeError(err)
			assert.Equal(t, domain.ErrorTypeFailedPrecondition, typ)
			assert.Empty(t, f.mariaDB.DeleteCalls())
			assert.Empty(t, f.mongoDB.DeleteCalls())
			_, err = f.svc.GetApplication(f.ctx, f.app.ID)
			assert.NoError(t, err, "application should not be deleted")
		})
	}
}

// Once the application and its variables are gone, nothing records the name of
// a database left undropped. Deletion must stop at the first failed drop.
func TestDeleteApplication_StopsWhenDatabaseDeletionFails(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		failMariaDB     bool
		failMongoDB     bool
		wantMongoDBCall bool
	}{
		{name: "mariadb", failMariaDB: true, wantMongoDBCall: false},
		{name: "mongodb", failMongoDB: true, wantMongoDBCall: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newAppFixture(t, "delete-app-db-fails-"+tt.name, appDatabases{MariaDB: true, MongoDB: true})
			failIf := func(fail bool) func(context.Context, domain.DeleteArgs) error {
				return func(context.Context, domain.DeleteArgs) error {
					if fail {
						return errors.New("drop failed")
					}
					return nil
				}
			}
			f.mariaDB.DeleteFunc = failIf(tt.failMariaDB)
			f.mongoDB.DeleteFunc = failIf(tt.failMongoDB)

			err := f.svc.DeleteApplication(f.ctx, f.app.ID)

			assert.Error(t, err)
			assert.Equal(t, tt.wantMongoDBCall, len(f.mongoDB.DeleteCalls()) > 0)
			_, err = f.svc.GetApplication(f.ctx, f.app.ID)
			assert.NoError(t, err, "application should not be deleted")
		})
	}
}

// An application whose database creation fails is left behind, and its owner
// can only clean it up by deleting it. A database that may have been created
// must already be recorded by a system variable, or deletion refuses to drop it.
func TestDeleteApplication_AfterDatabaseCreationFails(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		failMariaDB bool
		failMongoDB bool
	}{
		{name: "mariadb", failMariaDB: true},
		{name: "mongodb", failMongoDB: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newServiceFixture(t)
			failIf := func(fail bool) func(context.Context, domain.CreateArgs) error {
				return func(context.Context, domain.CreateArgs) error {
					if fail {
						return errors.New("create failed")
					}
					return nil
				}
			}
			f.mariaDB.CreateFunc = failIf(tt.failMariaDB)
			f.mongoDB.CreateFunc = failIf(tt.failMongoDB)
			_, err := f.createApp("create-db-fails-"+tt.name, appDatabases{MariaDB: true, MongoDB: true})
			require.Error(t, err)
			apps, err := f.svc.GetApplications(f.ctx, apiserver.GetAppScope{Scope: apiserver.GetAppScopeMine})
			require.NoError(t, err)
			require.Len(t, apps, 1)
			app := apps[0].App

			err = f.svc.DeleteApplication(f.ctx, app.ID)

			require.NoError(t, err)
			mariaDBCalls := f.mariaDB.DeleteCalls()
			require.Len(t, mariaDBCalls, 1)
			assert.Equal(t, domain.DBName(app.ID), mariaDBCalls[0].Args.Database)
			mongoDBCalls := f.mongoDB.DeleteCalls()
			require.Len(t, mongoDBCalls, 1)
			assert.Equal(t, domain.DBName(app.ID), mongoDBCalls[0].Args.Database)
		})
	}
}

func TestCreateApplication_GivesApplicationTheCreatedDatabases(t *testing.T) {
	t.Parallel()
	f := newAppFixture(t, "create-app-db-envs", appDatabases{MariaDB: true, MongoDB: true})
	dbName := domain.DBName(f.app.ID)
	var mariaDBCreates, mongoDBCreates []domain.CreateArgs
	for _, c := range f.mariaDB.CreateCalls() {
		mariaDBCreates = append(mariaDBCreates, c.Args)
	}
	for _, c := range f.mongoDB.CreateCalls() {
		mongoDBCreates = append(mongoDBCreates, c.Args)
	}

	tests := []struct {
		name                                        string
		hostKey, portKey, userKey, passKey, nameKey string
		wantHost                                    string
		wantPort                                    string
		createCalls                                 []domain.CreateArgs
	}{
		{
			name:    "mariadb",
			hostKey: domain.EnvMariaDBHostnameKey, portKey: domain.EnvMariaDBPortKey, userKey: domain.EnvMariaDBUserKey,
			passKey: domain.EnvMariaDBPasswordKey, nameKey: domain.EnvMariaDBDatabaseKey,
			wantHost: "mariadb.example.com", wantPort: "3306",
			createCalls: mariaDBCreates,
		},
		{
			name:    "mongodb",
			hostKey: domain.EnvMongoDBHostnameKey, portKey: domain.EnvMongoDBPortKey, userKey: domain.EnvMongoDBUserKey,
			passKey: domain.EnvMongoDBPasswordKey, nameKey: domain.EnvMongoDBDatabaseKey,
			wantHost: "mongodb.example.com", wantPort: "27017",
			createCalls: mongoDBCreates,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			systemValue := func(key string) string {
				t.Helper()
				env, ok := f.env(t, key)
				require.True(t, ok, "%s should be set", key)
				assert.True(t, env.System, "%s should be a system variable", key)
				return env.Value
			}

			assert.Equal(t, tt.wantHost, systemValue(tt.hostKey))
			assert.Equal(t, tt.wantPort, systemValue(tt.portKey))
			assert.Equal(t, dbName, systemValue(tt.userKey))
			assert.Equal(t, dbName, systemValue(tt.nameKey))
			password := systemValue(tt.passKey)
			assert.NotEmpty(t, password)
			require.Len(t, tt.createCalls, 1)
			assert.Equal(t, domain.CreateArgs{Database: dbName, Password: password}, tt.createCalls[0])
		})
	}
}

// Deletion refuses an application whose database names are not recorded, and
// use_mariadb / use_mongodb cannot be changed afterwards. An application whose
// variables were only partly recorded could therefore never be deleted.
func TestCreateApplication_LeavesNothingWhenRecordingDatabaseFails(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		failKey string
	}{
		{name: "before the first database name", failKey: domain.EnvMariaDBUserKey},
		{name: "after the first database name", failKey: domain.EnvMongoDBHostnameKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newServiceFixture(t)
			require.NoError(t, f.c.Invoke(func(db *sql.DB) error {
				_, err := db.ExecContext(t.Context(), fmt.Sprintf(
					"CREATE TRIGGER fail_env_insert BEFORE INSERT ON environments FOR EACH ROW "+
						"BEGIN IF NEW.`key` = '%s' THEN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'injected'; END IF; END",
					tt.failKey,
				))
				return err
			}))

			_, err := f.createApp("record-db-fails", appDatabases{MariaDB: true, MongoDB: true})

			require.Error(t, err)
			apps, err := f.svc.GetApplications(f.ctx, apiserver.GetAppScope{Scope: apiserver.GetAppScopeMine})
			require.NoError(t, err)
			assert.Empty(t, apps, "application should not be left behind")
			assert.Empty(t, f.mariaDB.CreateCalls())
			assert.Empty(t, f.mongoDB.CreateCalls())
		})
	}
}

func TestApplication_UsesOnlyConfiguredDatabase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		dbs           appDatabases
		unusedNameKey string
	}{
		{name: "mariadb only", dbs: appDatabases{MariaDB: true}, unusedNameKey: domain.EnvMongoDBDatabaseKey},
		{name: "mongodb only", dbs: appDatabases{MongoDB: true}, unusedNameKey: domain.EnvMariaDBDatabaseKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newAppFixture(t, "single-db-"+strings.ReplaceAll(tt.name, " ", "-"), tt.dbs)
			_, ok := f.env(t, tt.unusedNameKey)
			assert.False(t, ok, "%s should not be set", tt.unusedNameKey)

			err := f.svc.DeleteApplication(f.ctx, f.app.ID)

			require.NoError(t, err)
			callCount := func(used bool) int { return lo.Ternary(used, 1, 0) }
			assert.Len(t, f.mariaDB.CreateCalls(), callCount(tt.dbs.MariaDB))
			assert.Len(t, f.mariaDB.DeleteCalls(), callCount(tt.dbs.MariaDB))
			assert.Len(t, f.mongoDB.CreateCalls(), callCount(tt.dbs.MongoDB))
			assert.Len(t, f.mongoDB.DeleteCalls(), callCount(tt.dbs.MongoDB))
		})
	}
}

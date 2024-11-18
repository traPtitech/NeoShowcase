package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationConfig_Validate(t *testing.T) {
	tests := []struct {
		name       string
		deployType DeployType
		config     ApplicationConfig
		wantErr    bool
	}{
		{
			name:       "valid (runtime dockerfile)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeDockerfile{
					DockerfileName: "Dockerfile",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "./main",
					},
					BaseImage: "golang:1.20",
					BuildCmd:  "go build -o main",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with no build cmd (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "python3 main.py",
					},
					BaseImage: "python:3",
					BuildCmd:  "",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with scratch (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "./my-binary",
					},
					BaseImage: "",
					BuildCmd:  "",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty entrypoint cmd (runtime cmd, has base image)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "",
					},
					BaseImage: "php:7-apache",
					BuildCmd:  "",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty entrypoint cmd (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "",
					},
					BaseImage: "",
					BuildCmd:  "",
				},
			},
			wantErr: true,
		},
		{
			name:       "valid (static buildpack)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticBuildpack{
					StaticConfig: StaticConfig{
						ArtifactPath: "./dist",
					},
					Context: "",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty artifact path (static buildpack)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticBuildpack{
					StaticConfig: StaticConfig{
						ArtifactPath: "",
					},
					Context: "",
				},
			},
			wantErr: true,
		},
		{
			name:       "valid (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticDockerfile{
					StaticConfig: StaticConfig{
						ArtifactPath: "./dist",
					},
					DockerfileName: "Dockerfile",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty artifact path (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticDockerfile{
					StaticConfig: StaticConfig{
						ArtifactPath: "",
					},
					DockerfileName: "Dockerfile",
				},
			},
			wantErr: true,
		},
		{
			name:       "valid (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					StaticConfig: StaticConfig{
						ArtifactPath: "./dist",
					},
					BaseImage: "node:18",
					BuildCmd:  "yarn build",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with no build cmd (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					StaticConfig: StaticConfig{
						ArtifactPath: "./dist",
					},
					BaseImage: "alpine:latest",
					BuildCmd:  "",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with scratch (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					StaticConfig: StaticConfig{
						ArtifactPath: "./dist",
					},
					BaseImage: "",
					BuildCmd:  "",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty artifact path (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					StaticConfig: StaticConfig{
						ArtifactPath: "",
					},
					BaseImage: "",
					BuildCmd:  "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate(tt.deployType)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name+" (hash)", func(t *testing.T) {
			xxh3 := tt.config.Hash(nil)
			t.Logf("hash: %v", xxh3)
			assert.Len(t, xxh3, 16)
			for i := 0; i < 5; i++ {
				assert.Equal(t, xxh3, tt.config.Hash(nil))
			}
		})
	}
}

func TestApplication_SelfValidate(t *testing.T) {
	runtimeValidConfig := ApplicationConfig{
		BuildConfig: &BuildConfigRuntimeDockerfile{DockerfileName: "Dockerfile"},
	}
	require.NoError(t, runtimeValidConfig.Validate(DeployTypeRuntime))

	tests := []struct {
		name    string
		app     Application
		wantErr bool
	}{
		{
			name: "valid",
			app: Application{
				Name:         "test",
				RepositoryID: "abc",
				RefName:      "master",
				Commit:       EmptyCommit,
				DeployType:   DeployTypeRuntime,
				Running:      false,
				CurrentBuild: "",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Config:       runtimeValidConfig,
				Websites:     nil,
				OwnerIDs:     []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			app: Application{
				Name:         "",
				RepositoryID: "abc",
				RefName:      "master",
				Commit:       EmptyCommit,
				DeployType:   DeployTypeRuntime,
				Running:      false,
				CurrentBuild: "",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Config:       runtimeValidConfig,
				Websites:     nil,
				OwnerIDs:     []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "empty repository id",
			app: Application{
				Name:         "test",
				RepositoryID: "",
				RefName:      "master",
				Commit:       EmptyCommit,
				DeployType:   DeployTypeRuntime,
				Running:      false,
				CurrentBuild: "",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Config:       runtimeValidConfig,
				Websites:     nil,
				OwnerIDs:     []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "empty owners",
			app: Application{
				Name:         "test",
				RepositoryID: "abc",
				RefName:      "master",
				Commit:       EmptyCommit,
				DeployType:   DeployTypeRuntime,
				Running:      false,
				CurrentBuild: "",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Config:       runtimeValidConfig,
				Websites:     nil,
				OwnerIDs:     []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.app.SelfValidate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashWithEnv(t *testing.T) {
	config := ApplicationConfig{
		BuildConfig: &BuildConfigRuntimeDockerfile{DockerfileName: "Dockerfile"},
	}
	require.NoError(t, config.Validate(DeployTypeRuntime))

	testEnv := []*Environment{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	hash := config.Hash(testEnv)

	t.Run("hash should not be equal when env is different", func(t *testing.T) {
		testEnv2 := []*Environment{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
			{Key: "key3", Value: "value3"},
		}
		assert.NotEqual(t, hash, config.Hash(testEnv2))
	})

	t.Run("order should not matter", func(t *testing.T) {
		reversedTestEnv := []*Environment{
			{Key: "key2", Value: "value2"},
			{Key: "key1", Value: "value1"},
		}
		assert.Equal(t, hash, config.Hash(reversedTestEnv))
	})
}

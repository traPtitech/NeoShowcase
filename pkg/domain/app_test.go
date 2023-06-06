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
			name:       "empty entrypoint cmd (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigRuntimeCmd{
					RuntimeConfig: RuntimeConfig{
						Entrypoint: "",
					},
					BaseImage: "golang:1.20",
					BuildCmd:  "go build -o main",
				},
			},
			wantErr: true,
		},
		{
			name:       "valid (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticDockerfile{
					DockerfileName: "Dockerfile",
					ArtifactPath:   "./dist",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty artifact path (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticDockerfile{
					DockerfileName: "Dockerfile",
					ArtifactPath:   "",
				},
			},
			wantErr: true,
		},
		{
			name:       "valid (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "node:18",
					BuildCmd:     "yarn build",
					ArtifactPath: "./dist",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with no build cmd (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "alpine:latest",
					BuildCmd:     "",
					ArtifactPath: "./dist",
				},
			},
			wantErr: false,
		},
		{
			name:       "valid with scratch (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "",
					BuildCmd:     "",
					ArtifactPath: "./dist",
				},
			},
			wantErr: false,
		},
		{
			name:       "empty artifact path (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "",
					BuildCmd:     "",
					ArtifactPath: "",
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
				Name:          "test",
				RepositoryID:  "abc",
				RefName:       "master",
				DeployType:    DeployTypeRuntime,
				Running:       false,
				CurrentCommit: EmptyCommit,
				WantCommit:    EmptyCommit,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Config:        runtimeValidConfig,
				Websites:      nil,
				OwnerIDs:      []string{"abc"},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			app: Application{
				Name:          "",
				RepositoryID:  "abc",
				RefName:       "master",
				DeployType:    DeployTypeRuntime,
				Running:       false,
				CurrentCommit: EmptyCommit,
				WantCommit:    EmptyCommit,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Config:        runtimeValidConfig,
				Websites:      nil,
				OwnerIDs:      []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "empty repository id",
			app: Application{
				Name:          "test",
				RepositoryID:  "",
				RefName:       "master",
				DeployType:    DeployTypeRuntime,
				Running:       false,
				CurrentCommit: EmptyCommit,
				WantCommit:    EmptyCommit,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Config:        runtimeValidConfig,
				Websites:      nil,
				OwnerIDs:      []string{"abc"},
			},
			wantErr: true,
		},
		{
			name: "empty owners",
			app: Application{
				Name:          "test",
				RepositoryID:  "abc",
				RefName:       "master",
				DeployType:    DeployTypeRuntime,
				Running:       false,
				CurrentCommit: EmptyCommit,
				WantCommit:    EmptyCommit,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				Config:        runtimeValidConfig,
				Websites:      nil,
				OwnerIDs:      []string{},
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

package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func TestApplicationConfig_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		deployType DeployType
		config     ApplicationConfig
		want       bool
	}{
		{
			name:       "valid (runtime dockerfile)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildType: BuildTypeRuntimeDockerfile,
				BuildConfig: &BuildConfigRuntimeDockerfile{
					DockerfileName: "Dockerfile",
				},
			},
			want: true,
		},
		{
			name:       "valid (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildType: BuildTypeRuntimeCmd,
				BuildConfig: &BuildConfigRuntimeCmd{
					BaseImage: "golang:1.20",
					BuildCmd:  "go build -o main",
				},
				Entrypoint: "./main",
			},
			want: true,
		},
		{
			name:       "valid with no build cmd (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildType: BuildTypeRuntimeCmd,
				BuildConfig: &BuildConfigRuntimeCmd{
					BaseImage: "python:3",
					BuildCmd:  "",
				},
				Entrypoint: "python3 main.py",
			},
			want: true,
		},
		{
			name:       "valid with scratch (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildType: BuildTypeRuntimeCmd,
				BuildConfig: &BuildConfigRuntimeCmd{
					BaseImage: "",
					BuildCmd:  "",
				},
				Entrypoint: "./my-binary",
			},
			want: true,
		},
		{
			name:       "empty entrypoint cmd (runtime cmd)",
			deployType: DeployTypeRuntime,
			config: ApplicationConfig{
				BuildType: BuildTypeRuntimeCmd,
				BuildConfig: &BuildConfigRuntimeCmd{
					BaseImage: "golang:1.20",
					BuildCmd:  "go build -o main",
				},
				Entrypoint: "",
			},
			want: false,
		},
		{
			name:       "valid (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticDockerfile,
				BuildConfig: &BuildConfigStaticDockerfile{
					DockerfileName: "Dockerfile",
					ArtifactPath:   "./dist",
				},
			},
			want: true,
		},
		{
			name:       "empty artifact path (static dockerfile)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticDockerfile,
				BuildConfig: &BuildConfigStaticDockerfile{
					DockerfileName: "Dockerfile",
					ArtifactPath:   "",
				},
			},
			want: false,
		},
		{
			name:       "valid (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticCmd,
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "node:18",
					BuildCmd:     "yarn build",
					ArtifactPath: "./dist",
				},
			},
			want: true,
		},
		{
			name:       "valid with no build cmd (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticCmd,
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "alpine:latest",
					BuildCmd:     "",
					ArtifactPath: "./dist",
				},
			},
			want: true,
		},
		{
			name:       "valid with scratch (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticCmd,
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "",
					BuildCmd:     "",
					ArtifactPath: "./dist",
				},
			},
			want: true,
		},
		{
			name:       "empty artifact path (static cmd)",
			deployType: DeployTypeStatic,
			config: ApplicationConfig{
				BuildType: BuildTypeStaticCmd,
				BuildConfig: &BuildConfigStaticCmd{
					BaseImage:    "",
					BuildCmd:     "",
					ArtifactPath: "",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.IsValid(tt.deployType); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IsValid(t *testing.T) {
	runtimeValidConfig := ApplicationConfig{
		BuildType:   BuildTypeRuntimeDockerfile,
		BuildConfig: &BuildConfigRuntimeDockerfile{DockerfileName: "Dockerfile"},
	}
	require.True(t, runtimeValidConfig.IsValid(DeployTypeRuntime))

	tests := []struct {
		name string
		app  Application
		want bool
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
			want: true,
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
			want: false,
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
			want: false,
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
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.app.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		{"ok", "google.com", true},
		{"wildcard ng", "*.trap.show", false},
		{"multi wildcard ng", "*.*.trap.show", false},
		{"wildcard in middle", "trap.*.show", false},
		{"trailing dot ng", "google.com.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDomain(tt.domain); got != tt.want {
				t.Errorf("IsValidDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAvailableDomain_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		{"ok", "google.com", true},
		{"wildcard ok", "*.trap.show", true},
		{"multi wildcard ng", "*.*.trap.show", false},
		{"wildcard in middle", "trap.*.show", false},
		{"trailing dot ng", "google.com.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AvailableDomain{
				Domain: tt.domain,
			}
			if got := a.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAvailableDomain_match(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		target string
		want   bool
	}{
		{"ok", "google.com", "google.com", true},
		{"ng", "google.com", "example.com", false},
		{"wildcard ok", "*.google.com", "test.google.com", true},
		{"wildcard ok2", "*.google.com", "hello.test.google.com", true},
		{"wildcard ng", "*.google.com", "example.com", false},
		{"wildcard ng2", "*.google.com", "test.example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AvailableDomain{
				Domain: tt.domain,
			}
			got := a.match(tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAvailableDomainSlice_IsAvailable(t *testing.T) {
	tests := []struct {
		name string
		s    AvailableDomainSlice
		fqdn string
		want bool
	}{
		{
			name: "empty",
			s:    AvailableDomainSlice{},
			fqdn: "google.com",
			want: false,
		},
		{
			name: "empty (nil)",
			s:    nil,
			fqdn: "google.com",
			want: false,
		},
		{
			name: "ok",
			s:    AvailableDomainSlice{{Domain: "google.com", Available: true}},
			fqdn: "google.com",
			want: true,
		},
		{
			name: "subdomain ok",
			s:    AvailableDomainSlice{{Domain: "*.google.com", Available: true}},
			fqdn: "sub.google.com",
			want: true,
		},
		{
			name: "ng",
			s:    AvailableDomainSlice{{Domain: "google.com", Available: true}},
			fqdn: "yahoo.com",
			want: false,
		},
		{
			name: "specific subdomain ng 1",
			s:    AvailableDomainSlice{{Domain: "*.google.com", Available: true}, {Domain: "sub.google.com", Available: false}},
			fqdn: "sub.google.com",
			want: false,
		},
		{
			name: "specific subdomain ng 2",
			s:    AvailableDomainSlice{{Domain: "sub.google.com", Available: false}, {Domain: "*.google.com", Available: true}},
			fqdn: "sub.google.com",
			want: false,
		},
		{
			name: "specific wildcard subdomain ng 1",
			s:    AvailableDomainSlice{{Domain: "*.sub.google.com", Available: false}, {Domain: "*.google.com", Available: true}},
			fqdn: "test.sub.google.com",
			want: false,
		},
		{
			name: "specific wildcard subdomain ng 2",
			s:    AvailableDomainSlice{{Domain: "*.google.com", Available: true}, {Domain: "*.sub.google.com", Available: false}},
			fqdn: "test.sub.google.com",
			want: false,
		},
		{
			name: "specific wildcard subdomain ok 1",
			s:    AvailableDomainSlice{{Domain: "*.sub.google.com", Available: false}, {Domain: "*.google.com", Available: true}},
			fqdn: "sub.google.com",
			want: true,
		},
		{
			name: "specific wildcard subdomain ok 2",
			s:    AvailableDomainSlice{{Domain: "*.google.com", Available: true}, {Domain: "*.sub.google.com", Available: false}},
			fqdn: "sub.google.com",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.IsAvailable(tt.fqdn)
			assert.Equal(t, tt.want, got)
		})
	}
}

const validSSHKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMwwAAAJCX+efxl/nn
8QAAAAtzc2gtZWQyNTUxOQAAACAC1iAC54T1ooCQN545XcXDPdTxJEEDdt9TsO3MwoPMww
AAAEA+FzwWKIYduEDOqkEOZ2wmxZWPc2wpZeWj+J8e3Q6x0QLWIALnhPWigJA3njldxcM9
1PEkQQN231Ow7czCg8zDAAAADG1vdG9AbW90by13cwE=
-----END OPENSSH PRIVATE KEY-----`

func TestRepository_IsValid(t *testing.T) {
	tests := []struct {
		name string
		repo Repository
		want bool
	}{
		{
			name: "valid auth none (http)",
			repo: Repository{
				Name:     "test",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			want: true,
		},
		{
			name: "valid auth none (https)",
			repo: Repository{
				Name:     "test",
				URL:      "https://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			want: true,
		},
		{
			name: "valid auth basic",
			repo: Repository{
				Name: "test",
				URL:  "https://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method:   RepositoryAuthMethodBasic,
					Username: "username",
					Password: "password",
				}),
				OwnerIDs: []string{"abc"},
			},
			want: true,
		},
		{
			name: "valid auth ssh",
			repo: Repository{
				Name: "test",
				URL:  "git@github.com:traPtitech/NeoShowcase.git",
				Auth: optional.From(RepositoryAuth{
					Method: RepositoryAuthMethodSSH,
					SSHKey: validSSHKey,
				}),
				OwnerIDs: []string{"abc"},
			},
			want: true,
		},
		{
			name: "invalid name",
			repo: Repository{
				Name:     "",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			want: false,
		},
		{
			name: "invalid url",
			repo: Repository{
				Name:     "test",
				URL:      "ttp://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			want: false,
		},
		{
			name: "invalid scheme (auth none)",
			repo: Repository{
				Name:     "test",
				URL:      "git@github.com:traPtitech/NeoShowcase.git",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{"abc"},
			},
			want: false,
		},
		{
			name: "invalid scheme (auth basic)",
			repo: Repository{
				Name: "test",
				URL:  "http://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method:   RepositoryAuthMethodBasic,
					Username: "username",
					Password: "password",
				}),
				OwnerIDs: []string{"abc"},
			},
			want: false,
		},
		{
			name: "invalid scheme (auth ssh)",
			repo: Repository{
				Name: "test",
				URL:  "https://github.com/traPtitech/NeoShowcase",
				Auth: optional.From(RepositoryAuth{
					Method: RepositoryAuthMethodSSH,
					SSHKey: validSSHKey,
				}),
				OwnerIDs: []string{"abc"},
			},
			want: false,
		},
		{
			name: "invalid owners",
			repo: Repository{
				Name:     "test",
				URL:      "http://github.com/traPtitech/NeoShowcase",
				Auth:     optional.Of[RepositoryAuth]{},
				OwnerIDs: []string{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.repo.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepositoryAuth_IsValid(t *testing.T) {
	tests := []struct {
		name string
		auth RepositoryAuth
		want bool
	}{
		{
			name: "valid basic auth",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "root",
				Password: "password",
				SSHKey:   "",
			},
			want: true,
		},
		{
			name: "invalid username",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "",
				Password: "password",
				SSHKey:   "",
			},
			want: false,
		},
		{
			name: "invalid password",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodBasic,
				Username: "root",
				Password: "",
				SSHKey:   "",
			},
			want: false,
		},
		{
			name: "valid ssh auth",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   validSSHKey,
			},
			want: true,
		},
		{
			name: "valid ssh auth (uses default system key)",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   "",
			},
			want: true,
		},
		{
			name: "invalid ssh private key",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey: `-----BEGIN OPENSSH PRIVATE KEY------
-----END OPENSSH PRIVATE KEY-----`,
			},
			want: false,
		},
		{
			name: "invalid ssh auth (public key)",
			auth: RepositoryAuth{
				Method:   RepositoryAuthMethodSSH,
				Username: "",
				Password: "",
				SSHKey:   `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALWIALnhPWigJA3njldxcM91PEkQQN231Ow7czCg8zD`,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.auth.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebsite_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		website Website
		want    bool
	}{
		{"ok1", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: 80}, true},
		{"ok2", Website{FQDN: "google.com", PathPrefix: "/path/to/prefix", HTTPPort: 8080}, true},
		{"invalid fqdn1", Website{FQDN: "google.com.", PathPrefix: "/", HTTPPort: 80}, false},
		{"invalid fqdn2", Website{FQDN: "*.google.com", PathPrefix: "/", HTTPPort: 80}, false},
		{"invalid fqdn3", Website{FQDN: "google.*.com", PathPrefix: "/", HTTPPort: 80}, false},
		{"invalid fqdn4", Website{FQDN: "goo gle.com", PathPrefix: "/", HTTPPort: 80}, false},
		{"invalid fqdn5", Website{FQDN: "no space", PathPrefix: "/", HTTPPort: 80}, false},
		{"invalid path1", Website{FQDN: "google.com", PathPrefix: "", HTTPPort: 80}, false},
		{"invalid path2", Website{FQDN: "google.com", PathPrefix: "../test", HTTPPort: 80}, false},
		{"invalid path3", Website{FQDN: "google.com", PathPrefix: "/test/", HTTPPort: 80}, false},
		{"strip prefix ok1", Website{FQDN: "google.com", PathPrefix: "/", StripPrefix: false, HTTPPort: 80}, true},
		{"strip prefix ok2", Website{FQDN: "google.com", PathPrefix: "/test", StripPrefix: false, HTTPPort: 80}, true},
		{"strip prefix ng", Website{FQDN: "google.com", PathPrefix: "/", StripPrefix: true, HTTPPort: 80}, false},
		{"strip prefix ok3", Website{FQDN: "google.com", PathPrefix: "/test", StripPrefix: true, HTTPPort: 80}, true},
		{"invalid port1", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: -1}, false},
		{"invalid port2", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: 65536}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.website.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebsite_pathComponents(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{"top", "/", []string{}},
		{"first layer", "/test", []string{"test"}},
		{"multiple layers", "/path/to/prefix", []string{"path", "to", "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Website{
				PathPrefix: tt.path,
			}
			if got := w.pathComponents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pathComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebsite_ConflictsWith(t *testing.T) {
	pathTests := []struct {
		name     string
		target   string
		existing []string
		want     bool
	}{
		{"ok1", "/", []string{}, false},
		{"ok2", "/foo", []string{"/api", "/spa"}, false},
		{"ok3", "/api/v2", []string{"/api/v1", "/spa"}, false},
		{"ok4", "/api2", []string{"/api"}, false},
		{"ok5", "/api", []string{"/api2"}, false},
		{"ng1", "/", []string{"/"}, true},
		{"ng2", "/api", []string{"/"}, true},
		{"ng3", "/api/v2", []string{"/api"}, true},
		{"ng4", "/api", []string{"/api"}, true},
	}
	for _, tt := range pathTests {
		t.Run("path "+tt.name, func(t *testing.T) {
			w := &Website{
				PathPrefix: tt.target,
			}
			existingWebsites := lo.Map(tt.existing, func(ex string, i int) *Website {
				return &Website{PathPrefix: ex}
			})
			if got := w.ConflictsWith(existingWebsites); got != tt.want {
				t.Errorf("ConflictsWith() = %v, want %v", got, tt.want)
			}
		})
	}

	fullTests := []struct {
		name     string
		target   *Website
		existing []*Website
		want     bool
	}{
		{
			name:     "ng if same scheme",
			target:   &Website{PathPrefix: "/", HTTPS: false},
			existing: []*Website{{PathPrefix: "/", HTTPS: false}},
			want:     true,
		},
		{
			name:     "ok if different scheme",
			target:   &Website{PathPrefix: "/", HTTPS: true},
			existing: []*Website{{PathPrefix: "/", HTTPS: false}},
			want:     false,
		},
		{
			name:     "ok if different fqdn",
			target:   &Website{FQDN: "google.com", PathPrefix: "/", HTTPS: false},
			existing: []*Website{{FQDN: "yahoo.com", PathPrefix: "/", HTTPS: false}},
			want:     false,
		},
	}
	for _, tt := range fullTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.ConflictsWith(tt.existing); got != tt.want {
				t.Errorf("ConflictsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

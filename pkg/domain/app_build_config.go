package domain

import (
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/mattn/go-shellwords"
)

type BuildType int

const (
	BuildTypeRuntimeBuildpack BuildType = iota
	BuildTypeRuntimeCmd
	BuildTypeRuntimeDockerfile
	BuildTypeStaticCmd
	BuildTypeStaticDockerfile
)

func (b BuildType) DeployType() DeployType {
	switch b {
	case BuildTypeRuntimeBuildpack, BuildTypeRuntimeCmd, BuildTypeRuntimeDockerfile:
		return DeployTypeRuntime
	case BuildTypeStaticCmd, BuildTypeStaticDockerfile:
		return DeployTypeStatic
	default:
		panic(fmt.Sprintf("unknown build type: %v", b))
	}
}

type RuntimeConfig struct {
	UseMariaDB bool
	UseMongoDB bool
	Entrypoint string
	Command    string
}

func parseArgs(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	return shellwords.Parse(s)
}

func (rc *RuntimeConfig) Validate() error {
	if _, err := parseArgs(rc.Entrypoint); err != nil {
		return errors.Wrap(err, "invalid entrypoint")
	}
	if _, err := parseArgs(rc.Command); err != nil {
		return errors.Wrap(err, "invalid command")
	}
	return nil
}

func (rc *RuntimeConfig) MariaDB() bool {
	return rc.UseMariaDB
}

func (rc *RuntimeConfig) MongoDB() bool {
	return rc.UseMongoDB
}

func (rc *RuntimeConfig) EntrypointArgs() []string {
	args, _ := parseArgs(rc.Entrypoint)
	return args
}

func (rc *RuntimeConfig) CommandArgs() []string {
	args, _ := parseArgs(rc.Command)
	return args
}

type StaticConfig struct{}

func (sc *StaticConfig) MariaDB() bool {
	return false
}

func (sc *StaticConfig) MongoDB() bool {
	return false
}

func (sc *StaticConfig) EntrypointArgs() []string {
	panic("no entrypoint for static config")
}

func (sc *StaticConfig) CommandArgs() []string {
	panic("no command for static config")
}

type BuildConfig interface {
	isBuildConfig()
	BuildType() BuildType
	Validate() error
	MariaDB() bool
	MongoDB() bool
	EntrypointArgs() []string
	CommandArgs() []string
}

type buildConfigEmbed struct{}

func (buildConfigEmbed) isBuildConfig() {}

type BuildConfigRuntimeBuildpack struct {
	RuntimeConfig
	Context string
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeBuildpack) BuildType() BuildType {
	return BuildTypeRuntimeBuildpack
}

func (bc *BuildConfigRuntimeBuildpack) Validate() error {
	if err := bc.RuntimeConfig.Validate(); err != nil {
		return err
	}
	// NOTE: context is not necessary
	return nil
}

type BuildConfigRuntimeCmd struct {
	RuntimeConfig
	BaseImage     string
	BuildCmd      string
	BuildCmdShell bool
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeCmd) BuildType() BuildType {
	return BuildTypeRuntimeCmd
}

func (bc *BuildConfigRuntimeCmd) Validate() error {
	if err := bc.RuntimeConfig.Validate(); err != nil {
		return err
	}
	if bc.Entrypoint == "" && bc.Command == "" {
		return errors.New("entrypoint or command is required")
	}
	// NOTE: base image is not necessary (default: scratch)
	// NOTE: build cmd is not necessary
	return nil
}

type BuildConfigRuntimeDockerfile struct {
	RuntimeConfig
	DockerfileName string
	Context        string
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeDockerfile) BuildType() BuildType {
	return BuildTypeRuntimeDockerfile
}

func (bc *BuildConfigRuntimeDockerfile) Validate() error {
	if err := bc.RuntimeConfig.Validate(); err != nil {
		return err
	}
	if bc.DockerfileName == "" {
		return errors.New("dockerfile_name is required")
	}
	// NOTE: Runtime Dockerfile build could have no entrypoint/command but is impossible to catch only from config
	// (can only catch at runtime)
	return nil
}

type BuildConfigStaticCmd struct {
	StaticConfig
	BaseImage     string
	BuildCmd      string
	BuildCmdShell bool
	ArtifactPath  string
	buildConfigEmbed
}

func (bc *BuildConfigStaticCmd) BuildType() BuildType {
	return BuildTypeStaticCmd
}

func (bc *BuildConfigStaticCmd) Validate() error {
	// NOTE: base image is not necessary (default: scratch)
	// NOTE: build cmd is not necessary
	if bc.ArtifactPath == "" {
		return errors.New("artifact_path is required")
	}
	return nil
}

type BuildConfigStaticDockerfile struct {
	StaticConfig
	DockerfileName string
	Context        string
	ArtifactPath   string
	buildConfigEmbed
}

func (bc *BuildConfigStaticDockerfile) BuildType() BuildType {
	return BuildTypeStaticDockerfile
}

func (bc *BuildConfigStaticDockerfile) Validate() error {
	if bc.DockerfileName == "" {
		return errors.New("dockerfile_name is required")
	}
	if bc.ArtifactPath == "" {
		return errors.New("artifact_path is required")
	}
	return nil
}

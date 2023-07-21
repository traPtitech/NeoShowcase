package domain

import (
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/google/shlex"
	"github.com/samber/lo"
)

type BuildType int

const (
	BuildTypeRuntimeBuildpack BuildType = iota
	BuildTypeRuntimeCmd
	BuildTypeRuntimeDockerfile
	BuildTypeStaticBuildpack
	BuildTypeStaticCmd
	BuildTypeStaticDockerfile
)

func (b BuildType) DeployType() DeployType {
	switch b {
	case BuildTypeRuntimeBuildpack, BuildTypeRuntimeCmd, BuildTypeRuntimeDockerfile:
		return DeployTypeRuntime
	case BuildTypeStaticBuildpack, BuildTypeStaticCmd, BuildTypeStaticDockerfile:
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

const shellSpecialCharacters = "`" + `~#$&*()\|[]{};'"<>?!=`

func ParseArgs(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	args, err := shlex.Split(s)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse command")
	}
	shellSyntax := lo.ContainsBy(args, func(arg string) bool {
		return strings.ContainsAny(arg, shellSpecialCharacters)
	})
	if shellSyntax {
		return []string{"sh", "-c", s}, nil
	}
	return args, nil
}

func (rc *RuntimeConfig) Validate() error {
	if _, err := ParseArgs(rc.Entrypoint); err != nil {
		return errors.Wrap(err, "entrypoint")
	}
	if _, err := ParseArgs(rc.Command); err != nil {
		return errors.Wrap(err, "command")
	}
	return nil
}

func (rc *RuntimeConfig) MariaDB() bool {
	return rc.UseMariaDB
}

func (rc *RuntimeConfig) MongoDB() bool {
	return rc.UseMongoDB
}

func (rc *RuntimeConfig) GetRuntimeConfig() RuntimeConfig {
	return *rc
}

func (rc *RuntimeConfig) GetStaticConfig() StaticConfig {
	panic("not static config")
}

type StaticConfig struct {
	ArtifactPath string
	SPA          bool
}

func (sc *StaticConfig) Validate() error {
	if sc.ArtifactPath == "" {
		return errors.New("artifact_path is required for static builds")
	}
	return nil
}

func (sc *StaticConfig) MariaDB() bool {
	return false
}

func (sc *StaticConfig) MongoDB() bool {
	return false
}

func (sc *StaticConfig) GetRuntimeConfig() RuntimeConfig {
	panic("not runtime config")
}

func (sc *StaticConfig) GetStaticConfig() StaticConfig {
	return *sc
}

type BuildConfig interface {
	isBuildConfig()
	BuildType() BuildType
	Validate() error

	MariaDB() bool
	MongoDB() bool

	GetRuntimeConfig() RuntimeConfig
	GetStaticConfig() StaticConfig
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
	BaseImage string
	BuildCmd  string
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeCmd) BuildType() BuildType {
	return BuildTypeRuntimeCmd
}

func (bc *BuildConfigRuntimeCmd) Validate() error {
	if err := bc.RuntimeConfig.Validate(); err != nil {
		return err
	}
	// NOTE: Base image could have no entrypoint/command but is impossible to catch only from config
	if bc.BaseImage == "" && bc.Entrypoint == "" && bc.Command == "" {
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

type BuildConfigStaticBuildpack struct {
	StaticConfig
	Context string
	buildConfigEmbed
}

func (bc *BuildConfigStaticBuildpack) BuildType() BuildType {
	return BuildTypeStaticBuildpack
}

func (bc *BuildConfigStaticBuildpack) Validate() error {
	if err := bc.StaticConfig.Validate(); err != nil {
		return err
	}
	// NOTE: context is not necessary
	return nil
}

type BuildConfigStaticCmd struct {
	StaticConfig
	BaseImage string
	BuildCmd  string
	buildConfigEmbed
}

func (bc *BuildConfigStaticCmd) BuildType() BuildType {
	return BuildTypeStaticCmd
}

func (bc *BuildConfigStaticCmd) Validate() error {
	if err := bc.StaticConfig.Validate(); err != nil {
		return err
	}
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
	buildConfigEmbed
}

func (bc *BuildConfigStaticDockerfile) BuildType() BuildType {
	return BuildTypeStaticDockerfile
}

func (bc *BuildConfigStaticDockerfile) Validate() error {
	if err := bc.StaticConfig.Validate(); err != nil {
		return err
	}
	if bc.DockerfileName == "" {
		return errors.New("dockerfile_name is required")
	}
	return nil
}

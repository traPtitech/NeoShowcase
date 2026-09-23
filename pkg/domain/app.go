package domain

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/oops"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type ApplicationConfig struct {
	BuildConfig  BuildConfig
	AutoShutdown AutoShutdownConfig
}

func (c *ApplicationConfig) Validate(deployType DeployType) error {
	if c.BuildConfig.BuildType().DeployType() != deployType {
		return oops.New("deploy type doesn't match build type")
	}
	if err := c.BuildConfig.Validate(); err != nil {
		return oops.Wrapf(err, "invalid build_config")
	}
	if err := c.AutoShutdown.Validate(deployType); err != nil {
		return oops.Wrapf(err, "invalid auto_shutdown")
	}
	return nil
}

// AutoShutdownConfig is the setting to stop a runtime application while it receives no requests.
// Enabling it is only allowed for runtime applications.
type AutoShutdownConfig struct {
	Enabled bool
	// Startup must be set if Enabled is true.
	Startup StartupBehavior
}

func (c *AutoShutdownConfig) Validate(deployType DeployType) error {
	if !c.Enabled {
		return nil
	}
	if deployType != DeployTypeRuntime {
		return oops.New("auto shutdown is only available for runtime applications")
	}
	if c.Startup == StartupBehaviorUndefined {
		return oops.New("startup is required if auto shutdown is enabled")
	}
	return nil
}

// StartupBehavior represents the behavior of the application when it starts up.
// This is used to determine how to handle the request until the application is ready.
type StartupBehavior int

const (
	StartupBehaviorUndefined StartupBehavior = iota
	// StartupBehaviorLoadingPage is a strategy that shows a loading page until the application is ready.
	// This is suitable for web applications that are accessed from frontend.
	StartupBehaviorLoadingPage
	// StartupBehaviorBlocking is a strategy that blocks the request until the application is ready.
	// This is suitable for API servers.
	StartupBehaviorBlocking
)

// Hash returns a digest of the build inputs, which are the build config and env.
// Settings that do not affect builds, such as AutoShutdown, are not included
// so that changing them does not trigger a rebuild.
func (c *ApplicationConfig) Hash(env []*Environment) string {
	b := lo.Must(json.Marshal(c.BuildConfig))
	sort.SliceStable(env, func(i, j int) bool { return env[i].Key < env[j].Key })
	e := lo.Must(json.Marshal(env))
	b = append(b, e...)
	return hash.XXH3Hex(b)
}

type DeployType int

const (
	DeployTypeRuntime DeployType = iota
	DeployTypeStatic
)

var EmptyCommit = strings.Repeat("0", 40)

type Application struct {
	ID               string
	Name             string
	RepositoryID     string
	RefName          string
	Commit           string
	DeployType       DeployType
	Running          bool
	Container        ContainerState
	ContainerMessage string
	CurrentBuild     string
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Config           ApplicationConfig
	Websites         []*Website
	PortPublications []*PortPublication
	OwnerIDs         []string
}

func (a *Application) SelfValidate() error {
	if a.Name == "" {
		return oops.New("name is required")
	}
	if a.RepositoryID == "" {
		return oops.New("repository_id is required")
	}
	if a.RefName == "" {
		return oops.New("ref_name is required")
	}
	if err := a.Config.Validate(a.DeployType); err != nil {
		return oops.Wrapf(err, "invalid config")
	}
	for _, website := range a.Websites {
		if err := website.Validate(); err != nil {
			return oops.Wrapf(err, "invalid website")
		}
	}
	for _, p := range a.PortPublications {
		if err := p.Validate(); err != nil {
			return oops.Wrapf(err, "invalid port publication")
		}
	}
	if len(a.OwnerIDs) == 0 {
		return oops.New("owner_ids cannot be empty")
	}
	return nil
}

func (a *Application) Validate(
	actor *User,
	existingApps []*Application,
	domains AvailableDomainSlice,
	ports AvailablePortSlice,
) error {
	if err := a.SelfValidate(); err != nil {
		return err
	}

	// resource availability check
	for _, website := range a.Websites {
		if website.Authentication != AuthenticationTypeOff && !domains.IsAuthAvailable(website.FQDN) {
			return oops.Errorf("auth not available for domain %s", website.FQDN)
		}
		if !domains.IsAvailable(website.FQDN) {
			return oops.Errorf("domain %s not available", website.FQDN)
		}
	}
	for _, p := range a.PortPublications {
		if !ports.IsAvailable(p.InternetPort, p.Protocol) {
			return oops.Errorf("port %d/%s not available", p.InternetPort, p.Protocol)
		}
	}

	// resource conflict check
	// exclude self if contained
	existingApps = lo.Filter(existingApps, func(app *Application, _ int) bool { return app.ID != a.ID })
	if a.WebsiteConflicts(existingApps, actor) {
		return oops.New("website conflict")
	}
	for _, p := range a.PortPublications {
		if p.ConflictsWith(existingApps) {
			return oops.Errorf("port %d/%s conflicts with existing port publication", p.InternetPort, p.Protocol)
		}
	}

	return nil
}

func (a *Application) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(a.OwnerIDs, user.ID)
}

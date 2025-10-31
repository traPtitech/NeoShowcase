package domain

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type ApplicationConfig struct {
	BuildConfig BuildConfig
}

func (c *ApplicationConfig) Validate(deployType DeployType) error {
	if c.BuildConfig.BuildType().DeployType() != deployType {
		return errors.New("deploy type doesn't match build type")
	}
	if err := c.BuildConfig.Validate(); err != nil {
		return errors.Wrap(err, "invalid build_config")
	}
	return nil
}

func (c *ApplicationConfig) Hash(env []*Environment) string {
	b := lo.Must(json.Marshal(c))
	sort.SliceStable(env, func(i, j int) bool { return env[i].Key < env[j].Key })
	e := lo.Must(json.Marshal(env))
	b = append(b, e...)
	return hash.XXH3Hex(b)
}

type DeployType int

const (
	DeployTypeRuntime DeployType = iota
	DeployTypeStatic
	DeployTypeFunction
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
		return errors.New("name is required")
	}
	if a.RepositoryID == "" {
		return errors.New("repository_id is required")
	}
	if a.RefName == "" {
		return errors.New("ref_name is required")
	}
	if err := a.Config.Validate(a.DeployType); err != nil {
		return errors.Wrap(err, "invalid config")
	}
	for _, website := range a.Websites {
		if err := website.Validate(); err != nil {
			return errors.Wrap(err, "invalid website")
		}
	}
	for _, p := range a.PortPublications {
		if err := p.Validate(); err != nil {
			return errors.Wrap(err, "invalid port publication")
		}
	}
	if len(a.OwnerIDs) == 0 {
		return errors.New("owner_ids cannot be empty")
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
			return errors.Errorf("auth not available for domain %s", website.FQDN)
		}
		if !domains.IsAvailable(website.FQDN) {
			return errors.Errorf("domain %s not available", website.FQDN)
		}
	}
	for _, p := range a.PortPublications {
		if !ports.IsAvailable(p.InternetPort, p.Protocol) {
			return errors.Errorf("port %d/%s not available", p.InternetPort, p.Protocol)
		}
	}

	// resource conflict check
	// exclude self if contained
	existingApps = lo.Filter(existingApps, func(app *Application, _ int) bool { return app.ID != a.ID })
	if a.WebsiteConflicts(existingApps, actor) {
		return errors.New("website conflict")
	}
	for _, p := range a.PortPublications {
		if p.ConflictsWith(existingApps) {
			return errors.Errorf("port %d/%s conflicts with existing port publication", p.InternetPort, p.Protocol)
		}
	}

	return nil
}

func (a *Application) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(a.OwnerIDs, user.ID)
}

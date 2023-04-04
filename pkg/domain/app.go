package domain

import (
	"context"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/net/idna"

	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateDeploying
	ApplicationStateRunning
	ApplicationStateErrored
)

type AuthenticationType int

const (
	AuthenticationTypeOff AuthenticationType = iota
	AuthenticationTypeSoft
	AuthenticationTypeHard
)

type ApplicationConfig struct {
	UseMariaDB     bool
	UseMongoDB     bool
	BaseImage      string
	DockerfileName string
	ArtifactPath   string
	BuildCmd       string
	EntrypointCmd  string
	Authentication AuthenticationType
}

type BuildType int

const (
	BuildTypeRuntime BuildType = iota
	BuildTypeStatic
)

var EmptyCommit = strings.Repeat("0", 40)

type Application struct {
	ID            string
	Name          string
	RepositoryID  string
	RefName       string
	BuildType     BuildType
	State         ApplicationState
	CurrentCommit string
	WantCommit    string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Config   ApplicationConfig
	Websites []*Website
	OwnerIDs []string
}

type Artifact struct {
	ID        string
	BuildID   string
	Size      int64
	CreatedAt time.Time
	DeletedAt optional.Of[time.Time]
}

func NewArtifact(buildID string, size int64) *Artifact {
	return &Artifact{
		ID:        NewID(),
		BuildID:   buildID,
		Size:      size,
		CreatedAt: time.Now(),
	}
}

func IsValidDomain(domain string) bool {
	// 面倒なのでtrailing dotは無しで統一
	if strings.HasSuffix(domain, ".") {
		return false
	}
	_, err := idna.Lookup.ToUnicode(domain)
	return err == nil
}

type AvailableDomain struct {
	Domain    string
	Available bool
}

type AvailableDomainSlice []*AvailableDomain

func (a *AvailableDomain) IsValid() bool {
	domain := a.Domain
	domain = strings.TrimPrefix(domain, "*.")
	return IsValidDomain(domain)
}

func (a *AvailableDomain) match(fqdn string) bool {
	if fqdn == a.Domain {
		return true
	}
	if strings.HasPrefix(a.Domain, "*.") {
		baseDomain := strings.TrimPrefix(a.Domain, "*")
		if strings.HasSuffix(fqdn, baseDomain) {
			return true
		}
	}
	return false
}

func (s AvailableDomainSlice) IsAvailable(fqdn string) bool {
	for _, a := range s {
		if !a.Available && a.match(fqdn) {
			return false
		}
	}
	for _, a := range s {
		if a.Available && a.match(fqdn) {
			return true
		}
	}
	return false
}

type BuildStatus int

const (
	BuildStatusQueued BuildStatus = iota
	BuildStatusBuilding
	BuildStatusSucceeded
	BuildStatusFailed
	BuildStatusCanceled
	BuildStatusSkipped
)

func (t BuildStatus) IsFinished() bool {
	switch t {
	case BuildStatusSucceeded, BuildStatusFailed, BuildStatusCanceled, BuildStatusSkipped:
		return true
	default:
		return false
	}
}

type Build struct {
	ID            string
	Commit        string
	Status        BuildStatus
	ApplicationID string
	StartedAt     optional.Of[time.Time]
	UpdatedAt     optional.Of[time.Time]
	FinishedAt    optional.Of[time.Time]
	Retriable     bool
	Artifact      optional.Of[Artifact]
}

func NewBuild(applicationID string, commit string) *Build {
	return &Build{
		ID:            NewID(),
		Commit:        commit,
		Status:        BuildStatusQueued,
		ApplicationID: applicationID,
	}
}

type Environment struct {
	Key    string
	Value  string
	System bool
}

type Repository struct {
	ID       string
	Name     string
	URL      string
	Auth     optional.Of[RepositoryAuth]
	OwnerIDs []string
}

type RepositoryAuthMethod int

const (
	RepositoryAuthMethodBasic RepositoryAuthMethod = iota
	RepositoryAuthMethodSSH
)

type RepositoryAuth struct {
	Method   RepositoryAuthMethod
	Username string
	Password string
	SSHKey   string
}

type Website struct {
	ID          string
	FQDN        string
	PathPrefix  string
	StripPrefix bool
	HTTPS       bool
	HTTPPort    int
}

func (w *Website) IsValid() bool {
	if !IsValidDomain(w.FQDN) {
		return false
	}
	if !strings.HasPrefix(w.PathPrefix, "/") {
		return false
	}
	if w.PathPrefix != "/" && strings.HasSuffix(w.PathPrefix, "/") {
		return false
	}
	if w.StripPrefix && w.PathPrefix == "/" {
		return false
	}
	if !(0 <= w.HTTPPort && w.HTTPPort < 65536) {
		return false
	}
	return true
}

func (w *Website) pathComponents() []string {
	if w.PathPrefix == "/" {
		return []string{}
	}
	return strings.Split(w.PathPrefix, "/")[1:]
}

func (w *Website) pathContainedBy(target *Website) bool {
	this := w.pathComponents()
	other := target.pathComponents()
	if len(this) < len(other) {
		return false
	}
	for i := range other {
		if this[i] != other[i] {
			return false
		}
	}
	return true
}

// ConflictsWith checks whether this website's path prefix is contained in the existing websites' path prefixes.
func (w *Website) ConflictsWith(existing []*Website) bool {
	for _, ex := range existing {
		if w.FQDN != ex.FQDN {
			continue
		}
		if w.pathContainedBy(ex) {
			return true
		}
	}
	return false
}

func GetArtifactsInUse(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*Artifact, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		BuildType: optional.From(BuildTypeStatic),
	})
	if err != nil {
		return nil, err
	}

	commits := make(map[string]struct{}, 2*len(applications))
	for _, app := range applications {
		commits[app.WantCommit] = struct{}{}
		commits[app.CurrentCommit] = struct{}{}
	}
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{CommitIn: optional.From(lo.Keys(commits)), Status: optional.From(BuildStatusSucceeded)})
	if err != nil {
		return nil, err
	}

	// Last succeeded builds for each commit
	slices.SortFunc(builds, func(a, b *Build) bool { return a.StartedAt.ValueOrZero().Before(b.StartedAt.ValueOrZero()) })
	commitToBuild := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.Commit, b })

	artifacts := make([]*Artifact, len(commitToBuild))
	for _, build := range commitToBuild {
		if !build.Artifact.Valid {
			continue
		}
		artifacts = append(artifacts, &build.Artifact.V)
	}
	return artifacts, nil
}

func GetActiveStaticSites(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*StaticSite, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		BuildType: optional.From(BuildTypeStatic),
		State:     optional.From(ApplicationStateRunning),
	})
	if err != nil {
		return nil, err
	}

	commits := lo.Map(applications, func(app *Application, i int) string { return app.CurrentCommit })
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{CommitIn: optional.From(commits), Status: optional.From(BuildStatusSucceeded)})
	if err != nil {
		return nil, err
	}

	// Last succeeded builds for each commit
	slices.SortFunc(builds, func(a, b *Build) bool { return a.StartedAt.ValueOrZero().Before(b.StartedAt.ValueOrZero()) })
	commitToBuild := lo.SliceToMap(builds, func(b *Build) (string, *Build) { return b.Commit, b })

	var sites []*StaticSite
	for _, app := range applications {
		build, ok := commitToBuild[app.CurrentCommit]
		if !ok {
			continue
		}
		if !build.Artifact.Valid {
			continue
		}
		for _, website := range app.Websites {
			sites = append(sites, &StaticSite{
				Application: app,
				Website:     website,
				ArtifactID:  build.Artifact.V.ID,
			})
		}
	}
	return sites, nil
}

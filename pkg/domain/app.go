package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/net/idna"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ApplicationState int

const (
	ApplicationStateIdle ApplicationState = iota
	ApplicationStateDeploying
	ApplicationStateRunning
	ApplicationStateErrored
)

func (s ApplicationState) String() string {
	switch s {
	case ApplicationStateIdle:
		return "IDLE"
	case ApplicationStateDeploying:
		return "DEPLOYING"
	case ApplicationStateRunning:
		return "RUNNING"
	case ApplicationStateErrored:
		return "ERRORED"
	default:
		return ""
	}
}

func ApplicationStateFromString(str string) ApplicationState {
	switch str {
	case "IDLE":
		return ApplicationStateIdle
	case "DEPLOYING":
		return ApplicationStateDeploying
	case "RUNNING":
		return ApplicationStateRunning
	case "ERRORED":
		return ApplicationStateErrored
	default:
		panic(fmt.Sprintf("unknown application state: %v", str))
	}
}

type AuthenticationType int

const (
	AuthenticationTypeOff AuthenticationType = iota
	AuthenticationTypeSoft
	AuthenticationTypeHard
)

func (a AuthenticationType) String() string {
	switch a {
	case AuthenticationTypeOff:
		return "off"
	case AuthenticationTypeSoft:
		return "soft"
	case AuthenticationTypeHard:
		return "hard"
	default:
		return ""
	}
}

func AuthenticationTypeFromString(s string) AuthenticationType {
	switch s {
	case "off":
		return AuthenticationTypeOff
	case "soft":
		return AuthenticationTypeSoft
	case "hard":
		return AuthenticationTypeHard
	default:
		panic(fmt.Sprintf("unknown authentication type: %v", s))
	}
}

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

var EmptyCommit = strings.Repeat("0", 40)

type Application struct {
	ID            string
	Name          string
	RepositoryID  string
	BranchName    string
	BuildType     builder.BuildType
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
	Size      int64
	CreatedAt time.Time
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

type Build struct {
	ID            string
	Commit        string
	Status        builder.BuildStatus
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
		Status:        builder.BuildStatusQueued,
		ApplicationID: applicationID,
	}
}

type Environment struct {
	ID            string
	ApplicationID string
	Key           string
	Value         string
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

func (t RepositoryAuthMethod) String() string {
	switch t {
	case RepositoryAuthMethodBasic:
		return "basic"
	case RepositoryAuthMethodSSH:
		return "ssh"
	default:
		return ""
	}
}

func RepositoryAuthMethodFromString(s string) RepositoryAuthMethod {
	switch s {
	case "basic":
		return RepositoryAuthMethodBasic
	case "ssh":
		return RepositoryAuthMethodSSH
	default:
		panic(fmt.Sprintf("unknown auth type: %v", s))
	}
}

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

func GetActiveStaticSites(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*StaticSite, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		BuildType: optional.From(builder.BuildTypeStatic),
		State:     optional.From(ApplicationStateRunning),
	})
	if err != nil {
		return nil, err
	}

	commits := lo.Map(applications, func(app *Application, i int) string { return app.CurrentCommit })
	builds, err := buildRepo.GetBuilds(ctx, GetBuildCondition{CommitIn: optional.From(commits), Status: optional.From(builder.BuildStatusSucceeded)})
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

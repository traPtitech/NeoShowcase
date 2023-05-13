package domain

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/mattn/go-shellwords"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"golang.org/x/net/idna"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
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

type BuildConfig interface {
	isBuildConfig()
	BuildType() BuildType
	Validate() error
}

type buildConfigEmbed struct{}

func (buildConfigEmbed) isBuildConfig() {}

type BuildConfigRuntimeBuildpack struct {
	Context string
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeBuildpack) BuildType() BuildType {
	return BuildTypeRuntimeBuildpack
}

func (bc *BuildConfigRuntimeBuildpack) Validate() error {
	// NOTE: context is not necessary
	return nil
}

type BuildConfigRuntimeCmd struct {
	BaseImage     string
	BuildCmd      string
	BuildCmdShell bool
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeCmd) BuildType() BuildType {
	return BuildTypeRuntimeCmd
}

func (bc *BuildConfigRuntimeCmd) Validate() error {
	// NOTE: base image is not necessary (default: scratch)
	// NOTE: build cmd is not necessary
	return nil
}

type BuildConfigRuntimeDockerfile struct {
	DockerfileName string
	Context        string
	buildConfigEmbed
}

func (bc *BuildConfigRuntimeDockerfile) BuildType() BuildType {
	return BuildTypeRuntimeDockerfile
}

func (bc *BuildConfigRuntimeDockerfile) Validate() error {
	if bc.DockerfileName == "" {
		return errors.New("dockerfile_name is required")
	}
	return nil
}

type BuildConfigStaticCmd struct {
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

type ApplicationConfig struct {
	UseMariaDB  bool
	UseMongoDB  bool
	BuildConfig BuildConfig
	Entrypoint  string
	Command     string
}

func validateCommand(s string) error {
	_, err := shellwords.Parse(s)
	return err
}

func (c *ApplicationConfig) Validate(deployType DeployType) error {
	if c.BuildConfig.BuildType().DeployType() != deployType {
		return errors.New("build type doesn't match deploy type")
	}
	if err := c.BuildConfig.Validate(); err != nil {
		return errors.Wrap(err, "invalid build_config")
	}
	if c.BuildConfig.BuildType() == BuildTypeRuntimeCmd && c.Entrypoint == "" && c.Command == "" {
		return errors.New("entrypoint or command is required for runtime_cmd build type")
	}
	// NOTE: Runtime Dockerfile build could have no entrypoint/command but is impossible to catch only from config
	// (can only catch at runtime)
	if c.Entrypoint != "" {
		if err := validateCommand(c.Entrypoint); err != nil {
			return errors.Wrap(err, "invalid entrypoint")
		}
	}
	if c.Command != "" {
		if err := validateCommand(c.Command); err != nil {
			return errors.Wrap(err, "invalid command")
		}
	}
	return nil
}

func (c *ApplicationConfig) EntrypointArgs() []string {
	args, _ := shellwords.Parse(c.Entrypoint)
	return args
}

func (c *ApplicationConfig) CommandArgs() []string {
	args, _ := shellwords.Parse(c.Command)
	return args
}

type DeployType int

const (
	DeployTypeRuntime DeployType = iota
	DeployTypeStatic
)

var EmptyCommit = strings.Repeat("0", 40)

type Application struct {
	ID            string
	Name          string
	RepositoryID  string
	RefName       string
	DeployType    DeployType
	Running       bool
	Container     ContainerState
	CurrentCommit string
	WantCommit    string
	CreatedAt     time.Time
	UpdatedAt     time.Time

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
	ctx context.Context,
	actor *User,
	controller ControllerServiceClient,
	existingApps []*Application,
	domains AvailableDomainSlice,
	ports AvailablePortSlice,
) (validateErr error, err error) {
	if err = a.SelfValidate(); err != nil {
		return err, nil
	}

	for _, website := range a.Websites {
		if website.Authentication == AuthenticationTypeOff {
			continue
		}
		available, err := controller.AuthAvailable(ctx, website.FQDN)
		if err != nil {
			return nil, errors.Wrap(err, "checking auth availability")
		}
		if !available {
			return errors.Errorf("auth not available for domain %s", website.FQDN), nil
		}
	}

	for _, website := range a.Websites {
		if !domains.IsAvailable(website.FQDN) {
			return errors.Errorf("domain %s not available", website.FQDN), nil
		}
	}

	// exclude self if contained
	existingApps = lo.Filter(existingApps, func(app *Application, _ int) bool { return app.ID != a.ID })

	if a.WebsiteConflicts(existingApps, actor) {
		return errors.New("website conflict"), nil
	}

	for _, p := range a.PortPublications {
		if !ports.IsAvailable(p.InternetPort, p.Protocol) {
			return errors.Errorf("port %d/%s not available", p.InternetPort, p.Protocol), nil
		}
		if p.ConflictsWith(existingApps) {
			return errors.Errorf("port %d/%s conflicts with existing port publication", p.InternetPort, p.Protocol), nil
		}
	}

	return nil, nil
}

func (a *Application) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(a.OwnerIDs, user.ID)
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

func ValidateDomain(domain string) error {
	// 面倒なのでtrailing dotは無しで統一
	if strings.HasSuffix(domain, ".") {
		return errors.New("trailing dot not allowed in domain")
	}
	if strings.HasPrefix(domain, ".") {
		return errors.New("leading dot not allowed in domain")
	}
	_, err := idna.Lookup.ToUnicode(domain)
	if err != nil {
		return errors.Wrap(err, "invalid domain")
	}
	return nil
}

func ValidateWildcardDomain(domain string) error {
	if !strings.HasPrefix(domain, "*.") {
		return errors.New("wildcard domain needs to begin with *.")
	}
	baseDomain := strings.TrimPrefix(domain, "*.")
	return ValidateDomain(baseDomain)
}

func MatchDomain(source, target string) bool {
	if source == target {
		return true
	}
	if strings.HasPrefix(source, "*.") {
		baseDomain := strings.TrimPrefix(source, "*")
		if strings.HasSuffix(target, baseDomain) {
			return true
		}
	}
	return false
}

type AvailableDomain struct {
	Domain    string
	Available bool
}

type AvailableDomainSlice []*AvailableDomain

func (a *AvailableDomain) Validate() error {
	err := ValidateWildcardDomain(a.Domain)
	if err == nil {
		return nil
	}
	return ValidateDomain(a.Domain)
}

func (a *AvailableDomain) match(fqdn string) bool {
	return MatchDomain(a.Domain, fqdn)
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

type PortPublicationProtocol string

const (
	PortPublicationProtocolTCP PortPublicationProtocol = "tcp"
	PortPublicationProtocolUDP PortPublicationProtocol = "udp"
)

type AvailablePort struct {
	StartPort int
	EndPort   int
	Protocol  PortPublicationProtocol
}

type AvailablePortSlice []*AvailablePort

func isValidPort(port int) error {
	if !(0 <= port && port < 65536) {
		return errors.New("invalid port (needs to be within 0 to 65535)")
	}
	return nil
}

func (ap *AvailablePort) Validate() error {
	if err := isValidPort(ap.StartPort); err != nil {
		return errors.Wrap(err, "invalid start port")
	}
	if err := isValidPort(ap.EndPort); err != nil {
		return errors.Wrap(err, "invalid end port")
	}
	if ap.EndPort < ap.StartPort {
		return errors.New("end port comes before start port")
	}
	return nil
}

func (s AvailablePortSlice) IsAvailable(port int, protocol PortPublicationProtocol) bool {
	return lo.ContainsBy(s, func(ap *AvailablePort) bool {
		return ap.Protocol == protocol && ap.StartPort <= port && port <= ap.EndPort
	})
}

type UnavailablePort struct {
	Port     int
	Protocol PortPublicationProtocol
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
	QueuedAt      time.Time
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
		QueuedAt:      time.Now(),
	}
}

type Environment struct {
	ApplicationID string
	Key           string
	Value         string
	System        bool
}

type Repository struct {
	ID       string
	Name     string
	URL      string
	Auth     optional.Of[RepositoryAuth]
	OwnerIDs []string
}

func (r *Repository) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	ep, err := transport.NewEndpoint(r.URL)
	if err != nil {
		return errors.Wrap(err, "invalid url")
	}
	if !r.Auth.Valid {
		// URL is in http(s) format
		if ep.Protocol != "http" && ep.Protocol != "https" {
			return errors.New("url has to be http(s) protocol when auth is none")
		}
	} else if r.Auth.V.Method == RepositoryAuthMethodBasic {
		// URL is in https format
		if ep.Protocol != "https" {
			return errors.New("url has to be https protocol when auth is basic")
		}
	} else if r.Auth.V.Method == RepositoryAuthMethodSSH {
		// URL is in ssh format
		if ep.Protocol != "ssh" {
			return errors.New("url has to be ssh protocol when auth is ssh")
		}
	}
	if len(r.OwnerIDs) == 0 {
		return errors.New("owner_ids cannot be empty")
	}
	return nil
}

func (r *Repository) IsOwner(user *User) bool {
	return user.Admin || lo.Contains(r.OwnerIDs, user.ID)
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

func (r *RepositoryAuth) Validate() error {
	switch r.Method {
	case RepositoryAuthMethodBasic:
		if r.Username == "" {
			return errors.New("username cannot be empty")
		}
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
	case RepositoryAuthMethodSSH:
		if r.SSHKey != "" {
			_, err := ssh.NewPublicKeys("", []byte(r.SSHKey), "")
			if err != nil {
				return errors.Wrap(err, "invalid ssh private key")
			}
		}
	}
	return nil
}

type AuthenticationType int

const (
	AuthenticationTypeOff AuthenticationType = iota
	AuthenticationTypeSoft
	AuthenticationTypeHard
)

type Website struct {
	ID             string
	FQDN           string
	PathPrefix     string
	StripPrefix    bool
	HTTPS          bool
	H2C            bool
	HTTPPort       int
	Authentication AuthenticationType
}

func (w *Website) Validate() error {
	if err := ValidateDomain(w.FQDN); err != nil {
		return errors.Wrap(err, "invalid domain")
	}
	if !strings.HasPrefix(w.PathPrefix, "/") {
		return errors.New("path_prefix has to start with /")
	}
	if w.PathPrefix != "/" && strings.HasSuffix(w.PathPrefix, "/") {
		return errors.New("path_prefix requires no trailing slash")
	}
	if w.StripPrefix && w.PathPrefix == "/" {
		return errors.New("strip_prefix has to be false when path_prefix is /")
	}
	u, err := url.ParseRequestURI(w.PathPrefix)
	if err != nil {
		return errors.Wrap(err, "invalid path")
	}
	if u.EscapedPath() != w.PathPrefix {
		return errors.New("invalid path: either not escaped or contains non-path elements")
	}
	if err = isValidPort(w.HTTPPort); err != nil {
		return errors.Wrap(err, "invalid http port")
	}
	return nil
}

func (w *Website) pathComponents() []string {
	// NOTE: empty PathPrefix must not exist
	if w.PathPrefix == "/" {
		return []string{}
	}
	return strings.Split(w.PathPrefix[1:], "/")
}

func (w *Website) pathContainedBy(target *Website) bool {
	this := w.pathComponents()
	other := target.pathComponents()
	if len(this) < len(other) {
		return false
	}
	return ds.Equals(this[:len(other)], other)
}

func (w *Website) Equals(target *Website) bool {
	if w.FQDN != target.FQDN {
		return false
	}
	if w.HTTPS != target.HTTPS {
		return false
	}
	return w.PathPrefix == target.PathPrefix
}

func (w *Website) conflictsWith(target *Website) bool {
	if w.FQDN != target.FQDN {
		return false
	}
	if w.HTTPS != target.HTTPS {
		return false
	}
	return w.pathContainedBy(target) || target.pathContainedBy(w)
}

func (a *Application) WebsiteConflicts(existing []*Application, actor *User) bool {
	for _, w := range a.Websites {
		for _, ex := range existing {
			for _, w2 := range ex.Websites {
				if w.Equals(w2) {
					return true
				}
				if w.conflictsWith(w2) && !ex.IsOwner(actor) {
					return true
				}
			}
		}
	}
	return false
}

type PortPublication struct {
	InternetPort    int
	ApplicationPort int
	Protocol        PortPublicationProtocol
}

func (p *PortPublication) Validate() error {
	if err := isValidPort(p.InternetPort); err != nil {
		return errors.Wrap(err, "invalid internet port")
	}
	if err := isValidPort(p.ApplicationPort); err != nil {
		return errors.Wrap(err, "invalid application port")
	}
	return nil
}

func (p *PortPublication) ToUnavailablePort() *UnavailablePort {
	return &UnavailablePort{
		Port:     p.InternetPort,
		Protocol: p.Protocol,
	}
}

func (p *PortPublication) ConflictsWith(existing []*Application) bool {
	return lo.ContainsBy(existing, func(app *Application) bool {
		return lo.ContainsBy(app.PortPublications, func(used *PortPublication) bool {
			return used.InternetPort == p.InternetPort && used.Protocol == p.Protocol
		})
	})
}

func GetArtifactsInUse(ctx context.Context, appRepo ApplicationRepository, buildRepo BuildRepository) ([]*Artifact, error) {
	applications, err := appRepo.GetApplications(ctx, GetApplicationCondition{
		DeployType: optional.From(DeployTypeStatic),
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
		DeployType: optional.From(DeployTypeStatic),
		Running:    optional.From(true),
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

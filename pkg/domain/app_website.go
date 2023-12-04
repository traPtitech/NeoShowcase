package domain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"golang.org/x/net/idna"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

func ValidateDomain(domain string) error {
	// ドメインが大文字を含むときはエラー
	if domain != strings.ToLower(domain) {
		return errors.Errorf("domain %v must be lower case", domain)
	}
	// 面倒なのでtrailing dotは無しで統一
	if strings.HasSuffix(domain, ".") {
		return errors.Errorf("trailing dot not allowed in domain %v", domain)
	}
	if strings.HasPrefix(domain, ".") {
		return errors.Errorf("leading dot not allowed in domain %v", domain)
	}
	// allow underscore; Showcaseとのcompatibilityのため 本来はホスト名にunderscoreが入るのはダメ
	// https://stackoverflow.com/questions/2180465/can-domain-name-subdomains-have-an-underscore-in-it
	_, err := idna.Lookup.ToUnicode(strings.ReplaceAll(domain, "_", "-"))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid domain %v", domain))
	}
	return nil
}

func ValidateWildcardDomain(domain string) error {
	if !strings.HasPrefix(domain, "*.") {
		return errors.Errorf("wildcard domain needs to begin with *. (got %v)", domain)
	}
	baseDomain := strings.TrimPrefix(domain, "*.")
	return ValidateDomain(baseDomain)
}

func ValidatePossibleWildcardDomain(domain string) error {
	err := ValidateWildcardDomain(domain)
	if err == nil {
		return nil
	}
	return ValidateDomain(domain)
}

func ContainsDomain(source, target string) bool {
	if source == target {
		return true
	}
	if strings.HasPrefix(source, "*.") {
		baseSource := strings.TrimPrefix(source, "*")
		if strings.HasSuffix(target, baseSource) {
			return true
		}
	}
	return false
}

type AvailableDomain struct {
	Domain         string
	ExcludeDomains []string
	AuthAvailable  bool
	AlreadyBound   bool // Actual availability (whether domain is bound to a specific app or not)
}

type AvailableDomainSlice []*AvailableDomain

func (a *AvailableDomain) Validate() error {
	if err := ValidatePossibleWildcardDomain(a.Domain); err != nil {
		return err
	}
	for _, excludeDomain := range a.ExcludeDomains {
		if err := ValidatePossibleWildcardDomain(excludeDomain); err != nil {
			return err
		}
		if !ContainsDomain(a.Domain, excludeDomain) {
			return errors.Errorf("exclude domain %v is not contained within %v", excludeDomain, a.Domain)
		}
	}
	return nil
}

func (a *AvailableDomain) SetAlreadyBound(existing []*Application) {
	if strings.HasPrefix(a.Domain, "*.") {
		// Wildcard domain cannot be bound to one app, it has infinite number of subdomains
		a.AlreadyBound = false
	} else {
		a.AlreadyBound = lo.ContainsBy(existing, func(app *Application) bool {
			return lo.ContainsBy(app.Websites, func(w *Website) bool {
				return w.FQDN == a.Domain && w.PathPrefix == "/" // Intentional vague checking of http or https
			})
		})
	}
}

func (a *AvailableDomain) Match(fqdn string) bool {
	for _, excludeDomain := range a.ExcludeDomains {
		if ContainsDomain(excludeDomain, fqdn) {
			return false
		}
	}
	return ContainsDomain(a.Domain, fqdn)
}

func (s AvailableDomainSlice) IsAvailable(fqdn string) bool {
	return lo.ContainsBy(s, func(ad *AvailableDomain) bool {
		return ad.Match(fqdn)
	})
}

func (s AvailableDomainSlice) IsAuthAvailable(fqdn string) bool {
	return lo.ContainsBy(s, func(ad *AvailableDomain) bool {
		return ad.Match(fqdn) && ad.AuthAvailable
	})
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

func (w *Website) Compare(other *Website) bool {
	return strings.Compare(w.ID, other.ID) < 0
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

func (w *Website) Normalize() {
	w.FQDN = strings.ToLower(w.FQDN)
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
		// check with existing websites
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

		// check with self
		for _, w2 := range a.Websites {
			if w.ID == w2.ID {
				continue
			}
			if w.Equals(w2) {
				return true
			}
		}
	}
	return false
}

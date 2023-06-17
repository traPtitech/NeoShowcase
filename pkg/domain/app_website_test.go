package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{"ok 1", "google.com", false},
		{"ok 2", "hyphens-are-allowed.example.com", false},
		{"ok 3", "日本語.jp", false},
		{"ok 4", "underscore_allowed.example.com", false},
		{"invalid characters 1", "admin@example.com", true},
		{"invalid characters 2", "space not allowed.example.com", true},
		{"wildcard ng", "*.trap.show", true},
		{"multi wildcard ng", "*.*.trap.show", true},
		{"wildcard in middle", "trap.*.show", true},
		{"leading dot ng", ".example.com", true},
		{"trailing dot ng", "google.com.", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDomain(tt.domain)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAvailableDomain_Validate(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{"ok", "google.com", false},
		{"wildcard ok", "*.trap.show", false},
		{"multi wildcard ng", "*.*.trap.show", true},
		{"wildcard in middle", "trap.*.show", true},
		{"trailing dot ng", "google.com.", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AvailableDomain{
				Domain: tt.domain,
			}
			err := a.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAvailableDomain_Match(t *testing.T) {
	simpleTests := []struct {
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
	for _, tt := range simpleTests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AvailableDomain{
				Domain: tt.domain,
			}
			got := a.Match(tt.target)
			assert.Equal(t, tt.want, got)
		})
	}

	tests := []struct {
		name string
		ad   AvailableDomain
		fqdn string
		want bool
	}{
		{
			name: "ok",
			ad:   AvailableDomain{Domain: "google.com"},
			fqdn: "google.com",
			want: true,
		},
		{
			name: "subdomain ok",
			ad:   AvailableDomain{Domain: "*.google.com"},
			fqdn: "sub.google.com",
			want: true,
		},
		{
			name: "ng",
			ad:   AvailableDomain{Domain: "google.com"},
			fqdn: "yahoo.com",
			want: false,
		},
		{
			name: "specific subdomain ng 1",
			ad:   AvailableDomain{Domain: "*.google.com", ExcludeDomains: []string{"sub.google.com"}},
			fqdn: "sub.google.com",
			want: false,
		},
		{
			name: "specific wildcard subdomain ng 1",
			ad:   AvailableDomain{Domain: "*.google.com", ExcludeDomains: []string{"*.sub.google.com"}},
			fqdn: "test.sub.google.com",
			want: false,
		},
		{
			name: "specific wildcard subdomain ok 1",
			ad:   AvailableDomain{Domain: "*.google.com", ExcludeDomains: []string{"*.sub.google.com"}},
			fqdn: "sub.google.com",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, tt.ad.Validate())
			got := tt.ad.Match(tt.fqdn)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWebsite_Validate(t *testing.T) {
	tests := []struct {
		name    string
		website Website
		wantErr bool
	}{
		{"ok1", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: 80}, false},
		{"ok2", Website{FQDN: "google.com", PathPrefix: "/path/to/prefix", HTTPPort: 8080}, false},
		{"ok3", Website{FQDN: "google.com", PathPrefix: "/space%20is%20encoded", HTTPPort: 8080}, false},
		{"ok4", Website{FQDN: "trap.show", PathPrefix: "/~toki", HTTPPort: 8080}, false},
		{"ok5", Website{FQDN: "trap.show", PathPrefix: "/~toki/bot_converter", HTTPPort: 8080}, false},
		{"invalid fqdn1", Website{FQDN: "google.com.", PathPrefix: "/", HTTPPort: 80}, true},
		{"invalid fqdn2", Website{FQDN: "*.google.com", PathPrefix: "/", HTTPPort: 80}, true},
		{"invalid fqdn3", Website{FQDN: "google.*.com", PathPrefix: "/", HTTPPort: 80}, true},
		{"invalid fqdn4", Website{FQDN: "goo gle.com", PathPrefix: "/", HTTPPort: 80}, true},
		{"invalid fqdn5", Website{FQDN: "no space", PathPrefix: "/", HTTPPort: 80}, true},
		{"invalid path1", Website{FQDN: "google.com", PathPrefix: "", HTTPPort: 80}, true},
		{"invalid path2", Website{FQDN: "google.com", PathPrefix: "../test", HTTPPort: 80}, true},
		{"invalid path3", Website{FQDN: "google.com", PathPrefix: "/test/", HTTPPort: 80}, true},
		{"invalid path4", Website{FQDN: "google.com", PathPrefix: "/space not encoded", HTTPPort: 80}, true},
		{"invalid path5", Website{FQDN: "google.com", PathPrefix: "/query?", HTTPPort: 80}, true},
		{"invalid path6", Website{FQDN: "google.com", PathPrefix: "/query?foo", HTTPPort: 80}, true},
		{"invalid path7", Website{FQDN: "google.com", PathPrefix: "/query?foo=bar", HTTPPort: 80}, true},
		{"invalid path8", Website{FQDN: "google.com", PathPrefix: "https://google.com/test", HTTPPort: 80}, true},
		{"strip prefix ok1", Website{FQDN: "google.com", PathPrefix: "/", StripPrefix: false, HTTPPort: 80}, false},
		{"strip prefix ok2", Website{FQDN: "google.com", PathPrefix: "/test", StripPrefix: false, HTTPPort: 80}, false},
		{"strip prefix ng", Website{FQDN: "google.com", PathPrefix: "/", StripPrefix: true, HTTPPort: 80}, true},
		{"strip prefix ok3", Website{FQDN: "google.com", PathPrefix: "/test", StripPrefix: true, HTTPPort: 80}, false},
		{"invalid port1", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: -1}, true},
		{"invalid port2", Website{FQDN: "google.com", PathPrefix: "/", HTTPPort: 65536}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.website.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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

func TestWebsite_conflictsWith(t *testing.T) {
	pathTests := []struct {
		name     string
		target   string
		existing string
		want     bool
	}{
		{"ok1", "/foo", "/api", false},
		{"ok2", "/foo", "/spa", false},
		{"ok3", "/api/v2", "/api/v1", false},
		{"ok4", "/api2", "/api", false},
		{"ok5", "/api", "/api2", false},
		{"ng1", "/", "/", true},
		{"ng2", "/api", "/", true},
		{"ng3", "/", "/api", true},
		{"ng4", "/api/v2", "/api", true},
		{"ng5", "/api", "/api/v2", true},
		{"ng6", "/api", "/api", true},
	}
	for _, tt := range pathTests {
		t.Run("path "+tt.name, func(t *testing.T) {
			w := &Website{PathPrefix: tt.target}
			target := &Website{PathPrefix: tt.existing}
			if got := w.conflictsWith(target); got != tt.want {
				t.Errorf("conflictsWith() = %v, want %v", got, tt.want)
			}
		})
	}

	fullTests := []struct {
		name     string
		target   *Website
		existing *Website
		want     bool
	}{
		{
			name:     "ng if same scheme",
			target:   &Website{PathPrefix: "/", HTTPS: false},
			existing: &Website{PathPrefix: "/", HTTPS: false},
			want:     true,
		},
		{
			name:     "ok if different scheme",
			target:   &Website{PathPrefix: "/", HTTPS: true},
			existing: &Website{PathPrefix: "/", HTTPS: false},
			want:     false,
		},
		{
			name:     "ok if different fqdn",
			target:   &Website{FQDN: "google.com", PathPrefix: "/", HTTPS: false},
			existing: &Website{FQDN: "yahoo.com", PathPrefix: "/", HTTPS: false},
			want:     false,
		},
	}
	for _, tt := range fullTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.conflictsWith(tt.existing); got != tt.want {
				t.Errorf("conflictsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_WebsiteConflicts(t *testing.T) {
	u1 := &User{ID: "user1"}
	u2 := &User{ID: "user2"}
	u3 := &User{ID: "user3"}
	admin := &User{ID: "user4", Admin: true}
	existing := &Application{
		Websites: []*Website{{
			FQDN:       "bar.trap.games",
			PathPrefix: "/",
		}},
		OwnerIDs: []string{u2.ID, u3.ID},
	}
	tests := []struct {
		name     string
		target   *Application
		existing *Application
		actor    *User
		want     bool
	}{
		{
			name: "ok (different fqdn, no conflict)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "foo.trap.games",
					PathPrefix: "/api",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    u1,
			want:     false,
		},
		{
			name: "ng (conflict, no ownership of the other)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "bar.trap.games",
					PathPrefix: "/api",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    u1,
			want:     true,
		},
		{
			name: "ng (conflict, owner of the other, but same website)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "bar.trap.games",
					PathPrefix: "/",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    u3,
			want:     true,
		},
		{
			name: "ng (conflict, actor is admin, but same website)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "bar.trap.games",
					PathPrefix: "/",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    admin,
			want:     true,
		},
		{
			name: "ok (conflict, but owner of the other)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "bar.trap.games",
					PathPrefix: "/api",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    u3,
			want:     false,
		},
		{
			name: "ok (conflict, but actor is admin)",
			target: &Application{
				Websites: []*Website{{
					FQDN:       "bar.trap.games",
					PathPrefix: "/api",
				}},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    admin,
			want:     false,
		},
		{
			name: "ng (self conflict)",
			target: &Application{
				Websites: []*Website{
					{
						ID:         "w1",
						FQDN:       "foo.trap.games",
						PathPrefix: "/api",
					},
					{
						ID:         "w2",
						FQDN:       "foo.trap.games",
						PathPrefix: "/api",
					},
				},
				OwnerIDs: []string{u1.ID, u3.ID},
			},
			existing: existing,
			actor:    u1,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.WebsiteConflicts([]*Application{tt.existing}, tt.actor)
			assert.Equal(t, tt.want, got)
		})
	}
}

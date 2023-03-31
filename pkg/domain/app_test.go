package domain

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

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
			if got := a.match(tt.target); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsAvailable(tt.fqdn); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
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
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
}

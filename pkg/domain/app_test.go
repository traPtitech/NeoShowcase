package domain

import (
	"testing"

	"github.com/samber/lo"
)

func TestIsValidDomain(t *testing.T) {
	type args struct {
		domain string
	}
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

func TestAvailableDomain_Match(t *testing.T) {
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
			if got := a.Match(tt.target); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
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
		{"ng1", "/", []string{"/"}, true},
		{"ng2", "/api", []string{"/"}, true},
		{"ng3", "/api/v2", []string{"/api"}, true},
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

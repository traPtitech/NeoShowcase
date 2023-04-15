package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLSTargetDomain(t *testing.T) {
	tests := []struct {
		name    string
		fqdn    string
		domains WildcardDomains
		want    string
	}{
		{"wildcard match", "sub.google.com", WildcardDomains{"*.google.com"}, "*.google.com"},
		{"not wildcard match", "google.com", WildcardDomains{"*.google.com"}, "google.com"},
		{"recursive wildcard match", "grand.children.google.com", WildcardDomains{"*.google.com"}, "*.children.google.com"},
		{"domain txt record not in control", "test.yahoo.com", WildcardDomains{"*.google.com"}, "test.yahoo.com"},
		{"wildcard match (no wildcard)", "sub.google.com", WildcardDomains{}, "sub.google.com"},
		{"not wildcard match (no wildcard)", "google.com", WildcardDomains{}, "google.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.domains.TLSTargetDomain(&Website{FQDN: tt.fqdn}), "(%v).TLSTargetDomain(%v)", tt.domains, tt.fqdn)
		})
	}
}

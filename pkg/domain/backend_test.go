package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLSTargetDomain(t *testing.T) {
	tests := []struct {
		name     string
		wildcard bool
		fqdn     string
		ads      AvailableDomainSlice
		want     string
	}{
		{"wildcard match", true, "sub.google.com", AvailableDomainSlice{{Domain: "*.google.com", Available: true}}, "*.google.com"},
		{"not wildcard match", true, "sub.google.com", AvailableDomainSlice{{Domain: "sub.google.com", Available: true}}, "sub.google.com"},
		{"wildcard match (no wildcard)", false, "sub.google.com", AvailableDomainSlice{{Domain: "*.google.com", Available: true}}, "sub.google.com"},
		{"not wildcard match (no wildcard)", false, "sub.google.com", AvailableDomainSlice{{Domain: "sub.google.com", Available: true}}, "sub.google.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, TLSTargetDomain(tt.wildcard, &Website{FQDN: tt.fqdn}, tt.ads), "TLSTargetDomain(%v, %v, %v)", tt.wildcard, tt.fqdn, tt.ads)
		})
	}
}

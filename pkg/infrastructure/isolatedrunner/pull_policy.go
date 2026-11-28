// Package isolatedrunner holds what the implementations of builder.IsolatedRunner have in common.
package isolatedrunner

import "github.com/samber/oops"

// PullPolicy decides when the image of a job is pulled, with the same meaning as the imagePullPolicy of k8s.
type PullPolicy string

const (
	PullAlways       PullPolicy = "Always"
	PullIfNotPresent PullPolicy = "IfNotPresent"
	PullNever        PullPolicy = "Never"
)

// Validate reports whether p is one of the defined policies.
func (p PullPolicy) Validate() error {
	switch p {
	case PullAlways, PullIfNotPresent, PullNever:
		return nil
	default:
		return oops.Errorf("invalid pull policy %q (supported values: %s, %s, %s)", p, PullAlways, PullIfNotPresent, PullNever)
	}
}

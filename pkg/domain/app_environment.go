package domain

import (
	"regexp"

	"github.com/samber/oops"
)

type Environment struct {
	ApplicationID string
	Key           string
	Value         string
	System        bool
}

func (e *Environment) GetKV() (string, string) {
	return e.Key, e.Value
}

var environmentVariableKeyFormat = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

func (e *Environment) Validate() error {
	if !environmentVariableKeyFormat.MatchString(e.Key) {
		return oops.Errorf("bad key format: %s", e.Key)
	}
	return nil
}

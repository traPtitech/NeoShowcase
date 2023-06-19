package domain

type Environment struct {
	ApplicationID string
	Key           string
	Value         string
	System        bool
}

func (e *Environment) GetKV() (string, string) {
	return e.Key, e.Value
}

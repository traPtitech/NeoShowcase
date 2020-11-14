package generator

type Engine interface {
	Init() error
	Reconcile(sites []*Site) error
}

type Site struct {
	ID            string
	FQDN          string
	PathPrefix    string
	ArtifactID    string
	ApplicationID string
}

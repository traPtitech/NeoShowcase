package cli

var (
	version  string
	revision string
)

func SetVersion(ver, rev string) {
	version = ver
	revision = rev
}

func GetVersion() (ver, rev string) {
	return version, revision
}

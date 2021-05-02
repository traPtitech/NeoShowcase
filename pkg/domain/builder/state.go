package builder

type State int

const (
	StateUnknown State = iota
	StateUnavailable
	StateWaiting
	StateBuilding
)

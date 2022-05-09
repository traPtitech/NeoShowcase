package builder

import "fmt"

type BuildStatus int

const (
	BuildStatusBuilding BuildStatus = iota
	BuildStatusSucceeded
	BuildStatusFailed
	BuildStatusCanceled
	BuildStatusQueued
	BuildStatusSkipped
)

func (t BuildStatus) String() string {
	switch t {
	case BuildStatusBuilding:
		return "BUILDING"
	case BuildStatusSucceeded:
		return "SUCCEEDED"
	case BuildStatusFailed:
		return "FAILED"
	case BuildStatusCanceled:
		return "CANCELED"
	case BuildStatusQueued:
		return "QUEUED"
	case BuildStatusSkipped:
		return "SKIPPED"
	}
	return ""
}

func BuildStatusFromString(str string) BuildStatus {
	switch str {
	case "BUILDING":
		return BuildStatusBuilding
	case "SUCCEEDED":
		return BuildStatusSucceeded
	case "FAILED":
		return BuildStatusFailed
	case "CANCELED":
		return BuildStatusCanceled
	case "QUEUED":
		return BuildStatusQueued
	case "SKIPPED":
		return BuildStatusSkipped
	default:
		panic(fmt.Errorf("UNKNOWN BUILD STATUS: %s", str))
	}
}

func (t BuildStatus) IsFinished() bool {
	switch t {
	case BuildStatusSucceeded, BuildStatusFailed, BuildStatusCanceled, BuildStatusSkipped:
		return true
	default:
		return false
	}
}

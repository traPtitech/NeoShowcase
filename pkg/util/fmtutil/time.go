package fmtutil

import (
	"fmt"
	"time"
)

func DurationHuman(d time.Duration) string {
	switch {
	case d < time.Second:
		return "less than a second"
	case d < 2*time.Second:
		return "a second"
	case d < time.Minute:
		return fmt.Sprintf("%d seconds", d/time.Second)
	case d < 2*time.Minute:
		return "a minute"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes", d/time.Minute)
	case d < 2*time.Hour:
		return "an hour"
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours", d/time.Hour)
	case d < 48*time.Hour:
		return "a day"
	default:
		return fmt.Sprintf("%d days", d/time.Hour/24)
	}
}

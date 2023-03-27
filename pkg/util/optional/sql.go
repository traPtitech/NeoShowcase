package optional

import (
	"time"

	"github.com/volatiletech/null/v8"
)

func IntoTime(o Of[time.Time]) null.Time {
	return null.Time{Time: o.V, Valid: o.Valid}
}

func FromTime(t null.Time) Of[time.Time] {
	return Of[time.Time]{V: t.Time, Valid: t.Valid}
}

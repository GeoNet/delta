package meta

import (
	"time"
)

type Span struct {
	Start time.Time
	End   time.Time
}

func (s Span) before(span Span) bool {
	return s.Start.Before(span.Start)
}

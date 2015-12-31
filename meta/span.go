package meta

import (
	"time"
)

type Span struct {
	Start time.Time
	End   time.Time
}

func (s Span) Before(span Span) bool {
	return s.Start.Before(span.Start)
}

package meta

import (
	"testing"
	"time"
)

func TestTelemetry(t *testing.T) {

	t.Run("check telemetry", testListFunc("testdata/telemetries.csv", &TelemetryList{
		Telemetry{
			Station:  "KAVZ",
			Location: "10",
			Span: Span{
				Start: time.Date(2004, 12, 3, 22, 0, 0, 0, time.UTC),
				End:   time.Date(2004, 12, 14, 19, 0, 0, 0, time.UTC),
			},
			ScaleFactor: 1.0,
			factor:      "",
		},
	}))
}

package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var telemetryChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for overlapping telemeties": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			telemetries := set.Telemetries()
			for i := 0; i < len(telemetries); i++ {
				for j := i + 1; j < len(telemetries); j++ {
					if telemetries[i].Station != telemetries[j].Station {
						continue
					}
					if telemetries[i].Location != telemetries[j].Location {
						continue
					}
					if !telemetries[i].Span.Overlaps(telemetries[j].Span) {
						continue
					}

					t.Errorf("error: telemetry overlap for %q", telemetries[i].String())
				}
			}
		}
	},

	"check for zero gain telemetries": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, v := range set.Telemetries() {
				if v.Gain != 0.0 {
					continue
				}
				t.Errorf("error: telemetry with zero gain for %q", v.String())
			}
		}
	},
}

func TestTelemities(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range telemetryChecks {
		t.Run(k, v(set))
	}
}

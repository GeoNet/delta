package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var timingsChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate timings": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			timings := set.Timings()
			for i := 0; i < len(timings); i++ {
				for j := i + 1; j < len(timings); j++ {
					if timings[i].Station != timings[j].Station {
						continue
					}
					if timings[i].Location != timings[j].Location {
						continue
					}
					if !timings[i].Span.Overlaps(timings[j].Span) {
						continue
					}

					t.Errorf("timings overlap, %s %s %s", timings[i].Station, timings[i].Location, timings[i].Span)
				}
			}
		}
	},
	"check for unknown stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, v := range set.Timings() {
				if _, ok := set.Station(v.Station); !ok {
					t.Errorf("missing timings station: %s", v.Station)
				}
			}
		}
	},
	"check for unknown sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, v := range set.Timings() {
				if _, ok := set.Site(v.Station, v.Location); !ok {
					t.Errorf("missing timings site: %s %s", v.Station, v.Location)
				}
			}
		}
	},
}

func TestTimings(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range timingsChecks {
		t.Run(k, v(set))
	}
}

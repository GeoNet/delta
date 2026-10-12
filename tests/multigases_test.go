package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var multigasChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for calibration overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			multigases := set.Multigases()
			for i := 0; i < len(multigases); i++ {
				for j := i + 1; j < len(multigases); j++ {
					if multigases[i].Station != multigases[j].Station {
						continue
					}
					if multigases[i].Location != multigases[j].Location {
						continue
					}
					if multigases[i].Gas != multigases[j].Gas {
						continue
					}

					if multigases[i].End.Before(multigases[j].Start) {
						continue
					}
					if multigases[i].Start.After(multigases[j].End) {
						continue
					}
					if multigases[i].End.Equal(multigases[j].Start) {
						continue
					}
					if multigases[i].Start.Equal(multigases[j].End) {
						continue
					}

					t.Errorf("multigas calibration for %-5s/%-2s has overlap between %s and %s",
						multigases[i].Station, multigases[i].Location,
						multigases[i].Start.Format(meta.DateTimeFormat),
						multigases[i].End.Format(meta.DateTimeFormat))

				}
			}
		}
	},
	"check for missing calibration stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, m := range set.Multigases() {
				if _, ok := set.Site(m.Station, m.Location); !ok {
					t.Errorf("unable to find multigas site %-5s (%-2s)", m.Station, m.Location)
				}
			}
		}
	},
}

func TestMultigas(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range multigasChecks {
		t.Run(k, v(set))
	}
}

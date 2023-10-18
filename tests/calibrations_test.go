package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var calibrationChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for calibration overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.CalibrationList)
			for _, s := range set.Calibrations() {
				installs[s.Id()] = append(installs[s.Id()], s)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].End.Before(v[j].Start) {
							continue
						}
						if v[i].Start.After(v[j].End) {
							continue
						}
						if v[i].End.Equal(v[j].Start) {
							continue
						}
						if v[i].Start.Equal(v[j].End) {
							continue
						}

						t.Errorf("calibration %s/%s has number \"%d\" overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Number,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for missing assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Calibrations() {
				var found bool
				for _, a := range set.Assets() {
					if a.Model != s.Model {
						continue
					}
					if a.Serial != s.Serial {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find calibration asset: %s [%s]", s.Model, s.Serial)
			}
		}
	},
}

func TestCalibrations(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range calibrationChecks {
		t.Run(k, v(set))
	}
}

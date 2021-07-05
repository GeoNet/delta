package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testCalibrations = map[string]func([]meta.Calibration) func(t *testing.T){

	"check for calibration installation overlaps": func(installed []meta.Calibration) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.CalibrationList)
			for _, s := range installed {
				if _, ok := installs[s.Model]; !ok {
					installs[s.Model] = meta.CalibrationList{}
				}
				installs[s.Model] = append(installs[s.Model], s)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].Serial != v[j].Serial {
							continue
						}
						if v[i].Component != v[j].Component {
							continue
						}
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

						t.Errorf("calibration %s/%s has component %-2d overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Component,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
}

var testCalibrationsAssets = map[string]func([]meta.Calibration, []meta.Asset) func(t *testing.T){
	"check for missing assets": func(installed []meta.Calibration, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range installed {
				var found bool
				for _, a := range assets {
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
	var installed meta.CalibrationList
	loadListFile(t, "../install/calibrations.csv", &installed)

	for k, fn := range testCalibrations {
		t.Run(k, fn(installed))
	}

}

func TestCalibrations_Assets(t *testing.T) {
	var installed meta.CalibrationList
	loadListFile(t, "../install/calibrations.csv", &installed)

	var sensors meta.AssetList
	loadListFile(t, "../assets/sensors.csv", &sensors)

	for k, fn := range testCalibrationsAssets {
		t.Run(k, fn(installed, sensors))
	}
}

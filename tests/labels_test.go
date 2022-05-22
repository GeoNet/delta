package delta_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testLabels = map[string]func([]meta.Label) func(t *testing.T){
	"check for duplicated labels": func(labels []meta.Label) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(labels); i++ {
				for j := i + 1; j < len(labels); j++ {
					if labels[i].Type != labels[j].Type {
						continue
					}
					if labels[i].SamplingRate != labels[j].SamplingRate {
						continue
					}
					if labels[i].Azimuth != labels[j].Azimuth {
						continue
					}
					if labels[i].Dip != labels[j].Dip {
						continue
					}

					t.Errorf("label overlap: %s (%g) at %g/%g", labels[i].Type, labels[i].SamplingRate, labels[i].Azimuth, labels[i].Dip)
				}
			}
		}
	},

	"check for invalid codes": func(labels []meta.Label) func(t *testing.T) {
		return func(t *testing.T) {
			for _, l := range labels {
				if s := strings.ToUpper(l.Code); s != l.Code {
					t.Errorf("label code case issue: %s (%g) %s", l.Type, l.SamplingRate, l.Code)
				}
				if n := len(l.Code); n != 3 {
					t.Errorf("label code len issue: %s (%g) %s", l.Type, l.SamplingRate, l.Code)
				}
			}
		}
	},
	"check for invalid flags": func(labels []meta.Label) func(t *testing.T) {
		return func(t *testing.T) {
			for _, l := range labels {
				f := []byte(l.Flags)
				sort.Slice(f, func(i, j int) bool {
					return f[i] < f[j]
				})
				if s := string(f); s != l.Flags {
					t.Errorf("label flags sort issue: %s (%g) %s", l.Type, l.SamplingRate, l.Flags)
				}
			}
		}
	},
}

func TestLabels(t *testing.T) {

	var labels meta.LabelList
	loadListFile(t, "../response/labels.csv", &labels)

	for k, fn := range testLabels {
		t.Run(k, fn(labels))
	}
}

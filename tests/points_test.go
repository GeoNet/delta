package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var pointChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicated points": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			points := set.Points()
			for i := 0; i < len(points); i++ {
				for j := i + 1; j < len(points); j++ {
					if points[i].Sample != points[j].Sample {
						continue
					}
					if points[i].Location != points[j].Location {
						continue
					}
					t.Errorf("point duplication: %s/%s", points[i].Sample, points[i].Location)
				}
			}
		}
	},

	"check for duplicated sample points": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			samples := make(map[string]meta.Sample)
			for _, s := range set.Samples() {
				samples[s.Code] = s
			}
			for _, c := range set.Points() {
				if s, ok := samples[c.Sample]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: point start mismatch: " + strings.Join([]string{
							c.Sample,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: point end mismatch: " + strings.Join([]string{
							c.Sample,
							c.Location,
							c.End.String(),
							"after",
							s.End.String(),
						}, " "))
					}
				}
			}
		}
	},
}

func TestPoints(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range pointChecks {
		t.Run(k, v(set))
	}
}

package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var sessionChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check session overlap": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			sessions := set.Sessions()
			for i := 0; i < len(sessions); i++ {
				for j := i + 1; j < len(sessions); j++ {
					if sessions[i].Mark != sessions[j].Mark {
						continue
					}
					if sessions[i].Interval != sessions[j].Interval {
						continue
					}
					if sessions[i].End.Equal(sessions[j].Start) {
						t.Errorf("session start matches end: %s", strings.Join([]string{
							sessions[i].Mark,
							sessions[i].Interval.String(),
							sessions[i].Start.String(),
							sessions[i].End.String(),
							sessions[j].Interval.String(),
							sessions[j].Start.String(),
							sessions[j].End.String(),
						}, " "))
					}
					if sessions[i].Start.After(sessions[j].End) {
						continue
					}
					if !sessions[i].End.After(sessions[j].Start) {
						continue
					}
					t.Errorf("session overlap: %s", strings.Join([]string{
						sessions[i].Mark,
						sessions[i].Interval.String(),
						sessions[i].Start.String(),
						sessions[i].End.String(),
						sessions[j].Interval.String(),
						sessions[j].Start.String(),
						sessions[j].End.String(),
					}, " "))
				}
			}
		}
	},

	"check session spans": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Sessions() {
				if s.Start.After(s.End) {
					t.Errorf("session span mismatch: %s", strings.Join([]string{
						s.Mark,
						s.Interval.String(),
						s.Start.String(),
						"after",
						s.End.String(),
					}, " "))
				}
			}
		}
	},
	"check session satellite system": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Sessions() {
				switch s.SatelliteSystem {
				case "GPS":
				case "GPS+GLO":
				case "GPS+GLO+GAL+BDS+QZSS":
				default:
					t.Errorf("unknown satellite system: %s", s.SatelliteSystem)
				}
			}
		}
	},

	"check session marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			check := make(map[string]meta.Mark)
			for _, m := range set.Marks() {
				check[m.Code] = m
			}
			for _, s := range set.Sessions() {
				if _, ok := check[s.Mark]; !ok {
					t.Errorf("unknown session mark: %s", s.Mark)
				}
			}
		}
	},
}

func TestSessions(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range sessionChecks {
		t.Run(k, v(set))
	}
}

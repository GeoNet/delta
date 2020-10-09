package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testSessions = map[string]func([]meta.Session) func(t *testing.T){

	"check session overlap": func(sessions []meta.Session) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(sessions); i++ {
				for j := i + 1; j < len(sessions); j++ {
					if sessions[i].Mark != sessions[j].Mark {
						continue
					}
					if sessions[i].Interval != sessions[j].Interval {
						continue
					}
					if sessions[i].End.Equal(sessions[j].Start) {
						t.Errorf("session start matches end: " + strings.Join([]string{
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

	"check session spans": func(sessions []meta.Session) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range sessions {
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
	"check session satellite system": func(sessions []meta.Session) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range sessions {
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
}

var testSessionsMarks = map[string]func([]meta.Session, []meta.Mark) func(t *testing.T){

	"check session marks": func(sessions []meta.Session, marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {
			check := make(map[string]meta.Mark)
			for _, m := range marks {
				check[m.Code] = m
			}
			for _, s := range sessions {
				if _, ok := check[s.Mark]; !ok {
					t.Errorf("unknown session mark: %s", s.Mark)
				}
			}
		}
	},
}

func TestSessions(t *testing.T) {

	var sessions meta.SessionList
	loadListFile(t, "../install/sessions.csv", &sessions)

	for k, fn := range testSessions {
		t.Run(k, fn(sessions))
	}
}

func TestSessions_Marks(t *testing.T) {

	var sessions meta.SessionList
	loadListFile(t, "../install/sessions.csv", &sessions)

	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testSessionsMarks {
		t.Run(k, fn(sessions, marks))
	}
}

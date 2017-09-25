package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSessions(t *testing.T) {

	var sessions meta.SessionList
	loadListFile(t, "../install/sessions.csv", &sessions)

	t.Run("check for session overlaps", func(t *testing.T) {
		for i := 0; i < len(sessions); i++ {
			for j := i + 1; j < len(sessions); j++ {
				if sessions[i].Mark != sessions[j].Mark {
					continue
				}
				if sessions[i].Model != sessions[j].Model {
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
					continue
				}
				if sessions[i].Start.After(sessions[j].End) {
					continue
				}
				if !sessions[i].End.After(sessions[j].Start) {
					continue
				}
				t.Errorf("session overlap: " + strings.Join([]string{
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
	})

	t.Run("check for missing session marks", func(t *testing.T) {
		var list meta.MarkList
		loadListFile(t, "../network/marks.csv", &list)

		marks := make(map[string]meta.Mark)
		for _, m := range list {
			marks[m.Code] = m
		}
		for _, s := range sessions {
			if _, ok := marks[s.Mark]; !ok {
				t.Errorf("unknown session mark: " + s.Mark)
			}
		}
	})

	t.Run("check for session span mismatches", func(t *testing.T) {
		for _, s := range sessions {
			if s.Start.After(s.End) {
				t.Errorf("session span mismatch: " + strings.Join([]string{
					s.Mark,
					s.Interval.String(),
					s.Start.String(),
					"after",
					s.End.String(),
				}, " "))
			}
		}
	})

	t.Run("check for unknown session satellite systems", func(t *testing.T) {
		for _, s := range sessions {
			switch s.SatelliteSystem {
			case "GPS":
			case "GPS+GLO":
			case "GPS+GLO+GAL+BDS+QZSS":
			default:
				t.Error("unknown satellite system: " + s.SatelliteSystem)
			}
		}

	})
}

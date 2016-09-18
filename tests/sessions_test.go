package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSessions(t *testing.T) {

	var sessions meta.SessionList

	t.Log("Load sessions file")
	if err := meta.LoadList("../install/sessions.csv", &sessions); err != nil {
		t.Fatal(err)
	}

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

	marks := make(map[string]meta.Mark)
	{
		var list meta.MarkList
		t.Log("Load marks file")
		if err := meta.LoadList("../network/marks.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, m := range list {
			marks[m.Code] = m
		}
	}

	for _, s := range sessions {
		if _, ok := marks[s.Mark]; !ok {
			t.Log("unknown session mark: " + s.Mark)
		}
		if s.Start.After(s.End) {
			t.Log("session span mismatch: " + strings.Join([]string{
				s.Mark,
				s.Interval.String(),
				s.Start.String(),
				"after",
				s.End.String(),
			}, " "))
		}
		switch s.SatelliteSystem {
		case "GPS":
		case "GPS+GLO":
		case "GPS+GLO+GAL+BDS+QZSS":
		default:
			t.Error("unknown satellite system: " + s.SatelliteSystem)
		}
	}
}

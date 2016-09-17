package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestStreams(t *testing.T) {

	var streams meta.StreamList

	t.Log("Load streams file")
	if err := meta.LoadList("../install/streams.csv", &streams); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(streams); i++ {
		for j := i + 1; j < len(streams); j++ {
			if streams[i].Station != streams[j].Station {
				continue
			}
			if streams[i].Location != streams[j].Location {
				continue
			}
			if streams[i].Start.After(streams[j].End) {
				continue
			}
			if streams[i].End.Before(streams[j].Start) {
				continue
			}
			if streams[i].SamplingRate != streams[j].SamplingRate {
				continue
			}
			t.Errorf("stream overlap: " + strings.Join([]string{
				streams[i].Station,
				streams[i].Location,
				streams[i].Start.String(),
				streams[i].End.String(),
			}, " "))
		}
	}

	stas := make(map[string]meta.Station)
	{
		var list meta.StationList
		t.Log("Load stations file")
		if err := meta.LoadList("../network/stations.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			stas[s.Code] = s
		}
	}

	sites := make(map[string]map[string]meta.Site)
	{
		var list meta.SiteList
		t.Log("Load sites file")
		if err := meta.LoadList("../network/sites.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			if _, ok := sites[s.Station]; !ok {
				sites[s.Station] = make(map[string]meta.Site)
			}
			sites[s.Station][s.Location] = s
		}
	}

	for _, c := range streams {
		if _, ok := stas[c.Station]; !ok {
			t.Log("unknown stream station: " + c.Station)
		} else if s, ok := sites[c.Station]; !ok {
			t.Log("unknown stream station: " + c.Station)
		} else if _, ok := s[c.Location]; !ok {
			t.Log("unknown stream station/location: " + c.Station + "/" + c.Location)
		}
		if c.Start.After(c.End) {
			t.Log("stream span mismatch: " + strings.Join([]string{
				c.Station,
				c.Location,
				c.Start.String(),
				"after",
				c.End.String(),
			}, " "))
		}
	}
}

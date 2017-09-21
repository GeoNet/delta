package delta_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var recorderSamplingRates = []float64{50, 200}
var sensorSamplingRates = []float64{0.1, 1, 10, 50, 100, 200}

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

	for _, s := range streams {
		if s.SamplingRate == 0 {
			t.Errorf("invalid stream sample rate: " + strings.Join([]string{
				s.Station,
				s.Location,
				s.Start.String(),
				s.End.String(),
			}, " "))
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

	var assets = make(map[string]meta.Asset)
	{
		var list meta.AssetList
		t.Log("Load recorders assets file")
		if err := meta.LoadList("../assets/recorders.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, l := range list {
			assets[l.Model+":::"+l.Serial] = l
		}
		t.Log("Load sensors assets file")
		if err := meta.LoadList("../assets/sensors.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			assets[s.Model+":::"+s.Serial] = s
		}
	}

	var recorders []meta.InstalledRecorder
	{
		var list meta.InstalledRecorderList
		t.Log("Load recorders file")
		if err := meta.LoadList("../install/recorders.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, r := range list {
			recorders = append(recorders, r)
		}
	}

	var missing []meta.Stream

	for _, r := range recorders {
		if a, ok := assets[r.DataloggerModel+":::"+r.Serial]; !ok || a.Number == "" {
			continue
		}
		if r.End.Before(time.Now()) {
			continue
		}
		var handled bool
		for _, s := range streams {
			if s.Station != r.Station || r.Location != s.Location {
				continue
			}
			if r.Start.After(s.End) || r.End.Before(s.Start) {
				continue
			}
			handled = true
		}
		if !handled {
			for _, sps := range recorderSamplingRates {
				missing = append(missing, meta.Stream{
					Station:      r.Station,
					Location:     r.Location,
					SamplingRate: sps,
					Span: meta.Span{
						Start: r.Start,
						End:   r.End,
					},
				})
			}
			t.Errorf("no current stream defined for recorder: %s [%s/%s] %s %s", r.String(), r.Station, r.Location, r.Start, r.End)
		}
	}

	var sensors []meta.InstalledSensor
	{
		var list meta.InstalledSensorList
		t.Log("Load sensors file")
		if err := meta.LoadList("../install/sensors.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			sensors = append(sensors, s)
		}
	}

	for _, v := range sensors {
		if a, ok := assets[v.Model+":::"+v.Serial]; !ok || a.Number == "" {
			continue
		}
		if v.End.Before(time.Now()) {
			continue
		}
		var handled bool
		for _, s := range streams {
			if s.Station != v.Station || v.Location != s.Location {
				continue
			}
			if v.Start.After(s.End) || v.End.Before(s.Start) {
				continue
			}
			handled = true
		}
		if !handled {
			for _, sps := range sensorSamplingRates {
				missing = append(missing, meta.Stream{
					Station:      v.Station,
					Location:     v.Location,
					SamplingRate: sps,
					Span: meta.Span{
						Start: v.Start,
						End:   v.End,
					},
				})
			}
			t.Errorf("no current stream defined for sensor: %s [%s/%s] %s %s", v.String(), v.Station, v.Location, v.Start, v.End)
		}
	}

	if len(missing) > 0 {
		sort.Sort(meta.StreamList(missing))
		t.Log("\n" + string(meta.MarshalList(meta.StreamList(missing))))
	}

}

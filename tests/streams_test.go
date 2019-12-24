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
	loadListFile(t, "../install/streams.csv", &streams)

	t.Run("check for overlapping streams", func(t *testing.T) {
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
	})

	stas := make(map[string]meta.Station)
	t.Run("load stations file", func(t *testing.T) {
		var list meta.StationList
		loadListFile(t, "../network/stations.csv", &list)
		for _, s := range list {
			stas[s.Code] = s
		}
	})

	sites := make(map[string]map[string]meta.Site)
	t.Run("load sites file", func(t *testing.T) {
		var list meta.SiteList
		loadListFile(t, "../network/sites.csv", &list)
		for _, s := range list {
			if _, ok := sites[s.Station]; !ok {
				sites[s.Station] = make(map[string]meta.Site)
			}
			sites[s.Station][s.Location] = s
		}
	})

	t.Run("check for invalid stream sample rates", func(t *testing.T) {
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
	})

	t.Run("check for invalid stream spans", func(t *testing.T) {
		for _, c := range streams {
			if c.Start.After(c.End) {
				t.Error("stream span mismatch: " + strings.Join([]string{
					c.Station,
					c.Location,
					c.Start.String(),
					"after",
					c.End.String(),
				}, " "))
			}
		}
	})

	t.Run("check for invalid stream stations", func(t *testing.T) {
		for _, c := range streams {
			if _, ok := stas[c.Station]; !ok {
				t.Error("unknown stream station: " + c.Station)
			}
		}
	})

	t.Run("check for invalid dates: stream within station", func(t *testing.T) {
		for _, c := range streams {
			if s, ok := stas[c.Station]; ok {
				switch {
				case c.Start.Before(s.Start):
					t.Error("error: stream span mismatch: " + strings.Join([]string{
						c.Station,
						c.Location,
						c.Start.String(),
						"before",
						s.Start.String(),
					}, " "))
				case s.End.Before(time.Now()) && c.End.After(s.End):
					t.Error("error: stream span mismatch: " + strings.Join([]string{
						c.Station,
						c.Location,
						c.End.String(),
						"after",
						s.End.String(),
					}, " "))
				}
			}
		}
	})

	t.Run("check for invalid stream sites", func(t *testing.T) {
		for _, c := range streams {
			if _, ok := sites[c.Station]; !ok {
				t.Error("unknown stream station: " + c.Station)
			}
		}
	})

	t.Run("check for invalid stream locations", func(t *testing.T) {
		for _, c := range streams {
			if s, ok := sites[c.Station]; ok {
				if _, ok := s[c.Location]; !ok {
					t.Error("unknown stream station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	})

	t.Run("check for invalid dates: stream within site", func(t *testing.T) {
		for _, c := range streams {
			if s, ok := sites[c.Station]; ok {
				if l, ok := s[c.Location]; ok {
					switch {
					case c.Start.Before(l.Start):
						t.Error("error: stream span start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							l.Start.String(),
						}, " "))
					case l.End.Before(time.Now()) && c.End.After(l.End):
						t.Error("error: stream span end mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.End.String(),
							"after",
							l.End.String(),
						}, " "))
					}
				}
			}
		}
	})

	t.Run("check for invalid stream spans", func(t *testing.T) {
		for _, c := range streams {
			if c.Start.After(c.End) {
				t.Error("stream span mismatch: " + strings.Join([]string{
					c.Station,
					c.Location,
					c.Start.String(),
					"after",
					c.End.String(),
				}, " "))
			}
		}
	})

	var missing []meta.Stream

	t.Run("check for missing recorder streams", func(t *testing.T) {

		var list meta.AssetList
		loadListFile(t, "../assets/recorders.csv", &list)

		var assets = make(map[string]meta.Asset)
		for _, l := range list {
			assets[l.Model+":::"+l.Serial] = l
		}

		var recorders meta.InstalledRecorderList
		loadListFile(t, "../install/recorders.csv", &recorders)

		for _, r := range recorders {
			a, ok := assets[r.DataloggerModel+":::"+r.Serial]
			if !ok {
				continue
			}
			if a.Number == "" {
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
			if handled {
				continue
			}
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
			t.Errorf("no current stream defined for recorder: %s [%s/%s] %s %s",
				r.String(), r.Station, r.Location, r.Start, r.End)
		}
	})

	t.Run("check for missing sensor streams", func(t *testing.T) {
		var sensors meta.InstalledSensorList
		loadListFile(t, "../install/sensors.csv", &sensors)

		var list meta.AssetList
		loadListFile(t, "../assets/sensors.csv", &list)

		var assets = make(map[string]meta.Asset)
		for _, l := range list {
			assets[l.Model+":::"+l.Serial] = l
		}

		for _, v := range sensors {
			a, ok := assets[v.Model+":::"+v.Serial]
			if !ok {
				continue
			}
			if a.Number == "" {
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
			if handled {
				continue
			}
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
			t.Errorf("no current stream defined for sensor: %s [%s/%s] %s %s",
				v.String(), v.Station, v.Location, v.Start, v.End)
		}
	})

	sort.Sort(meta.StreamList(missing))

	if len(missing) > 0 {
		res, err := meta.MarshalList(meta.StreamList(missing))
		if err != nil {
			t.Fatal(err)
		}
		t.Error("\n" + string(res))
	}

	t.Run("check for invalid axial labels", func(t *testing.T) {
		for _, s := range streams {
			switch s.Axial {
			case "true", "false":
			case "ZNE", "Z12", "XYZ":
			default:
				t.Errorf("invalid stream axial code: %s", s.Axial)
			}
		}
	})
}

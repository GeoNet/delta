package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var testStreams = map[string]func([]meta.Stream) func(t *testing.T){
	"check for invalid axial labels": func(streams []meta.Stream) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range streams {
				switch s.Axial {
				case "true", "false":
				case "ZNE", "Z12", "XYZ":
				default:
					t.Errorf("invalid stream axial code: %s", s.Axial)
				}
			}
		}
	},

	"check for invalid stream span overlaps": func(streams []meta.Stream) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},

	"check for invalid stream spans": func(streams []meta.Stream) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},

	"check for invalid stream sample rates": func(streams []meta.Stream) func(t *testing.T) {

		return func(t *testing.T) {
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
		}
	},
}

var testStreamsStations = map[string]func([]meta.Stream, []meta.Station) func(t *testing.T){
	"check for invalid stream stations": func(streams []meta.Stream, list []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]interface{})
			for _, s := range list {
				stas[s.Code] = true
			}

			for _, c := range streams {
				if _, ok := stas[c.Station]; !ok {
					t.Error("unknown stream station: " + c.Station)
				}
			}
		}
	},
	"check for invalid dates: stream within station": func(streams []meta.Stream, list []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range list {
				stas[s.Code] = s
			}

			for _, c := range streams {
				if s, ok := stas[c.Station]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: stream span mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: stream span mismatch: " + strings.Join([]string{
							c.Station,
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

var testStreamsSites = map[string]func([]meta.Stream, []meta.Site) func(t *testing.T){
	"check for invalid stream sites": func(streams []meta.Stream, list []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[string]interface{})
			for _, s := range list {
				sites[s.Station] = true
			}

			for _, c := range streams {
				if _, ok := sites[c.Station]; !ok {
					t.Error("unknown stream station: " + c.Station)
				}
			}
		}
	},

	"check for invalid stream locations": func(streams []meta.Stream, list []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[struct {
				s string
				l string
			}]interface{})
			for _, s := range list {
				sites[struct {
					s string
					l string
				}{s: s.Station, l: s.Location}] = true
			}

			for _, c := range streams {
				if _, ok := sites[struct {
					s string
					l string
				}{s: c.Station, l: c.Location}]; !ok {
					t.Error("unknown stream station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	},

	"check for invalid dates: stream within site": func(streams []meta.Stream, list []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[struct {
				s string
				l string
			}]meta.Site)
			for _, s := range list {
				sites[struct {
					s string
					l string
				}{s: s.Station, l: s.Location}] = s
			}

			for _, c := range streams {
				if s, ok := sites[struct {
					s string
					l string
				}{s: c.Station, l: c.Location}]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: stream span start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: stream span end mismatch: " + strings.Join([]string{
							c.Station,
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

var testStreamsInstalledSensorsAssets = map[string]func([]meta.Stream, []meta.InstalledSensor, []meta.Asset) func(t *testing.T){
	"check for invalid stream sensor sites": func(streams []meta.Stream, sensors []meta.InstalledSensor, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			var list = make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range assets {
				list[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, v := range sensors {
				a, ok := list[struct {
					m string
					s string
				}{m: v.Model, s: v.Serial}]
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

				t.Errorf("no current stream defined for sensor: %s [%s/%s] %s %s",
					v.String(), v.Station, v.Location, v.Start, v.End)
			}
		}
	},
}

var testStreamsInstalledRecordersAssets = map[string]func([]meta.Stream, []meta.InstalledRecorder, []meta.Asset) func(t *testing.T){
	"check for invalid stream sites": func(streams []meta.Stream, recorders []meta.InstalledRecorder, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			var list = make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range assets {
				list[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, r := range recorders {
				a, ok := list[struct {
					m string
					s string
				}{m: r.DataloggerModel, s: r.Serial}]
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

				t.Errorf("no current stream defined for recorder: %s [%s/%s] %s %s",
					r.String(), r.Station, r.Location, r.Start, r.End)
			}
		}
	},
}

func TestStreams(t *testing.T) {

	var streams meta.StreamList
	loadListFile(t, "../install/streams.csv", &streams)

	for k, fn := range testStreams {
		t.Run(k, fn(streams))
	}
}

func TestStreams_Stations(t *testing.T) {

	var streams meta.StreamList
	loadListFile(t, "../install/streams.csv", &streams)

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testStreamsStations {
		t.Run(k, fn(streams, stations))
	}
}

func TestStreams_Sites(t *testing.T) {

	var streams meta.StreamList
	loadListFile(t, "../install/streams.csv", &streams)

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	for k, fn := range testStreamsSites {
		t.Run(k, fn(streams, sites))
	}
}

func TestStreams_InstalledSensorsAssets(t *testing.T) {
	var streams meta.StreamList
	loadListFile(t, "../install/streams.csv", &streams)

	var sensors meta.InstalledSensorList
	loadListFile(t, "../install/sensors.csv", &sensors)

	var assets meta.AssetList
	loadListFile(t, "../assets/sensors.csv", &assets)

	for k, fn := range testStreamsInstalledSensorsAssets {
		t.Run(k, fn(streams, sensors, assets))
	}
}

func TestStreams_InstalledRecordersAssets(t *testing.T) {
	var streams meta.StreamList
	loadListFile(t, "../install/streams.csv", &streams)

	var recorders meta.InstalledRecorderList
	loadListFile(t, "../install/recorders.csv", &recorders)

	var assets meta.AssetList
	loadListFile(t, "../assets/recorders.csv", &assets)

	for k, fn := range testStreamsInstalledRecordersAssets {
		t.Run(k, fn(streams, recorders, assets))
	}
}

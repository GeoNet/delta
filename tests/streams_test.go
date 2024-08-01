package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var streamChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for invalid axial labels": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Streams() {
				switch s.Axial {
				case "true", "false":
				case "ZNE", "Z12", "XYZ":
				default:
					t.Errorf("invalid stream axial code: %s", s.Axial)
				}
			}
		}
	},

	"check for invalid stream span overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Streams() {
				if c.Start.After(c.End) {
					t.Error("stream span mismatch: " + strings.Join([]string{
						c.Station,
						c.Location,
						c.Source,
						c.Start.String(),
						"after",
						c.End.String(),
					}, " "))
				}
			}
		}
	},

	"check for invalid stream spans": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			streams := set.Streams()
			for i := 0; i < len(streams); i++ {
				for j := i + 1; j < len(streams); j++ {
					if streams[i].Station != streams[j].Station {
						continue
					}
					if streams[i].Location != streams[j].Location {
						continue
					}
					if streams[i].Source != streams[j].Source {
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
						streams[i].Source,
						streams[i].Start.String(),
						streams[i].End.String(),
					}, " "))
				}
			}
		}
	},

	"check for invalid stream sample rates": func(set *meta.Set) func(t *testing.T) {

		return func(t *testing.T) {
			for _, s := range set.Streams() {
				if s.SamplingRate == 0 {
					t.Errorf("invalid stream sample rate: " + strings.Join([]string{
						s.Station,
						s.Location,
						s.Source,
						s.Start.String(),
						s.End.String(),
					}, " "))
				}
			}
		}
	},

	"check for invalid stream stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]interface{})
			for _, s := range set.Stations() {
				stas[s.Code] = true
			}

			for _, c := range set.Streams() {
				if _, ok := stas[c.Station]; !ok {
					t.Error("unknown stream station: " + c.Station)
				}
			}
		}
	},

	"check for invalid dates: stream within station": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stas[s.Code] = s
			}

			for _, c := range set.Streams() {
				if s, ok := stas[c.Station]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: stream span mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Source,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: stream span mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Source,
							c.End.String(),
							"after",
							s.End.String(),
						}, " "))
					}
				}
			}
		}
	},

	"check for invalid stream sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[string]interface{})
			for _, s := range set.Sites() {
				sites[s.Station] = true
			}

			for _, c := range set.Streams() {
				if _, ok := sites[c.Station]; !ok {
					t.Error("unknown stream station: " + c.Station)
				}
			}
		}
	},

	"check for invalid stream locations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[struct {
				s string
				l string
			}]interface{})
			for _, s := range set.Sites() {
				sites[struct {
					s string
					l string
				}{s: s.Station, l: s.Location}] = true
			}

			for _, c := range set.Streams() {
				if _, ok := sites[struct {
					s string
					l string
				}{s: c.Station, l: c.Location}]; !ok {
					t.Error("unknown stream station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	},

	"check for invalid dates: stream within site": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			sites := make(map[struct {
				s string
				l string
			}]meta.Site)
			for _, s := range set.Sites() {
				sites[struct {
					s string
					l string
				}{s: s.Station, l: s.Location}] = s
			}

			for _, c := range set.Streams() {
				if s, ok := sites[struct {
					s string
					l string
				}{s: c.Station, l: c.Location}]; ok {
					switch {
					case c.Start.Before(s.Start):
						t.Log("warning: stream span start mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Source,
							c.Start.String(),
							"before",
							s.Start.String(),
						}, " "))
					case s.End.Before(time.Now()) && c.End.After(s.End):
						t.Log("warning: stream span end mismatch: " + strings.Join([]string{
							c.Station,
							c.Location,
							c.Source,
							c.End.String(),
							"after",
							s.End.String(),
						}, " "))
					}
				}
			}
		}
	},

	"check for invalid stream sensor sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			var list = make(map[struct {
				m string
				s string
			}]meta.Asset)

			for _, a := range set.Assets() {
				list[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, v := range set.InstalledSensors() {
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
				for _, s := range set.Streams() {
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

	"check for invalid stream recorder sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			var list = make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range set.Assets() {
				list[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, r := range set.InstalledRecorders() {
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
				for _, s := range set.Streams() {
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

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range streamChecks {
		t.Run(k, v(set))
	}
}

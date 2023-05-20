package delta_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var connectionChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for connection overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			connections := set.Connections()
			for i := 0; i < len(connections); i++ {
				for j := i + 1; j < len(connections); j++ {
					if connections[i].Station != connections[j].Station {
						continue
					}
					if connections[i].Location != connections[j].Location {
						continue
					}
					if connections[i].Number != connections[j].Number {
						continue
					}
					if connections[i].Start.After(connections[j].End) {
						continue
					}
					if connections[i].End.Before(connections[j].Start) {
						continue
					}
					t.Errorf("connection overlap: " + strings.Join([]string{
						connections[i].Station,
						connections[i].Location,
						strconv.Itoa(connections[i].Number),
						connections[i].Start.String(),
						connections[i].End.String(),
					}, " "))
				}
			}
		}
	},

	"check for connection span mismatch": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Connections() {
				if c.Start.After(c.End) {
					t.Error("connection span mismatch: " + strings.Join([]string{
						c.Station,
						c.Location,
						strconv.Itoa(c.Number),
						c.Start.String(),
						"after",
						c.End.String(),
					}, " "))
				}
			}
		}
	},

	"check for missing connection stations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range set.Stations() {
				stas[s.Code] = s
			}
			for _, c := range set.Connections() {
				if _, ok := stas[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
		}
	},

	"check for missing connection sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			stations := make(map[string]meta.Site)
			for _, s := range set.Sites() {
				stations[s.Station] = s
			}

			for _, c := range set.Connections() {
				if _, ok := stations[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
		}
	},

	"check for missing connection site locations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			sites := make(map[string]map[string]meta.Site)
			for _, s := range set.Sites() {
				if _, ok := sites[s.Station]; !ok {
					sites[s.Station] = make(map[string]meta.Site)
				}
				sites[s.Station][s.Location] = s
			}

			for _, c := range set.Connections() {
				if _, ok := sites[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
			for _, c := range set.Connections() {
				if _, ok := sites[c.Station]; !ok {
					continue
				}

				if _, ok := sites[c.Station][c.Location]; !ok {
					t.Error("unknown connection station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	},

	"check for missing datalogger places": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			places := make(map[string]string)
			for _, d := range set.DeployedDataloggers() {
				switch d.Role {
				case "":
					places[d.Place] = d.Place
				default:
					places[d.Place+"/"+d.Role] = d.Place
				}
			}
			for _, c := range set.Connections() {
				switch c.Role {
				case "":
					if _, ok := places[c.Place]; !ok {
						t.Error("error: unknown datalogger place: " + c.Place)
					}
				default:
					if _, ok := places[c.Place+"/"+c.Role]; !ok {
						t.Error("error: unknown datalogger place/role: " + c.Place + "/" + c.Role)
					}
				}
			}
		}
	},

	"check for missing sensor connections": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			assets := make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range set.Assets() {
				assets[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, s := range set.InstalledSensors() {
				if a, ok := assets[struct {
					m string
					s string
				}{m: s.Model, s: s.Serial}]; !ok || a.Number == "" {
					continue
				}
				if s.End.Before(time.Now()) {
					continue
				}
				var handled bool
				for _, c := range set.Connections() {
					if c.Station != s.Station || c.Location != s.Location {
						continue
					}
					if c.Start.After(s.End) || c.End.Before(s.Start) {
						continue
					}
					handled = true
				}
				if handled {
					continue
				}

				t.Errorf("no current connection defined for sensor: %s [%s/%s] %s %s",
					s.String(), s.Station, s.Location, s.Start, s.End)
			}
		}
	},

	"check for missing datalogger connections": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			assets := make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range set.Assets() {
				assets[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, d := range set.DeployedDataloggers() {

				if a, ok := assets[struct {
					m string
					s string
				}{m: d.Model, s: d.Serial}]; !ok || a.Number == "" {
					continue
				}

				if d.End.Before(time.Now()) {
					continue
				}
				var handled bool
				for _, c := range set.Connections() {
					if c.Place != d.Place || c.Role != d.Role {
						continue
					}
					if c.Start.After(d.End) || c.End.Before(d.Start) {
						continue
					}
					handled = true
				}
				if handled {
					continue
				}

				t.Errorf("no current connection defined for datalogger: %s [%s/%s] %s %s",
					d.String(), d.Place, d.Role, d.Start, d.End)
			}
		}
	},
}

func TestConnections(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range connectionChecks {
		t.Run(k, v(set))
	}
}

package delta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var testConnections = map[string]func([]meta.Connection) func(t *testing.T){

	"check for connection overlaps": func(connections []meta.Connection) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(connections); i++ {
				for j := i + 1; j < len(connections); j++ {
					if connections[i].Station != connections[j].Station {
						continue
					}
					if connections[i].Location != connections[j].Location {
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
						connections[i].Start.String(),
						connections[i].End.String(),
					}, " "))
				}
			}
		}
	},
	"check for connection span mismatch": func(connections []meta.Connection) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range connections {
				if c.Start.After(c.End) {
					t.Error("connection span mismatch: " + strings.Join([]string{
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
}

var testConnectionsStations = map[string]func([]meta.Connection, []meta.Station) func(t *testing.T){
	"check for missing connection stations": func(connections []meta.Connection, list []meta.Station) func(t *testing.T) {
		return func(t *testing.T) {
			stas := make(map[string]meta.Station)
			for _, s := range list {
				stas[s.Code] = s
			}
			for _, c := range connections {
				if _, ok := stas[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
		}
	},
}

var testConnectionsSites = map[string]func([]meta.Connection, []meta.Site) func(t *testing.T){
	"check for missing connection stations": func(connections []meta.Connection, list []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {

			stations := make(map[string]meta.Site)
			for _, s := range list {
				stations[s.Station] = s
			}

			for _, c := range connections {
				if _, ok := stations[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
		}
	},
	"check for missing connection site locations": func(connections []meta.Connection, list []meta.Site) func(t *testing.T) {
		return func(t *testing.T) {

			sites := make(map[string]map[string]meta.Site)
			for _, s := range list {
				if _, ok := sites[s.Station]; !ok {
					sites[s.Station] = make(map[string]meta.Site)
				}
				sites[s.Station][s.Location] = s
			}

			for _, c := range connections {
				if _, ok := sites[c.Station]; !ok {
					t.Error("unknown connection station: " + c.Station)
				}
			}
			for _, c := range connections {
				if _, ok := sites[c.Station]; !ok {
					continue
				}

				if _, ok := sites[c.Station][c.Location]; !ok {
					t.Error("unknown connection station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	},
}

var testConnectionsDeployedDataloggers = map[string]func([]meta.Connection, []meta.DeployedDatalogger) func(t *testing.T){
	"check for missing datalogger places": func(connections []meta.Connection, list []meta.DeployedDatalogger) func(t *testing.T) {
		return func(t *testing.T) {

			places := make(map[string]string)
			for _, d := range list {
				switch d.Role {
				case "":
					places[d.Place] = d.Place
				default:
					places[d.Place+"/"+d.Role] = d.Place
				}
			}
			for _, c := range connections {
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
}

var testConnectionsInstalledSensorAssets = map[string]func([]meta.Connection, []meta.InstalledSensor, []meta.Asset) func(t *testing.T){
	"check for missing sensor connections": func(connections []meta.Connection, sensors []meta.InstalledSensor, list []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			assets := make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range list {
				assets[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, s := range sensors {
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
				for _, c := range connections {
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
}

var testConnectionsDeployedDataloggersAssets = map[string]func([]meta.Connection, []meta.DeployedDatalogger, []meta.Asset) func(t *testing.T){
	"check for missing datalogger connections": func(connections []meta.Connection, dataloggers []meta.DeployedDatalogger, list []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			assets := make(map[struct {
				m string
				s string
			}]meta.Asset)
			for _, a := range list {
				assets[struct {
					m string
					s string
				}{m: a.Model, s: a.Serial}] = a
			}

			for _, d := range dataloggers {

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
				for _, c := range connections {
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

	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	for k, fn := range testConnections {
		t.Run(k, fn(connections))
	}
}

func TestConnections_Stations(t *testing.T) {

	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	var stations meta.StationList
	loadListFile(t, "../network/stations.csv", &stations)

	for k, fn := range testConnectionsStations {
		t.Run(k, fn(connections, stations))
	}
}

func TestConnections_Sites(t *testing.T) {

	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	var sites meta.SiteList
	loadListFile(t, "../network/sites.csv", &sites)

	for k, fn := range testConnectionsSites {
		t.Run(k, fn(connections, sites))
	}
}

func TestConnections_DeployedDataloggers(t *testing.T) {
	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	var dataloggers meta.DeployedDataloggerList
	loadListFile(t, "../install/dataloggers.csv", &dataloggers)

	for k, fn := range testConnectionsDeployedDataloggers {
		t.Run(k, fn(connections, dataloggers))
	}
}

func TestConnections_DeployedDataloggersAssets(t *testing.T) {
	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	var dataloggers meta.DeployedDataloggerList
	loadListFile(t, "../install/dataloggers.csv", &dataloggers)

	var assets meta.AssetList
	loadListFile(t, "../assets/dataloggers.csv", &assets)

	for k, fn := range testConnectionsDeployedDataloggersAssets {
		t.Run(k, fn(connections, dataloggers, assets))
	}
}

func TestConnections_InstalledSensorAssets(t *testing.T) {
	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	var sensors meta.InstalledSensorList
	loadListFile(t, "../install/sensors.csv", &sensors)

	var assets meta.AssetList
	loadListFile(t, "../assets/sensors.csv", &assets)

	for k, fn := range testConnectionsInstalledSensorAssets {
		t.Run(k, fn(connections, sensors, assets))
	}
}

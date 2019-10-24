package delta_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

func TestConnections(t *testing.T) {

	var connections meta.ConnectionList
	loadListFile(t, "../install/connections.csv", &connections)

	t.Run("check for connection overlaps", func(t *testing.T) {
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
	})

	stas := make(map[string]meta.Station)
	t.Run("load stations list", func(t *testing.T) {
		var list meta.StationList
		loadListFile(t, "../network/stations.csv", &list)
		for _, s := range list {
			stas[s.Code] = s
		}
	})

	sites := make(map[string]map[string]meta.Site)
	t.Run("load sites list", func(t *testing.T) {
		var list meta.SiteList
		loadListFile(t, "../network/sites.csv", &list)
		for _, s := range list {
			if _, ok := sites[s.Station]; !ok {
				sites[s.Station] = make(map[string]meta.Site)
			}
			sites[s.Station][s.Location] = s
		}
	})

	t.Run("check for missing connection stations", func(t *testing.T) {
		for _, c := range connections {
			if _, ok := stas[c.Station]; !ok {
				t.Error("unknown connection station: " + c.Station)
			}
		}
	})

	t.Run("check for missing connection sites", func(t *testing.T) {
		for _, c := range connections {
			if _, ok := sites[c.Station]; !ok {
				t.Error("unknown connection station: " + c.Station)
			}
		}
	})

	t.Run("check for missing connection locations", func(t *testing.T) {
		for _, c := range connections {
			if s, ok := sites[c.Station]; ok {
				if _, ok := s[c.Location]; !ok {
					t.Error("unknown connection station/location: " + c.Station + "/" + c.Location)
				}
			}
		}
	})

	t.Run("check for connection span mismatch", func(t *testing.T) {
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
	})

	t.Run("check for missing connection places", func(t *testing.T) {
		var list meta.DeployedDataloggerList
		loadListFile(t, "../install/dataloggers.csv", &list)

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
	})

	var missing []meta.Connection

	t.Run("check for missing sensor connections", func(t *testing.T) {
		var list meta.AssetList
		loadListFile(t, "../assets/sensors.csv", &list)

		var assets = make(map[string]meta.Asset)
		for _, s := range list {
			assets[s.Model+":::"+s.Serial] = s
		}

		var sensors meta.InstalledSensorList
		loadListFile(t, "../install/sensors.csv", &sensors)

		for _, s := range sensors {
			if a, ok := assets[s.Model+":::"+s.Serial]; !ok || a.Number == "" {
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

			var place string
			if p, ok := stas[s.Station]; ok {
				place = p.Name
			}
			missing = append(missing, meta.Connection{
				Station:  s.Station,
				Location: s.Location,
				Place:    place,
				Span: meta.Span{
					Start: s.Start,
					End:   s.End,
				},
			})
			t.Errorf("no current connection defined for sensor: %s [%s/%s] %s %s",
				s.String(), s.Station, s.Location, s.Start, s.End)
		}
	})

	t.Run("check for missing datalogger connections", func(t *testing.T) {

		var list meta.AssetList
		loadListFile(t, "../assets/dataloggers.csv", &list)

		var assets = make(map[string]meta.Asset)
		for _, d := range list {
			assets[d.Model+":::"+d.Serial] = d
		}

		var dataloggers meta.DeployedDataloggerList
		loadListFile(t, "../install/dataloggers.csv", &dataloggers)

		for _, d := range dataloggers {
			if a, ok := assets[d.Model+":::"+d.Serial]; !ok || a.Number == "" {
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

			sta, loc := "XXXX", "LL"
			for k, v := range stas {
				if v.Name == d.Place {
					sta = k
				}
			}

			missing = append(missing, meta.Connection{
				Station:  sta,
				Location: loc,
				Place:    d.Place,
				Role:     d.Role,
				Span: meta.Span{
					Start: d.Start,
					End:   d.End,
				},
			})

			t.Errorf("no current connection defined for datalogger: %s [%s/%s] %s %s",
				d.String(), d.Place, d.Role, d.Start, d.End)
		}
	})

	sort.Sort(meta.ConnectionList(missing))

	if len(missing) > 0 {
		res, err := meta.MarshalList(meta.ConnectionList(missing))
		if err != nil {
			t.Fatal(err)
		}
		t.Error("\n" + string(res))
	}

}

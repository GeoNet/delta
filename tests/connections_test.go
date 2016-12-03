package delta_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

func TestConnections(t *testing.T) {
	var recorders meta.InstalledRecorderList

	if err := meta.LoadList("../install/recorders.csv", &recorders); err != nil {
		t.Fatal(err)
	}

	var connections meta.ConnectionList

	if err := meta.LoadList("../install/connections.csv", &connections); err != nil {
		t.Fatal(err)
	}

	var sensors []meta.InstalledSensor
	{
		var list meta.InstalledSensorList
		t.Log("Load stations file")
		if err := meta.LoadList("../install/sensors.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			sensors = append(sensors, s)
		}
	}

	var dataloggers []meta.DeployedDatalogger
	places := make(map[string]string)
	{
		var list meta.DeployedDataloggerList
		t.Log("Load stations file")
		if err := meta.LoadList("../install/dataloggers.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, d := range list {
			dataloggers = append(dataloggers, d)
			if d.Role != "" {
				places[d.Place+"/"+d.Role] = d.Place
			} else {
				places[d.Place] = d.Place
			}
		}
	}

	stations := make(map[string]meta.Station)
	{
		var list meta.StationList
		t.Log("Load stations file")
		if err := meta.LoadList("../network/stations.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			stations[s.Code] = s
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

	var assets = make(map[string]meta.Asset)
	{
		var list meta.AssetList
		t.Log("Load datalogger assets file")
		if err := meta.LoadList("../assets/dataloggers.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, d := range list {
			assets[d.Model+":::"+d.Serial] = d
		}
		t.Log("Load sensor assets file")
		if err := meta.LoadList("../assets/sensors.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, s := range list {
			assets[s.Model+":::"+s.Serial] = s
		}
	}

	t.Run("check for connection overlap", func(t *testing.T) {
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

	t.Run("check connection settings", func(t *testing.T) {
		for _, c := range connections {
			if _, ok := stations[c.Station]; !ok {
				t.Log("unknown connection station: " + c.Station)
			} else if s, ok := sites[c.Station]; !ok {
				t.Log("unknown connection station: " + c.Station)
			} else if _, ok := s[c.Location]; !ok {
				t.Log("unknown connection station/location: " + c.Station + "/" + c.Location)
			}
			if c.Start.After(c.End) {
				t.Log("connection span mismatch: " + strings.Join([]string{
					c.Station,
					c.Location,
					c.Start.String(),
					"after",
					c.End.String(),
				}, " "))
			}
		}
	})

	t.Run("check connection place and roles", func(t *testing.T) {
		for _, c := range connections {
			if c.Role != "" {
				if _, ok := places[c.Place+"/"+c.Role]; !ok {
					t.Log("warning: unknown datalogger place/role: " + c.Place + "/" + c.Role)
				}
			} else {
				if _, ok := places[c.Place]; !ok {
					t.Log("warning: unknown datalogger place: " + c.Place)
				}
			}
		}
	})

	t.Run("check for missing connections", func(t *testing.T) {
		var missing []meta.Connection

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
			if !handled {
				var place string
				if p, ok := stations[s.Station]; ok {
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
				t.Errorf("no current connection defined for sensor: %s [%s/%s] %s %s", s.String(), s.Station, s.Location, s.Start, s.End)
			}
		}

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
			if !handled {
				s, l := "XXXX", "LL"
				for k, v := range stations {
					if v.Name == d.Place {
						s = k
					}
				}

				missing = append(missing, meta.Connection{
					Station:  s,
					Location: l,
					Place:    d.Place,
					Role:     d.Role,
					Span: meta.Span{
						Start: d.Start,
						End:   d.End,
					},
				})

				t.Errorf("no current connection defined for datalogger: %s [%s/%s] %s %s", d.String(), d.Place, d.Role, d.Start, d.End)
			}
		}

		if len(missing) > 0 {
			sort.Sort(meta.ConnectionList(missing))
			t.Log("\n" + string(meta.MarshalList(meta.ConnectionList(missing))))
		}
	})

}

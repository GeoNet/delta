package delta_test

import (
	"strings"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestConnections(t *testing.T) {

	var connections meta.ConnectionList

	t.Log("Load connections file")
	if err := meta.LoadList("../install/connections.csv", &connections); err != nil {
		t.Fatal(err)
	}

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

	for _, c := range connections {
		if _, ok := stas[c.Station]; !ok {
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

	places := make(map[string]string)
	{
		var list meta.DeployedDataloggerList
		t.Log("Load installed dataloggers file")
		if err := meta.LoadList("../install/dataloggers.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, d := range list {
			if d.Role != "" {
				places[d.Place+"/"+d.Role] = d.Place
			} else {
				places[d.Place] = d.Place
			}
		}
	}

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

	/*
		Station Code,Location Code,Datalogger Place,Datalogger Role,Pre Amp,Gain,Start Date,End Date
		A11,10,Alfredton,Short Period,false,0,1993-02-01T01:23:00Z,1993-02-21T08:48:00Z
		ABAZ,10,Whangaparaoa Navy Base,,false,0,2008-10-13T04:00:00Z,9999-01-01T00:00:00Z
		AC1A,10,Wards Pass #1,,false,0,2001-09-12T03:00:01Z,2002-01-22T22:00:00Z
		AC2A,10,Wards Pass #2,,false,0,2001-09-13T02:00:01Z,2002-01-22T22:00:00Z
		AC3A,10,Acheron #3,,false,0,2001-07-31T23:00:00Z,2002-01-22T23:00:00Z
		AC4A,10,Acheron #4,Short Period,false,0,2001-08-01T01:00:00Z,2002-01-22T23:00:00Z
		AC5A,10,Wards Pass #5,,false,0,2001-08-01T04:00:00Z,2001-11-27T23:00:00Z
		AGA,10,Angora Road,Short Period,false,0,1990-02-20T22:33:00Z,1992-04-07T17:01:00Z
		AHAA,10,Ahaura,Short Period,false,0,1995-11-19T17:31:00Z,1995-12-12T01:28:00Z
	*/

	/*
	 */
}

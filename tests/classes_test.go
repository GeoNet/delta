package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestClasses(t *testing.T) {

	var classes meta.ClassList
	t.Log("Load site classes file")
	{
		if err := meta.LoadList("../network/classes.csv", &classes); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(classes); i++ {
		for j := i + 1; j < len(classes); j++ {
			if classes[i].Station == classes[j].Station {
				t.Errorf("class site duplication: " + classes[i].Station + " <=> " + classes[j].Station)
			}
		}
	}

	stations := make(map[string]meta.Station)
	t.Log("Load stations file")
	{
		var list meta.StationList
		if err := meta.LoadList("../network/stations.csv", &list); err != nil {
			t.Fatal(err)
		}
		for _, l := range list {
			stations[l.Code] = l
		}
	}

	for _, c := range classes {
		if _, ok := stations[c.Station]; !ok {
			t.Errorf("class station missing: " + c.Station)
		}
	}

	for _, c := range classes {
		switch c.Class {
		case "A", "B", "C", "D", "E":
		default:
			t.Errorf("class invalid class: " + c.Station + " " + c.Class)
		}
		switch c.QVs30 {
		case "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Vs30: " + c.Station + " " + c.QVs30)
		}
		switch c.QTsite {
		case "I", "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Tsite: " + c.Station + " " + c.QTsite)
		}
		switch c.DTsite {
		case "I", "Ms", "Mw", "Mn", "Mu", "Ma":
		default:
			t.Errorf("class invalid D_Tsite: " + c.Station + " " + c.DTsite)
		}

		switch c.QZb {
		case "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Zb: " + c.Station + " " + c.QZb)
		}

	}

}

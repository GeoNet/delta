package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestVs30s(t *testing.T) {

	var vs30s meta.Vs30List
	t.Log("Load site vs30 file")
	{
		if err := meta.LoadList("../environment/vs30.csv", &vs30s); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(vs30s); i++ {
		for j := i + 1; j < len(vs30s); j++ {
			if vs30s[i].Station == vs30s[j].Station {
				t.Errorf("class site duplication: " + vs30s[i].Station + " <=> " + vs30s[j].Station)
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

	for _, v := range vs30s {
		if _, ok := stations[v.Station]; !ok {
			t.Errorf("class station missing: " + v.Station)
		}
	}

	for _, v := range vs30s {
		switch v.QVs30 {
		case "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Vs30: " + v.Station + " " + v.QVs30)
		}
		switch v.QTsite {
		case "I", "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Tsite: " + v.Station + " " + v.QTsite)
		}
		switch v.DTsite {
		case "I", "Ms", "Mw", "Mn", "Mu", "Ma":
		default:
			t.Errorf("class invalid D_Tsite: " + v.Station + " " + v.DTsite)
		}

		switch v.QZb {
		case "Q3", "Q2", "Q1":
		default:
			t.Errorf("class invalid Q_Zb: " + v.Station + " " + v.QZb)
		}

	}

}

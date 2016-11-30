package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
	"io/ioutil"
	"os"
	"github.com/GeoNet/delta"
	"github.com/golang/protobuf/proto"
)

const apiDir = "../.tmp/api/delta"

func TestMarks(t *testing.T) {

	var marks meta.MarkList
	t.Log("Load network marks file")
	{
		if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(marks); i++ {
		for j := i + 1; j < len(marks); j++ {
			if marks[i].Code == marks[j].Code {
				t.Errorf("mark duplication: " + marks[i].Code)
			}
		}
	}

}

// TestMarksProto creates a Protobuf file of Marks.
 func TestMarksProto(t *testing.T) {
	var networks meta.NetworkList

	if err := meta.LoadList("../network/networks.csv", &networks); err != nil {
		t.Error(err)
	}

	 var net = make(map[string]*delta.Network)

	 for _, v := range networks {
		 n := delta.Network{
			 Code: v.Code,
			 External: v.External,
			 Description: v.Description,
			 Restricted: v.Restricted,
		 }

		 net[v.Code] = &n
	 }

	var marks meta.MarkList
	if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
		t.Error(err)
	}

	if len(marks) == 0 {
		t.Error("zero length mark list.")
	}

	var m delta.Marks

	m.Marks = make(map[string]*delta.Mark)

	for _, v := range marks {
		pt := delta.Point{
			Longitude:v.Longitude,
			Latitude:v.Latitude,
			Elevation:v.Elevation,
			Datum:v.Datum,
		}

		s := delta.Span{
			Start:v.Start.Unix(),
			End:v.End.Unix(),
		}

		mk := delta.Mark{
			Code:v.Code,
			Name:v.Name,
			Network: net[v.Network],
			Point:&pt,
			Span:&s,
		}

		m.Marks[mk.Code] = &mk
	}

	b, err := proto.Marshal(&m)
	if err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(apiDir, 0777); err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(apiDir + "/marks.pb", b, 0644); err != nil {
		t.Error(err)
	}
}

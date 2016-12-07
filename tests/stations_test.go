package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta"
	"github.com/golang/protobuf/proto"
	"os"
	"io/ioutil"
	"fmt"
	"bytes"
)

func TestStations(t *testing.T) {

	var stations meta.StationList
	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < len(stations); i++ {
		for j := i + 1; j < len(stations); j++ {
			if stations[i].Code == stations[j].Code {
				t.Errorf("station duplication: " + stations[i].Code)
			}
		}
	}

}

// TestStationsProto creates a Protobuf file of Stations.
// Creates a GeoJSON file of gauge locations.
func TestStationsProto(t *testing.T) {
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

	var stations meta.StationList
	if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
		t.Error(err)
	}

	if len(stations) == 0 {
		t.Error("zero length stations list.")
	}

	var s delta.Stations
	s.Stations = make(map[string]*delta.Station)

	for _, v := range stations {
		pt := delta.Point{
			Longitude:v.Longitude,
			Latitude:v.Latitude,
			Elevation:v.Elevation,
			Datum:v.Datum,
		}

		sp := delta.Span{
			Start:v.Start.Unix(),
			End:v.End.Unix(),
		}

		st := delta.Station{
			Code:v.Code,
			Name:v.Name,
			Network: net[v.Network],
			Point:&pt,
			Span:&sp,
		}

		s.Stations[st.Code] = &st
	}

	b, err := proto.Marshal(&s)
	if err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(apiDir, 0777); err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(apiDir + "/stations.pb", b, 0644); err != nil {
		t.Error(err)
	}

	// GeoJSON for sea level gauges.
	// Output the Station long/lat.  This is not the same as the
	// sensor locations.  Will need to consider adding this as a
	// separate layer (station versus location).
	// Assume that the network code is enough to ID the gauges.
	var bf bytes.Buffer

	bf.WriteString(`{"type": "FeatureCollection","features": [`)

	f := false
	for _, v := range s.Stations {
		if v.GetSpan() == nil || v.GetPoint() == nil || v.GetNetwork() == nil || v.Network.Restricted {
			continue
		}

		if v.Network.Code != "TG" {
			continue
		}

		if f {
			bf.WriteString(`,`)
		}
		f = true

		bf.WriteString(`{"type":"Feature","geometry":{"type": "Point","coordinates": [`)
		bf.WriteString(fmt.Sprintf("%f,%f", v.Point.Longitude, v.Point.Latitude))
		bf.WriteString(`]},"properties":{`)
		bf.WriteString(fmt.Sprintf("\"code\":\"%s\"", v.Code))
		bf.WriteString(fmt.Sprintf(",\"name\":\"%s\"", v.Name))
		bf.WriteString(fmt.Sprintf(",\"datum\":\"%s\"", v.Point.Datum))
		bf.WriteString(fmt.Sprintf(",\"elevation\":%f", v.Point.Elevation))
		bf.WriteString(fmt.Sprintf(",\"start\":\"%d\"", v.Span.Start))
		bf.WriteString(fmt.Sprintf(",\"end\":\"%d\"", v.Span.End))
		bf.WriteString(fmt.Sprintf(",\"network\":\"%s\"", v.Network.External))
		bf.WriteString(`}}`)
	}

	bf.WriteString(`]}`)

	if err := ioutil.WriteFile(apiDir + "/gauges.geojson", bf.Bytes(), 0644); err != nil {
		t.Error(err)
	}
}

package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta"
	"github.com/golang/protobuf/proto"
	"os"
	"io/ioutil"
	"github.com/GeoNet/delta/resp"
	"bytes"
	"fmt"
	"strings"
	"encoding/json"
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

// TestStationsProto creates:
// Protobuf file of Stations
// JSON file of Stations
// GeoJSON of Site for each sensor type.
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

	var sites meta.SiteList
	if err := meta.LoadList("../network/sites.csv", &sites); err != nil {
		t.Error(err)
	}

	for _, v := range sites {
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
		si := delta.Site{
			Location: v.Location,
			Survey: v.Survey,
			Span: &sp,
			Point: &pt,
		}

		if _, ok := s.Stations[v.Station]; ok {
			if s.Stations[v.Station].Sites == nil {
				s.Stations[v.Station].Sites = make(map[string]*delta.Site)
			}
			s.Stations[v.Station].Sites[si.Location] = &si
		}
	}

	var installed meta.InstalledSensorList
	if err := meta.LoadList("../install/sensors.csv", &installed); err != nil {
		t.Fatal(err)
	}

	for _, v := range installed {
		st := resp.SensorModels[v.Model]

		e := delta.Equipment{
			Make: v.Make,
			Model: v.Model,
			Serial: v.Serial,
			Type: st.Type,
		}
		sp := delta.Span{
			Start:v.Start.Unix(),
			End:v.End.Unix(),
		}
		o := delta.Orientation{
			Dip: v.Dip,
			Azimuth: v.Azimuth,
		}
		off := delta.Offset{
			North: v.North,
			East: v.East,
			Vertical: v.Vertical,
		}
		sc := delta.Scale{
			Bias: v.Bias,
			Factor: v.Factor,
		}
		is := delta.InstalledSensor{
			Equipment: &e,
			Span: &sp,
			Orientation: &o,
			Offset: &off,
			Scale: &sc,
		}

		if _, ok := s.Stations[v.Station]; ok {
			if _, ok := s.Stations[v.Station].Sites[v.Location]; ok {
				s.Stations[v.Station].Sites[v.Location].InstalledSensor = append(s.Stations[v.Station].Sites[v.Location].InstalledSensor, &is)
			}
		}
	}

	// strong motion recorders that are a sensor and a datalogger as a package.
	var recorders meta.InstalledRecorderList
	if err := meta.LoadList("../install/recorders.csv", &recorders); err != nil {
		t.Fatal(err)
	}

	for _, v := range recorders {
		st, ok := resp.SensorModels[v.InstalledSensor.Model]
		if !ok {
			// in resp/auto.go some models have " SENSOR" appended.
			st, ok = resp.SensorModels[v.InstalledSensor.Model + " SENSOR"]
		}

		e := delta.Equipment{
			Make: v.Make,
			Model: v.Model,
			Serial: v.Serial,
			Type: st.Type,
		}
		sp := delta.Span{
			Start:v.Start.Unix(),
			End:v.End.Unix(),
		}
		o := delta.Orientation{
			Dip: v.Dip,
			Azimuth: v.Azimuth,
		}
		off := delta.Offset{
			North: v.North,
			East: v.East,
			Vertical: v.Vertical,
		}
		sc := delta.Scale{
			Bias: v.Bias,
			Factor: v.Factor,
		}
		is := delta.InstalledSensor{
			Equipment: &e,
			Span: &sp,
			Orientation: &o,
			Offset: &off,
			Scale: &sc,
		}

		if _, ok := s.Stations[v.Station]; ok {
			if _, ok := s.Stations[v.Station].Sites[v.Location]; ok {
				s.Stations[v.Station].Sites[v.Location].InstalledSensor = append(s.Stations[v.Station].Sites[v.Location].InstalledSensor, &is)
			}
		}
	}

	// protobuf of all station information

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

	// json of all station information

	b, err = json.Marshal(&s)
	if err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(apiDir, 0777); err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(apiDir + "/stations.json", b, 0644); err != nil {
		t.Error(err)
	}

	// GeoJSON files of site for each sensor type.

	var out = make(map[string]*bytes.Buffer)

	for _, v:= range resp.SensorModels {
		if v.Type != "" {
			out[v.Type] = &bytes.Buffer{}
		}
	}

	for k := range out {
		out[k].WriteString(`{"type": "FeatureCollection","features": [`)
	}

	for _, station := range s.Stations {
		if station.GetSpan() == nil || station.GetPoint() == nil || station.GetNetwork() == nil || station.Network.Restricted {
			continue
		}
		if station.Network.External != "NZ" || station.Network.Code == "XX" {
			continue
		}

		for _, site := range station.Sites {
			for _, v := range site.InstalledSensor {
				_, ok := out[v.Equipment.Type]

				switch  {
				case v.Equipment.Type == "":
					t.Logf("%s.%s no sensor model for %s, site not classified.",
						station.Code, site.Location, v.Equipment.Model)
				case ok:
					writeProps(station, site, v, out[v.Equipment.Type])
				case !ok:
					t.Logf("unclassified type %s", v.Equipment.Type)
				}
			}
		}
	}

	for k := range out {
		out[k].WriteString(`]}`)

		fn := strings.Replace(strings.ToLower(k), " ", "", -1)
		fn = apiDir + "/" + fn + ".geojson"

		if err := ioutil.WriteFile(fn, out[k].Bytes(), 0644); err != nil {
			t.Error(err)
		}
	}
}

func writeProps(station *delta.Station, site *delta.Site, sensor *delta.InstalledSensor, b *bytes.Buffer) {

	if !bytes.HasSuffix(b.Bytes(), []byte("[")) {
		b.WriteString(",")
	}

	b.WriteString(`{"type":"Feature","geometry":{"type": "Point","coordinates": [`)
	b.WriteString(fmt.Sprintf("%f,%f", site.Point.Longitude, site.Point.Latitude))
	b.WriteString(`]},"properties":{`)
	b.WriteString(fmt.Sprintf("\"network\":\"%s\"", station.Network.External))
	b.WriteString(fmt.Sprintf(",\"station\":\"%s\"", station.Code))
	b.WriteString(fmt.Sprintf(",\"location\":\"%s\"", site.Location))
	b.WriteString(fmt.Sprintf(",\"name\":\"%s\"", station.Name))
	b.WriteString(fmt.Sprintf(",\"datum\":\"%s\"", site.Point.Datum))
	b.WriteString(fmt.Sprintf(",\"elevation\":%f", site.Point.Elevation))
	b.WriteString(fmt.Sprintf(",\"vertical_offset\":%f", sensor.Offset.Vertical))
	b.WriteString(fmt.Sprintf(",\"start\":\"%d\"", sensor.Span.Start))
	b.WriteString(fmt.Sprintf(",\"end\":\"%d\"", sensor.Span.End))

	b.WriteString(`}}`)
}

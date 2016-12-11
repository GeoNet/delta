package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
	"io/ioutil"
	"os"
	"github.com/GeoNet/delta"
	"github.com/golang/protobuf/proto"
	"encoding/json"
	"fmt"
	"bytes"
	"time"
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

// TestMarksProto creates protobuf and JSON files of Marks.
// These are pushed to S3 (by Travis) for use in api.geonet.org.nz
// Three files are created:
// marks.pb - fully hydrated protobuf will all GNSS Mark information.
// marks.json - JSON version of marks.pb (for use in browsers).
// marks.geojson - GeoJSON of Mark locations.
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

	 var monuments meta.MonumentList
	 if err := meta.LoadList("../network/monuments.csv", &monuments); err != nil {
		 t.Error(err)
	 }

	 for _,v := range monuments {
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }

		 mn := delta.Monument{
			 DomesNumber:v.DomesNumber,
			 MarkType:v.MarkType,
			 Type     :v.Type,
			 GroundRelationship:v.GroundRelationship,
			 FoundationType:v.FoundationType,
			 FoundationDepth:v.FoundationDepth,
			 Bedrock: v.Bedrock,
			 Geology: v.Geology,
			 Span: &s,
		 }
		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].Monument = &mn
		 }
	 }

	 var antennas meta.InstalledAntennaList
	 if err := meta.LoadList("../install/antennas.csv", &antennas); err != nil {
		 t.Error(err)
	 }

	 for _, v := range antennas {
		 e := delta.Equipment{
			 Make: v.Make,
			 Model: v.Model,
			 Serial: v.Serial,
		 }
		 o := delta.Offset{
			 Vertical:v.Vertical,
			 North:v.North,
			 East:v.East,
		 }
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }
		 a := delta.InstalledAntenna{
			 Equipment:&e,
			 Offset:&o,
			 Span:&s,
			 Azimuth:v.Azimuth,
		 }

		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].InstalledAntenna = append(m.Marks[v.Mark].InstalledAntenna, &a)
		 }
	 }

	 var radomes meta.InstalledRadomeList
	 if err := meta.LoadList("../install/radomes.csv", &radomes); err != nil {
		 t.Error(err)
	 }

	 for _, v := range radomes {
		 e := delta.Equipment{
			 Make: v.Make,
			 Model: v.Model,
			 Serial: v.Serial,
		 }
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }
		 r := delta.InstalledRadome{
			 Equipment: &e,
			 Span: &s,
		 }
		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].InstalledRadome = append(m.Marks[v.Mark].InstalledRadome, &r)
		 }
	 }

	 var firmwares meta.FirmwareHistoryList
	 if err := meta.LoadList("../install/firmware.csv", &firmwares); err != nil {
		 t.Error(err)
	 }

	 var fw = make(map[delta.Equipment]delta.Receiver)

	 for _, v := range firmwares {
		 e := delta.Equipment{
			 Make: v.Make,
			 Model: v.Model,
			 Serial: v.Serial,
		 }
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }
		 f := delta.Firmware{
			 Version: v.Version,
			 Notes: v.Notes,
			 Span: &s,
		 }

		 rx := fw[e]
		 rx.Equipment =&e
		 rx.Firmware = append(rx.Firmware, &f)
		 fw[e] = rx
	 }

	 var receivers meta.DeployedReceiverList
	 if err := meta.LoadList("../install/receivers.csv", &receivers); err != nil {
		 t.Error(err)
	 }

	 for _,v := range receivers {
		 e := delta.Equipment{
			 Make: v.Make,
			 Model: v.Model,
			 Serial: v.Serial,
		 }
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }

		 rx, ok := fw[e]
		 if !ok {
			 t.Errorf("no firware for %v", e)
		 }

		 d := delta.DeployedReceiver{
			 Receiver: &rx,
			 Span: &s,
		 }

		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].DeployedReceiver = append(m.Marks[v.Mark].DeployedReceiver, &d)
		 }
	 }

	 var sessions meta.SessionList
	 if err := meta.LoadList("../install/sessions.csv", &sessions); err != nil {
		 t.Error(err)
	 }

	 for _, v := range sessions {
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }
		 se := delta.Session{
			 Operator: v.Operator,
			 Agency          :v.Agency,
			 Model           : v.Model,
			 SatelliteSystem : v.SatelliteSystem,
			 Interval        :v.Interval.Nanoseconds(),
			 ElevationMask: v.ElevationMask,
			 HeaderComment: v.HeaderComment,
			 Span: &s,
		 }
		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].Session = append(m.Marks[v.Mark].Session, &se)
		 }
	 }

	 var metsensors meta.InstalledMetSensorList
	 if err := meta.LoadList("../install/metsensors.csv", &metsensors); err != nil {
		 t.Error(err)
	 }

	 for _, v := range metsensors {
		 e := delta.Equipment{
			 Make: v.Make,
			 Model: v.Model,
			 Serial: v.Serial,
		 }
		 s := delta.Span{
			 Start:v.Start.Unix(),
			 End:v.End.Unix(),
		 }
		 p := delta.Point{
			 Longitude:v.Longitude,
			 Latitude:v.Latitude,
			 Elevation:v.Elevation,
			 Datum:v.Datum,
		 }
		 ms := delta.InstalledMetSensor{
			 Equipment: &e,
			 Span: &s,
			 Point: &p,
			 IMSComment:v.IMSComment,
		 }
		 if _, ok := m.Marks[v.Mark]; ok {
			 m.Marks[v.Mark].InstalledMetSensor = append(m.Marks[v.Mark].InstalledMetSensor, &ms)
		 }
	 }

	 // output files

	 if err := os.MkdirAll(apiDir, 0777); err != nil {
		 t.Error(err)
	 }

	 b, err := proto.Marshal(&m)
	 if err != nil {
		 t.Error(err)
	 }

	 if err := ioutil.WriteFile(apiDir + "/marks.pb", b, 0644); err != nil {
		 t.Error(err)
	 }

	 b, err = json.Marshal(&m)
	 if err != nil {
		 t.Error(err)
	 }

	 if err := ioutil.WriteFile(apiDir + "/marks.json", b, 0644); err != nil {
		 t.Error(err)
	 }

	 var bf bytes.Buffer

	 bf.WriteString(`{"type": "FeatureCollection","features": [`)

	 f := false
	 for _, v := range m.Marks {
		 if v.GetSpan() == nil || v.GetPoint() == nil || v.GetNetwork() == nil {
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
		 bf.WriteString(fmt.Sprintf(",\"start\":\"%s\"", time.Unix(v.Span.Start, 0).UTC().Format(time.RFC3339)))
		 bf.WriteString(fmt.Sprintf(",\"end\":\"%s\"", time.Unix(v.Span.End, 0).UTC().Format(time.RFC3339)))

		 var group string

		 switch v.Network.Code {
		 case "LI", "IG":
			 group = "LINZ"
		 case "CG", "SA":
			 group = "GeoNet"
		 default:
			 group = "other"
		 }

		 bf.WriteString(fmt.Sprintf(",\"group\":\"%s\"", group))

		 bf.WriteString(`}}`)
	 }

	 bf.WriteString(`]}`)

	 if err := ioutil.WriteFile(apiDir + "/marks.geojson", bf.Bytes(), 0644); err != nil {
		 t.Error(err)
	 }
 }

package gloria_pb

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"testing"
	"time"
)

// TestProto creates an example Gloria proto file for TAUP.
func TestProto(t *testing.T) {
	m := Mark{
		Code:        "TAUP",
		DomesNumber: "50217M001",
		Point: &Point{
			Longitude: 176.0809947,
			Latitude:  -38.74271665,
			Elevation: 427.0279,
		},
		DeployedReceiver: []*DeployedReceiver{
			&DeployedReceiver{
				Receiver: &Receiver{
					Model:        "TRIMBLE NETR9",
					SerialNumber: "5033K69574",
					Firmware: []*Firmware{
						&Firmware{Version: "5.15",
							Span: &Span{
								Start: time.Date(2016, 12, 20, 22, 21, 26, 0, time.UTC).Unix(),
								End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
				},
				Span: &Span{
					Start: time.Date(2015, 9, 19, 17, 25, 1, 0, time.UTC).Unix(),
					End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		InstalledAntenna: []*InstalledAntenna{
			&InstalledAntenna{
				Antenna: &Antenna{
					Model:        "TRM57971.00",
					SerialNumber: "1441031450",
				},
				Offset: &Offset{
					East:     0.0,
					North:    0.0,
					Vertical: 0.0550,
				},
				Span: &Span{
					Start: time.Date(2015, 9, 19, 17, 25, 1, 0, time.UTC).Unix(),
					End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		Comment: `Data supplied by the GeoNet project.  GeoNet is core
funded by EQC, with support from LINZ, and is
operated by GNS on behalf of EQC and all New Zealanders.
Contact: www.geonet.org.nz  email: info@geonet.org.nz`,
	}

	b, err := proto.Marshal(&m)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("testdata/taup.pb", b, 0644)
	if err != nil {
		t.Error(err)
	}
}

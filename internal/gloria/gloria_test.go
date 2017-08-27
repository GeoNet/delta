package gloria_test

import (
	"github.com/GeoNet/delta/internal/gloria"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"testing"
	"time"
)

// TestProto creates an example Gloria proto file for TAUP.
func TestProto(t *testing.T) {
	m := gloria.Mark{
		Code:        "TAUP",
		DomesNumber: "50217M001",
		Point: &gloria.Point{
			Longitude: 176.0809947,
			Latitude:  -38.74271665,
			Elevation: 427.0279,
		},
		DeployedReceiver: []*gloria.DeployedReceiver{
			&gloria.DeployedReceiver{
				Receiver: &gloria.Receiver{
					Model:        "TRIMBLE NETR9",
					SerialNumber: "5033K69574",
					Firmware: []*gloria.Firmware{
						&gloria.Firmware{Version: "5.15",
							Span: &gloria.Span{
								Start: time.Date(2016, 12, 20, 22, 21, 26, 0, time.UTC).Unix(),
								End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
							},
						},
					},
				},
				Span: &gloria.Span{
					Start: time.Date(2015, 9, 19, 17, 25, 1, 0, time.UTC).Unix(),
					End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
		InstalledAntenna: []*gloria.InstalledAntenna{
			&gloria.InstalledAntenna{
				Antenna: &gloria.Antenna{
					Model:        "TRM57971.00",
					SerialNumber: "1441031450",
				},
				Offset: &gloria.Offset{
					East:     0.0,
					North:    0.0,
					Vertical: 0.0550,
				},
				Span: &gloria.Span{
					Start: time.Date(2015, 9, 19, 17, 25, 1, 0, time.UTC).Unix(),
					End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
		},
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

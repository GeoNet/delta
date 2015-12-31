package meta

import (
	"io/ioutil"
	//	"path/filepath"
	"testing"
	"time"
)

func TestList(t *testing.T) {

	//	var liststart, _ = time.Parse(DateTimeFormat, "2010-01-01T12:00:00Z")
	//	var listend, _ = time.Parse(DateTimeFormat, "2012-01-01T12:00:00Z")

	var listtests = []struct {
		f string
		l List
	}{
		/*
			{
				"testdata/networks.csv",
				Networks{
					Network{
						Code:        "AA",
						External:    "XX",
						Description: "A Name",
						Restricted:  false,
					},
					Network{
						Code:        "BB",
						External:    "XX",
						Description: "B Name",
						Restricted:  true,
					},
				},
				func() List { return &Networks{} },
			},
			{
				"testdata/stations.csv",
				Stations{
					Station{
						Code:      "AAAA",
						Network:   "AA",
						Name:      "A Name",
						Latitude:  -41.5,
						Longitude: 173.5,
						StartTime: liststart,
						EndTime:   listend,
					},
					Station{
						Code:      "BBBB",
						Network:   "BB",
						Name:      "B Name",
						Latitude:  -42.5,
						Longitude: 174.5,
						StartTime: liststart.Add(time.Hour),
						EndTime:   listend.Add(time.Hour),
					},
				},
				func() List { return &Stations{} },
			},
		*/
		{
			"testdata/assets.csv",
			&Assets{
				{
					Equipment: Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220063995",
					},
					Manufacturer: "Trimble Navigation Ltd.",
				},
				{
					Equipment: Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220066912",
					},
					Manufacturer: "Trimble Navigation Ltd.",
				},
			},
		},
		{
			"testdata/antennas.csv",
			&InstalledAntennas{
				{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220063995",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2000-08-02T23:59:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Offset: Offset{
						Height: 0.0015,
						North:  0.0,
						East:   0.0,
					},
					MarkCode: "CNCL",
				},
				{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220066912",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2000-08-14T23:59:52Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2011-02-07T22:35:00Z")
								return v
							}(),
						},
					},
					Offset: Offset{
						Height: 0.0013,
						North:  0.0,
						East:   0.0,
					},
					MarkCode: "MTJO",
				},
			},
		},
		{
			"testdata/dataloggers.csv",
			&DeployedDataloggers{
				DeployedDatalogger{
					Install: Install{
						Equipment: Equipment{
							Make:   "GNSScience",
							Model:  "EARSS/3",
							Serial: "152",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2001-01-18T13:22:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2001-02-10T10:50:00Z")
								return v
							}(),
						},
					},
					Place: "Pukeroa",
					Role:  "Short Period",
				},
				DeployedDatalogger{
					Install: Install{
						Equipment: Equipment{
							Make:   "Kinemetrics",
							Model:  "Q330/3",
							Serial: "2216",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2009-02-10T23:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Place: "Turoa Road End",
				},
			},
		},
		{
			"testdata/gauges.csv",
			&InstalledGauges{
				InstalledGauge{
					Install: Install{
						Equipment: Equipment{
							Make:   "GESensing",
							Model:  "Druck PDCR-1830",
							Serial: "2427881",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2007-03-06T00:00:02Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2007-05-22T23:30:00Z")
								return v
							}(),
						},
					},
					Offset: Offset{},
					Orientation: Orientation{
						Dip: 90.0,
					},
					CableLength:  20.0,
					StationCode:  "WLGT",
					LocationCode: "41",
				},
				InstalledGauge{
					Install: Install{
						Equipment: Equipment{
							Make:   "GESensing",
							Model:  "Druck PTX-1830",
							Serial: "2504328",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2007-03-06T00:00:02Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2012-05-29T22:55:00Z")
								return v
							}(),
						},
					},
					Offset: Offset{},
					Orientation: Orientation{
						Dip: 90.0,
					},
					CableLength:  20.0,
					StationCode:  "WLGT",
					LocationCode: "40",
				},
			},
		},
		{
			"testdata/metsensors.csv",
			&InstalledMetSensors{
				InstalledMetSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65123",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "1998-07-09T23:59:59Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Point: Point{
						Latitude:  -41.2351,
						Longitude: 174.917,
						Elevation: 26,
						Datum:     "NZGD2000",
					},
					MarkCode: "GRAC",
				},
				InstalledMetSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65125",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2000-08-15T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Point: Point{
						Latitude:  -43.9857,
						Longitude: 170.4649,
						Elevation: 1044,
						Datum:     "NZGD2000",
					},
					MarkCode: "MTJO",
				},
			},
		},
		{
			"testdata/radomes.csv",
			&InstalledRadomes{
				InstalledRadome{
					Install: Install{
						Equipment: Equipment{
							Make:   "LeicaGeosystems",
							Model:  "LEIS Radome",
							Serial: "0220148020",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "1999-09-27T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2000-01-21T00:00:00Z")
								return v
							}(),
						},
					},
					MarkCode: "MQZG",
				},
				InstalledRadome{
					Install: Install{
						Equipment: Equipment{
							Make:   "Thales",
							Model:  "SCIS Radome",
							Serial: "0220063995",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2000-08-03T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					MarkCode: "CNCL",
				},
			},
		},
		{
			"testdata/receivers.csv",
			&DeployedReceivers{
				DeployedReceiver{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "5700",
							Serial: "220280300",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2002-12-31T01:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2012-08-31T15:00:01Z")
								return v
							}(),
						},
					},
					Place: "Mount Hodgkinson",
				},
				DeployedReceiver{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "NetR9",
							Serial: "5014K66721",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2010-10-12T00:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2014-07-27T22:00:00Z")
								return v
							}(),
						},
					},
					Place: "Methven",
				},
			},
		},
		{
			"testdata/sensors.csv",
			&InstalledSensors{
				InstalledSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "AppliedGeomechanics",
							Model:  "Lily tiltmeter",
							Serial: "N7935",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2009-11-15T01:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2013-05-24T00:00:00Z")
								return v
							}(),
						},
					},
					Orientation: Orientation{
						Azimuth: 233.0,
						Dip:     0.0,
					},
					Offset: Offset{
						Height: -64.0,
					},
					StationCode:  "COVZ",
					LocationCode: "90",
				},
				InstalledSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Guralp",
							Model:  "CMG-3ESPC",
							Serial: "T36194",
						},
						Span: Span{
							Start: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2012-05-21T11:00:04Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(DateTimeFormat, "2013-07-10T23:00:00Z")
								return v
							}(),
						},
					},
					Orientation: Orientation{
						Azimuth: 0.0,
						Dip:     0.0,
					},
					Offset: Offset{
						Height: 0.0,
					},
					StationCode:  "INZ",
					LocationCode: "10",
				},
			},
		},
	}

	for _, tt := range listtests {
		res := MarshalList(tt.l)

		t.Log("Compare raw list file: " + tt.f)
		{
			b, err := ioutil.ReadFile(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(b) {
				t.Errorf("list file text mismatch: %s [\n%s\n]", tt.f, diff(string(res), string(b)))
			}
		}
		t.Log("Check encode/decode list: " + tt.f)
		{
			if err := UnmarshalList(res, tt.l); err != nil {
				t.Fatal(err)
			}

			s := MarshalList(tt.l)
			if string(res) != string(s) {
				t.Errorf("list encode/reencode mismatch: %s [\n%s\n]", tt.f, diff(string(res), string(s)))
			}
		}

		t.Log("Check list file: " + tt.f)
		{
			if err := LoadList(tt.f, tt.l); err != nil {
				t.Fatal(err)
			}

			s := MarshalList(tt.l)
			if string(res) != string(s) {
				t.Errorf("list file list mismatch: %s [\n%s\n]", tt.f, diff(string(res), string(s)))
			}
		}
	}
}

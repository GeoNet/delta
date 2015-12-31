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

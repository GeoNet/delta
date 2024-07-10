package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mustParse(str string) time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSensor(t *testing.T) {

	example := Network{
		Stations: []Station{
			{
				Code:        "ABAZ",
				Name:        "Army Bay",
				Network:     "NZ",
				Description: "Auckland volcano seismic network",
				StartDate:   mustParse("2008-10-13T00:00:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),
				Latitude:    -36.600224003,
				Longitude:   174.832332909,
				Elevation:   74,

				Sites: []Site{
					{
						Code:      "10",
						Latitude:  -36.600224003,
						Longitude: 174.832332909,
						Elevation: 74,
						Datum:     "WGS84",
						Survey:    "External GPS Device",
						StartDate: mustParse("2008-10-13T04:00:00Z"),
						EndDate:   mustParse("2023-07-28T01:08:24Z"),

						Sensors: []Sensor{
							{
								Model:     "L4C-3D",
								Make:      "Sercel",
								Type:      "Short Period Seismometer",
								Channels:  "EHZ",
								StartDate: mustParse("2008-10-13T04:00:00Z"),
								EndDate:   mustParse("2010-03-15T02:00:00Z"),
							},
						},
					},
				},
			},
			{
				Code:        "AUCT",
				Name:        "Auckland",
				Network:     "NZ",
				Description: "National tsunami gauge network",
				Latitude:    -36.8314371,
				Longitude:   174.7865372,
				StartDate:   mustParse("2009-03-26T00:00:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),

				Sites: []Site{
					{
						Code:      "40",
						Latitude:  -36.8314371,
						Longitude: 174.7865372,
						Datum:     "WGS84",
						Survey:    "Unknown",
						StartDate: mustParse("2009-03-26T00:00:00Z"),
						EndDate:   mustParse("9999-01-01T00:00:00Z"),

						Sensors: []Sensor{
							{
								Model:     "Druck PTX-1830",
								Make:      "General Electric",
								Type:      "Coastal Pressure Sensor",
								Channels:  "BTT",
								StartDate: mustParse("2009-03-26T02:30:00Z"),
								EndDate:   mustParse("2012-08-28T01:00:00Z"),
							},
						},
					},
				},
			},
			{
				Code:        "ALS1R",
				Name:        "Auckland Scenic Drive Waiatarua",
				Network:     "NZ",
				Description: "GeoNet environmental data using low rate data collection platform network",
				StartDate:   mustParse("2023-02-24T01:09:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),
				Latitude:    -36.93285,
				Longitude:   174.576436,
				Elevation:   370.17,

				Sites: []Site{
					{
						Code:      "01",
						Latitude:  -36.93285,
						Longitude: 174.576436,
						Elevation: 370.17,
						Datum:     "WGS84",
						Survey:    "External GPS Device",
						StartDate: mustParse("2023-02-24T01:09:00Z"),
						EndDate:   mustParse("9999-01-01T00:00:00Z"),

						Sensors: []Sensor{
							{
								Model:     "Tipping bucket rain gauge TB3 0.2mm",
								Make:      "Hyquest",
								Type:      "Environmental Sensor",
								StartDate: mustParse("2023-02-24T01:09:00Z"),
								EndDate:   mustParse("9999-01-01T00:00:00Z"),
							},
						},
					},
				},
			},
		},
		Marks: []Mark{
			{
				Code:           "2004",
				Name:           "Puketapu Road",
				Network:        "NZ",
				Description:    "Temporary sites",
				Latitude:       -38.672344067,
				Longitude:      175.804587325,
				Elevation:      514.959,
				MarkType:       "Forced Centering",
				FoundationType: "Reinforced Concrete",
				StartDate:      mustParse("2000-06-19T00:00:00Z"),
				EndDate:        mustParse("9999-01-01T00:00:00Z"),

				Antennas: []Sensor{
					{
						Model:     "TRM57971.00",
						Make:      "Trimble Navigation Ltd.",
						Type:      "GNSS Antenna",
						StartDate: mustParse("2022-12-06T01:37:08Z"),
						EndDate:   mustParse("9999-01-01T00:00:00Z"),
					},
				},
				Receivers: []Sensor{
					{
						Model:     "TRIMBLE ALLOY",
						Make:      "Trimble Navigation Ltd.",
						Type:      "GNSS Receiver",
						StartDate: mustParse("2022-12-06T01:37:08Z"),
						EndDate:   mustParse("9999-01-01T00:00:00Z"),
					},
				},
			},
		},
		Buoys: []Station{
			{
				Code:        "NZA",
				Name:        "Offshore Wellington Hikurangi",
				Network:     "NZ",
				Description: "National tsunami deep ocean network",
				StartDate:   mustParse("2019-12-21T17:00:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),
				Latitude:    -42.371,
				Longitude:   176.911,
				Depth:       2690,
				Datum:       "WGS84",

				Sites: []Site{
					{
						Code:      "40",
						Latitude:  -42.3707,
						Longitude: 176.9109,
						Depth:     2690,
						Datum:     "WGS84",
						Survey:    "Unknown",
						StartDate: mustParse("2019-12-21T17:00:00Z"),
						EndDate:   mustParse("2021-12-17T01:00:00Z"),

						Sensors: []Sensor{
							{
								Model:     "BPR Subsystem",
								Make:      "SAIC",
								Type:      "DART Bottom Pressure Recorder",
								StartDate: mustParse("2019-12-22T00:00:00Z"),
								EndDate:   mustParse("2021-12-17T23:59:59Z"),
							},
						},
					},
				},
			},
		},
		Mounts: []Mount{
			{
				Code:        "DISC",
				Name:        "Discovery Lodge",
				Network:     "VC",
				Description: "Volcano monitoring camera network",
				Latitude:    -39.156143224,
				Longitude:   175.491281443,
				Elevation:   9999,
				Datum:       "WGS84",
				StartDate:   mustParse("2020-07-28T02:20:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),

				Views: []View{
					{
						Code:        "01",
						Label:       "Ruapehu North",
						Description: "Images of Mount Ruapehu from the volcano camera situated at Discovery Lodge.",
						Azimuth:     155,
						StartDate:   mustParse("2020-07-28T02:20:00Z"),
						EndDate:     mustParse("9999-01-01T00:00:00Z"),

						Sensors: []Sensor{
							{
								Model:     "M12 3MP",
								Make:      "Mobotix AG",
								Type:      "Camera",
								Azimuth:   144,
								Dip:       6,
								StartDate: mustParse("2020-07-28T02:20:00Z"),
								EndDate:   mustParse("2021-02-04T21:10:00Z"),
							},
						},
					},
				},
			},
			{

				Code:        "RUD01",
				Name:        "Tukino Skifield",
				Network:     "NZ",
				Description: "GeoNet environmental data using low rate data collection platform network",
				Latitude:    -39.2775,
				Longitude:   175.6087,
				Elevation:   1748,
				Datum:       "WGS84",
				StartDate:   mustParse("2021-05-05T05:13:00Z"),
				EndDate:     mustParse("9999-01-01T00:00:00Z"),

				Views: []View{
					{
						Code:        "01",
						Label:       "Tukino Skifield",
						Description: "SO2 from Ruapehu from Tukino Skifield",
						Azimuth:     266,
						Dip:         -60,
						StartDate:   mustParse("2022-05-05T05:13:00Z"),
						EndDate:     mustParse("9999-01-01T00:00:00Z"),

						Sensors: []Sensor{
							{
								Model:     "Avaspec-Mini2048CL",
								Make:      "Avaspec CompactLine",
								Type:      "DOAS",
								Azimuth:   266,
								Dip:       -60,
								StartDate: mustParse("2022-05-05T05:13:00Z"),
								EndDate:   mustParse("9999-01-01T00:00:00Z"),
							},
						},
					},
				},
			},
		},
		Samples: []Station{
			{
				Code:        "WI228",
				Name:        "White Island Yellow Duck Stream",
				Network:     "MC",
				Description: "Manually collected volcano monitoring data",
				StartDate:   mustParse("2017-06-29T00:00:00Z"),
				EndDate:     mustParse("2017-06-29T00:00:00Z"),
				Latitude:    -37.520946,
				Longitude:   177.186671,
				Elevation:   24,
				Datum:       "WGS84",

				Sites: []Site{
					{
						Code:      "MC01",
						Latitude:  -37.520946,
						Longitude: 177.186671,
						Elevation: 24,
						Datum:     "WGS84",
						Survey:    "External GPS Device",
						StartDate: mustParse("2017-06-29T00:00:00Z"),
						EndDate:   mustParse("2017-06-29T00:00:00Z"),

						Sensors: []Sensor{
							{
								Code:      "01",
								Property:  "Al-conc",
								Type:      "Manual Collection",
								Aspect:    "stream",
								StartDate: mustParse("2017-06-29T00:00:00Z"),
								EndDate:   mustParse("2017-06-29T00:00:00Z"),
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer

	if err := example.EncodeXML(&buf, "", "  "); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile("testdata/sensor.xml")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), raw) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(buf.String(), string(raw)))
	}
}

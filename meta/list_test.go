package meta_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

func TestList(t *testing.T) {

	var listtests = []struct {
		f string
		l meta.List
	}{
		{
			"testdata/networks.csv",
			&meta.NetworkList{
				meta.Network{
					NetworkCode:  "AK",
					ExternalCode: "NZ",
					Description:  "Auckland volcano seismic network",
					Restricted:   false,
				},
				meta.Network{
					NetworkCode:  "CB",
					ExternalCode: "NZ",
					Description:  "Canterbury regional seismic network",
					Restricted:   false,
				},
			},
		},
		{
			"testdata/stations.csv",
			&meta.StationList{
				meta.Station{
					Reference: meta.Reference{
						Code:    "DFE",
						Network: "TR",
						Name:    "Dawson Falls",
					},
					Point: meta.Point{
						Latitude:  -39.325743417,
						Longitude: 174.103863732,
						Elevation: 880.0,
						Datum:     "WGS84",
					},
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "1993-12-14T00:00:00Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2010-02-23T00:00:00Z")
							return v
						}(),
					},
				},
				meta.Station{
					Reference: meta.Reference{
						Code:    "TBAS",
						Network: "SM",
						Name:    "Tolaga Bay Area School",
					},
					Point: meta.Point{
						Latitude:  -38.372803703,
						Longitude: 178.300778623,
						Elevation: 8.0,
						Datum:     "WGS84",
					},
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2002-03-05T00:00:00Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
							return v
						}(),
					},
					Notes: "Is located in the Kiln Shed next to the hall",
				},
			},
		},
		{
			"testdata/sites.csv",
			&meta.SiteList{
				meta.Site{
					Point: meta.Point{
						Latitude:  -39.198244208,
						Longitude: 175.547981982,
						Elevation: 1116.0,
						Datum:     "WGS84",
					},
					StationCode:  "CNZ",
					LocationCode: "12",
				},
				meta.Site{
					Point: meta.Point{
						Latitude:  -45.091369824,
						Longitude: 169.411775594,
						Elevation: 701.0,
						Datum:     "WGS84",
					},
					StationCode:  "MSCZ",
					LocationCode: "10",
				},
			},
		},
		{
			"testdata/marks.csv",
			&meta.MarkList{
				meta.Mark{
					Reference: meta.Reference{
						Code:    "AHTI",
						Network: "CG",
						Name:    "Ahititi",
					},
					Point: meta.Point{
						Latitude:  -38.411447554,
						Longitude: 178.046002897,
						Elevation: 563.221,
						Datum:     "NZGD2000",
					},
					MarkType:           "Forced Centering",
					MonumentType:       "Deep Braced",
					GroundRelationship: 0.0,
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2009-01-01T00:00:00Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
							return v
						}(),
					},
				},
				meta.Mark{
					Reference: meta.Reference{
						Code:    "DUND",
						Network: "LI",
						Name:    "Dunedin",
					},
					Point: meta.Point{
						Latitude:  -45.88366604,
						Longitude: 170.5971706,
						Elevation: 386.964,
						Datum:     "NZGD2000",
					},
					GroundRelationship: -1.25,
					MarkType:           "Forced Centering",
					MonumentType:       "Short Braced",
					DomeNumber:         "50212M003",
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2005-08-10T00:00:00Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
							return v
						}(),
					},
				},
			},
		},
		{
			"testdata/assets.csv",
			&meta.AssetList{
				{
					Equipment: meta.Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220063995",
					},
					Manufacturer: "Trimble Navigation Ltd.",
					AssetNumber:  "100",
				},
				{
					Equipment: meta.Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220066912",
					},
					Manufacturer: "Trimble Navigation Ltd.",
					AssetNumber:  "101",
				},
			},
		},
		{
			"testdata/antennas.csv",
			&meta.InstalledAntennaList{
				{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220063995",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2000-08-02T23:59:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Offset: meta.Offset{
						Height: 0.0015,
						North:  0.0,
						East:   0.0,
					},
					MarkCode: "CNCL",
				},
				{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220066912",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2000-08-14T23:59:52Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2011-02-07T22:35:00Z")
								return v
							}(),
						},
					},
					Offset: meta.Offset{
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
			&meta.DeployedDataloggerList{
				meta.DeployedDatalogger{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "GNSScience",
							Model:  "EARSS/3",
							Serial: "152",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2001-01-18T13:22:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2001-02-10T10:50:00Z")
								return v
							}(),
						},
					},
					Place: "Pukeroa",
					Role:  "Short Period",
				},
				meta.DeployedDatalogger{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Kinemetrics",
							Model:  "Q330/3",
							Serial: "2216",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2009-02-10T23:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
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
			&meta.InstalledGaugeList{
				meta.InstalledGauge{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "GESensing",
							Model:  "Druck PDCR-1830",
							Serial: "2427881",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2007-03-06T00:00:02Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2007-05-22T23:30:00Z")
								return v
							}(),
						},
					},
					Offset: meta.Offset{},
					Orientation: meta.Orientation{
						Dip: 90.0,
					},
					CableLength:  20.0,
					StationCode:  "WLGT",
					LocationCode: "41",
				},
				meta.InstalledGauge{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "GESensing",
							Model:  "Druck PTX-1830",
							Serial: "2504328",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2007-03-06T00:00:02Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2012-05-29T22:55:00Z")
								return v
							}(),
						},
					},
					Offset: meta.Offset{},
					Orientation: meta.Orientation{
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
			&meta.InstalledMetSensorList{
				meta.InstalledMetSensor{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65123",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "1998-07-09T23:59:59Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Point: meta.Point{
						Latitude:  -41.2351,
						Longitude: 174.917,
						Elevation: 26,
						Datum:     "NZGD2000",
					},
					MarkCode: "GRAC",
				},
				meta.InstalledMetSensor{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65125",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2000-08-15T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
								return v
							}(),
						},
					},
					Point: meta.Point{
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
			&meta.InstalledRadomeList{
				meta.InstalledRadome{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "LeicaGeosystems",
							Model:  "LEIS Radome",
							Serial: "0220148020",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "1999-09-27T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2000-01-21T00:00:00Z")
								return v
							}(),
						},
					},
					MarkCode: "MQZG",
				},
				meta.InstalledRadome{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Thales",
							Model:  "SCIS Radome",
							Serial: "0220063995",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2000-08-03T00:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
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
			&meta.DeployedReceiverList{
				meta.DeployedReceiver{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Trimble",
							Model:  "5700",
							Serial: "220280300",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2002-12-31T01:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2012-08-31T15:00:01Z")
								return v
							}(),
						},
					},
					Mark: "HORN",
				},
				meta.DeployedReceiver{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Trimble",
							Model:  "NetR9",
							Serial: "5014K66721",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2010-10-12T00:00:01Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2014-07-27T22:00:00Z")
								return v
							}(),
						},
					},
					Mark: "METH",
				},
			},
		},
		{
			"testdata/sensors.csv",
			&meta.InstalledSensorList{
				meta.InstalledSensor{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "AppliedGeomechanics",
							Model:  "Lily tiltmeter",
							Serial: "N7935",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2009-11-15T01:00:00Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2013-05-24T00:00:00Z")
								return v
							}(),
						},
					},
					Orientation: meta.Orientation{
						Azimuth: 233.0,
						Dip:     0.0,
					},
					Offset: meta.Offset{
						Height: -64.0,
					},
					StationCode:  "COVZ",
					LocationCode: "90",
				},
				meta.InstalledSensor{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   "Guralp",
							Model:  "CMG-3ESPC",
							Serial: "T36194",
						},
						Span: meta.Span{
							Start: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2012-05-21T11:00:04Z")
								return v
							}(),
							End: func() time.Time {
								v, _ := time.Parse(meta.DateTimeFormat, "2013-07-10T23:00:00Z")
								return v
							}(),
						},
					},
					Orientation: meta.Orientation{
						Azimuth: 0.0,
						Dip:     0.0,
					},
					Offset: meta.Offset{
						Height: 0.0,
					},
					StationCode:  "INZ",
					LocationCode: "10",
				},
			},
		},
		{
			"testdata/recorders.csv",
			&meta.InstalledRecorderList{
				meta.InstalledRecorder{
					InstalledSensor: meta.InstalledSensor{
						Install: meta.Install{
							Equipment: meta.Equipment{
								Make:   "CSI",
								Model:  "CUSP3A",
								Serial: "3A-040001",
							},
							Span: meta.Span{
								Start: func() time.Time {
									v, _ := time.Parse(meta.DateTimeFormat, "2004-11-27T00:00:00Z")
									return v
								}(),
								End: func() time.Time {
									v, _ := time.Parse(meta.DateTimeFormat, "2010-03-25T00:30:00Z")
									return v
								}(),
							},
						},
						Orientation: meta.Orientation{
							Azimuth: 266.0,
							Dip:     0.0,
						},
						Offset: meta.Offset{
							Height: 0.0,
						},
						StationCode:  "AMBC",
						LocationCode: "20",
					},
					DataloggerModel: "CUSP3A",
				},
				meta.InstalledRecorder{
					InstalledSensor: meta.InstalledSensor{
						Install: meta.Install{
							Equipment: meta.Equipment{
								Make:   "Kinemetrics",
								Model:  "FBA-ES-T-DECK",
								Serial: "1275",
							},
							Span: meta.Span{
								Start: func() time.Time {
									v, _ := time.Parse(meta.DateTimeFormat, "2014-04-17T00:10:00Z")
									return v
								}(),
								End: func() time.Time {
									v, _ := time.Parse(meta.DateTimeFormat, "2014-07-29T00:00:00Z")
									return v
								}(),
							},
						},
						Orientation: meta.Orientation{
							Azimuth: 210.0,
							Dip:     0.0,
						},
						Offset: meta.Offset{
							Height: 0.0,
						},
						StationCode:  "EKS3",
						LocationCode: "20",
					},
					DataloggerModel: "BASALT",
				},
			},
		},
		{
			"testdata/connections.csv",
			&meta.ConnectionList{
				meta.Connection{
					StationCode:  "APZ",
					LocationCode: "10",
					Place:        "The Paps",
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2006-05-07T03:23:54Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
							return v
						}(),
					},
				},
				meta.Connection{
					StationCode:  "BSWZ",
					LocationCode: "10",
					Place:        "Blackbirch Station",
					Span: meta.Span{
						Start: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "2003-12-09T00:00:00Z")
							return v
						}(),
						End: func() time.Time {
							v, _ := time.Parse(meta.DateTimeFormat, "9999-01-01T00:00:00Z")
							return v
						}(),
					},
				},
			},
		},
	}

	for _, tt := range listtests {
		res := meta.MarshalList(tt.l)

		t.Log("Compare raw list file: " + tt.f)
		{
			b, err := ioutil.ReadFile(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(b) {
				t.Errorf("list file text mismatch: %s [\n%s\n]", tt.f, diff(res, b))
			}
		}
		t.Log("Check encode/decode list: " + tt.f)
		{
			if err := meta.UnmarshalList(res, tt.l); err != nil {
				t.Fatal(err)
			}

			s := meta.MarshalList(tt.l)
			if string(res) != string(s) {
				t.Errorf("list encode/reencode mismatch: %s [\n%s\n]", tt.f, diff(res, s))
			}
		}

		t.Log("Check list file: " + tt.f)
		{
			if err := meta.LoadList(tt.f, tt.l); err != nil {
				t.Fatal(err)
			}

			s := meta.MarshalList(tt.l)
			if string(res) != string(s) {
				t.Errorf("list file list mismatch: %s [\n%s\n]", tt.f, diff(res, s))
			}
		}
	}
}

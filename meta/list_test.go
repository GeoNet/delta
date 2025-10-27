package meta

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func testListFunc(path string, list List) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := MarshalList(list)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("compare raw list file: "+path, func(t *testing.T) {
			check, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(check) {
				t.Errorf("unexpected %s content -got/+exp\n%s", path, cmp.Diff(res, check))
			}
		})
		t.Run("check encode/decode list file: "+path, func(t *testing.T) {
			if err := UnmarshalList(res, list); err != nil {
				t.Fatal(err)
			}
			check, err := MarshalList(list)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(check) {
				t.Errorf("unexpected %s content -got/+exp\n%s", path, cmp.Diff(res, check))
			}
		})
		t.Run("check list file: "+path, func(t *testing.T) {
			if err := LoadList(path, list); err != nil {
				t.Fatal(err)
			}
			check, err := MarshalList(list)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(check) {
				t.Errorf("unexpected %s content -got/+exp\n%s", path, cmp.Diff(res, check))
			}
		})
	}
}

func TestList(t *testing.T) {

	var listtests = []struct {
		f string
		l List
	}{
		{
			"testdata/networks.csv",
			&NetworkList{
				Network{
					Code:        "AK",
					External:    "NZ",
					Description: "Auckland volcano seismic network",
					Restricted:  false,
				},
				Network{
					Code:        "CB",
					External:    "NZ",
					Description: "Canterbury regional seismic network",
					Restricted:  false,
				},
			},
		},
		{
			"testdata/stations.csv",
			&StationList{
				Station{
					Reference: Reference{
						Code:    "DFE",
						Network: "TR",
						Name:    "Dawson Falls",
					},
					Position: Position{
						Latitude:  -39.325743417,
						Longitude: 174.103863732,
						Elevation: 880.0,
						Depth:     0,
						Datum:     "WGS84",

						latitude:  "-39.325743417",
						longitude: "174.103863732",
						elevation: "880",
						depth:     "0",
					},
					Span: Span{
						Start: time.Date(1993, time.December, 14, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2010, time.February, 23, 0, 0, 0, 0, time.UTC),
					},
				},
				Station{
					Reference: Reference{
						Code:    "TBAS",
						Network: "SM",
						Name:    "Tolaga Bay Area School",
					},
					Position: Position{
						Latitude:  -38.372803703,
						Longitude: 178.300778623,
						Elevation: 8.0,
						Depth:     0,
						Datum:     "WGS84",

						latitude:  "-38.372803703",
						longitude: "178.300778623",
						elevation: "8",
						depth:     "0",
					},
					Span: Span{
						Start: time.Date(2002, time.March, 5, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					//Notes: "Is located in the Kiln Shed next to the hall",
				},
			},
		},
		{
			"testdata/mounts.csv",
			&MountList{
				Mount{
					Reference: Reference{
						Code: "MTSR",
						Name: "Ruapehu South",
					},
					Position: Position{
						Latitude:  -39.384607843,
						Longitude: 175.470410324,
						Elevation: 840,
						Datum:     "WGS84",

						latitude:  "-39.384607843",
						longitude: "175.470410324",
						elevation: "840",
					},
					Span: Span{
						Start: time.Date(2011, time.September, 8, 0, 10, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Description: "Images of Mount Ruapehu from the volcano camera situated at Mangateitei.",
				},
				Mount{
					Reference: Reference{
						Code: "RIMM",
						Name: "Raoul Island",
					},
					Position: Position{
						Latitude:  -29.267332,
						Longitude: -177.907235,
						Elevation: 490,
						Datum:     "WGS84",

						latitude:  "-29.267332",
						longitude: "-177.907235",
						elevation: "490",
					},
					Span: Span{
						Start: time.Date(2009, time.May, 18, 0, 0, 2, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Description: "Images looking into Green Lake on Raoul Island from the volcano camera situated on Mount Moumoukai.",
				},
			},
		},
		{
			"testdata/views.csv",
			&ViewList{
				View{
					Mount: "MTSR",
					Code:  "01",
					Label: "Ruapehu South",
					Orientation: Orientation{
						Dip:     10.0,
						Azimuth: 180.0,

						dip:     "10",
						azimuth: "180",
					},
					Description: "Images of Mount Ruapehu from the volcano camera situated at Mangateitei.",
					Span: Span{
						Start: time.Date(2011, time.September, 8, 0, 10, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},

				View{
					Mount: "RIMM",
					Code:  "02",
					Label: "Raoul Island",
					Orientation: Orientation{
						Dip:     -10.0,
						Azimuth: 280.0,

						dip:     "-10",
						azimuth: "280",
					},
					Description: "Images looking into Green Lake on Raoul Island from the volcano camera situated on Mount Moumoukai.",
					Span: Span{
						Start: time.Date(2009, time.May, 18, 0, 0, 2, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},

		{
			"testdata/sites.csv",
			&SiteList{
				Site{
					Position: Position{
						Latitude:  -39.198244208,
						Longitude: 175.547981982,
						Elevation: 1116.0,
						Depth:     0,
						Datum:     "WGS84",
						latitude:  "-39.198244208",
						longitude: "175.547981982",
						elevation: "1116",
						depth:     "0",
					},
					Survey: "GPS",
					Span: Span{
						Start: time.Date(2014, time.May, 16, 0, 0, 15, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Station:  "CNZ",
					Location: "12",
				},
				Site{
					Position: Position{
						Latitude:  -45.091369824,
						Longitude: 169.411775594,
						Elevation: 701.0,
						Depth:     0,
						Datum:     "WGS84",
						latitude:  "-45.091369824",
						longitude: "169.411775594",
						elevation: "701",
						depth:     "0",
					},
					Survey: "Map",
					Span: Span{
						Start: time.Date(1986, time.December, 9, 20, 10, 0, 0, time.UTC),
						End:   time.Date(1996, time.May, 1, 21, 38, 0, 0, time.UTC),
					},
					Station:  "MSCZ",
					Location: "10",
				},
			},
		},
		{
			"testdata/features.csv",
			&FeatureList{
				Feature{
					Span: Span{
						Start: time.Date(2014, time.May, 16, 0, 0, 15, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Station:     "CNZ",
					Location:    "12",
					Sublocation: "01",
					Property:    "Toto",
					Description: "Over the rainbow",
					Aspect:      "oklahoma",
				},
				Feature{
					Span: Span{
						Start: time.Date(1986, time.December, 9, 20, 10, 0, 0, time.UTC),
						End:   time.Date(1996, time.May, 1, 21, 38, 0, 0, time.UTC),
					},
					Station:     "MSCZ",
					Location:    "10",
					Sublocation: "02",
					Property:    "Tin",
					Description: "Somewhere up on high",
					Aspect:      "tulsa",
				},
			},
		},
		{
			"testdata/marks.csv",
			&MarkList{
				Mark{
					Reference: Reference{
						Code:    "AHTI",
						Network: "CG",
						Name:    "Ahititi",
					},
					Igs: false,
					Position: Position{
						Latitude:  -38.411447554,
						Longitude: 178.046002897,
						Elevation: 563.221,
						Datum:     "NZGD2000",

						latitude:  "-38.411447554",
						longitude: "178.046002897",
						elevation: "563.221",
					},
					Span: Span{
						Start: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Mark{
					Reference: Reference{
						Code:    "DUND",
						Network: "LI",
						Name:    "Dunedin",
					},
					Igs: true,
					Position: Position{
						Latitude:  -45.88366604,
						Longitude: 170.5971706,
						Elevation: 386.964,
						Datum:     "NZGD2000",

						latitude:  "-45.88366604",
						longitude: "170.5971706",
						elevation: "386.964",
					},
					Span: Span{
						Start: time.Date(2005, time.August, 10, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/monuments.csv",
			&MonumentList{
				Monument{
					Mark:               "CLIM",
					DomesNumber:        "",
					MarkType:           "Forced Centering",
					Type:               "Deep Braced",
					GroundRelationship: -1.00,
					groundRelationship: "-1",
					FoundationType:     "Steel Rods",
					FoundationDepth:    10.0,
					foundationDepth:    "10",
					Bedrock:            "Greywacke",
					Geology:            "",
					Span: Span{
						Start: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Monument{
					Mark:               "TAUP",
					DomesNumber:        "50217M001",
					MarkType:           "Forced Centering",
					Type:               "Pillar",
					GroundRelationship: -1.25,
					groundRelationship: "-1.25",
					FoundationType:     "Concrete",
					FoundationDepth:    2.0,
					foundationDepth:    "2",
					Bedrock:            "Rhyolite",
					Geology:            "",
					Span: Span{
						Start: time.Date(2005, time.August, 10, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/visibility.csv",
			&VisibilityList{
				Visibility{
					Mark:          "AHTI",
					SkyVisibility: "good",
					Span: Span{
						Start: time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Visibility{
					Mark:          "DUND",
					SkyVisibility: "clear to NW",
					Span: Span{
						Start: time.Date(2005, time.August, 10, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/assets.csv",
			&AssetList{
				{
					Equipment: Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220063995",
					},
					Number: "100",
				},
				{
					Equipment: Equipment{
						Make:   "Trimble",
						Model:  "Chokering Model 29659.00",
						Serial: "0220066912",
					},
					Number: "101",
				},
			},
		},
		{
			"testdata/antennas.csv",
			&InstalledAntennaList{
				{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220063995",
						},
						Span: Span{
							Start: time.Date(2000, time.August, 2, 23, 59, 1, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Offset: Offset{
						Vertical: 0.0015,
						North:    0.0,
						East:     0.0,

						vertical: "0.0015",
						north:    "0",
						east:     "0",
					},
					Mark:    "CNCL",
					Azimuth: 0.0,
					azimuth: "0",
				},
				{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "Chokering Model 29659.00",
							Serial: "0220066912",
						},
						Span: Span{
							Start: time.Date(2000, time.August, 14, 23, 59, 52, 0, time.UTC),
							End:   time.Date(2011, time.February, 7, 22, 35, 0, 0, time.UTC),
						},
					},
					Offset: Offset{
						Vertical: 0.0013,
						North:    0.0,
						East:     0.0,

						vertical: "0.0013",
						north:    "0",
						east:     "0",
					},
					Mark:    "MTJO",
					Azimuth: 10.0,
					azimuth: "10",
				},
			},
		},
		{
			"testdata/dataloggers.csv",
			&DeployedDataloggerList{
				DeployedDatalogger{
					Install: Install{
						Equipment: Equipment{
							Make:   "GNS Science",
							Model:  "EARSS/3",
							Serial: "152",
						},
						Span: Span{
							Start: time.Date(2001, time.January, 18, 13, 22, 0, 0, time.UTC),
							End:   time.Date(2001, time.February, 10, 10, 50, 0, 0, time.UTC),
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
							Start: time.Date(2009, time.February, 10, 23, 0, 1, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Place: "Turoa Road End",
				},
			},
		},
		{
			"testdata/metsensors.csv",
			&InstalledMetSensorList{
				InstalledMetSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65123",
						},
						Span: Span{
							Start: time.Date(1998, time.July, 9, 23, 59, 59, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Position: Position{
						Latitude:  -41.2351,
						Longitude: 174.917,
						Elevation: 26,
						Datum:     "NZGD2000",

						latitude:  "-41.2351",
						longitude: "174.917",
						elevation: "26",
					},
					Mark: "GRAC",
					Accuracy: MetSensorAccuracy{
						Humidity:    2.0,
						Pressure:    0.5,
						Temperature: 1,

						humidity:    "2",
						pressure:    "0.5",
						temperature: "1",
					},
				},
				InstalledMetSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Paroscientific",
							Model:  "Paroscientific meterological sensor",
							Serial: "65125",
						},
						Span: Span{
							Start: time.Date(2000, time.August, 15, 0, 0, 0, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Position: Position{
						Latitude:  -43.9857,
						Longitude: 170.4649,
						Elevation: 1044,
						Datum:     "NZGD2000",

						latitude:  "-43.9857",
						longitude: "170.4649",
						elevation: "1044",
					},
					Mark: "MTJO",
					Accuracy: MetSensorAccuracy{
						Humidity:    2.0,
						Pressure:    0.5,
						Temperature: 1,

						humidity:    "2",
						pressure:    "0.5",
						temperature: "1",
					},
				},
			},
		},
		{
			"testdata/cameras.csv",
			&InstalledCameraList{
				InstalledCamera{
					Install: Install{
						Equipment: Equipment{
							Make:   "Axis Communications AB",
							Model:  "AXIS-221",
							Serial: "00408C6DC9E1",
						},
						Span: Span{
							Start: time.Date(2006, time.February, 24, 14, 0, 0, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Dip:     0.0,
						Azimuth: 20.0,

						dip:     "0",
						azimuth: "20",
					},
					Offset: Offset{
						Vertical: -3.0,

						vertical: "-3",
						north:    "0",
						east:     "0",
					},
					Mount: "WHWI",
					View:  "01",
					Notes: "Looking at White Island",
				},
				InstalledCamera{
					Install: Install{
						Equipment: Equipment{
							Make:   "Mobotix AG",
							Model:  "M12 3MP",
							Serial: "0003c5041fc7",
						},
						Span: Span{
							Start: time.Date(2009, time.March, 3, 2, 0, 0, 0, time.UTC),
							End:   time.Date(2009, time.September, 18, 1, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Dip:     0.0,
						Azimuth: 280.0,

						dip:     "0",
						azimuth: "280",
					},
					Offset: Offset{
						Vertical: -10.0,

						vertical: "-10",
						north:    "0",
						east:     "0",
					},
					Mount: "K",
					View:  "01",
					Notes: "Bearing is magnetic.",
				},
			},
		},
		{
			"testdata/radomes.csv",
			&InstalledRadomeList{
				InstalledRadome{
					Install: Install{
						Equipment: Equipment{
							Make:   "LeicaGeosystems",
							Model:  "LEIS Radome",
							Serial: "0220148020",
						},
						Span: Span{
							Start: time.Date(1999, time.September, 27, 0, 0, 0, 0, time.UTC),
							End:   time.Date(2000, time.January, 21, 0, 0, 0, 0, time.UTC),
						},
					},
					Mark: "MQZG",
				},
				InstalledRadome{
					Install: Install{
						Equipment: Equipment{
							Make:   "Thales",
							Model:  "SCIS Radome",
							Serial: "0220063995",
						},
						Span: Span{
							Start: time.Date(2000, time.August, 3, 0, 0, 0, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Mark: "CNCL",
				},
			},
		},
		{
			"testdata/receivers.csv",
			&DeployedReceiverList{
				DeployedReceiver{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "5700",
							Serial: "220280300",
						},
						Span: Span{
							Start: time.Date(2002, time.December, 31, 1, 0, 1, 0, time.UTC),
							End:   time.Date(2012, time.August, 31, 15, 0, 1, 0, time.UTC),
						},
					},
					Mark: "HORN",
				},
				DeployedReceiver{
					Install: Install{
						Equipment: Equipment{
							Make:   "Trimble",
							Model:  "NetR9",
							Serial: "5014K66721",
						},
						Span: Span{
							Start: time.Date(2010, time.October, 12, 0, 0, 1, 0, time.UTC),
							End:   time.Date(2014, time.July, 27, 22, 0, 0, 0, time.UTC),
						},
					},
					Mark: "METH",
				},
			},
		},
		{
			"testdata/sensors.csv",
			&InstalledSensorList{
				InstalledSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "AppliedGeomechanics",
							Model:  "Lily tiltmeter",
							Serial: "N7935",
						},
						Span: Span{
							Start: time.Date(2009, time.November, 15, 1, 0, 0, 0, time.UTC),
							End:   time.Date(2013, time.May, 24, 0, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Azimuth: 233.0,
						Dip:     0.0,

						azimuth: "233",
						dip:     "0",
					},
					Offset: Offset{
						Vertical: -64.0,

						vertical: "64",
						north:    "0",
						east:     "0",
					},
					Scale: Scale{
						Factor: 1.0,
						Bias:   0.0,

						factor: "1",
						bias:   "0",
					},
					Station:  "COVZ",
					Location: "90",
				},
				InstalledSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   "Guralp",
							Model:  "CMG-3ESPC",
							Serial: "T36194",
						},
						Span: Span{
							Start: time.Date(2012, time.May, 21, 11, 0, 4, 0, time.UTC),
							End:   time.Date(2013, time.July, 10, 23, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Azimuth: 0.0,
						Dip:     0.0,

						azimuth: "0",
						dip:     "0",
					},
					Offset: Offset{
						Vertical: 0.0,

						vertical: "0",
						north:    "0",
						east:     "0",
					},
					Scale: Scale{
						Factor: 1.0,
						Bias:   0.0,

						factor: "1",
						bias:   "0",
					},
					Station:  "INZ",
					Location: "10",
				},
			},
		},
		{
			"testdata/recorders.csv",
			&InstalledRecorderList{
				InstalledRecorder{
					InstalledSensor: InstalledSensor{
						Install: Install{
							Equipment: Equipment{
								Make:   "Canterbury Seismic Instruments",
								Model:  "CUSP3A",
								Serial: "3A-040001",
							},
							Span: Span{
								Start: time.Date(2004, time.November, 27, 0, 0, 0, 0, time.UTC),
								End:   time.Date(2010, time.March, 25, 0, 30, 0, 0, time.UTC),
							},
						},
						Orientation: Orientation{
							Azimuth: 266.0,
							Dip:     0.0,

							azimuth: "266",
							dip:     "0",
						},
						Offset: Offset{
							Vertical: 0.0,

							vertical: "0",
						},
						Station:  "AMBC",
						Location: "20",
					},
					DataloggerModel: "CUSP3A",
				},
				InstalledRecorder{
					InstalledSensor: InstalledSensor{
						Install: Install{
							Equipment: Equipment{
								Make:   "Kinemetrics",
								Model:  "FBA-ES-T-DECK",
								Serial: "1275",
							},
							Span: Span{
								Start: time.Date(2014, time.April, 17, 0, 10, 0, 0, time.UTC),
								End:   time.Date(2014, time.July, 29, 0, 0, 0, 0, time.UTC),
							},
						},
						Orientation: Orientation{
							Azimuth: 210.0,
							Dip:     0.0,

							azimuth: "210",
							dip:     "0",
						},
						Offset: Offset{
							Vertical: 0.0,

							vertical: "0",
						},
						Station:  "EKS3",
						Location: "20",
					},
					DataloggerModel: "BASALT",
				},
			},
		},
		{
			"testdata/connections.csv",
			&ConnectionList{
				Connection{
					Station:  "APZ",
					Location: "10",
					Place:    "The Paps",
					Number:   3,
					Span: Span{
						Start: time.Date(2006, time.May, 7, 3, 23, 54, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					number: "3",
				},
				Connection{
					Station:  "BSWZ",
					Location: "10",
					Place:    "Blackbirch Station",
					Span: Span{
						Start: time.Date(2003, time.December, 9, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/sessions.csv",
			&SessionList{
				Session{
					Mark:            "TAUP",
					Operator:        "GeoNet",
					Agency:          "GNS",
					Model:           "TRIMBLE NETRS",
					SatelliteSystem: "GPS",
					Interval:        time.Second * 30,
					ElevationMask:   0,
					elevationMask:   "0",
					HeaderComment:   "linz",
					Format:          "trimble_5700 x5",
					Span: Span{
						Start: time.Date(2002, time.March, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/streams.csv",
			&StreamList{
				Stream{
					Station:      "AKSS",
					Location:     "20",
					Band:         "B",
					Source:       "N",
					SamplingRate: 50.0,
					samplingRate: "50",
					Axial:        "true",
					Triggered:    true,
					Span: Span{
						Start: time.Date(2011, time.August, 25, 0, 25, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Stream{
					Station:      "APZ",
					Location:     "20",
					Band:         "H",
					Source:       "N",
					SamplingRate: 200.0,
					samplingRate: "200",
					Axial:        "false",
					Triggered:    false,
					Span: Span{
						Start: time.Date(2007, time.May, 2, 22, 0, 1, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			"testdata/gauges.csv",
			&GaugeList{
				Gauge{
					Span: Span{
						Start: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Reference: Reference{
						Code:    "AUCT",
						Network: "TG",
					},
					Number:   "363",
					TimeZone: 180.0,
					timeZone: "180",
					Position: Position{
						Latitude:  36.5,
						Longitude: 174.47,

						latitude:  "36.5",
						longitude: "174.47",
					},
					Crex: "-3683144 17478654 AUCT",
				},
				Gauge{
					Span: Span{
						Start: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Reference: Reference{
						Code:    "CPIT",
						Network: "TG",
					},
					Number:   "313",
					TimeZone: 180.0,
					timeZone: "180",
					Position: Position{
						Latitude:  40.55,
						Longitude: 176.13,

						latitude:  "40.55",
						longitude: "176.13",
					},
					Crex: "-4089929 17623168 CPIT",
				},
			},
		},
		{
			"testdata/constituents.csv",
			&ConstituentList{
				Constituent{
					Span: Span{
						Start: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Gauge:     "AUCT",
					Number:    1,
					Name:      "Z0",
					Amplitude: 186.2448,
					Lag:       0,

					amplitude: "186.2448",
					lag:       "0",
				},
				Constituent{
					Span: Span{
						Start: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
					Gauge:     "AUCT",
					Location:  "41",
					Number:    2,
					Name:      "SA",
					Amplitude: 3.8781,
					Lag:       112.03,

					amplitude: "3.8781",
					lag:       "112.03",
				},
			},
		},
		{
			"testdata/doases.csv",
			&InstalledDoasList{
				InstalledDoas{
					Install: Install{
						Equipment: Equipment{
							Make:   "Avaspec CompactLine",
							Model:  "Avaspec-Mini2048CL",
							Serial: "2002125M1",
						},
						Span: Span{
							Start: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Dip:     -60.0,
						Azimuth: 169.0,

						dip:     "-60",
						azimuth: "169",
					},
					Offset: Offset{
						vertical: "0",
						north:    "0",
						east:     "0",
					},
					Mount: "TOD01",
					View:  "01",
				},
				InstalledDoas{
					Install: Install{
						Equipment: Equipment{
							Make:   "Avaspec CompactLine",
							Model:  "Avaspec-Mini2048CL",
							Serial: "2002127M1",
						},
						Span: Span{
							Start: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
							End:   time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Orientation: Orientation{
						Dip:     -60.0,
						Azimuth: 250.0,

						dip:     "-60",
						azimuth: "250",
					},
					Offset: Offset{
						vertical: "0",
						north:    "0",
						east:     "0",
					},
					Mount: "TOD02",
					View:  "01",
				},
			},
		},
		{
			"testdata/classes.csv",
			&ClassList{
				Class{
					Station:     "WHAS",
					SiteClass:   "C",
					Vs30:        270,
					Vs30Quality: "Q3",
					Tsite: Range{
						Value: 0.4,
					},
					TsiteMethod:   "I",
					TsiteQuality:  "Q3",
					BasementDepth: 40,
					DepthQuality:  "Q3",
					Citations:     []string{"Perrin2015a"},
					Notes:         "Perrin et al. 2015",
				},
				Class{
					Station:     "WKZ",
					SiteClass:   "B",
					Vs30:        1000,
					Vs30Quality: "Q3",
					Tsite: Range{
						Compare: LessThan,
						Value:   0.1,
					},
					TsiteMethod:   "I",
					TsiteQuality:  "Q3",
					BasementDepth: 0,
					DepthQuality:  "Q3",
					Citations:     []string{"Kaiser2017", "Perrin2015a"},
					Notes:         "Perrin et al. 2015",
				},
			},
		},
	}

	for _, tt := range listtests {
		res, err := MarshalList(tt.l)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("Compare raw list file: " + tt.f)
		{
			b, err := os.ReadFile(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(b) {
				t.Errorf("unexpected %s content -got/+exp\n%s", tt.f, cmp.Diff(string(res), string(b)))
			}
		}
		t.Log("Check encode/decode list: " + tt.f)
		{
			if err := UnmarshalList(res, tt.l); err != nil {
				t.Fatal(err)
			}

			s, err := MarshalList(tt.l)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(s) {
				t.Errorf("unexpected %s content -got/+exp\n%s", tt.f, cmp.Diff(string(res), string(s)))
			}
		}

		t.Log("Check list file: " + tt.f)
		{
			if err := LoadList(tt.f, tt.l); err != nil {
				t.Fatal(err)
			}

			s, err := MarshalList(tt.l)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(s) {
				t.Errorf("unexpected %s content -got/+exp\n%s", tt.f, cmp.Diff(string(res), string(s)))
			}
		}
	}
}

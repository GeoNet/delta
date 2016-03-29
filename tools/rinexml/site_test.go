package main

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestSiteXML_Marshal(t *testing.T) {

	var tests = []struct {
		s string
		x SiteXML
	}{
		{
			strings.Join([]string{
				"<SITE>",
				"  <mark>",
				"    <geodetic-code>YALD</geodetic-code>",
				"    <domes-number></domes-number>",
				"  </mark>",
				"  <location>",
				"    <latitude>-43.490766992</latitude>",
				"    <longitude>172.481129517</longitude>",
				"    <height>64.65</height>",
				"    <datum>WGS84</datum>",
				"  </location>",
				"  <cgps-session>",
				"    <start-time>2015-08-31T20:40:01</start-time>",
				"    <stop-time>open</stop-time>",
				`    <observation-interval unit="s">30</observation-interval>`,
				"    <operator>",
				"      <name>GeoNet</name>",
				"      <agency>GNS</agency>",
				"    </operator>",
				"    <rinex>",
				"      <header-comment-name>geonet</header-comment-name>",
				"      <header-comment-text>Data supplied by the GeoNet project. GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders. Contact: www.geonet.org.nz  email: info@geonet.org.nz.</header-comment-text>",
				"    </rinex>",
				"    <data-format>trimble_netr9</data-format>",
				`    <download-name-format type="long">`,
				"      <year>x4 A4 x*</year>",
				"      <month>x8 A2 x*</month>",
				"      <day>x10 A2 x*</day>",
				"      <hour>x12 A2 x*</hour>",
				"      <minute>x14 A2 x*</minute>",
				"      <second>x16 A2 x*</second>",
				"    </download-name-format>",
				"    <receiver>",
				"      <serial-number>5307K50971</serial-number>",
				"      <igs-designation>TRIMBLE NETR9</igs-designation>",
				"      <firmware-history>",
				"        <start-time>2015-08-31T20:40:01</start-time>",
				"        <stop-time>open</stop-time>",
				"        <version>4.85</version>",
				"      </firmware-history>",
				"      <firmware-history>",
				"        <start-time>2014-12-16T00:00:01</start-time>",
				"        <stop-time>2015-08-31T20:40:00</stop-time>",
				"        <version>4.82</version>",
				"      </firmware-history>",
				"      <firmware-history>",
				"        <start-time>2014-05-07T00:00:00</start-time>",
				"        <stop-time>2014-12-16T00:00:00</stop-time>",
				"        <version>4.81</version>",
				"      </firmware-history>",
				"    </receiver>",
				"    <installed-cgps-antenna>",
				`      <height unit="m">0.035</height>`,
				`      <offset-east unit="m">0</offset-east>`,
				`      <offset-north unit="m">0</offset-north>`,
				"      <radome>NONE</radome>",
				"      <cgps-antenna>",
				"        <serial-number>1441040153</serial-number>",
				"        <igs-designation>TRM57971.00</igs-designation>",
				"      </cgps-antenna>",
				"    </installed-cgps-antenna>",
				"  </cgps-session>",
				"</SITE>",
				"",
			}, "\n"),
			SiteXML{
				XMLName: xml.Name{Local: "SITE"},
				Mark: MarkXML{
					GeodeticCode: "YALD",
					DomesNumber:  "",
				},
				Location: LocationXML{
					Latitude:  -43.490766992,
					Longitude: 172.481129517,
					Height:    64.65,
					Datum:     "WGS84",
				},
				Sessions: []CGPSSessionXML{
					CGPSSessionXML{
						StartTime: "2015-08-31T20:40:01",
						StopTime:  "open",
						ObservationInterval: Number{
							Units: "s",
							Value: 30,
						},
						Operator: OperatorXML{
							Name:   "GeoNet",
							Agency: "GNS",
						},
						Rinex: RinexXML{
							HeaderCommentName: "geonet",
							HeaderCommentText: "Data supplied by the GeoNet project. GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders. Contact: www.geonet.org.nz  email: info@geonet.org.nz.",
						},
						DataFormat: "trimble_netr9",
						DownloadNameFormat: DownloadNameFormatXML{
							Type:   "long",
							Year:   "x4 A4 x*",
							Month:  "x8 A2 x*",
							Day:    "x10 A2 x*",
							Hour:   "x12 A2 x*",
							Minute: "x14 A2 x*",
							Second: "x16 A2 x*",
						},
						Receiver: ReceiverXML{
							SerialNumber:   "5307K50971",
							IGSDesignation: "TRIMBLE NETR9",
							FirmwareHistories: []FirmwareHistoryXML{
								FirmwareHistoryXML{
									StartTime: "2015-08-31T20:40:01",
									StopTime:  "open",
									Version:   "4.85",
								},
								FirmwareHistoryXML{
									StartTime: "2014-12-16T00:00:01",
									StopTime:  "2015-08-31T20:40:00",
									Version:   "4.82",
								},
								FirmwareHistoryXML{
									StartTime: "2014-05-07T00:00:00",
									StopTime:  "2014-12-16T00:00:00",
									Version:   "4.81",
								},
							},
						},
						InstalledCGPSAntenna: InstalledCGPSAntennaXML{
							Height:      Number{Units: "m", Value: 0.0350},
							OffsetEast:  Number{Units: "m", Value: 0.0000},
							OffsetNorth: Number{Units: "m", Value: 0.0000},
							Radome:      "NONE",
							CGPSAntenna: CGPSAntennaXML{
								SerialNumber:   "1441040153",
								IGSDesignation: "TRM57971.00",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s, err := test.x.Marshal()
		if err != nil {
			t.Error(err)
		}

		if (string)(s) != test.s {
			t.Error(strings.Join([]string{"marshalling mismatch:", (string)(s), test.s, ""}, "\n=========\n"))
		}
	}
}

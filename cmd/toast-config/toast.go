package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

const ToastTimeFormat = "2006-01-02T15:04:05.999999Z07:00"
const PublicIdTimeFormat = "20060102150405.000000"

const ToastSchema = "http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.11"
const ToastSchemaVersion = "0.11"

var streamRe = regexp.MustCompile("^(BTZ|UTZ)$")

var count int

func parameterCount() int {
	count++
	return count
}

func publicId() string {
	return fmt.Sprintf("%s.%04d", time.Now().UTC().Format(PublicIdTimeFormat), parameterCount())
}

type Parameter struct {
	PublicId string `xml:"publicID,attr"`

	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type ParameterSet struct {
	PublicId string `xml:"publicID,attr"`
	Created  string `xml:"created,attr"`

	ModuleId   string      `xml:"moduleID"`
	Parameters []Parameter `xml:"parameter"`
}

type Setup struct {
	Name    string `xml:"name,attr"`
	Enabled bool   `xml:"enabled,attr"`

	ParameterSetID string `xml:"parameterSetID"`
}

type Station struct {
	PublicId    string `xml:"publicID,attr"`
	NetworkCode string `xml:"networkCode,attr"`
	StationCode string `xml:"stationCode,attr"`
	Enabled     bool   `xml:"enabled,attr"`

	AgencyID     string `xml:"creationInfo>agencyID"`
	Author       string `xml:"creationInfo>author"`
	CreationTime string `xml:"creationInfo>creationTime"`

	Setup *Setup `xml:"setup,omitempty"`
}

type Module struct {
	PublicId string `xml:"publicID,attr"`
	Name     string `xml:"name,attr"`
	Enabled  bool   `xml:"enabled,attr"`

	Stations []Station `xml:"station"`
}

type Config struct {
	ParameterSets []ParameterSet `xml:"parameterSet"`
	Module        Module         `xml:"module"`
}

type Toast struct {
	XMLName   xml.Name `xml:"seiscomp"`
	NameSpace string   `xml:"xmlns,attr"`
	Version   string   `xml:"version,attr"`

	Config Config `xml:"Config"`
}

func NewToast(set *meta.Set, agency string, networks ...string) (*Toast, error) {

	var stations []Station
	var parameterSets []ParameterSet

	for _, list := range networks {
		for _, code := range strings.Split(list, ",") {
			if code = strings.TrimSpace(code); code == "" {
				continue
			}
			net, ok := set.Network(code)
			if !ok {
				continue
			}
			for _, stn := range set.Stations() {
				if stn.Network != code {
					continue
				}

				var setup *Setup
				created := time.Now().UTC()
				for _, site := range set.Sites() {
					if site.Station != stn.Code {
						continue
					}
					if time.Since(site.End) > 0 {
						continue
					}

					for _, coll := range set.Collections(site) {
						if !streamRe.MatchString(coll.Code()) {
							continue
						}

						parameterSets = append(parameterSets, ParameterSet{
							PublicId: fmt.Sprintf("ParameterSet/trunk/Station/%s/%s/default", net.External, stn.Code),
							Created:  created.Format(ToastTimeFormat),
							ModuleId: "Config/trunk",
							Parameters: []Parameter{
								{
									//PublicId: fmt.Sprintf("Parameter/20200727001055.014792.4500",
									PublicId: fmt.Sprintf("Parameter/%s", publicId()),
									//20200727001055.014792.4500",
									Name:  "detecLocid",
									Value: site.Location,
								},
								{
									//PublicId: "Parameter/20200727001055.014751.4499",
									Name:  "detecStream",
									Value: coll.Code(),
								},
							},
						})

						setup = &Setup{
							Name:           "default",
							Enabled:        true,
							ParameterSetID: fmt.Sprintf("ParameterSet/trunk/Station/%s/%s/default", net.External, stn.Code),
						}
						break
					}
					break
				}
				stations = append(stations, Station{
					PublicId:     fmt.Sprintf("Config/trunk/%s/%s", net.External, stn.Code),
					NetworkCode:  net.External,
					StationCode:  stn.Code,
					Enabled:      true,
					AgencyID:     agency,
					Author:       "trunk",
					CreationTime: created.Format(ToastTimeFormat),
					Setup:        setup,
				})

			}
		}
	}

	toast := Toast{
		XMLName: xml.Name{
			Space: ToastSchema,
			Local: "seiscomp",
		},
		NameSpace: ToastSchema,
		Version:   ToastSchemaVersion,

		Config: Config{
			ParameterSets: parameterSets,
			Module: Module{
				PublicId: "Config/trunk",
				Name:     "trunk",
				Enabled:  true,
				Stations: stations,
			},
		},
	}

	return &toast, nil
}

func (t *Toast) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, t)
}

func (t *Toast) Marshal() ([]byte, error) {
	raw, err := xml.Marshal(t)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(raw, '\n')...), nil
}

func (t *Toast) MarshalIndent(a, b string) ([]byte, error) {
	raw, err := xml.MarshalIndent(t, a, b)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(raw, '\n')...), nil
}

func (t *Toast) Encode(wr io.Writer) error {
	data, err := t.MarshalIndent("", "  ")
	if err != nil {
		return err
	}
	if _, err := wr.Write(data); err != nil {
		return err
	}
	return nil
}

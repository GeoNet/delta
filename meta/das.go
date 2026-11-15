package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	dasVersion int = iota
	dasNetwork
	dasLocation
	dasCountry
	dasPrincipalInvestigatorName
	dasPointOfContact
	dasPointOfContactEmail
	dasPointOfContactAddress
	dasStart
	dasEnd
	dasLatitude
	dasLongitude
	dasElevation
	dasDepth
	dasDatum
	dasFundingAgency
	dasProjectNumber
	dasDigitalObjectIdentifier
	dasPurposeOfDataCollection
	dasPrincipalInvestigatorEmail
	dasPrincipalInvestigatorAddress
)

var dasHeaders Header = map[string]int{
	"Network Code": dasNetwork,
	"Location":     dasLocation,
	"Country":      dasCountry,
	// "Latitude":   stationLatitude,
	// "Longitude":  stationLongitude,
	"Principal Investigator":     dasPrincipalInvestigatorName,
	"Point Of Contact":           dasPointOfContact,
	"Point Of Contact Email":     dasPointOfContactEmail,
	"Point Of Contact Address":   dasPointOfContactAddress,
	"Funding Agency":             dasFundingAgency,
	"Project Number":             dasProjectNumber,
	"Digital Object Identifier":  dasDigitalObjectIdentifier,
	"Purpose Of Data Collection": dasPurposeOfDataCollection,
	// "Name":                       stationName,

	// "Elevation":  dasElevation,
	// "Depth":      dasDepth,
	// "Datum":      dasDatum,
	"Start Date": dasStart,
	"End Date":   dasEnd,
}

var DASTable Table = Table{
	name:    "DAS",
	headers: dasHeaders,
	primary: []string{"Network Code", "Location", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
		// "Network Code":               "dasNetworkCode",
		// "Principal Investigator":     "dasPrincipalInvestigator",
		// "Point Of Contact":           "dasPointOfContact",
		// "Point Of Contact Email":     "daspointOfContactEmail",
		// "Point Of Contact Address":   "dasPointOfContactAddress",
		// "Start Date":                 "dasStart",
		// "End Date":                   "dasEnd",
		// "Funding Agency":             "dasFundingAgency",
		// "Project Number":             "dasProjectNumber",
		// "Digital Object Identifier":  "dasDigitalObjectIdentifier",
		// "Purpose Of Data Collection": "dasPurposeOfDataCollection",
	},
	start: "Start Date",
	end:   "End Date",
}

type DASStation struct {
	// Reference
	Position
	Span

	// Version                 string `json:"version"`
	Network                 string                   `json:"network_code"`
	Location                string                   `json:"location"`
	Country                 string                   `json:"country"`
	PrincipalInvestigator   DASPrincipalInvestigator `json:"principal_investigator"`
	PointOfContact          string                   `json:"point_of_contact"`
	PointOfContactEmail     string                   `json:"point_of_contact_email"`
	PointOfContactAddress   string                   `json:"point_of_contact_address"`
	FundingAgency           string                   `json:"funding_agency"`
	ProjectNumber           string                   `json:"project_number"`
	DigitalObjectIdentifier string                   `json:"digital_object_identifier"`
	PurposeOfDataCollection string                   `json:"purpose_of_data_collection"`
}

type DASPrincipalInvestigator struct {
	Name    string
	Email   string
	Address string
}

type DASList []DASStation

func (dl DASList) Len() int      { return len(dl) }
func (dl DASList) Swap(i, j int) { dl[i], dl[j] = dl[j], dl[i] }
func (dl DASList) Less(i, j int) bool {
	switch {
	case strings.ToLower(dl[i].Network) < strings.ToLower(dl[j].Network):
		return true
	case strings.ToLower(dl[i].Network) > strings.ToLower(dl[j].Network):
		return false
	default:
		return dl[i].Start.Before(dl[j].Start)
	}
}

func (dl DASList) encode() [][]string {
	var data [][]string

	data = append(data, dasHeaders.Columns())

	for _, row := range dl {
		data = append(data, []string{
			// strings.TrimSpace(row.Version),
			strings.TrimSpace(row.Network),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Country),
			strings.TrimSpace(row.PrincipalInvestigator.Name),
			strings.TrimSpace(row.PointOfContact),
			strings.TrimSpace(row.PointOfContactEmail),
			strings.TrimSpace(row.PointOfContactAddress),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
			strings.TrimSpace(row.FundingAgency),
			strings.TrimSpace(row.ProjectNumber),
			strings.TrimSpace(row.DigitalObjectIdentifier),
			strings.TrimSpace(row.PurposeOfDataCollection),
		})
	}

	return data
}

func (dl *DASList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var dasList []DASStation

	fields := dasHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)
		start, err := time.Parse(DateTimeFormat, d[dasStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[dasEnd])
		if err != nil {
			return err
		}
		var dasPrincipalInvestigator DASPrincipalInvestigator
		if d[dasPrincipalInvestigatorName] == d[dasPointOfContact] {
			dasPrincipalInvestigator = DASPrincipalInvestigator{
				Name:    d[dasPrincipalInvestigatorName],
				Email:   d[dasPointOfContactEmail],
				Address: d[dasPointOfContactAddress],
			}
		} else {
			dasPrincipalInvestigator = DASPrincipalInvestigator{
				Name:    d[dasPrincipalInvestigatorName],
				Email:   "",
				Address: "",
			}
			fmt.Printf("Principle investigator %s is different to point of contact for das deployment %s, metadata may be incomplete\n", d[dasPrincipalInvestigatorName], d[dasNetwork])
		}

		dasList = append(dasList, DASStation{
			Span: Span{
				Start: start,
				End:   end,
			},
			// Version:                 strings.TrimSpace(d[dasVersion]),
			Network:                 strings.TrimSpace(d[dasNetwork]),
			Location:                strings.TrimSpace(d[dasLocation]),
			Country:                 strings.TrimSpace(d[dasCountry]),
			PrincipalInvestigator:   dasPrincipalInvestigator,
			PointOfContact:          strings.TrimSpace(d[dasPointOfContact]),
			PointOfContactEmail:     strings.TrimSpace(d[dasPointOfContactEmail]),
			PointOfContactAddress:   strings.TrimSpace(d[dasPointOfContactAddress]),
			FundingAgency:           strings.TrimSpace(d[dasFundingAgency]),
			ProjectNumber:           strings.TrimSpace(d[dasProjectNumber]),
			DigitalObjectIdentifier: strings.TrimSpace(d[dasDigitalObjectIdentifier]),
			PurposeOfDataCollection: strings.TrimSpace(d[dasPurposeOfDataCollection]),
		})
		// fmt.Print("succeded parsing das", "\n")
	}

	*dl = DASList(dasList)

	return nil
}

func LoadDAS(path string) ([]DASStation, error) {
	var dl []DASStation

	if err := LoadList(path, (*DASList)(&dl)); err != nil {
		return nil, err
	}

	sort.Sort(DASList(dl))

	return dl, nil
}

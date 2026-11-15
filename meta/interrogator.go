package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	interrogatorID int = iota
	interrogatorNetworkCode
	interrogatorManufacturer
	interrogatorModel
	interrogatorSerialNumber
	interrogatorFirmwareVersion
	interrogatorComment
	interrogatorAcquisitions
	interrogatorStart
	interrogatorEnd
	interrogatorLast
)

var interrogatorHeaders Header = map[string]int{
	"Manufacturer":    interrogatorManufacturer,
	"Model":           interrogatorModel,
	"Serial Number":   interrogatorSerialNumber,
	"Interrogator ID": interrogatorID,
	"Network Code":    interrogatorNetworkCode,
	// "Firmware Version": interrogatorFirmwareVersion,
	// "Comment":      interrogatorComment,
	// "Acquisitions": interrogatorAcquisitions,
	"Start Date": interrogatorStart,
	"End Date":   interrogatorEnd,
}

var InterrogatorTable Table = Table{
	name:    "Interrogator",
	headers: interrogatorHeaders,
	primary: []string{"Interrogator ID", "Network Code", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledInterrogator struct {
	Span

	InterrogatorID string
	Manufacturer   string
	Model          string
	SerialNumber   string
	NetworkCode    string
	// FirmwareVersion string `json:"firmware_version"`
	// Comment      string `json:"comment"`
	// Acquisitions int    `json:"acquisitions"`
}

type InterrogatorList []InstalledInterrogator

func (il InterrogatorList) Len() int      { return len(il) }
func (il InterrogatorList) Swap(i, j int) { il[i], il[j] = il[j], il[i] }
func (il InterrogatorList) Less(i, j int) bool {
	switch {
	case strings.ToLower(il[i].InterrogatorID) < strings.ToLower(il[j].InterrogatorID):
		return true
	case strings.ToLower(il[i].InterrogatorID) > strings.ToLower(il[j].InterrogatorID):
		return false
	default:
		return il[i].Start.Before(il[j].Start)
	}
}

func (il InterrogatorList) encode() [][]string {
	var data [][]string

	data = append(data, interrogatorHeaders.Columns())

	for _, row := range il {
		data = append(data, []string{
			strings.TrimSpace(row.InterrogatorID),
			strings.TrimSpace(row.Manufacturer),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.SerialNumber),
			// strings.TrimSpace(row.FirmwareVersion),
			// strings.TrimSpace(row.Comment),
			// strconv.Itoa(row.Acquisitions),
			strings.TrimSpace(row.NetworkCode),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (il *InterrogatorList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var interrogators []InstalledInterrogator

	fields := interrogatorHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		// acquisitions, err := strconv.Atoi(strings.TrimSpace(d[interrogatorAcquisitions]))
		// if err != nil {
		// return err
		// }

		start, err := time.Parse(DateTimeFormat, d[interrogatorStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[interrogatorEnd])
		if err != nil {
			return err
		}

		interrogators = append(interrogators, InstalledInterrogator{
			Span: Span{
				Start: start,
				End:   end,
			},
			InterrogatorID: strings.TrimSpace(d[interrogatorID]),
			Manufacturer:   strings.TrimSpace(d[interrogatorManufacturer]),
			Model:          strings.TrimSpace(d[interrogatorModel]),
			SerialNumber:   strings.TrimSpace(d[interrogatorSerialNumber]),
			NetworkCode:    strings.TrimSpace(d[interrogatorNetworkCode]),
			// FirmwareVersion: strings.TrimSpace(d[interrogatorFirmwareVersion]),
			// Comment:         strings.TrimSpace(d[interrogatorComment]),
			// Acquisitions:    acquisitions,
		})
		// fmt.Print("succeded parsing Interrogators", "\n")
	}

	*il = InterrogatorList(interrogators)

	return nil
}

func LoadInterrogators(path string) ([]InstalledInterrogator, error) {
	var il []InstalledInterrogator

	if err := LoadList(path, (*InterrogatorList)(&il)); err != nil {
		return nil, err
	}

	sort.Sort(InterrogatorList(il))

	return il, nil
}

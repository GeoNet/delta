package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	cableNetworkCode int = iota
	cableID
	cableModel
	cableBoundingBox
	cableOwner
	cableCharacteristics
	cableOutsideDiameter
	cableEnvironment
	cableInstallationEnvironment
	cableInstallationDate
	cableRemovalDate

	// cableOutsideDiameterUnit
	// cableComment
	// cableFibers
	cableStart
	cableEnd
	// cableLast
)

var cableHeaders Header = map[string]int{
	"Network Code":                   cableNetworkCode,
	"Cable ID":                       cableID,
	"Cable Model":                    cableModel,
	"Bounding Box":                   cableBoundingBox,
	"Cable Owner":                    cableOwner,
	"Cable Characteristics":          cableCharacteristics,
	"Cable Outside Diameter":         cableOutsideDiameter,
	"Cable Environment":              cableEnvironment,
	"Cable Installation Environment": cableInstallationEnvironment,
	// "Cable Outside Diameter Unit":    cableOutsideDiameterUnit,
	"Start Date": cableStart,
	"End Date":   cableEnd,
}

var CableTable Table = Table{
	name:    "Cable",
	headers: cableHeaders,
	primary: []string{"Network Code", "Cable ID", "Start Date"},
	native:  []string{"Cable Outside Diameter"},
	foreign: map[string][]string{},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledCable struct {
	Span
	NetworkCode string
	CableID     string

	BoundingBox []float64
	Owner       string
	// InstallationDate        time.Time `json:"cable_installation_date"`
	// RemovalDate             time.Time `json:"cable_removal_date"`
	Characteristics         string
	Environment             string
	InstallationEnvironment string
	Model                   string
	OutsideDiameter         float64
	OutsideDiameterUnit     string
	// Comment                 string    `json:"comment"`
	// Fibers string `json:"fibers"`

	outsideDiameter string // shadow value used to retain formatting
}

type CableList []InstalledCable

func (cl CableList) Len() int      { return len(cl) }
func (cl CableList) Swap(i, j int) { cl[i], cl[j] = cl[j], cl[i] }
func (cl CableList) Less(i, j int) bool {
	switch {
	case strings.ToLower(cl[i].CableID) < strings.ToLower(cl[j].CableID):
		return true
	case strings.ToLower(cl[i].CableID) > strings.ToLower(cl[j].CableID):
		return false
	default:
		return cl[i].Start.Before(cl[j].Start)
	}
}

func (cl CableList) encode() [][]string {
	var data [][]string

	data = append(data, cableHeaders.Columns())

	for _, row := range cl {
		data = append(data, []string{
			strings.TrimSpace(row.NetworkCode),
			strings.TrimSpace(row.CableID),
			strings.TrimSpace(row.Model),
			strconv.Quote(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(row.BoundingBox)), ","), "[]")), // converts slice of floats to quoted string

			strings.TrimSpace(row.Owner),
			// row.InstallationDate.Format(DateTimeFormat),
			// row.RemovalDate.Format(DateTimeFormat),
			strings.TrimSpace(row.Characteristics),
			strings.TrimSpace(row.outsideDiameter),
			strings.TrimSpace(row.Environment),
			strings.TrimSpace(row.InstallationEnvironment),

			// strings.TrimSpace(row.OutsideDiameterUnit),
			// strings.TrimSpace(row.Comment),
			// strings.TrimSpace(row.Fibers),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (cl *CableList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var cables []InstalledCable

	fields := cableHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[cableStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[cableEnd])
		if err != nil {
			return err
		}

		outsideDiameter, err := strconv.ParseFloat(strings.TrimSpace(d[cableOutsideDiameter]), 64)
		if err != nil {
			return err
		}
		var boundingBox []float64
		for _, chunk := range strings.Split(d[cableBoundingBox], ",") {
			f, err := strconv.ParseFloat(chunk, 64)
			if err != nil {
				return err
			}
			boundingBox = append(boundingBox, f)
		}
		// fmt.Print(d[cableBoundingBox], boundingBox, "\n")

		cables = append(cables, InstalledCable{
			Span: Span{
				Start: start,
				End:   end,
			},
			NetworkCode: strings.TrimSpace(d[networkCode]),
			CableID:     strings.TrimSpace(d[cableID]),
			BoundingBox: boundingBox,
			// BoundingBox: strings.TrimSpace(d[cableBoundingBox]),
			Owner:                   strings.TrimSpace(d[cableOwner]),
			Characteristics:         strings.TrimSpace(d[cableCharacteristics]),
			Environment:             strings.TrimSpace(d[cableEnvironment]),
			InstallationEnvironment: strings.TrimSpace(d[cableInstallationEnvironment]),
			Model:                   strings.TrimSpace(d[cableModel]),
			OutsideDiameter:         outsideDiameter,
			outsideDiameter:         strings.TrimSpace(d[cableOutsideDiameter]), //shadow to retain formatting
			// OutsideDiameterUnit:     strings.TrimSpace(d[cableOutsideDiameterUnit]),
			// Comment:                 strings.TrimSpace(d[cableComment]),
			// Fibers:                  strings.TrimSpace(d[cableFibers]),

		})
		// fmt.Print("succeded parsing Cable", "\n")
		// fmt.Print(cables[0], "\n\n")
	}

	*cl = CableList(cables)

	return nil
}

func LoadCables(path string) ([]InstalledCable, error) {
	var cl []InstalledCable

	if err := LoadList(path, (*CableList)(&cl)); err != nil {
		return nil, err
	}

	sort.Sort(CableList(cl))

	return cl, nil
}

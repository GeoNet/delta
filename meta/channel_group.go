package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	channelGroupID int = iota
	ChannelAcquisitionID
	channelGroupCableID
	channelGroupFiberID
	channelGroupCoordinateGenerationDate
	channelGroupCoordinateSystem
	channelGroupReferenceFrame
	channelGroupLocationMethod
	// channelGroupDistanceAlongFiberUnit
	// channelGroupXCoordinateUnit
	channelGroupUncertaintyInXCoordinate
	// channelGroupUncertaintyInXCoordinateUnit
	// channelGroupYCoordinateUnit
	channelGroupUncertaintyInYCoordinate
	// channelGroupUncertaintyInYCoordinateUnit
	// channelGroupElevationAboveSeaLevelUnit
	channelGroupUncertaintyInDepth
	channelGroupUncertaintyInElevation
	// channelGroupUncertaintyInElevationUnit
	// channelGroupStrikeUnit
	channelGroupUncertaintyInStrike
	// channelGroupDipUnit
	channelGroupUncertaintyInDip
	// channelGroupUncertaintyInDipUnit
	channelGroupFirstUsableChannelID
	channelGroupLastUsableChannelID
	// channelGroupChannels
	channelGroupStart
	channelGroupEnd
	// channelGroupLast
)

var channelGroupHeaders Header = map[string]int{
	"Channel Group ID":           channelGroupID,
	"Acquisition ID":             ChannelAcquisitionID,
	"Cable ID":                   channelGroupCableID,
	"Fiber ID":                   channelGroupFiberID,
	"Coordinate Generation Date": channelGroupCoordinateGenerationDate,
	"Coordinate System":          channelGroupCoordinateSystem,
	"Reference Frame":            channelGroupReferenceFrame,
	"Location Method":            channelGroupLocationMethod,
	// "Distance Along Fiber Unit":        channelGroupDistanceAlongFiberUnit,
	// "X Coordinate Unit":                channelGroupXCoordinateUnit,
	"X Uncertainty": channelGroupUncertaintyInXCoordinate,
	// "Uncertainty In X Coordinate Unit": channelGroupUncertaintyInXCoordinateUnit,
	// "Y Coordinate Unit":                channelGroupYCoordinateUnit,
	"Y Uncertainty": channelGroupUncertaintyInYCoordinate,
	// "Uncertainty In Y Coordinate Unit": channelGroupUncertaintyInYCoordinateUnit,
	"Depth Uncertainty": channelGroupUncertaintyInDepth,
	// "Elevation Above Sea Level Unit":   channelGroupElevationAboveSeaLevelUnit,
	"Elevation Uncertainty": channelGroupUncertaintyInElevation,
	// "Uncertainty In Elevation Unit":    channelGroupUncertaintyInElevationUnit,
	// "Strike Unit":                      channelGroupStrikeUnit,
	"Strike Uncertainty": channelGroupUncertaintyInStrike,
	// "Dip Unit":                         channelGroupDipUnit,
	"Dip Uncertainty": channelGroupUncertaintyInDip,
	// "Uncertainty In Dip Unit":          channelGroupUncertaintyInDipUnit,
	"First Usable Channel ID": channelGroupFirstUsableChannelID,
	"Last Usable Channel ID":  channelGroupLastUsableChannelID,
	// "Channels":                channelGroupChannels,
	"Start Date": channelGroupStart,
	"End Date":   channelGroupEnd,
}

var ChannelGroupTable Table = Table{
	name:    "Channel Group",
	headers: channelGroupHeaders,
	primary: []string{"Channel Group ID", "Start Date"},
	native: []string{
		"Uncertainty In X Coordinate",
		"Uncertainty In Y Coordinate",
		"Uncertainty In Elevation",
		"Uncertainty In Strike",
		"Uncertainty In Dip",
		"First Usable Channel ID",
		"Last Usable Channel ID",
	},
	foreign: map[string][]string{},
	remap: map[string]string{
		"Channel Group ID":                 "ChannelGroupID",
		"Cable ID":                         "CableID",
		"Fiber ID":                         "FiberID",
		"Coordinate Generation Date":       "CoordinateGenerationDate",
		"Coordinate System":                "CoordinateSystem",
		"Reference Frame":                  "ReferenceFrame",
		"Location Method":                  "LocationMethod",
		"Distance Along Fiber Unit":        "DistanceAlongFiberUnit",
		"X Coordinate Unit":                "XCoordinateUnit",
		"Uncertainty In X Coordinate":      "UncertaintyInXCoordinate",
		"Uncertainty In X Coordinate Unit": "UncertaintyInXCoordinateUnit",
		"Y Coordinate Unit":                "YCoordinateUnit",
		"Uncertainty In Y Coordinate":      "UncertaintyInYCoordinate",
		"Uncertainty In Y Coordinate Unit": "UncertaintyInYCoordinateUnit",
		"Elevation Above Sea Level Unit":   "ElevationAboveSeaLevelUnit",
		"Uncertainty In Elevation":         "UncertaintyInElevation",
		"Uncertainty In Elevation Unit":    "UncertaintyInElevationUnit",
		"Strike Unit":                      "StrikeUnit",
		"Uncertainty In Strike":            "UncertaintyInStrike",
		"Dip Unit":                         "DipUnit",
		"Uncertainty In Dip":               "UncertaintyInDip",
		"Uncertainty In Dip Unit":          "UncertaintyInDipUnit",
		"First Usable Channel ID":          "FirstUsableChannelID",
		"Last Usable Channel ID":           "LastUsableChannelID",
		"Start Date":                       "Start",
		"End Date":                         "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type DASChannelGroup struct {
	Span

	ChannelGroupID           string `json:"channel_group_id"`
	AcquisitionID            string
	CableID                  string    `json:"cable_id"`
	FiberID                  string    `json:"fiber_id"`
	CoordinateGenerationDate time.Time `json:"coordinate_generation_date"`
	CoordinateSystem         string    `json:"coordinate_system"`
	ReferenceFrame           string    `json:"reference_frame"`
	LocationMethod           string    `json:"location_method"`
	// DistanceAlongFiberUnit       string    `json:"distance_along_fiber_unit"`
	XCoordinateUnit          string  `json:"x_coordinate_unit"`
	UncertaintyInXCoordinate float64 `json:"uncertainty_in_x_coordinate"`
	// UncertaintyInXCoordinateUnit string    `json:"uncertainty_in_x_coordinate_unit"`
	YCoordinateUnit          string  `json:"y_coordinate_unit"`
	UncertaintyInYCoordinate float64 `json:"uncertainty_in_y_coordinate"`
	// UncertaintyInYCoordinateUnit string    `json:"uncertainty_in_y_coordinate_unit"`
	UncertaintyInDepth float64
	// ElevationAboveSeaLevelUnit   string    `json:"elevation_above_sea_level_unit"`
	UncertaintyInElevation float64 `json:"uncertainty_in_elevation"`
	// UncertaintyInElevationUnit   string    `json:"uncertainty_in_elevation_unit"`
	StrikeUnit          string  `json:"strike_unit"`
	UncertaintyInStrike float64 `json:"uncertainty_in_strike"`
	// DipUnit             string  `json:"dip_unit"`
	UncertaintyInDip float64 `json:"uncertainty_in_dip"`
	// UncertaintyInDipUnit         string    `json:"uncertainty_in_dip_unit"`
	FirstUsableChannelID int    `json:"first_usable_channel_id"`
	LastUsableChannelID  int    `json:"last_usable_channel_id"`
	Channels             string `json:"channels"`

	uncertaintyInXCoordinateStr string // shadow value used to retain formatting
	uncertaintyInYCoordinateStr string // shadow value used to retain formatting
	uncertaintyInDepthStr       string // shadow value used to retain formatting
	uncertaintyInElevationStr   string // shadow value used to retain formatting
	uncertaintyInStrikeStr      string // shadow value used to retain formatting
	uncertaintyInDipStr         string // shadow value used to retain formatting
}

type DASChannelGroupList []DASChannelGroup

func (cgl DASChannelGroupList) Len() int      { return len(cgl) }
func (cgl DASChannelGroupList) Swap(i, j int) { cgl[i], cgl[j] = cgl[j], cgl[i] }
func (cgl DASChannelGroupList) Less(i, j int) bool {
	switch {
	case strings.ToLower(cgl[i].ChannelGroupID) < strings.ToLower(cgl[j].ChannelGroupID):
		return true
	case strings.ToLower(cgl[i].ChannelGroupID) > strings.ToLower(cgl[j].ChannelGroupID):
		return false
	default:
		return cgl[i].Start.Before(cgl[j].Start)
	}
}

func (cgl DASChannelGroupList) encode() [][]string {
	var data [][]string

	data = append(data, channelGroupHeaders.Columns())

	for _, row := range cgl {
		data = append(data, []string{
			strings.TrimSpace(row.ChannelGroupID),
			strings.TrimSpace(row.AcquisitionID),
			strings.TrimSpace(row.CableID),
			strings.TrimSpace(row.FiberID),
			row.CoordinateGenerationDate.Format(DateTimeFormat),
			strings.TrimSpace(row.CoordinateSystem),
			strings.TrimSpace(row.ReferenceFrame),
			strings.TrimSpace(row.LocationMethod),
			// strings.TrimSpace(row.DistanceAlongFiberUnit),
			strings.TrimSpace(row.XCoordinateUnit),
			strings.TrimSpace(row.uncertaintyInXCoordinateStr),
			// strings.TrimSpace(row.UncertaintyInXCoordinateUnit),
			strings.TrimSpace(row.YCoordinateUnit),
			strings.TrimSpace(row.uncertaintyInYCoordinateStr),
			// strings.TrimSpace(row.UncertaintyInYCoordinateUnit),
			strings.TrimSpace(row.uncertaintyInDepthStr),
			// strings.TrimSpace(row.ElevationAboveSeaLevelUnit),
			strings.TrimSpace(row.uncertaintyInElevationStr),
			// strings.TrimSpace(row.UncertaintyInElevationUnit),
			strings.TrimSpace(row.StrikeUnit),
			strings.TrimSpace(row.uncertaintyInStrikeStr),
			// strings.TrimSpace(row.DipUnit),
			strings.TrimSpace(row.uncertaintyInDipStr),
			// strings.TrimSpace(row.UncertaintyInDipUnit),
			strconv.Itoa(row.FirstUsableChannelID),
			strconv.Itoa(row.LastUsableChannelID),
			strings.TrimSpace(row.Channels),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (cgl *DASChannelGroupList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var channelGroups []DASChannelGroup

	fields := channelGroupHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		coordinateGenerationDate, err := time.Parse(DateTimeFormat, d[channelGroupCoordinateGenerationDate])
		if err != nil {
			return err
		}

		uncertaintyInXCoordinate, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInXCoordinate]), 64)
		if err != nil {
			return err
		}
		uncertaintyInYCoordinate, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInYCoordinate]), 64)
		if err != nil {
			return err
		}
		UncertaintyInDepth, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInDepth]), 64)
		if err != nil {
			return err
		}
		uncertaintyInElevation, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInElevation]), 64)
		if err != nil {
			return err
		}
		uncertaintyInStrike, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInStrike]), 64)
		if err != nil {
			return err
		}
		uncertaintyInDip, err := strconv.ParseFloat(strings.TrimSpace(d[channelGroupUncertaintyInDip]), 64)
		if err != nil {
			return err
		}

		firstUsableChannelID, err := strconv.Atoi(strings.TrimSpace(d[channelGroupFirstUsableChannelID]))
		if err != nil {
			return err
		}
		lastUsableChannelID, err := strconv.Atoi(strings.TrimSpace(d[channelGroupLastUsableChannelID]))
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[channelGroupStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[channelGroupEnd])
		if err != nil {
			return err
		}

		channelGroups = append(channelGroups, DASChannelGroup{
			Span: Span{
				Start: start,
				End:   end,
			},
			ChannelGroupID:           strings.TrimSpace(d[channelGroupID]),
			AcquisitionID:            strings.TrimSpace(d[ChannelAcquisitionID]),
			CableID:                  strings.TrimSpace(d[channelGroupCableID]),
			FiberID:                  strings.TrimSpace(d[channelGroupFiberID]),
			CoordinateGenerationDate: coordinateGenerationDate,
			CoordinateSystem:         strings.TrimSpace(d[channelGroupCoordinateSystem]),
			ReferenceFrame:           strings.TrimSpace(d[channelGroupReferenceFrame]),
			LocationMethod:           strings.TrimSpace(d[channelGroupLocationMethod]),
			// DistanceAlongFiberUnit:       strings.TrimSpace(d[channelGroupDistanceAlongFiberUnit]),
			// XCoordinateUnit:              strings.TrimSpace(d[channelGroupXCoordinateUnit]),
			UncertaintyInXCoordinate: uncertaintyInXCoordinate,
			// UncertaintyInXCoordinateUnit: strings.TrimSpace(d[channelGroupUncertaintyInXCoordinateUnit]),
			// YCoordinateUnit:              strings.TrimSpace(d[channelGroupYCoordinateUnit]),
			UncertaintyInYCoordinate: uncertaintyInYCoordinate,
			// UncertaintyInYCoordinateUnit: strings.TrimSpace(d[channelGroupUncertaintyInYCoordinateUnit]),
			UncertaintyInDepth: UncertaintyInDepth,
			// ElevationAboveSeaLevelUnit:   strings.TrimSpace(d[channelGroupElevationAboveSeaLevelUnit]),
			UncertaintyInElevation: uncertaintyInElevation,
			// UncertaintyInElevationUnit: strings.TrimSpace(d[channelGroupUncertaintyInElevationUnit]),
			// StrikeUnit:                   strings.TrimSpace(d[channelGroupStrikeUnit]),
			UncertaintyInStrike: uncertaintyInStrike,
			// DipUnit:                      strings.TrimSpace(d[channelGroupDipUnit]),
			UncertaintyInDip: uncertaintyInDip,
			// UncertaintyInDipUnit:         strings.TrimSpace(d[channelGroupUncertaintyInDipUnit]),
			FirstUsableChannelID: firstUsableChannelID,
			LastUsableChannelID:  lastUsableChannelID,
			// Channels:                     strings.TrimSpace(d[channelGroupChannels]),

			uncertaintyInXCoordinateStr: strings.TrimSpace(d[channelGroupUncertaintyInXCoordinate]),
			uncertaintyInYCoordinateStr: strings.TrimSpace(d[channelGroupUncertaintyInYCoordinate]),
			uncertaintyInElevationStr:   strings.TrimSpace(d[channelGroupUncertaintyInElevation]),
			uncertaintyInStrikeStr:      strings.TrimSpace(d[channelGroupUncertaintyInStrike]),
			uncertaintyInDipStr:         strings.TrimSpace(d[channelGroupUncertaintyInDip]),
		})
		// fmt.Print("succeded parsing Channel Groups", "\n")
	}

	*cgl = DASChannelGroupList(channelGroups)

	return nil
}

func LoadChannelGroups(path string) ([]DASChannelGroup, error) {
	var cgl []DASChannelGroup

	if err := LoadList(path, (*DASChannelGroupList)(&cgl)); err != nil {
		return nil, err
	}

	sort.Sort(DASChannelGroupList(cgl))

	return cgl, nil
}

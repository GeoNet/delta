package meta

import (
	"sort"
	"strconv"
	"strings"
)

const (
	dasChannelID int = iota
	dasChannelGroupID
	dasChannelDistanceAlongFiber
	dasChannelXCoordinate
	dasChannelYCoordinate
	dasChannelElevation
	dasChannelDepth
	dasChannelStrike
	dasChannelDip
	dasChannelStart
	dasChannelEnd
)

var dasChannelHeaders Header = map[string]int{
	"Channel ID":           dasChannelID,
	"Channel Group ID":     dasChannelGroupID,
	"Distance Along Fiber": dasChannelDistanceAlongFiber,
	"X Coordinate":         dasChannelXCoordinate,
	"Y Coordinate":         dasChannelYCoordinate,
	"Elevation":            dasChannelElevation,
	"Depth":                dasChannelDepth,
	"Strike":               dasChannelStrike,
	"Dip":                  dasChannelDip,
	"Start Date":           dasChannelStart,
	"End Date":             dasChannelEnd,
}

var DASChannelTable Table = Table{
	name:    "DAS Channel",
	headers: dasChannelHeaders,
	primary: []string{"Channel ID"},
	native:  []string{},
	foreign: map[string][]string{},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type DASChannel struct {
	ChannelID              string `json:"channel_id"`
	ChannelGroupID         string
	DistanceAlongFiber     float64 `json:"distance_along_fiber"`
	XCoordinate            float64 `json:"x_coordinate"`
	YCoordinate            float64 `json:"y_coordinate"`
	ElevationAboveSeaLevel float64 `json:"elevation_above_sea_level"`
	DepthBelowSurface      float64 `json:"depth_below_surface"`
	Strike                 float64 `json:"strike"`
	Dip                    float64 `json:"dip"`

	distanceAlongFiber     string // shadow value used to retain formatting
	xCoordinate            string // shadow value used to retain formatting
	yCoordinate            string // shadow value used to retain formatting
	elevationAboveSeaLevel string // shadow value used to retain formatting
	depthBelowSurface      string // shadow value used to retain formatting
	strike                 string // shadow value used to retain formatting
	dip                    string // shadow value used to retain formatting
}

type DASChannelList []DASChannel

func (cl DASChannelList) Len() int      { return len(cl) }
func (cl DASChannelList) Swap(i, j int) { cl[i], cl[j] = cl[j], cl[i] }
func (cl DASChannelList) Less(i, j int) bool {
	return strings.ToLower(cl[i].ChannelID) < strings.ToLower(cl[j].ChannelID)
}

func (cl DASChannelList) encode() [][]string {
	var data [][]string

	data = append(data, dasChannelHeaders.Columns())

	for _, row := range cl {
		data = append(data, []string{
			strings.TrimSpace(row.ChannelID),
			strings.TrimSpace(row.ChannelGroupID),
			strings.TrimSpace(row.distanceAlongFiber),
			strings.TrimSpace(row.xCoordinate),
			strings.TrimSpace(row.yCoordinate),
			strings.TrimSpace(row.elevationAboveSeaLevel),
			strings.TrimSpace(row.depthBelowSurface),
			strings.TrimSpace(row.strike),
			strings.TrimSpace(row.dip),
		})
	}

	return data
}

func (cl *DASChannelList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var channels []DASChannel

	fields := dasChannelHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)
		// fmt.Print(channelGroupHeaders, "\n")
		// fmt.Print(fields, "\n")
		// fmt.Print(row, "\n")
		// fmt.Print(d, "\n")

		distanceAlongFiber, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelDistanceAlongFiber]), 64)
		if err != nil {
			return err
		}
		xCoordinate, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelXCoordinate]), 64)
		if err != nil {
			return err
		}
		yCoordinate, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelYCoordinate]), 64)
		if err != nil {
			return err
		}
		elevationAboveSeaLevel, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelElevation]), 64)
		if err != nil {
			return err
		}
		depthBelowSurface, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelDepth]), 64)
		if err != nil {
			return err
		}
		strike, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelStrike]), 64)
		if err != nil {
			return err
		}
		dip, err := strconv.ParseFloat(strings.TrimSpace(d[dasChannelDip]), 64)
		if err != nil {
			return err
		}

		channels = append(channels, DASChannel{
			ChannelID:              strings.TrimSpace(d[dasChannelID]),
			ChannelGroupID:         strings.TrimSpace(d[dasChannelGroupID]),
			DistanceAlongFiber:     distanceAlongFiber,
			XCoordinate:            xCoordinate,
			YCoordinate:            yCoordinate,
			ElevationAboveSeaLevel: elevationAboveSeaLevel,
			DepthBelowSurface:      depthBelowSurface,
			Strike:                 strike,
			Dip:                    dip,

			distanceAlongFiber:     strings.TrimSpace(d[dasChannelDistanceAlongFiber]),
			xCoordinate:            strings.TrimSpace(d[dasChannelXCoordinate]),
			yCoordinate:            strings.TrimSpace(d[dasChannelYCoordinate]),
			elevationAboveSeaLevel: strings.TrimSpace(d[dasChannelElevation]),
			depthBelowSurface:      strings.TrimSpace(d[dasChannelDepth]),
			strike:                 strings.TrimSpace(d[dasChannelStrike]),
			dip:                    strings.TrimSpace(d[dasChannelDip]),
		})
	}

	*cl = DASChannelList(channels)

	return nil
}

func LoadDASChannels(path string) ([]DASChannel, error) {
	var cl []DASChannel

	if err := LoadList(path, (*DASChannelList)(&cl)); err != nil {
		return nil, err
	}

	sort.Sort(DASChannelList(cl))

	return cl, nil
}

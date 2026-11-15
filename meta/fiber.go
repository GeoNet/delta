package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	fiberID int = iota
	fiberCableID
	fiberGeometry
	fiberMode
	fiberRefractionIndex
	fiberWindingAngle
	// fiberWindingAngleUnit
	fiberStartLocation
	// fiberStartLocationUnit
	fiberEndLocation
	// fiberEndLocationUnit
	fiberOpticLength
	// fiberOpticLengthUnit
	fiberOneWayAttenuation
	// fiberOneWayAttenuationUnit
	fiberStart
	fiberEnd
	fiberLast
)

var fiberHeaders Header = map[string]int{
	"Fiber ID":               fiberID,
	"Cable ID":               fiberCableID,
	"Fiber Geometry":         fiberGeometry,
	"Fiber Mode":             fiberMode,
	"Fiber Refraction Index": fiberRefractionIndex,
	"Fiber Winding Angle":    fiberWindingAngle,
	// "Fiber Winding Angle Unit":       fiberWindingAngleUnit,
	"Fiber Start Location": fiberStartLocation,
	// "Fiber Start Location Unit":      fiberStartLocationUnit,
	"Fiber End Location": fiberEndLocation,
	// "Fiber End Location Unit":        fiberEndLocationUnit,
	"Fiber Optic Length": fiberOpticLength,
	// "Fiber Optic Length Unit":        fiberOpticLengthUnit,
	"Fiber One Way Attenuation": fiberOneWayAttenuation,
	// "Fiber One Way Attenuation Unit": fiberOneWayAttenuationUnit,
	"Start Date": fiberStart,
	"End Date":   fiberEnd,
}

var FiberTable Table = Table{
	name:    "Fiber",
	headers: fiberHeaders,
	primary: []string{"Fiber ID", "Cable ID", "Start Date"},
	native: []string{
		"Fiber Refraction Index",
		"Fiber Winding Angle",
		"Fiber Start Location",
		"Fiber End Location",
		"Fiber Optic Length",
		"Fiber One Way Attenuation",
	},
	foreign: map[string][]string{},
	remap: map[string]string{
		// "Fiber ID":                       "FiberID",
		// "Fiber Geometry":                 "Geometry",
		// "Fiber Mode":                     "Mode",
		// "Fiber Refraction Index":         "RefractionIndex",
		// "Fiber Winding Angle":            "WindingAngle",
		// "Fiber Winding Angle Unit":       "WindingAngleUnit",
		// "Fiber Start Location":           "StartLocation",
		// "Fiber Start Location Unit":      "StartLocationUnit",
		// "Fiber End Location":             "EndLocation",
		// "Fiber End Location Unit":        "EndLocationUnit",
		// "Fiber Optic Length":             "OpticLength",
		// "Fiber Optic Length Unit":        "OpticLengthUnit",
		// "Fiber One Way Attenuation":      "OneWayAttenuation",
		// "Fiber One Way Attenuation Unit": "OneWayAttenuationUnit",
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledFiber struct {
	Span

	FiberID               string
	CableID               string
	Geometry              string
	Mode                  string
	RefractionIndex       float64
	WindingAngle          float64
	WindingAngleUnit      string
	StartLocation         float64
	StartLocationUnit     string
	EndLocation           float64
	EndLocationUnit       string
	OpticLength           float64
	OpticLengthUnit       string
	OneWayAttenuation     float64
	OneWayAttenuationUnit string

	refractionIndex   string // shadow value used to retain formatting
	windingAngle      string // shadow value used to retain formatting
	startLocation     string // shadow value used to retain formatting
	endLocation       string // shadow value used to retain formatting
	opticLength       string // shadow value used to retain formatting
	oneWayAttenuation string // shadow value used to retain formatting
}

type FiberList []InstalledFiber

func (fl FiberList) Len() int      { return len(fl) }
func (fl FiberList) Swap(i, j int) { fl[i], fl[j] = fl[j], fl[i] }
func (fl FiberList) Less(i, j int) bool {
	switch {
	case strings.ToLower(fl[i].FiberID) < strings.ToLower(fl[j].FiberID):
		return true
	case strings.ToLower(fl[i].FiberID) > strings.ToLower(fl[j].FiberID):
		return false
	default:
		return fl[i].Start.Before(fl[j].Start)
	}
}

func (fl FiberList) encode() [][]string {
	var data [][]string

	data = append(data, fiberHeaders.Columns())

	for _, row := range fl {
		data = append(data, []string{
			strings.TrimSpace(row.FiberID),
			strings.TrimSpace(row.CableID),
			strings.TrimSpace(row.Geometry),
			strings.TrimSpace(row.Mode),
			strings.TrimSpace(row.refractionIndex),
			strings.TrimSpace(row.windingAngle),
			strings.TrimSpace(row.WindingAngleUnit),
			strings.TrimSpace(row.startLocation),
			strings.TrimSpace(row.StartLocationUnit),
			strings.TrimSpace(row.endLocation),
			strings.TrimSpace(row.EndLocationUnit),
			strings.TrimSpace(row.opticLength),
			strings.TrimSpace(row.OpticLengthUnit),
			strings.TrimSpace(row.oneWayAttenuation),
			strings.TrimSpace(row.OneWayAttenuationUnit),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (fl *FiberList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var fibers []InstalledFiber

	fields := fiberHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		refractionIndex, err := strconv.ParseFloat(strings.TrimSpace(d[fiberRefractionIndex]), 64)
		if err != nil {
			return err
		}
		windingAngle, err := strconv.ParseFloat(strings.TrimSpace(d[fiberWindingAngle]), 64)
		if err != nil {
			return err
		}
		startLocation, err := strconv.ParseFloat(strings.TrimSpace(d[fiberStartLocation]), 64)
		if err != nil {
			return err
		}
		endLocation, err := strconv.ParseFloat(strings.TrimSpace(d[fiberEndLocation]), 64)
		if err != nil {
			return err
		}
		opticLength, err := strconv.ParseFloat(strings.TrimSpace(d[fiberOpticLength]), 64)
		if err != nil {
			return err
		}
		oneWayAttenuation, err := strconv.ParseFloat(strings.TrimSpace(d[fiberOneWayAttenuation]), 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[fiberStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[fiberEnd])
		if err != nil {
			return err
		}

		fibers = append(fibers, InstalledFiber{
			Span: Span{
				Start: start,
				End:   end,
			},
			FiberID:         strings.TrimSpace(d[fiberID]),
			CableID:         strings.TrimSpace(d[fiberCableID]),
			Geometry:        strings.TrimSpace(d[fiberGeometry]),
			Mode:            strings.TrimSpace(d[fiberMode]),
			RefractionIndex: refractionIndex,
			WindingAngle:    windingAngle,
			// WindingAngleUnit:      strings.TrimSpace(d[fiberWindingAngleUnit]),
			StartLocation: startLocation,
			// StartLocationUnit:     strings.TrimSpace(d[fiberStartLocationUnit]),
			EndLocation: endLocation,
			// EndLocationUnit:       strings.TrimSpace(d[fiberEndLocationUnit]),
			OpticLength: opticLength,
			// OpticLengthUnit:       strings.TrimSpace(d[fiberOpticLengthUnit]),
			OneWayAttenuation: oneWayAttenuation,
			// OneWayAttenuationUnit: strings.TrimSpace(d[fiberOneWayAttenuationUnit]),

			refractionIndex:   strings.TrimSpace(d[fiberRefractionIndex]),
			windingAngle:      strings.TrimSpace(d[fiberWindingAngle]),
			startLocation:     strings.TrimSpace(d[fiberStartLocation]),
			endLocation:       strings.TrimSpace(d[fiberEndLocation]),
			opticLength:       strings.TrimSpace(d[fiberOpticLength]),
			oneWayAttenuation: strings.TrimSpace(d[fiberOneWayAttenuation]),
		})
	}

	*fl = FiberList(fibers)

	return nil
}

func LoadFibers(path string) ([]InstalledFiber, error) {
	var fl []InstalledFiber

	if err := LoadList(path, (*FiberList)(&fl)); err != nil {
		return nil, err
	}

	sort.Sort(FiberList(fl))

	return fl, nil
}

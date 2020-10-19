package meta

import (
	"time"
)

// DateTimeFormat is the standard date and time storage format used
// in the CSV files, it is assumed to have resolution of one second
// and that times are in UTC.
const DateTimeFormat = "2006-01-02T15:04:05Z"

// Reference describes a location where measurements can be taken.
type Reference struct {
	// Code is used to identify the measurement location.
	Code string
	// Network can be used to group multiple measurement locations.
	Network string
	// Name is used to label the measurement location.
	Name string
}

// Point describes a measurement location geographically.
type Point struct {
	// Latitued represents the location latitude, with negative values representing southern latitudes.
	Latitude float64
	// Longitude represents the location longitude, with negative values representing western longitudes.
	Longitude float64
	// Elevation represents the location height relative to the given datum.
	Elevation float64
	// Datum can be used to indicate the location measurement reference.
	Datum string

	latitude  string // shadow value used to retain formatting
	longitude string // shadow value used to retain formatting
	elevation string // shadow value used to retain formatting
}

// Orientation is used to describe how a piece of installed equipment is aligned.
type Orientation struct {
	// Dip represents the vertical deployment, with a zero value representing a horizontal installation,
	// a positive value indicating a installation downwards, whereas a negative value indicates an upward
	// facing installation.
	Dip float64
	// Azimuth represents an equipment installation bearing, ideally with reference to true north.
	Azimuth float64
	// Method can be used to indicate the method or measuring the azimuth.
	Method string

	dip     string // shadow value used to retain formatting
	azimuth string // shadow value used to retain formatting
}

// Offset can be used to adjust an equipment installation relative to a given Point.
type Offset struct {
	// Vertical represents an adjustment up or down, the exact interpretation will depend on the use case,
	// although it is assumed to have units of meters.
	Vertical float64
	// North can be used to offset the installation to northwards, it is asusmed to have units of meters.
	North float64
	// East can be used to offset the installation to eastwards, it is asusmed to have units of meters.
	East float64

	vertical string // shadow value used to retain formatting
	north    string // shadow value used to retain formatting
	east     string // shadow value used to retain formatting
}

// Scale can be used to represent a non-linear installation, such as a pressure sensor installed in sea water
// rather than fresh water.
type Scale struct {
	// Factor can be used to represent a change of scale of the recorded value.
	Factor float64
	// Bias can be used to represent an offset to the recorded value.
	Bias float64

	factor string // shadow value used to retain formatting
	bias   string // shadow value used to retain formatting
}

// Span represents a time window.
type Span struct {
	// Start represents the beginning of the time window.
	Start time.Time
	// End represents the stop time of the window.
	End time.Time
}

// Equipment represents an indiviual piece of hardware.
type Equipment struct {
	// Make describes the manufacturer or equipment maker.
	Make string
	// Model describes the manufacturer's model name.
	Model string
	// Serial describes the manufacturer's identification of the device.
	Serial string
}

func (e Equipment) String() string {
	return e.Make + " " + e.Model + " [" + e.Serial + "]"
}

// Less compares Equipment structs suitable for sorting.
func (e Equipment) Less(eq Equipment) bool {

	switch {
	case e.Make < eq.Make:
		return true
	case e.Make > eq.Make:
		return false
	case e.Model < eq.Model:
		return true
	case e.Model > eq.Model:
		return false
	}

	return e.Serial < eq.Serial
}

// Install is a compounded struct the represents the installation of a
// piece of equipment over a given time period.
type Install struct {
	// Equipment respresents the actual installed equipment.
	Equipment
	// Span describes the installed time period.
	Span
}

// Less compares Install structs suitable for sorting.
func (i Install) Less(in Install) bool {
	switch {
	case i.Equipment.Less(in.Equipment):
		return true
	case in.Equipment.Less(i.Equipment):
		return false
	default:
		return i.Start.Before(in.Start)
	}
}

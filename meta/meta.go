package meta

import (
	"strconv"
	"strings"
	"time"
)

// DateTimeFormat is the standard date and time storage format used
// in the CSV files, it is assumed to have resolution of one second
// and that times are in UTC.
const DateTimeFormat = "2006-01-02T15:04:05Z"

// Format outputs a Time using the DateTimeFormat format.
func Format(t time.Time) string {
	return t.Format(DateTimeFormat)
}

type Compare int

const (
	EqualTo Compare = iota
	LessThan
	GreaterThan
)

// Reference describes a location where measurements can be taken.
type Reference struct {
	// Code is used to identify the measurement location.
	Code string `json:"code"`
	// Network can be used to group multiple measurement locations.
	Network string `json:"network"`
	// Name is used to label the measurement location.
	Name string `json:"name"`
}

// Position describes a measurement location geographically.
type Position struct {
	// Latitude represents the location latitude, with negative values representing southern latitudes.
	Latitude float64 `json:"latitude"`
	// Longitude represents the location longitude, with negative values representing western longitudes.
	Longitude float64 `json:"longitude"`
	// Elevation represents the location height relative to the given datum.
	Elevation float64 `json:"elevation"`
	// Datum can be used to indicate the location measurement reference.
	Datum string `json:"datum,omitempty"`
	// Depth measures the depth of water at the measurement position, if appropriate.
	Depth float64 `json:"depth"`

	latitude  string // shadow value used to retain formatting
	longitude string // shadow value used to retain formatting
	elevation string // shadow value used to retain formatting
	depth     string // shadow value used to retain formatting
}

// ElevationOk returns the Elevation and whether it has been set.
func (p Position) ElevationOk() (float64, bool) {
	if p.elevation != "" {
		return p.Elevation, true
	}
	return 0.0, false
}

// DepthOk returns the Depth and whether it has been set.
func (p Position) DepthOk() (float64, bool) {
	if p.depth != "" {
		return p.Depth, true
	}
	return 0.0, false
}

// Orientation is used to describe how a piece of installed equipment is aligned.
type Orientation struct {
	// Dip represents the vertical deployment, with a zero value representing a horizontal installation,
	// a positive value indicating a installation downwards, whereas a negative value indicates an upward
	// facing installation.
	Dip float64 `json:"dip"`
	// Azimuth represents an equipment installation bearing, ideally with reference to true north.
	Azimuth float64 `json:"azimuth"`
	// Method can be used to indicate the method or measuring the azimuth.
	Method string `json:"method,omitempty"`

	dip     string // shadow value used to retain formatting
	azimuth string // shadow value used to retain formatting
}

// DipOk returns the Dip and whether it has been set.
func (o Orientation) DipOk() (float64, bool) {
	if o.dip != "" {
		return o.Dip, true
	}
	return 0.0, false
}

// AzimuthOk returns the Azimuth and whether it has been set.
func (o Orientation) AzimuthOk() (float64, bool) {
	if o.azimuth != "" {
		return o.Azimuth, true
	}
	return 0.0, false
}

// Offset can be used to adjust an equipment installation relative to a given Position.
type Offset struct {
	// Vertical represents an adjustment up or down, the exact interpretation will depend on the use case,
	// although it is assumed to have units of meters.
	Vertical float64 `json:"vertical,omitempty"`
	// North can be used to offset the installation to northwards, it is asusmed to have units of meters.
	North float64 `json:"north,omitempty"`
	// East can be used to offset the installation to eastwards, it is asusmed to have units of meters.
	East float64 `json:"east,omitempty"`

	vertical string // shadow value used to retain formatting
	north    string // shadow value used to retain formatting
	east     string // shadow value used to retain formatting
}

// VerticalOk returns the Vertical offset and whether it has been set.
func (o Offset) VerticalOk() (float64, bool) {
	if o.vertical != "" {
		return o.Vertical, true
	}
	return 0.0, false
}

// NorthOk returns the North offset and whether it has been set.
func (o Offset) NorthOk() (float64, bool) {
	if o.north != "" {
		return o.North, true
	}
	return 0.0, false
}

// EastOk returns the East offset and whether it has been set.
func (o Offset) EastOk() (float64, bool) {
	if o.east != "" {
		return o.East, true
	}
	return 0.0, false
}

// Scale can be used to represent a non-linear installation, such as a pressure sensor installed in sea water
// rather than fresh water.
type Scale struct {
	// Factor can be used to represent a change of scale of the recorded value.
	Factor float64 `json:"factor"`
	// Bias can be used to represent an offset to the recorded value.
	Bias float64 `json:"bias,omitempty"`

	factor string // shadow value used to retain formatting
	bias   string // shadow value used to retain formatting
}

// FactorOk returns the Factor and whether it has been set.
func (s Scale) FactortOk() (float64, bool) {
	if s.factor != "" {
		return s.Factor, true
	}
	return 0.0, false
}

// BiasOk returns the Bias and whether it has been set.
func (s Scale) BiasOk() (float64, bool) {
	if s.bias != "" {
		return s.Bias, true
	}
	return 0.0, false
}

// Span represents a time window.
type Span struct {
	// Start represents the beginning of the time window.
	Start time.Time `json:"start"`

	// End represents the stop time of the window.
	End time.Time `json:"end"`
}

// Overlaps returns whether the time range of the Span overlaps the given Span.
// It is assumed that the End time for a span can overlap with the start of
// another Span if they exactly match.
func (s Span) Overlaps(span Span) bool {
	switch {
	case s.Start.After(span.End):
		return false
	case s.End.Before(span.Start):
		return false
	default:
		return true
	}
}

// Extent returns the Span that is the sum of the given overlapping Span values,
// the extra boolean return value will be false if no window could be found. It
// is assumed that the end must be greater than the start of the resultant Span.
func (s Span) Extent(spans ...Span) (Span, bool) {

	clip := s

	for _, span := range spans {
		if span.Start.Before(clip.Start) {
			continue
		}
		clip.Start = span.Start
	}
	for _, span := range spans {
		if span.End.After(clip.End) {
			continue
		}
		clip.End = span.End
	}

	if clip.Start.After(clip.End) {
		return Span{}, false
	}

	if !clip.Overlaps(s) {
		return Span{}, false
	}

	for _, span := range spans {
		if span.Overlaps(clip) {
			continue
		}
		return Span{}, false
	}

	if clip.Start.Equal(clip.End) {
		return Span{}, false
	}

	return clip, true
}

type Range struct {
	Value   float64
	Compare Compare
}

func NewRange(s string) (Range, error) {
	switch {
	case strings.HasPrefix(s, "<"):
		v, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value:   v,
			Compare: LessThan,
		}, nil
	case strings.HasPrefix(s, ">"):
		v, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value:   v,
			Compare: GreaterThan,
		}, nil
	default:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value: v,
		}, nil
	}
}

func (r Range) String() string {
	switch r.Compare {
	case LessThan:
		return "<" + strconv.FormatFloat(r.Value, 'g', -1, 64)
	case GreaterThan:
		return ">" + strconv.FormatFloat(r.Value, 'g', -1, 64)
	default:
		return strconv.FormatFloat(r.Value, 'g', -1, 64)
	}
}

// Equipment represents an indiviual piece of hardware.
type Equipment struct {
	// Make describes the manufacturer or equipment maker.
	Make string `json:"make"`
	// Model describes the manufacturer's model name.
	Model string `json:"model"`
	// Serial describes the manufacturer's identification of the device.
	Serial string `json:"serial"`
}

func (e Equipment) String() string {
	return e.Make + " " + e.Model + " [" + e.Serial + "]"
}

// Less compares Equipment structs suitable for sorting.
func (e Equipment) Less(eq Equipment) bool {

	switch {
	case strings.ToLower(e.Make) < strings.ToLower(eq.Make):
		return true
	case strings.ToLower(e.Make) > strings.ToLower(eq.Make):
		return false
	case strings.ToLower(e.Model) < strings.ToLower(eq.Model):
		return true
	case strings.ToLower(e.Model) > strings.ToLower(eq.Model):
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

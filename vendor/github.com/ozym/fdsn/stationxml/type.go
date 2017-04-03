package stationxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// The type of data this channel collects. Corresponds to
// channel flags in SEED blockette 52. The SEED volume producer could
// use the first letter of an Output value as the SEED channel flag.
type Type uint

const (
	TypeUnknown Type = iota
	TypeTriggered
	TypeContinuous
	TypeHealth
	TypeGeophysical
	TypeWeather
	TypeFlag
	TypeSynthesized
	TypeInput
	TypeExperimental
	TypeMaintenance
	TypeBeam
	TypeEnd
)

var typeLookup = [...]string{
	TypeUnknown:      "UNKNOWN",
	TypeTriggered:    "TRIGGERED",
	TypeContinuous:   "CONTINUOUS",
	TypeHealth:       "HEALTH",
	TypeGeophysical:  "GEOPHYSICAL",
	TypeWeather:      "WEATHER",
	TypeFlag:         "FLAG",
	TypeSynthesized:  "SYNTHESIZED",
	TypeInput:        "INPUT",
	TypeExperimental: "EXPERIMENTAL",
	TypeMaintenance:  "MAINTENANCE",
	TypeBeam:         "BEAM",
}

var typeMap = map[string]Type{
	"UNKNOWN":      TypeUnknown,
	"TRIGGERED":    TypeTriggered,
	"CONTINUOUS":   TypeContinuous,
	"HEALTH":       TypeHealth,
	"GEOPHYSICAL":  TypeGeophysical,
	"WEATHER":      TypeWeather,
	"FLAG":         TypeFlag,
	"SYNTHESIZED":  TypeSynthesized,
	"INPUT":        TypeInput,
	"EXPERIMENTAL": TypeExperimental,
	"MAINTENANCE":  TypeMaintenance,
	"BEAM":         TypeBeam,
}

func (t Type) String() string {

	if t < TypeEnd {
		return typeLookup[t]
	}

	return ""
}

func (t Type) IsValid() error {

	if !(t < TypeEnd) {
		return fmt.Errorf("invalid type: %d", t)
	}

	return nil
}

func (t Type) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !(t < TypeEnd) {
		return fmt.Errorf("invalid type: %d", t)
	}
	return e.EncodeElement(typeLookup[t], start)
}

func (t *Type) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	if _, ok := typeMap[s]; !ok {
		return fmt.Errorf("invalid type: %s", s)
	}

	*t = typeMap[s]

	return nil
}

func (t Type) MarshalJSON() ([]byte, error) {
	if !(t < TypeEnd) {
		return nil, fmt.Errorf("invalid type: %d", t)
	}
	return json.Marshal(typeLookup[t])
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	s := string(b)

	if _, ok := typeMap[s]; !ok {
		return fmt.Errorf("invalid type: %s", s)
	}

	*t = typeMap[s]

	return nil
}

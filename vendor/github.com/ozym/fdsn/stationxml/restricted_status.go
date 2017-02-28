package stationxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type RestrictedStatus uint

const (
	StatusUnknown RestrictedStatus = iota
	StatusOpen
	StatusClosed
	StatusPartial
	statusEnd
)

var restrictedStatus = [...]string{
	StatusUnknown: "unknown",
	StatusOpen:    "open",
	StatusClosed:  "closed",
	StatusPartial: "partial",
}

var restrictedStatusMap = map[string]RestrictedStatus{
	"unknown": StatusUnknown,
	"open":    StatusOpen,
	"closed":  StatusClosed,
	"partial": StatusPartial,
}

func (r RestrictedStatus) String() string {

	if r < statusEnd {
		return restrictedStatus[r]
	}

	return ""
}

func (r RestrictedStatus) IsValid() error {

	if !(r < statusEnd) {
		return fmt.Errorf("invalid restricted value: %d", r)
	}

	return nil
}

func (r RestrictedStatus) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if !(r < statusEnd) {
		return xml.Attr{}, fmt.Errorf("invalid restricted value: %d", r)
	}
	return xml.Attr{Name: name, Value: restrictedStatus[r]}, nil
}

func (r *RestrictedStatus) UnmarshalXMLAttr(attr xml.Attr) error {

	if _, ok := restrictedStatusMap[attr.Value]; !ok {
		return fmt.Errorf("invalid restricted value: %s", attr.Value)
	}

	*r = restrictedStatusMap[attr.Value]

	return nil
}

func (r RestrictedStatus) MarshalJSON() ([]byte, error) {
	if !(r < statusEnd) {
		return nil, fmt.Errorf("invalid restricted value: %d", r)
	}
	return json.Marshal(restrictedStatus[r])
}

func (r *RestrictedStatus) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	s := string(b)

	if _, ok := restrictedStatusMap[s]; !ok {
		return fmt.Errorf("invalid restricted value: %s", s)
	}

	*r = restrictedStatusMap[s]

	return nil
}

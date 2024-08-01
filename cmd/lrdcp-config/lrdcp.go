package main

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/GeoNet/delta/meta"
)

type Lrdcp struct {
	Features []meta.Feature
	Sensors  []meta.InstalledSensor
	Gains    []meta.Gain
}

func (l *Lrdcp) Load(set *meta.Set) error {

	l.Features = set.Features()
	l.Sensors = make([]meta.InstalledSensor, 0)

	// only include the sensor/gain when the station is defined in features
	for _, f := range l.Features {
		// check sensors
		for _, s := range set.InstalledSensors() {
			if s.Station == f.Station && s.Location == f.Location {
				l.Sensors = append(l.Sensors, s)
				break
			}
		}
		// check gains
		for _, g := range set.Gains() {
			if g.Station == f.Station && g.Location == f.Location && g.Sublocation == f.Sublocation {
				l.Gains = append(l.Gains, g)
				break
			}
		}
	}

	return nil
}

func (l *Lrdcp) Marshal(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(l); err != nil {
		return err
	}
	return nil
}

func (l *Lrdcp) MarshalIndent(wr io.Writer, prefix, indent string) error {
	enc := xml.NewEncoder(wr)
	enc.Indent(prefix, indent)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(l); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(wr, "\n"); err != nil {
		return err
	}
	return nil
}

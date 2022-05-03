package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	codenameType = iota
	codenameSamplingRate
	codenameDip
	codenameAzimuth
	codenameCode
	codenameAxial
	codenameFlags
	codenameLast
)

// Codename is used to describe a channel stream.
type Codename struct {
	Type  string
	Code  string
	Axial string
	Flags string

	SamplingRate float64
	Dip          float64
	Azimuth      float64

	samplingRate string
	dip          string
	azimuth      string
	axial        string
}

// Less compares Codename structs suitable for sorting.
func (c Codename) Less(codename Codename) bool {

	switch {
	case strings.ToLower(c.Type) < strings.ToLower(codename.Type):
		return true
	case strings.ToLower(c.Type) > strings.ToLower(codename.Type):
		return false
	case c.SamplingRate < codename.SamplingRate:
		return true
	case c.SamplingRate > codename.SamplingRate:
		return false
	case c.Dip < codename.Dip:
		return true
	case c.Dip > codename.Dip:
		return false
	case c.Azimuth < codename.Azimuth:
		return true
	default:
		return false
	}
}

type CodenameList []Codename

func (s CodenameList) Len() int           { return len(s) }
func (s CodenameList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s CodenameList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s CodenameList) encode() [][]string {
	data := [][]string{{
		"Type",
		"Sampling Rate",
		"Dip",
		"Azimuth",
		"Code",
		"Axial",
		"Flags",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.samplingRate),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.axial),
			strings.TrimSpace(v.Flags),
		})
	}
	return data
}

func (s *CodenameList) decode(data [][]string) error {
	var codenames []Codename
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != codenameLast {
				return fmt.Errorf("incorrect number of codename fields")
			}

			rate, err := strconv.ParseFloat(d[codenameSamplingRate], 64)
			if err != nil {
				return err
			}
			if rate < 0 {
				rate = -1.0 / rate
			}

			dip, err := strconv.ParseFloat(d[codenameDip], 64)
			if err != nil {
				return err
			}
			azimuth, err := strconv.ParseFloat(d[codenameAzimuth], 64)
			if err != nil {
				return err
			}

			axial := strings.TrimSpace(d[codenameCode])
			if s := strings.TrimSpace(d[codenameAxial]); s != "" {
				axial = s
			}

			codenames = append(codenames, Codename{
				Type:  strings.TrimSpace(d[codenameType]),
				Code:  strings.TrimSpace(d[codenameCode]),
				Flags: strings.TrimSpace(d[codenameFlags]),
				Axial: axial,

				SamplingRate: rate,
				Dip:          dip,
				Azimuth:      azimuth,

				samplingRate: strings.TrimSpace(d[codenameSamplingRate]),
				dip:          strings.TrimSpace(d[codenameDip]),
				azimuth:      strings.TrimSpace(d[codenameAzimuth]),
				axial:        strings.TrimSpace(d[codenameAxial]),
			})
		}

		*s = CodenameList(codenames)
	}
	return nil
}

func LoadCodenames(path string) ([]Codename, error) {
	var s []Codename

	if err := LoadList(path, (*CodenameList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(CodenameList(s))

	return s, nil
}

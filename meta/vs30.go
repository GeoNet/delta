package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	vs30Station = iota
	vs30Vs30
	vs30Tsite
	vs30Zb
	vs30QVs30
	vs30QTsite
	vs30DTsite
	vs30QZb
	vs30References
	vs30Notes
	vs30Last
)

type Vs30 struct {
	Station    string
	Vs30       float64
	Tsite      Range
	Zb         float64
	QVs30      string
	QTsite     string
	DTsite     string
	QZb        string
	References string
	Notes      string
}

func (c Vs30) Less(vs30 Vs30) bool {
	switch {
	case c.Station < vs30.Station:
		return true
	default:
		return false
	}
}

type Vs30List []Vs30

func (c Vs30List) Len() int           { return len(c) }
func (c Vs30List) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Vs30List) Less(i, j int) bool { return c[i].Less(c[j]) }

func (c Vs30List) encode() [][]string {
	data := [][]string{{
		"Station",
		"Vs30",
		"Tsite",
		"Zb",
		"Q_Vs30",
		"Q_Tsite",
		"D_Tsite",
		"Q_Zb",
		"References",
		"Notes",
	}}
	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strconv.FormatFloat(v.Vs30, 'g', -1, 64),
			strings.TrimSpace(v.Tsite.String()),
			strconv.FormatFloat(v.Zb, 'g', -1, 64),
			strings.TrimSpace(v.QVs30),
			strings.TrimSpace(v.QTsite),
			strings.TrimSpace(v.DTsite),
			strings.TrimSpace(v.QZb),
			strings.TrimSpace(v.References),
			strings.TrimSpace(v.Notes),
		})
	}
	return data
}

func (c *Vs30List) decode(data [][]string) error {
	var vs30s []Vs30
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != vs30Last {
				return fmt.Errorf("incorrect number of installed vs30 fields")
			}
			var err error

			var vs30, zb float64
			if vs30, err = strconv.ParseFloat(d[vs30Vs30], 64); err != nil {
				return err
			}
			if zb, err = strconv.ParseFloat(d[vs30Zb], 64); err != nil {
				return err
			}

			var r Range
			if r, err = NewRange(d[vs30Tsite]); err != nil {
				return err
			}

			vs30s = append(vs30s, Vs30{
				Station:    strings.TrimSpace(d[vs30Station]),
				Vs30:       vs30,
				Tsite:      r,
				Zb:         zb,
				QVs30:      strings.TrimSpace(d[vs30QVs30]),
				QTsite:     strings.TrimSpace(d[vs30QTsite]),
				DTsite:     strings.TrimSpace(d[vs30DTsite]),
				QZb:        strings.TrimSpace(d[vs30QZb]),
				References: strings.TrimSpace(d[vs30References]),
				Notes:      strings.TrimSpace(d[vs30Notes]),
			})
		}

		*c = Vs30List(vs30s)
	}
	return nil
}

func LoadVs30s(path string) ([]Vs30, error) {
	var c []Vs30

	if err := LoadList(path, (*Vs30List)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(Vs30List(c))

	return c, nil
}

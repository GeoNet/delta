package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	classStation = iota
	classClass
	classVs30
	classTsite
	classZb
	classQVs30
	classQTsite
	classDTsite
	classQZb
	classReferences
	classLast
)

type Class struct {
	Station    string
	Class      string
	Vs30       float64
	Tsite      Range
	Zb         float64
	QVs30      string
	QTsite     string
	DTsite     string
	QZb        string
	References string
}

func (c Class) Less(class Class) bool {
	switch {
	case c.Station < class.Station:
		return true
	default:
		return false
	}
}

type ClassList []Class

func (c ClassList) Len() int           { return len(c) }
func (c ClassList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ClassList) Less(i, j int) bool { return c[i].Less(c[j]) }

func (c ClassList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Class",
		"Vs30",
		"Tsite",
		"Zb",
		"Q_Vs30",
		"Q_Tsite",
		"D_Tsite",
		"Q_Zb",
		"References",
	}}
	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Class),
			strconv.FormatFloat(v.Vs30, 'g', -1, 64),
			strings.TrimSpace(v.Tsite.String()),
			strconv.FormatFloat(v.Zb, 'g', -1, 64),
			strings.TrimSpace(v.QVs30),
			strings.TrimSpace(v.QTsite),
			strings.TrimSpace(v.DTsite),
			strings.TrimSpace(v.QZb),
			strings.TrimSpace(v.References),
		})
	}
	return data
}

func (c *ClassList) decode(data [][]string) error {
	var classes []Class
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != classLast {
				return fmt.Errorf("incorrect number of installed class fields")
			}
			var err error

			var vs30, zb float64
			if vs30, err = strconv.ParseFloat(d[classVs30], 64); err != nil {
				return err
			}
			if zb, err = strconv.ParseFloat(d[classZb], 64); err != nil {
				return err
			}

			var r Range
			if r, err = NewRange(d[classTsite]); err != nil {
				return err
			}

			classes = append(classes, Class{
				Station:    strings.TrimSpace(d[classStation]),
				Class:      strings.TrimSpace(d[classClass]),
				Vs30:       vs30,
				Tsite:      r,
				Zb:         zb,
				QVs30:      strings.TrimSpace(d[classQVs30]),
				QTsite:     strings.TrimSpace(d[classQTsite]),
				DTsite:     strings.TrimSpace(d[classDTsite]),
				QZb:        strings.TrimSpace(d[classQZb]),
				References: strings.TrimSpace(d[classReferences]),
			})
		}

		*c = ClassList(classes)
	}
	return nil
}

func LoadClasses(path string) ([]Class, error) {
	var c []Class

	if err := LoadList(path, (*ClassList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(ClassList(c))

	return c, nil
}

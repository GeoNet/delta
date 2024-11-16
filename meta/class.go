package meta

import (
	"sort"
	"strconv"
	"strings"
)

const (
	classStation = iota
	classSiteClass
	classVs30
	classVs30Quality
	classTsite
	classTsiteMethod
	classTsiteQuality
	classBasementDepth
	classDepthQuality
	classLink
	classCitations
	classNotes
	classLast
)

var classHeaders Header = map[string]int{
	"Station":        classStation,
	"Site Class":     classSiteClass,
	"Vs30":           classVs30,
	"Vs30 Quality":   classVs30Quality,
	"Tsite":          classTsite,
	"Tsite Method":   classTsiteMethod,
	"Tsite Quality":  classTsiteQuality,
	"Basement Depth": classBasementDepth,
	"Depth Quality":  classDepthQuality,
	"Link":           classLink,
	"Citations":      classCitations,
	"Notes":          classNotes,
}

var ClassTable Table = Table{
	name:     "Class",
	headers:  classHeaders,
	primary:  []string{"Station"},
	native:   []string{"Vs30", "Basement Depth"},
	foreign:  map[string]map[string]string{},
	nullable: []string{"Citations", "Link", "Notes"},
	remap: map[string]string{
		"Site Class":     "SiteClass",
		"Vs30 Quality":   "Vs30Quality",
		"Tsite Method":   "TsiteMethod",
		"Tsite Quality":  "TsiteQuality",
		"Basement Depth": "BasementDepth",
		"Depth Quality":  "DepthQuality",
	},
}

type Class struct {
	Station       string
	SiteClass     string
	Vs30          float64
	Vs30Quality   string
	Tsite         Range
	TsiteMethod   string
	TsiteQuality  string
	BasementDepth float64
	DepthQuality  string
	Link          string
	Citations     []string
	Notes         string
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
	var data [][]string

	data = append(data, classHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.SiteClass),
			strconv.FormatFloat(row.Vs30, 'g', -1, 64),
			strings.TrimSpace(row.Vs30Quality),
			strings.TrimSpace(row.Tsite.String()),
			strings.TrimSpace(row.TsiteMethod),
			strings.TrimSpace(row.TsiteQuality),
			strconv.FormatFloat(row.BasementDepth, 'g', -1, 64),
			strings.TrimSpace(row.DepthQuality),
			strings.TrimSpace(row.Link),
			strings.Join(row.Citations, " "),
			strings.TrimSpace(row.Notes),
		})
	}

	return data
}

func (c *ClassList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var classes []Class

	fields := classHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		vs30, err := strconv.ParseFloat(d[classVs30], 64)
		if err != nil {
			return err
		}
		zb, err := strconv.ParseFloat(d[classBasementDepth], 64)
		if err != nil {
			return err
		}

		r, err := NewRange(d[classTsite])
		if err != nil {
			return err
		}

		citations := strings.Fields(strings.TrimSpace(d[classCitations]))

		sort.Strings(citations)

		classes = append(classes, Class{
			Station:       strings.TrimSpace(d[classStation]),
			SiteClass:     strings.TrimSpace(d[classSiteClass]),
			Vs30:          vs30,
			Vs30Quality:   strings.TrimSpace(d[classVs30Quality]),
			Tsite:         r,
			TsiteMethod:   strings.TrimSpace(d[classTsiteMethod]),
			TsiteQuality:  strings.TrimSpace(d[classTsiteQuality]),
			BasementDepth: zb,
			DepthQuality:  strings.TrimSpace(d[classDepthQuality]),
			Link:          strings.TrimSpace(d[classLink]),
			Citations:     citations,
			Notes:         strings.TrimSpace(d[classNotes]),
		})
	}

	*c = ClassList(classes)

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

package meta

import (
	"sort"
	"strconv"
)

type Serial string

func (s Serial) Less(serial Serial) bool {
	i, err := strconv.Atoi(string(s))
	if err != nil {
		return s < serial
	}
	j, err := strconv.Atoi(string(serial))
	if err != nil {
		return s < serial
	}
	return i < j
}

func (s Serial) Greater(serial Serial) bool {
	i, err := strconv.Atoi(string(s))
	if err != nil {
		return s > serial
	}
	j, err := strconv.Atoi(string(serial))
	if err != nil {
		return s > serial
	}
	return i > j
}

func (s Serial) Equal(serial Serial) bool {
	i, err := strconv.Atoi(string(s))
	if err != nil {
		return s == serial
	}
	j, err := strconv.Atoi(string(serial))
	if err != nil {
		return s == serial
	}
	return i == j
}

type Equipment struct {
	Manufacturer string `csv:"Manufacturer"`
	Make         string `csv:"Make"`
	Model        string `csv:"Model"`
	Serial       string `csv:"Serial Number"`
	Notes        string `csv:"Notes"`
}

type Equipments []Equipment

func (e Equipments) Len() int      { return len(e) }
func (e Equipments) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Equipments) Less(i, j int) bool {

	switch {
	case e[i].Make < e[j].Make:
		return true
	case e[i].Make > e[j].Make:
		return false
	case e[i].Model < e[j].Model:
		return true
	case e[i].Model > e[j].Model:
		return false
	case Serial(e[i].Serial).Less(Serial(e[j].Serial)):
		return true
	default:
		return false
	}
}

func (e Equipments) List()      {}
func (e Equipments) Sort() List { sort.Sort(e); return e }

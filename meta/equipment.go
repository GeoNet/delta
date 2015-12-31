package meta

import (
	//	"sort"
	"strconv"
)

/*
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
*/

type Equipment struct {
	Make   string
	Model  string
	Serial string
}

func (e Equipment) less(eq Equipment) bool {

	// check by make & model first
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

	// use a numerical serial compare if possible
	if a, err := strconv.Atoi(e.Serial); err == nil {
		if b, err := strconv.Atoi(eq.Serial); err == nil {
			return a < b
		}
	}

	// otherwise a string compare
	return e.Serial < eq.Serial
}

/*
type Equipments []Equipment

func (e Equipments) Len() int           { return len(e) }
func (e Equipments) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e Equipments) Less(i, j int) bool { return e[i].less(e[j]) }

func (e Equipments) decode() [][]string {
	return nil
}
func (e *Equipments) encode(data [][]string) error {
	return nil
}

func LoadEquipment(path string) ([]Equipment, error) {
	var e []Equipment

	if err := LoadList(path, (*Equipments)(&e)); err != nil {
		return nil, err
	}

	sort.Sort(Equipments(e))

	return e, nil
}
*/

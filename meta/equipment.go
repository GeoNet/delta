package meta

import (
	"strconv"
)

type Equipment struct {
	Make   string
	Model  string
	Serial string
}

func (e Equipment) Less(eq Equipment) bool {

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

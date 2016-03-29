package stationxml

import (
	"fmt"
	"regexp"
)

type Email struct {
	Address string `xml:",chardata"`
}

func (e Email) IsValid() error {

	if !(len(e.Address) > 0) {
		return fmt.Errorf("empty email")
	}

	if !(regexp.MustCompile(`^[\w\.\-_]+@[\w\.\-_]+$`).MatchString(e.Address)) {
		return fmt.Errorf("bad email address: %s", e)
	}

	return nil
}

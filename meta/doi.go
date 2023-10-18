package meta

import (
	"bytes"
	"errors"
	"regexp"
)

const doiMatch = `^((?:https?://)?(?:[^/\s]+\.)*doi\.org/10|10|DOI:10)\.([^\s]+)+$`

var doiRe = regexp.MustCompile(doiMatch)

var ErrInvalidDOI = errors.New("invalid DOI reference")

type Doi string

// NewDoi returns a Doi value from a string after checking.
func NewDoi(doi string) (Doi, error) {
	var d Doi
	if err := d.UnmarshalText([]byte(doi)); err != nil {
		return "", err
	}
	return d, nil
}

// MustDoi returns a Doi value, or panics. This is useful for testing with known values.
func MustDoi(doi string) Doi {
	d, err := NewDoi(doi)
	if err != nil {
		panic(err)
	}
	return d
}

func (d Doi) String() string {
	return string(d)
}

func (d Doi) Equal(doi Doi) bool {
	return bytes.Equal([]byte(d), []byte(doi))
}

// UnmarshalText implements the TextUnmarshaler interface.
func (d *Doi) UnmarshalText(data []byte) error {
	if !doiRe.Match(data) {
		return ErrInvalidDOI
	}

	*d = Doi(data)

	return nil
}

// MarshalText implements the TextMarshaler interface.
func (d Doi) MarshalText() ([]byte, error) {
	return []byte(d), nil
}

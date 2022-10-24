package build

import (
	"bytes"
	"encoding/gob"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func clone(a, b *stationxml.ResponseStageType) error {

	var buff bytes.Buffer

	if err := gob.NewEncoder(&buff).Encode(a); err != nil {
		return err
	}
	if err := gob.NewDecoder(&buff).Decode(b); err != nil {
		return err
	}

	return nil
}

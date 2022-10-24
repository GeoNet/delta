package build

import (
	"bytes"
	"encoding/gob"
)

func clone(a, b interface{}) error {

	var buff bytes.Buffer

	if err := gob.NewEncoder(&buff).Encode(a); err != nil {
		return err
	}
	if err := gob.NewDecoder(&buff).Decode(b); err != nil {
		return err
	}

	return nil
}

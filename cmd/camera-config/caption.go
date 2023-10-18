package main

import (
	"encoding/json"
	"io"
	"os"
)

type Caption struct {
	Mount string `json:"mount"`
	View  string `json:"view"`
	Label string `json:"label"`
}

type Captions []Caption

func (c Captions) EncodeFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.Encode(file)
}

func (c Captions) Encode(wr io.Writer) error {
	return json.NewEncoder(wr).Encode(c)
}

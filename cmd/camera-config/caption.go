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
	file, err := os.Create(path) //nolint:gosec // disable G304
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := c.Encode(file); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (c Captions) Encode(wr io.Writer) error {
	return json.NewEncoder(wr).Encode(c)
}

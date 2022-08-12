package main

import (
	"encoding/json"
	"io"
	"os"
	"sort"
)

type Volcano struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type Volcam struct {
	Id        string    `json:"id"`
	Mount     string    `json:"mount"`
	View      string    `json:"view"`
	Title     string    `json:"title"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Datum     string    `json:"datum"`
	Azimuth   float64   `json:"azimuth"`
	Height    float64   `json:"height"`
	Ground    float64   `json:"ground"`
	Volcanoes []Volcano `json:"volcanoes"`
}

type Volcams []Volcam

func (v Volcams) EncodeFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return v.Encode(file)
}

func (v Volcams) Encode(wr io.Writer) error {
	for _, n := range v {
		sort.Slice(n.Volcanoes, func(i, j int) bool {
			return n.Volcanoes[i].Id < n.Volcanoes[j].Id
		})
	}
	sort.Slice(v, func(i, j int) bool {
		switch {
		case v[i].Mount < v[j].Mount:
			return true
		case v[i].Mount > v[j].Mount:
			return false
		case v[i].View < v[j].View:
			return true
		default:
			return false
		}
	})

	if err := json.NewEncoder(wr).Encode(v); err != nil {
		return err
	}

	return nil
}

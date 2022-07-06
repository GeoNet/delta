package main

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {

	for k, v := range map[string]string{
		"":             "",
		"u":            "u",
		"U":            "u",
		"Up":           "up",
		"UpDown":       "up_down",
		"UpDownUp":     "up_down_up",
		"_UpDownUp_":   "_up_down_up_",
		"_Up_Down_Up_": "_up_down_up_",
		"UPDownUP":     "up_down_up",
		"__UPDownUP__": "_up_down_up_",
	} {
		if s := toSnakeCase(k); s != v {
			t.Errorf("incorrect snake case conversion for \"%s\", expected \"%s\" got \"%s\"", k, v, s)
		}
	}
}

func TestFileName(t *testing.T) {

	for k, v := range map[string]string{
		"":             ".go",
		"A":            "a.go",
		"AB":           "ab.go",
		"AB/C":         "c.go",
		"AB/CD":        "cd.go",
		"AB/C_D":       "c_d.go",
		"AB/C_DE":      "c_de.go",
		"dir/Type":     "type.go",
		"dir/TestType": "test_type.go",
		"dir/TESTType": "test_type.go",
	} {
		if s := FileName(k, ".go"); s != v {
			t.Errorf("incorrect file name conversion for \"%s\", expected \"%s\" got \"%s\"", k, v, s)
		}
	}
}

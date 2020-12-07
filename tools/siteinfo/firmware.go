package main

import (
	"time"
)

var firmwareMap = map[string]string{
	"unknown":           "0.00",
	"NP 7.19 / SP 3.04": "7.19",
	"5.10/3.013":        "5.10",
	"3.2.32.8":          "3.20",
	"1.2.5":             "1.20",
	"1F39":              "8.32",
	"1F50":              "8.35",
	"1F60":              "8.36",
	"CC00":              "9.10",
	"CD00":              "9.20",
	"ZC00":              "9.92",
	"4.17A":             "4.17",
}

type FirmwareHistory struct {
	Start   time.Time
	End     time.Time
	Version string
}

func (f FirmwareHistory) Less(firmware FirmwareHistory) bool {
	return f.Start.Before(firmware.Start)
}

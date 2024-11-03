package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	firmwareHistoryMake int = iota
	firmwareHistoryModel
	firmwareHistorySerial
	firmwareHistoryVersion
	firmwareHistoryStart
	firmwareHistoryEnd
	firmwareHistoryNotes
	firmwareHistoryLast
)

var firmwareHistoryHeaders Header = map[string]int{
	"Make":       firmwareHistoryMake,
	"Model":      firmwareHistoryModel,
	"Serial":     firmwareHistorySerial,
	"Version":    firmwareHistoryVersion,
	"Start Date": firmwareHistoryStart,
	"End Date":   firmwareHistoryEnd,
	"Notes":      firmwareHistoryNotes,
}

type FirmwareHistory struct {
	Install

	Version string
	Notes   string
}

var FirmwareHistoryTable Table = Table{
	name:    "Firmware",
	headers: firmwareHistoryHeaders,
	primary: []string{"Make", "Model", "Serial", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{
		"Asset": {"Make", "Model", "Serial"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type FirmwareHistoryList []FirmwareHistory

func (fh FirmwareHistoryList) Len() int           { return len(fh) }
func (fh FirmwareHistoryList) Swap(i, j int)      { fh[i], fh[j] = fh[j], fh[i] }
func (fh FirmwareHistoryList) Less(i, j int) bool { return fh[i].Install.Less(fh[j].Install) }

func (fh FirmwareHistoryList) encode() [][]string {
	var data [][]string

	data = append(data, firmwareHistoryHeaders.Columns())

	for _, row := range fh {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Version),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
			strings.TrimSpace(row.Notes),
		})
	}

	return data
}

func (fh *FirmwareHistoryList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var histories []FirmwareHistory

	fields := firmwareHistoryHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[firmwareHistoryStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[firmwareHistoryEnd])
		if err != nil {
			return err
		}

		histories = append(histories, FirmwareHistory{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[firmwareHistoryMake]),
					Model:  strings.TrimSpace(d[firmwareHistoryModel]),
					Serial: strings.TrimSpace(d[firmwareHistorySerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Version: strings.TrimSpace(d[firmwareHistoryVersion]),
			Notes:   strings.TrimSpace(d[firmwareHistoryNotes]),
		})
	}

	*fh = FirmwareHistoryList(histories)

	return nil
}

func LoadFirmwareHistory(path string) ([]FirmwareHistory, error) {
	var fh []FirmwareHistory

	if err := LoadList(path, (*FirmwareHistoryList)(&fh)); err != nil {
		return nil, err
	}

	sort.Sort(FirmwareHistoryList(fh))

	return fh, nil
}

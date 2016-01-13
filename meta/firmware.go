package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type FirmwareHistory struct {
	Install

	Version string
	Notes   string
}

type FirmwareHistoryList []FirmwareHistory

func (f FirmwareHistoryList) Len() int           { return len(f) }
func (f FirmwareHistoryList) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FirmwareHistoryList) Less(i, j int) bool { return f[i].Install.less(f[j].Install) }

func (f FirmwareHistoryList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial Number",
		"Version",
		"Start Date",
		"End Date",
		"Notes",
	}}
	for _, v := range f {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Version),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
			strings.TrimSpace(v.Notes),
		})
	}
	return data
}

func (f *FirmwareHistoryList) decode(data [][]string) error {
	var histories []FirmwareHistory
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 7 {
				return fmt.Errorf("incorrect number of firmware history fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[4]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[5]); err != nil {
				return err
			}

			histories = append(histories, FirmwareHistory{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[0]),
						Model:  strings.TrimSpace(d[1]),
						Serial: strings.TrimSpace(d[2]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Version: strings.TrimSpace(d[3]),
				Notes:   strings.TrimSpace(d[6]),
			})
		}

		*f = FirmwareHistoryList(histories)
	}
	return nil
}

func LoadFirmwareHistory(path string) ([]FirmwareHistory, error) {
	var f []FirmwareHistory

	if err := LoadList(path, (*FirmwareHistoryList)(&f)); err != nil {
		return nil, err
	}

	sort.Sort(FirmwareHistoryList(f))

	return f, nil
}

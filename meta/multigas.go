package meta

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "15:04:05"

const (
	multigasStation int = iota
	multigasLocation
	multigasGas
	multigasConcentration
	multigasFrequency
	multigasDay
	multigasCalibration
	multigasZero
	multigasStart
	multigasEnd
)

var multigasHeaders Header = map[string]int{
	"Station":          multigasStation,
	"Location":         multigasLocation,
	"Gas":              multigasGas,
	"Concentration":    multigasConcentration,
	"Frequency":        multigasFrequency,
	"Day Of Week":      multigasDay,
	"Calibration Time": multigasCalibration,
	"Zero Time":        multigasZero,
	"Start Date":       multigasStart,
	"End Date":         multigasEnd,
}

var ErrInvalidCalibrationFrequency = errors.New("invalid calibration frequency")

type CalibrationFrequency int

const (
	CalibrationUnknown = iota
	CalibrationDaily
	CalibrationWeekly
)

var frquencyName = map[CalibrationFrequency]string{
	CalibrationUnknown: "",
	CalibrationDaily:   "D",
	CalibrationWeekly:  "W",
}

func NewCalibrationFrequency(str string) (CalibrationFrequency, error) {
	for freq, name := range frquencyName {
		if name != str {
			continue
		}
		return freq, nil
	}
	return CalibrationUnknown, ErrInvalidCalibrationFrequency
}

func (cf CalibrationFrequency) String() string {
	return frquencyName[cf]
}

var ErrInvalidWeekday = errors.New("invalid weekday")

func parseWeekday(str string) (time.Weekday, error) {
	for day := time.Sunday; day <= time.Saturday; day++ {
		if day.String() != str {
			continue
		}
		return day, nil
	}
	return 0, ErrInvalidWeekday
}

var MultigasTable Table = Table{
	name:    "Multigas",
	headers: multigasHeaders,
	primary: []string{"Station", "Location", "Gas", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{
		"Station": {"Station"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Multigas struct {
	Span

	Station       string
	Location      string
	Gas           string
	Concentration float64
	Frequency     CalibrationFrequency
	Day           time.Weekday
	Calibration   time.Time
	Zero          time.Time

	concentration string
}

func (mg Multigas) Compare(multigas Multigas) int {
	if cmp := strings.Compare(mg.Station, multigas.Station); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(mg.Location, multigas.Location); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(mg.Gas, multigas.Gas); cmp != 0 {
		return cmp
	}
	if cmp := mg.Start.Compare(mg.End); cmp != 0 {
		return cmp
	}
	return 0
}

func (mg Multigas) Less(multigas Multigas) bool {
	return mg.Compare(multigas) < 0
}

type MultigasList []Multigas

func (mg MultigasList) Len() int           { return len(mg) }
func (mg MultigasList) Swap(i, j int)      { mg[i], mg[j] = mg[j], mg[i] }
func (mg MultigasList) Less(i, j int) bool { return mg[i].Less(mg[j]) }

func (mg MultigasList) encode() [][]string {
	var data [][]string

	data = append(data, multigasHeaders.Columns())

	for _, row := range mg {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Gas),
			strings.TrimSpace(row.concentration),
			strings.TrimSpace(row.Frequency.String()),
			strings.TrimSpace(row.Day.String()),
			row.Calibration.Format(TimeFormat),
			row.Zero.Format(TimeFormat),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	//log.Println("->", data)

	return data
}

func (mg *MultigasList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}

	var multigas []Multigas

	fields := multigasHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		concentration, err := strconv.ParseFloat(d[multigasConcentration], 64)
		if err != nil {
			return err
		}

		frequency, err := NewCalibrationFrequency(d[multigasFrequency])
		if err != nil {
			return err
		}

		day, err := parseWeekday(d[multigasDay])
		if err != nil {
			return err
		}

		calibration, err := time.ParseInLocation(TimeFormat, d[multigasCalibration], time.UTC)
		if err != nil {
			return err
		}

		zero, err := time.ParseInLocation(TimeFormat, d[multigasZero], time.UTC)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[multigasStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[multigasEnd])
		if err != nil {
			return err
		}

		multigas = append(multigas, Multigas{
			Span: Span{
				Start: start,
				End:   end,
			},
			Station:       strings.TrimSpace(d[multigasStation]),
			Location:      strings.TrimSpace(d[multigasLocation]),
			Gas:           strings.TrimSpace(d[multigasGas]),
			Concentration: concentration,
			Frequency:     frequency,
			Day:           day,
			Calibration:   calibration,
			Zero:          zero,

			concentration: strings.TrimSpace(d[multigasConcentration]),
		})

	}

	*mg = MultigasList(multigas)

	return nil
}

func LoadMultigas(path string) ([]Multigas, error) {
	var mg []Multigas

	if err := Load(path, (*MultigasList)(&mg)); err != nil {
		return nil, err
	}

	sort.Sort(MultigasList(mg))

	return mg, nil
}

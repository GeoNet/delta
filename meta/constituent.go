package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	constituentGauge = iota
	constituentNumber
	constituentName
	constituentAmplitude
	constituentLag
	constituentStart
	constituentEnd
	constituentLast
)

var constituentHeaders Header = map[string]int{
	"Gauge":       constituentGauge,
	"Number":      constituentNumber,
	"Constituent": constituentName,
	"Amplitude":   constituentAmplitude,
	"Lag":         constituentLag,
	"Start Date":  constituentStart,
	"End Date":    constituentEnd,
}

type Constituent struct {
	Span

	Gauge     string
	Number    int
	Name      string
	Amplitude float64
	Lag       float64

	amplitude string // shadow variable to maintain formatting
	lag       string // shadow variable to maintain formatting
}

type ConstituentList []Constituent

func (c ConstituentList) Len() int      { return len(c) }
func (c ConstituentList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ConstituentList) Less(i, j int) bool {
	switch {
	case c[i].Gauge < c[j].Gauge:
		return true
	case c[i].Gauge > c[j].Gauge:
		return false
	case c[i].Start.Before(c[j].Start):
		return true
	case c[i].Start.After(c[j].Start):
		return false
	case c[i].Number < c[j].Number:
		return true
	default:
		return false
	}
}

func (c ConstituentList) encode() [][]string {
	var data [][]string

	data = append(data, constituentHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Gauge),
			strconv.Itoa(row.Number),
			strings.TrimSpace(row.Name),
			strings.TrimSpace(row.amplitude),
			strings.TrimSpace(row.lag),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (c *ConstituentList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var constituents []Constituent

	fields := constituentHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		num, err := ParseInt(d[constituentNumber])
		if err != nil {
			return err
		}

		amp, err := strconv.ParseFloat(d[constituentAmplitude], 64)
		if err != nil {
			return err
		}

		lag, err := strconv.ParseFloat(d[constituentLag], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[constituentStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[constituentEnd])
		if err != nil {
			return err
		}

		constituents = append(constituents, Constituent{
			Span: Span{
				Start: start,
				End:   end,
			},

			Gauge:     d[constituentGauge],
			Number:    num,
			Name:      d[constituentName],
			Amplitude: amp,
			Lag:       lag,

			amplitude: d[constituentAmplitude],
			lag:       d[constituentLag],
		})
	}

	*c = ConstituentList(constituents)

	return nil
}

func LoadConstituents(path string) ([]Constituent, error) {
	var c []Constituent

	if err := LoadList(path, (*ConstituentList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(ConstituentList(c))

	return c, nil
}

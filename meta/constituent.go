package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	constituentGauge = iota
	constituentNumber
	constituentName
	constituentAmplitude
	constituentLag
	constituentLast
)

type Constituent struct {
	Gauge     string
	Number    int
	Name      string
	Amplitude float64
	Lag       float64
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
	default:
		return c[i].Number < c[j].Number
	}
}

func (c ConstituentList) encode() [][]string {
	data := [][]string{{
		"Gauge",
		"Number",
		"Constituent",
		"Amplitude",
		"Lag",
	}}
	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Gauge),
			strconv.Itoa(v.Number),
			strings.TrimSpace(v.Name),
			strconv.FormatFloat(v.Amplitude, 'g', -1, 64),
			strconv.FormatFloat(v.Lag, 'g', -1, 64),
		})
	}
	return data
}

func (c *ConstituentList) decode(data [][]string) error {
	var constituents []Constituent
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != constituentLast {
				return fmt.Errorf("incorrect number of installed constituent fields")
			}
			var err error

			var num int
			if num, err = strconv.Atoi(d[constituentNumber]); err != nil {
				return err
			}

			var amp, lag float64
			if amp, err = strconv.ParseFloat(d[constituentAmplitude], 64); err != nil {
				return err
			}
			if lag, err = strconv.ParseFloat(d[constituentLag], 64); err != nil {
				return err
			}

			constituents = append(constituents, Constituent{
				Gauge:     d[constituentGauge],
				Number:    num,
				Name:      d[constituentName],
				Amplitude: amp,
				Lag:       lag,
			})
		}

		*c = ConstituentList(constituents)
	}
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

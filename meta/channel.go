package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	channelMake = iota
	channelModel
	channelType
	channelNumber
	channelSamplingRate
	channelResponse
	channelLast
)

// Channel is used to describe a generic recording from a Datalogger.
type Channel struct {
	Make         string
	Model        string
	Type         string
	SamplingRate float64
	Response     string
	Number       int

	number       string
	samplingRate string
	gain         string
}

// Less compares Channel structs suitable for sorting.
func (c Channel) Less(comp Channel) bool {

	switch {
	case strings.ToLower(c.Make) < strings.ToLower(comp.Make):
		return true
	case strings.ToLower(c.Make) > strings.ToLower(comp.Make):
		return false
	case strings.ToLower(c.Model) < strings.ToLower(comp.Model):
		return true
	case strings.ToLower(c.Model) > strings.ToLower(comp.Model):
		return false
	case c.Number < comp.Number:
		return true
	case c.Number > comp.Number:
		return false
	case c.SamplingRate < comp.SamplingRate:
		return true
	default:
		return false
	}
}

type ChannelList []Channel

func (s ChannelList) Len() int           { return len(s) }
func (s ChannelList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ChannelList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s ChannelList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Type",
		"Number",
		"SamplingRate",
		"Response",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.number),
			strings.TrimSpace(v.samplingRate),
			strings.TrimSpace(v.gain),
			strings.TrimSpace(v.Response),
		})
	}

	return data
}
func (s *ChannelList) decode(data [][]string) error {
	var channels []Channel

	if !(len(data) > 1) {
		return nil
	}

	for _, d := range data[1:] {
		if len(d) != channelLast {
			return fmt.Errorf("incorrect number of installed channel fields")
		}

		var number int
		if s := strings.TrimSpace(d[channelNumber]); s != "" {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			number = v
		}

		var samplingRate float64
		if s := strings.TrimSpace(d[channelSamplingRate]); s != "" {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			switch {
			case v < 0.0:
				samplingRate = -1.0 / v
			default:
				samplingRate = v
			}
		}

		channels = append(channels, Channel{
			Make:     strings.TrimSpace(d[channelMake]),
			Model:    strings.TrimSpace(d[channelModel]),
			Type:     strings.TrimSpace(d[channelType]),
			Response: strings.TrimSpace(d[channelResponse]),

			Number:       number,
			SamplingRate: samplingRate,

			samplingRate: strings.TrimSpace(d[channelSamplingRate]),
			number:       strings.TrimSpace(d[channelNumber]),
		})
	}

	*s = ChannelList(channels)

	return nil
}

func LoadChannels(path string) ([]Channel, error) {
	var s []Channel

	if err := LoadList(path, (*ChannelList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(ChannelList(s))

	return s, nil
}

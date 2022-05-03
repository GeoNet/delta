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
	channelDescription
	channelType
	channelPin
	channelSamplingRate
	channelGain
	channelResponse
	channelLast
)

type Channel struct {
	Make        string
	Model       string
	Description string
	Type        string
	Response    string

	Pin          int
	SamplingRate float64
	Gain         float64

	pin          string
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
	case c.Pin < comp.Pin:
		return true
	case c.Pin > comp.Pin:
		return false
	case c.SamplingRate < comp.SamplingRate:
		return true
	case c.SamplingRate > comp.SamplingRate:
		return false
	case c.Gain < comp.Gain:
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
		"Description",
		"Type",
		"Pin",
		"SamplingRate",
		"Gain",
		"Response",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Description),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.pin),
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

		var pin int
		if v := strings.TrimSpace(d[channelPin]); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			pin = n
		}

		var gain float64
		if v := strings.TrimSpace(d[channelGain]); v != "" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			gain = f
		}

		samplingRate, err := strconv.ParseFloat(d[channelSamplingRate], 64)
		if err != nil {
			return err
		}

		channels = append(channels, Channel{
			Make:        strings.TrimSpace(d[channelMake]),
			Model:       strings.TrimSpace(d[channelModel]),
			Description: strings.TrimSpace(d[channelDescription]),
			Type:        strings.TrimSpace(d[channelType]),
			Response:    strings.TrimSpace(d[channelResponse]),

			Pin:          pin,
			Gain:         gain,
			SamplingRate: samplingRate,

			samplingRate: strings.TrimSpace(d[channelSamplingRate]),
			pin:          strings.TrimSpace(d[channelPin]),
			gain:         strings.TrimSpace(d[channelGain]),
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

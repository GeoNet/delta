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

var channelHeaders Header = map[string]int{
	"Make":         channelMake,
	"Model":        channelModel,
	"Type":         channelType,
	"Number":       channelNumber,
	"SamplingRate": channelSamplingRate,
	"Response":     channelResponse,
}

var ChannelTable Table = Table{
	name:    "Channel",
	headers: channelHeaders,
	primary: []string{"Make", "Model", "Number", "SamplingRate"},
	native:  []string{"Number", "SamplingRate"},
	foreign: map[string][]string{
		"Network": {"Network"},
	},
}

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
}

// Description returns a short label for the channel model family.
func (c Channel) Description() string {
	return fmt.Sprintf("%s %s %s", c.Make, strings.Split(strings.Fields(c.Model)[0], "/")[0], c.Type)
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
	case c.SamplingRate > comp.SamplingRate:
		return false
	default:
		return true
	}
}

type ChannelList []Channel

func (c ChannelList) Len() int           { return len(c) }
func (c ChannelList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ChannelList) Less(i, j int) bool { return c[i].Less(c[j]) }

func (c ChannelList) encode() [][]string {
	var data [][]string

	data = append(data, channelHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Type),
			strings.TrimSpace(row.number),
			strings.TrimSpace(row.samplingRate),
			strings.TrimSpace(row.Response),
		})
	}

	return data
}

func (c *ChannelList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var channels []Channel

	fields := channelHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		number, err := ParseInt(d[channelNumber])
		if err != nil {
			return err
		}

		samplingRate, err := strconv.ParseFloat(d[channelSamplingRate], 64)
		if err != nil {
			return err
		}
		if samplingRate < 0.0 {
			samplingRate = -1.0 / samplingRate
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

	*c = ChannelList(channels)

	return nil
}

func LoadChannels(path string) ([]Channel, error) {
	var c []Channel

	if err := LoadList(path, (*ChannelList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(ChannelList(c))

	return c, nil
}

package meta

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	acquisitionID int = iota
	acquisistionInterrogatorID
	acquisitionChannelGroup
	acquisitionStartTime
	acquisitionEndTime
	acquisitionSampleRate
	// acquisitionSampleRateUnit
	acquisitionGaugeLength
	// acquisitionGaugeLengthUnit
	acquisitionUnitOfMeasure
	acquisitionNumberOfChannels
	acquisitionSpatialSamplingInterval
	acquisitionSpatialSamplingIntervalUnits
	acquisitionPulseRate
	// acquisitionPulseRateUnit
	acquisitionPulseWidth
	// acquisitionPulseWidthUnit
	acquisitionComment
	acquisitionNativeHeaders
	acquisitionLast
)

var acquisitionHeaders Header = map[string]int{
	"Acquisition ID":  acquisitionID,
	"Interrogator ID": acquisistionInterrogatorID,
	"Channel Group":   acquisitionChannelGroup,
	"Sample Rate":     acquisitionSampleRate,
	// "Acquisition Sample Rate Unit":    acquisitionSampleRateUnit,
	"Gauge Length": acquisitionGaugeLength,
	// "Gauge Length Unit":               acquisitionGaugeLengthUnit,
	"Unit Of Measure":                 acquisitionUnitOfMeasure,
	"Number Of Channels":              acquisitionNumberOfChannels,
	"Spatial Sampling Interval":       acquisitionSpatialSamplingInterval,
	"Spatial Sampling Interval Units": acquisitionSpatialSamplingIntervalUnits,
	"Pulse Rate":                      acquisitionPulseRate,
	// "Pulse Rate Unit":                 acquisitionPulseRateUnit,
	"Pulse Width": acquisitionPulseWidth,
	// "Pulse Width Unit":                acquisitionPulseWidthUnit,
	// "Comment":                         acquisitionComment,
	"Native Headers": acquisitionNativeHeaders,

	"Start Date": acquisitionStartTime,
	"End Date":   acquisitionEndTime,
}

var AcquisitionTable Table = Table{
	name:    "Acquisition",
	headers: acquisitionHeaders,
	primary: []string{"Acquisition ID", "Acquisition Start Time"},
	native: []string{
		"Acquisition Sample Rate",
		"Gauge Length",
		"Number Of Channels",
		"Spatial Sampling Interval",
		"Pulse Rate",
		"Pulse Width",
		"Channel Groups",
	},
	foreign: map[string][]string{
		// "Interrogator ID": {"Interrogator ID"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Acquisition struct {
	Span

	AcquisitionID  string `json:"acquisition_id"`
	InterrogatorID string
	SampleRate     float64 `json:"acquisition_sample_rate"`
	// SampleRateUnit               string  `json:"acquisition_sample_rate_unit"`
	GaugeLength float64 `json:"gauge_length"`
	// GaugeLengthUnit              string  `json:"gauge_length_unit"`
	UnitOfMeasure                string  `json:"unit_of_measure"`
	NumberOfChannels             int     `json:"number_of_channels"`
	SpatialSamplingInterval      float64 `json:"spatial_sampling_interval"`
	SpatialSamplingIntervalUnits string  `json:"spatial_sampling_interval_units"`
	PulseRate                    float64 `json:"pulse_rate"`
	// PulseRateUnit                string  `json:"pulse_rate_unit"`
	PulseWidth float64 `json:"pulse_width"`
	// PulseWidthUnit               string  `json:"pulse_width_unit"`
	// Comment                      string  `json:"comment"`
	NativeHeaders map[string]interface{}
	ChannelGroups string `json:"channel_groups"`

	sampleRate              string // shadow value used to retain formatting
	gaugeLength             string // shadow value used to retain formatting
	spatialSamplingInterval string // shadow value used to retain formatting
	pulseRate               string // shadow value used to retain formatting
	pulseWidth              string // shadow value used to retain formatting
}

type AcquisitionList []Acquisition

func (al AcquisitionList) Len() int      { return len(al) }
func (al AcquisitionList) Swap(i, j int) { al[i], al[j] = al[j], al[i] }
func (al AcquisitionList) Less(i, j int) bool {
	switch {
	case strings.ToLower(al[i].AcquisitionID) < strings.ToLower(al[j].AcquisitionID):
		return true
	case strings.ToLower(al[i].AcquisitionID) > strings.ToLower(al[j].AcquisitionID):
		return false
	default:
		return al[i].Start.Before(al[j].Start)
	}
}

func (al AcquisitionList) encode() [][]string {
	var data [][]string

	data = append(data, acquisitionHeaders.Columns())

	for _, row := range al {
		data = append(data, []string{
			strings.TrimSpace(row.AcquisitionID),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
			strings.TrimSpace(row.sampleRate),
			// strings.TrimSpace(row.SampleRateUnit),
			strings.TrimSpace(row.gaugeLength),
			// strings.TrimSpace(row.GaugeLengthUnit),
			strings.TrimSpace(row.UnitOfMeasure),
			strconv.Itoa(row.NumberOfChannels),
			strings.TrimSpace(row.spatialSamplingInterval),
			strings.TrimSpace(row.SpatialSamplingIntervalUnits),
			strings.TrimSpace(row.pulseRate),
			// strings.TrimSpace(row.PulseRateUnit),
			strings.TrimSpace(row.pulseWidth),
			// strings.TrimSpace(row.PulseWidthUnit),
			// strings.TrimSpace(row.Comment),
			fmt.Sprint(row.NativeHeaders),
			strings.TrimSpace(row.ChannelGroups),
		})
	}

	return data
}

func (al *AcquisitionList) decode(data [][]string) error {
	if len(data) < 2 {
		return nil
	}
	// fmt.Print("\n\nhere\n\n")
	var acquisitions []Acquisition

	fields := acquisitionHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[acquisitionStartTime])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[acquisitionEndTime])
		if err != nil {
			return err
		}

		sampleRate, err := strconv.ParseFloat(d[acquisitionSampleRate], 64)
		if err != nil {
			return err
		}

		gaugeLength, err := strconv.ParseFloat(d[acquisitionGaugeLength], 64)
		if err != nil {
			return err
		}

		numberOfChannels, err := strconv.Atoi(strings.TrimSpace(d[acquisitionNumberOfChannels]))
		if err != nil {
			// fmt.Print("fail  ", d[acquisitionNumberOfChannels], "\n\n")
			return err
		}

		spatialSamplingInterval, err := strconv.ParseFloat(d[acquisitionSpatialSamplingInterval], 64)
		if err != nil {
			return err
		}

		pulseRate, err := strconv.ParseFloat(d[acquisitionPulseRate], 64)
		if err != nil {
			return err
		}
		// fmt.Printf("Pulse %s %s \n", d[acquisitionPulseRate], pulseRate)
		pulseWidth, err := strconv.ParseFloat(d[acquisitionPulseWidth], 64)
		if err != nil {
			return err
		}
		nativeHeaders := make(map[string]interface{})
		for _, s := range strings.Split(strings.Trim(d[acquisitionNativeHeaders], "{}"), ",") {
			if s != "" {
				k := strings.Split(s, ":")
				if len(k) != 2 {
					err := fmt.Sprintf("Invalid Native Header for acquisition: %s", d[acquisitionID])
					return errors.New(err)
				} else {
					nativeHeaders[strings.TrimSpace(k[0])] = strings.TrimSpace(k[1])
				}
			}
		}

		acquisitions = append(acquisitions, Acquisition{
			Span: Span{
				Start: start,
				End:   end,
			},
			AcquisitionID:  strings.TrimSpace(d[acquisitionID]),
			InterrogatorID: strings.TrimSpace(d[acquisistionInterrogatorID]),
			SampleRate:     sampleRate,
			// SampleRateUnit:               strings.TrimSpace(d[acquisitionSampleRateUnit]),
			GaugeLength: gaugeLength,
			// GaugeLengthUnit:              strings.TrimSpace(d[acquisitionGaugeLengthUnit]),
			UnitOfMeasure:                strings.TrimSpace(d[acquisitionUnitOfMeasure]),
			NumberOfChannels:             numberOfChannels,
			SpatialSamplingInterval:      spatialSamplingInterval,
			SpatialSamplingIntervalUnits: strings.TrimSpace(d[acquisitionSpatialSamplingIntervalUnits]),
			PulseRate:                    pulseRate,
			// PulseRateUnit:                strings.TrimSpace(d[acquisitionPulseRateUnit]),
			PulseWidth: pulseWidth,
			// PulseWidthUnit:               strings.TrimSpace(d[acquisitionPulseWidthUnit]),
			// Comment:                      strings.TrimSpace(d[acquisitionComment]),
			NativeHeaders: nativeHeaders,
			ChannelGroups: strings.TrimSpace(d[acquisitionChannelGroup]),

			sampleRate:              strings.TrimSpace(d[acquisitionSampleRate]),
			gaugeLength:             strings.TrimSpace(d[acquisitionGaugeLength]),
			spatialSamplingInterval: strings.TrimSpace(d[acquisitionSpatialSamplingInterval]),
			pulseRate:               strings.TrimSpace(d[acquisitionPulseRate]),
			pulseWidth:              strings.TrimSpace(d[acquisitionPulseWidth]),
		})
		// fmt.Print("here", acquisitions[0], "\n\n")
	}

	*al = AcquisitionList(acquisitions)

	return nil
}

func LoadAcquisitions(path string) ([]Acquisition, error) {
	var al []Acquisition

	if err := LoadList(path, (*AcquisitionList)(&al)); err != nil {
		return nil, err
	}

	sort.Sort(AcquisitionList(al))

	return al, nil
}

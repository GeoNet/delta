package metadb

import (
	"time"

	"github.com/GeoNet/delta/resp"
)

type Channel struct {
	Network  string
	External string
	Station  string
	Location string
	Code     string

	SampleRate float64
	Start      time.Time
	End        time.Time
}

func (m *MetaDB) Channels(sta string) ([]Channel, error) {
	var channels []Channel

	station, err := m.Station(sta)
	if err != nil {
		return nil, err
	}

	network, err := m.Network(station.Network)
	if err != nil {
		return nil, err
	}

	installations, err := m.Installations(sta)
	if err != nil {
		return nil, err
	}

	for _, installation := range installations {
		location, err := m.Site(station.Code, installation.Location)
		if err != nil {
			return nil, err
		}
		if location == nil {
			continue
		}

		for _, response := range resp.Streams(installation.Datalogger.Model, installation.Sensor.Model) {

			stream, err := m.StationLocationSamplingRateStartStream(
				station.Code,
				installation.Location,
				response.Datalogger.SampleRate,
				installation.Start)
			if err != nil {
				return nil, err
			}
			if stream == nil {
				continue
			}

			lookup := response.Channels(stream.Axial)
			for pin, _ := range response.Components {
				if !(pin < len(lookup)) {
					continue
				}

				channels = append(channels, Channel{
					Network:    network.Code,
					External:   network.External,
					Station:    station.Code,
					Location:   installation.Location,
					Code:       lookup[pin],
					SampleRate: response.SampleRate,
					Start:      installation.Start,
					End:        installation.End,
				})

			}

		}
	}

	return channels, nil
}

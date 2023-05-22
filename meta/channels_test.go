package meta

import (
	"testing"
)

func TestChannelList(t *testing.T) {
	t.Run("check channels", testListFunc("testdata/channels.csv", &ChannelList{
		Channel{
			Make:         "Nanometrics",
			Model:        "Centaur CTR4-6S",
			Type:         "Datalogger",
			SamplingRate: 200,
			Response:     "datalogger_nanometrics_centaur_200_response",

			samplingRate: "200",
		},
		Channel{
			Make:         "Quanterra",
			Model:        "Q330HR/6",
			Type:         "Datalogger",
			Number:       0,
			SamplingRate: 200,
			Response:     "datalogger_quanterra_q330_highgain_200_response",

			number:       "0",
			samplingRate: "200",
		},
		Channel{
			Make:         "Quanterra",
			Model:        "Q330HR/6",
			Type:         "Datalogger",
			Number:       3,
			SamplingRate: 200,
			Response:     "datalogger_quanterra_q330_200_response",

			number:       "3",
			samplingRate: "200",
		},
	}))
}

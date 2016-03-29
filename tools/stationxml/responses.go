package main

var Responses []Response = []Response{
	Response{

		Sensors: []Sensor{
			Sensor{
				Sensors:  []string{"2 Hz Duke Malin Seismometer"},
				Filters:  []string{"2_Hz_Duke_Malin_Seismometer"},
				Channels: "Z12",
				Reversed: true,
			},
		},

		Dataloggers: []Datalogger{
			Datalogger{
				Dataloggers:   []string{"Q330S/3", "Q330S/6"},
				Type:          "CG",
				Label:         "EH",
				Rate:          100.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"Q330S+_FLbelow100-100-PREAMP32"},
			},
		},
	},
	Response{

		Sensors: []Sensor{
			Sensor{
				Sensors:  []string{"L4C"},
				Filters:  []string{"L4C"},
				Channels: "Z",
				Reversed: false,
			},

			Sensor{
				Sensors:  []string{"L4C-3D"},
				Filters:  []string{"L4C"},
				Channels: "ZNE",
				Reversed: false,
			},
			Sensor{
				Sensors:  []string{"LE-3Dlite"},
				Filters:  []string{"LE-3Dlite"},
				Channels: "ZNE",
				Reversed: false,
			},
			Sensor{
				Sensors:  []string{"LE-3DliteMkII"},
				Filters:  []string{"LE-3DliteMkII"},
				Channels: "ZNE",
				Reversed: false,
			},
		},

		Dataloggers: []Datalogger{
			Datalogger{
				Dataloggers:   []string{"Q330/3", "Q330/6"},
				Type:          "CG",
				Label:         "EH",
				Rate:          100.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"Q330_FLbelow100-100"},
				Skip:          "KQ03",
			},
			Datalogger{
				Dataloggers:   []string{"Q330/3", "Q330/6"},
				Type:          "CG",
				Label:         "EH",
				Rate:          200.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"Q330_FLbelow100-200"},
				Match:         "KQ03",
			},
			Datalogger{
				Dataloggers:   []string{"EARSS/3"},
				Type:          "TG",
				Label:         "EH",
				Rate:          50.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"EARSS-50"},
			},
		},
	},
	Response{

		Sensors: []Sensor{
			Sensor{
				Sensors:  []string{"270-600/12V"},
				Filters:  []string{"270-600/12V"},
				Channels: "F",
				Reversed: false,
			},
			Sensor{
				Sensors:  []string{"270-600/24V"},
				Filters:  []string{"270-600/24V"},
				Channels: "F",
				Reversed: false,
			},
			Sensor{
				Sensors:  []string{"270-800/12V"},
				Filters:  []string{"270-800/12V"},
				Channels: "F",
				Reversed: false,
			},
		},

		Dataloggers: []Datalogger{
			Datalogger{
				Dataloggers:   []string{"Q330S/3", "Q330S/6"},
				Type:          "CW",
				Label:         "LD",
				Rate:          1.0,
				Frequency:     0.1,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"Q330S+_FLbelow100-1"},
			},
			Datalogger{
				Dataloggers:   []string{"Q330S/3", "Q330S/6"},
				Type:          "CW",
				Label:         "LD",
				Rate:          100.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"Q330S+_FLbelow100-100"},
			},
		},
	},

	Response{
		Sensors: []Sensor{
			Sensor{
				Sensors:  []string{"FBA-ES-T-DECK"},
				Filters:  []string{"FBA-ES-T-DECK"},
				Channels: "Z12",
				Reversed: false,
			},
			Sensor{
				Sensors:  []string{"FBA-23-DECK"},
				Filters:  []string{"FBA-23-DECK"},
				Channels: "Z12",
				Reversed: true,
			},
		},

		Dataloggers: []Datalogger{
			Datalogger{
				Dataloggers:   []string{"ETNA"},
				Type:          "TG",
				Label:         "HN",
				Rate:          200.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"ALTUS-200"},
			},
		},
	},
	Response{

		Sensors: []Sensor{
			Sensor{
				Sensors:  []string{"FBA-ES-T-DECK"},
				Filters:  []string{"FBA-ES-T-DECK"},
				Channels: "ZNE",
				Reversed: false,
				Match:    "^[A-Z][A-Z][A-Z]$",
			},
			Sensor{
				Sensors:  []string{"FBA-ES-T-DECK"},
				Filters:  []string{"FBA-ES-T-DECK"},
				Channels: "Z12",
				Reversed: false,
				Skip:     "^[A-Z][A-Z][A-Z]$",
			},
			Sensor{
				Sensors:  []string{"FBA-ES-T-BASALT"},
				Filters:  []string{"FBA-ES-T-BASALT"},
				Channels: "Z12",
				Reversed: false,
			},
		},

		Dataloggers: []Datalogger{
			Datalogger{
				Dataloggers:   []string{"BASALT"},
				Type:          "CG",
				Label:         "BN",
				Rate:          50.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"BASALT-50"},
			},
			Datalogger{
				Dataloggers:   []string{"BASALT"},
				Type:          "TG",
				Label:         "HN",
				Rate:          200.0,
				Frequency:     15.0,
				StorageFormat: "Steim2",
				ClockDrift:    0.0001,
				Filters:       []string{"BASALT-200"},
			},
		},
	},
}

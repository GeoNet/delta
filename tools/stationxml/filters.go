package main

var Filters map[string][]ResponseStage = map[string][]ResponseStage{

	"CMG-3ESP-GN": []ResponseStage{
		ResponseStage{
			Type:        "paz",
			Lookup:      "CMG-3ESP-GN",
			Frequency:   1.0,
			Gain:        2000,
			Scale:       1.0,
			InputUnits:  "M/S",
			OutputUnits: "V",
		},
	},

	"FBA-ES-T-BASALT": []ResponseStage{
		ResponseStage{
			Type:        "paz",
			Lookup:      "FBA-ES-T",
			Frequency:   1.0,
			Gain:        0.254712175,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"FBA-ES-T-DECK": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "FBA-ES-T",
			Frequency:   1.0,
			Gain:        0.1273560875,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"FBA-23-DECK": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "FBA-ES-T",
			Frequency:   1.0,
			Gain:        0.1273560875,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"SDP": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "SDP",
			Frequency:   1.0,
			Gain:        0.1273560875,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"SSA-320": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "SSA-320",
			Frequency:   1.0,
			Gain:        0.1273560875,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"CUSP SENSOR": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "CUSP",
			Frequency:   1.0,
			Gain:        0.10197162129779282425,
			Scale:       1.0,
			InputUnits:  "M/S**2",
			OutputUnits: "V",
		},
	},
	"270-600/12V": []ResponseStage{
		ResponseStage{
			Type:        "poly",
			Lookup:      "270-600/12V",
			Scale:       1.0,
			Frequency:   1.0,
			InputUnits:  "hPa",
			OutputUnits: "V",
		},
	},
	"270-600/24V": []ResponseStage{
		ResponseStage{
			Type:        "poly",
			Lookup:      "270-600/24V",
			Scale:       1.0,
			Frequency:   1.0,
			InputUnits:  "hPa",
			OutputUnits: "V",
		},
	},
	"L4C": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "L4C",
			Frequency:   15.0,
			Gain:        177.8,
			Scale:       1.0,
			InputUnits:  "M/S",
			OutputUnits: "V",
		},
	},
	"LE-3Dlite": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "LE-3Dlite",
			Frequency:   15.0,
			Gain:        400.0,
			Scale:       1.0,
			InputUnits:  "M/S",
			OutputUnits: "V",
		},
	},
	"LE-3DliteMkII": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "LE-3DliteMkII",
			Frequency:   15.0,
			Gain:        400.0,
			Scale:       1.0,
			InputUnits:  "M/S",
			OutputUnits: "V",
		},
	},
	"2_Hz_Duke_Malin_Seismometer": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "MALIN",
			Frequency:   15.0,
			Gain:        62.2,
			Scale:       1.0,
			InputUnits:  "M/S",
			OutputUnits: "V",
		},
	},
	//,
	// dataloggers,
	//,
	"Q330_FLbelow100-1": []ResponseStage{
		ResponseStage{
			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  1.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330_FLbelow100-1",
			SampleRate:  1.0,
			Delay:       15.930462,
			Correction:  15.930462,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330_FLbelow100-100": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			SampleRate:  100.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330_FLbelow100-100",
			SampleRate:  100.0,
			Delay:       0.04167,
			Correction:  0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330_FLbelow100-200-PREAMP": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "Q330-PREAMP",
			Frequency:   15.0,
			Gain:        30.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  200.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330_FLbelow100-200",
			SampleRate:  200.0,
			Delay:       0.020462,
			Correction:  0.020462,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330_FLbelow100-200": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  200.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330_FLbelow100-200",
			SampleRate:  200.0,
			Delay:       0.020462,
			Correction:  0.020462,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330S+_FLbelow100-100-PREAMP32": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "Q330-PREAMP",
			Frequency:   15.0,
			Gain:        32.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  100.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330S+_FLbelow100-100",
			SampleRate:  100.0,
			Delay:       0.04167,
			Correction:  0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330S+_FLbelow100-1": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  1.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330S+_FLbelow100-1",
			SampleRate:  1.0,
			Delay:       15.930462,
			Correction:  15.930462,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"Q330S+_FLbelow100-100": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  100.0,
			Gain:        419430.4,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:        "fir",
			Lookup:      "Q330S+_FLbelow100-100",
			SampleRate:  100.0,
			Delay:       0.04167,
			Correction:  0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"EARSS-50": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "EARSS-GAIN",
			Frequency:   15.0,
			Gain:        1.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "paz",
			Lookup:      "EARSS-50Hz",
			Frequency:   15.0,
			Gain:        1.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  50.0,
			Gain:        104857.60,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		},
	},
	"EARSS-100": []ResponseStage{
		ResponseStage{

			Type:        "paz",
			Lookup:      "EARSS-GAIN",
			Frequency:   15.0,
			Gain:        1.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "paz",
			Lookup:      "EARSS-100Hz",
			Frequency:   15.0,
			Gain:        1.0,
			Scale:       1.0,
			InputUnits:  "V",
			OutputUnits: "V",
		}, ResponseStage{
			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  100.0,
			Gain:        104857.60,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		},
	},
	"ALTUS-200": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  2000.0,
			Gain:        3355443.2,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "ALTUS_A200",
			SampleRate: 400.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "ALTUS_BNC",
			SampleRate: 200.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"CUSP-200": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  200.0,
			Gain:        1000000.0,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		},
	},
	"BASALT-50": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  30000.0,
			Gain:        1677721.6,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A5-50-S5C",
			SampleRate: 6000.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A3-50",
			SampleRate: 2000.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A4-50",
			SampleRate: 500.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A5-50",
			SampleRate: 100.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_B2-80",
			SampleRate: 50.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
	"BASALT-200": []ResponseStage{
		ResponseStage{

			Type:        "a2d",
			Lookup:      "A2D",
			Decimate:    1,
			SampleRate:  30000.0,
			Gain:        1677721.6,
			InputUnits:  "V",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A5-50-S5C",
			SampleRate: 6000.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A3-50",
			SampleRate: 2000.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_A5-50",
			SampleRate: 400.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		}, ResponseStage{
			Type:       "fir",
			Lookup:     "BASALT_B2-80",
			SampleRate: 200.0,
			//Delay      : 0.04167,
			//Correction : 0.04167,
			InputUnits:  "COUNTS",
			OutputUnits: "COUNTS",
		},
	},
}

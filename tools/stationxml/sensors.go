package main

var SensorModels map[string]SensorModel = map[string]SensorModel{
	"2 Hz Duke Malin Seismometer": SensorModel{
		Type:        "Short Period Borehole Seismometer",
		Description: "Duke Malin",
		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"270-600/12V": SensorModel{
		Type:        "Barometer",
		Description: "Setra 270",
		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"270-600/24V": SensorModel{
		Type:        "Barometer",
		Description: "Setra 270",
		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"270-800/12V": SensorModel{
		Type:        "Barometer",
		Description: "Setra 270",
		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"4.5 Hz Duke Malin Seismometer": SensorModel{
		Type:        "Short Period Borehole Seismometer",
		Description: "Duke Malin",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3ESP": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-3ESP",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3ESP-Z": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-3ESP",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3ESPC": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-3ESP",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3T": SensorModel{
		Type:         "Broadband Seismometer",
		Description:  "CMG-3T",
		Manufacturer: "Guralp",
		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3TB": SensorModel{
		Type:         "Broadband Seismometer",
		Description:  "CMG-3T",
		Manufacturer: "Guralp",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-3TB-GN": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-3T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-40T-30S": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-40T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-40T-60S": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "CMG-40T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CMG-6T": SensorModel{
		Type:        "Accelerometer",
		Description: "CMG-6T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"CUSP-Me SENSOR": SensorModel{
		Type:        "Building Array Sensor",
		Description: "CUSP Me",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSP3A SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP 3A",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSP3B SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP 3B",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSP3C SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP 3C",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSP3C3 SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP 3C",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSP3D SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP 3D",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"CUSPM SENSOR": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "CUSP M",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"Druck PTX-1830": SensorModel{
		Type:        "Pressure Sensor",
		Description: "Druck PTX",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"Druck PTX-1830-LAND": SensorModel{
		Type:        "Pressure Sensor",
		Description: "Druck PTX",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"FBA-23-DECK": SensorModel{
		Type:        "Accelerometer",
		Description: "FBA-23",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 180.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     90.0,
			},
			SensorComponent{
				Azimuth: 270.0,
				Dip:     0.0,
			},
		},
	},
	"FBA-ES-T": SensorModel{
		Type:        "Accelerometer",
		Description: "FBA-ES-T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"FBA-ES-T-BASALT": SensorModel{
		Type:        "Accelerometer",
		Description: "FBA-ES-T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"FBA-ES-T-DECK": SensorModel{
		Type:        "Accelerometer",
		Description: "FBA-ES-T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"FBA-ES-T-ISO": SensorModel{
		Type:        "Accelerometer",
		Description: "FBA-ES-T",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"GS-11D seismometer": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "GS-11D",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"IESE HS-1-LT Mini 4.5Hz": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "HS1-1-LT",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"IESE S10g-4.5 (with preamp)": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "S10g-4.5",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"IESE S10g-4.5 (without preamp)": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "S10g-4.5",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"IESE S31f-15 (with preamp)": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "S31f-15",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"InfraBSU microphone": SensorModel{
		Type:        "Mirophone",
		Description: "InfraBSU",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"Kinemetrics SBEPI": SensorModel{
		Type:        "Accelerometer",
		Description: "SBEPI",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"L4C": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "L4C",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"L4C-3D": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "L4C-3D",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"LE-3Dlite": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "LE-3Dlite",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"LE-3DliteMkII": SensorModel{
		Type:        "Short Period Seismometer",
		Description: "LE-3DliteMkII",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"LM35": SensorModel{
		Description: "LM35",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
	"Nanometrics Trillium 120QA": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "Trillium 120",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"SDP": SensorModel{
		Description: "SDP",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"SSA-320": SensorModel{
		Type:        "Strong Motion Sensor",
		Description: "SSA-320",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
		},
	},
	"STS-2": SensorModel{
		Type:        "Broadband Seismometer",
		Description: "STS-2",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     -90.0,
			},
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
			SensorComponent{
				Azimuth: 90.0,
				Dip:     0.0,
			},
		},
	},
	"Scout Hydrophone": SensorModel{
		Type:        "Hydrophone",
		Description: "Scout",

		Components: []SensorComponent{
			SensorComponent{
				Azimuth: 0.0,
				Dip:     0.0,
			},
		},
	},
}

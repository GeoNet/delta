package main

// list sensors that can be used
var blessedSensorList = []string{
	"FBA-ES-T",
	"FBA-ES-T-ISO",
	"FBA-ES-T-BASALT",
	"FBA-ES-T-OBSIDIAN",
	"FBA-ES-T-ETNA-2",
	"CUSP3C",
	"CUSP3D",
	"CMG-3ESP",
	"CMG-3ESPC",
	"CMG-3TB",
	"CMG-3TB-GN",
	"STS-2",
	"Trillium 120QA",
	"Trillium Compact 120PH-2",
	"L4C-3D",
	"L4C",
	"LE-3Dlite",
	"LE-3DliteMkII",
	"Titan",
	"EQR120",
}

// list dataloggers that can be used
var blessedDataloggerList = []string{
	"Q330/3",
	"Q330/6",
	"Q330HR/6",
	"Q330S/3",
	"Q330S/6",
	"Q330HRS/6",
	"BASALT",
	"OBSIDIAN",
	"ETNA 2",
	"CUSP3D",
	"CUSP3C",
	"Obsidian 4X Datalogger",
	"TitanSMA",
	"EQR120",
}

var blessedSensors = make(map[string]bool)
var blessedDataloggers = make(map[string]bool)

func init() {
	for _, s := range blessedSensorList {
		blessedSensors[s] = true
	}
	for _, s := range blessedDataloggerList {
		blessedDataloggers[s] = true
	}
}

func isBlessedSensor(model string) bool {
	if _, ok := blessedSensors[model]; ok {
		return true
	}
	return false
}

func isBlessedDatalogger(model string) bool {
	if _, ok := blessedDataloggers[model]; ok {
		return true
	}
	return false
}

package main

import (
	"log"
	"os"
)

func main() {

	generate := Generate{
		Fields: map[string]struct {
			Key string
		}{
			"assets":              {"Asset"},
			"calibrations":        {"Calibration"},
			"connections":         {"Connection"},
			"constituents":        {"Constituent"},
			"deployedDataloggers": {"DeployedDatalogger"},
			"deployedReceivers":   {"DeployedReceiver"},
			"doases":              {"InstalledDoas"},
			"domains":             {"Domain"},
			"features":            {"Feature"},
			"firmwareHistory":     {"FirmwareHistory"},
			"gains":               {"Gain"},
			"gauges":              {"Gauge"},
			"installedAntennas":   {"InstalledAntenna"},
			"installedCameras":    {"InstalledCamera"},
			"installedMetSensors": {"InstalledMetSensor"},
			"installedRadomes":    {"InstalledRadome"},
			"installedRecorders":  {"InstalledRecorder"},
			"installedSensors":    {"InstalledSensor"},
			"marks":               {"Mark"},
			"monuments":           {"Monument"},
			"mounts":              {"Mount"},
			"networks":            {"Network"},
			"placenames":          {"Placename"},
			"sessions":            {"Session"},
			"sites":               {"Site"},
			"stations":            {"Station"},
			"streams":             {"Stream"},
			"views":               {"View"},
			"visibilities":        {"Visibility"},
		},
		Lookup: map[string]struct {
			Key    string
			Fields []string
		}{
			"assets":     {"Asset", []string{"make", "model", "serial"}},
			"marks":      {"Mark", []string{"code"}},
			"monuments":  {"Monument", []string{"mark"}},
			"mounts":     {"Mount", []string{"code"}},
			"networks":   {"Network", []string{"code"}},
			"placenames": {"Placename", []string{"name"}},
			"sites":      {"Site", []string{"station", "location"}},
			"stations":   {"Station", []string{"code"}},
			"views":      {"View", []string{"mount", "code"}},
		},
	}

	if err := generate.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

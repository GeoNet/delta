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
			"channels":            {"Channel"},
			"classes":             {"Class"},
			"components":          {"Component"},
			"connections":         {"Connection"},
			"citations":           {"Citation"},
			"constituents":        {"Constituent"},
			"darts":               {"Dart"},
			"datasets":            {"Dataset"},
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
			"points":              {"Point"},
			"polarities":          {"Polarity"},
			"preamps":             {"Preamp"},
			"samples":             {"Sample"},
			"sessions":            {"Session"},
			"sites":               {"Site"},
			"stations":            {"Station"},
			"streams":             {"Stream"},
			"telemetries":         {"Telemetry"},
			"timings":             {"Timing"},
			"views":               {"View"},
			"visibilities":        {"Visibility"},
		},
		Lookup: map[string]struct {
			Key    string
			Fields []string
		}{
			"assets":       {"Asset", []string{"make", "model", "serial"}},
			"citations":    {"Citation", []string{"key"}},
			"classes":      {"Class", []string{"station"}},
			"constituents": {"Constituent", []string{"gauge", "location"}},
			"darts":        {"Dart", []string{"station"}},
			"datasets":     {"Dataset", []string{"domain", "network"}},
			"domains":      {"Domain", []string{"name"}},
			"marks":        {"Mark", []string{"code"}},
			"monuments":    {"Monument", []string{"mark"}},
			"mounts":       {"Mount", []string{"code"}},
			"networks":     {"Network", []string{"code"}},
			"placenames":   {"Placename", []string{"name"}},
			"points":       {"Point", []string{"sample", "location"}},
			"samples":      {"Sample", []string{"code"}},
			"sites":        {"Site", []string{"station", "location"}},
			"stations":     {"Station", []string{"code"}},
			"views":        {"View", []string{"mount", "code"}},
		},
	}

	if err := generate.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

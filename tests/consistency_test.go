package delta_test

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GeoNet/delta/meta"
)

var testConsistency = map[string]func(path string, list meta.List) func(t *testing.T){

	"check file consistency": func(path string, list meta.List) func(t *testing.T) {
		return func(t *testing.T) {

			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("unable to read %s file: %v", path, err)
			}

			var buf bytes.Buffer
			if err := csv.NewWriter(&buf).WriteAll(meta.EncodeList(list)); err != nil {
				t.Fatalf("unable to decode %s file: %v", path, err)
			}

			if string(raw) != buf.String() {
				t.Errorf("unexpected %s content -got/+exp\n%s", filepath.Base(path), cmp.Diff(string(raw), buf.String()))
			}
		}
	},
}

func TestConsistency(t *testing.T) {

	files := map[string]struct {
		f string
		l meta.List
	}{
		"darts":     {f: "../network/darts.csv", l: &meta.DartList{}},
		"marks":     {f: "../network/marks.csv", l: &meta.MarkList{}},
		"monuments": {f: "../network/monuments.csv", l: &meta.MonumentList{}},
		"mounts":    {f: "../network/mounts.csv", l: &meta.MountList{}},
		"networks":  {f: "../network/networks.csv", l: &meta.NetworkList{}},
		"points":    {f: "../network/points.csv", l: &meta.PointList{}},
		"samples":   {f: "../network/samples.csv", l: &meta.SampleList{}},
		"sites":     {f: "../network/sites.csv", l: &meta.SiteList{}},
		"stations":  {f: "../network/stations.csv", l: &meta.StationList{}},
		"views":     {f: "../network/views.csv", l: &meta.ViewList{}},

		"antennas":     {f: "../install/antennas.csv", l: &meta.InstalledAntennaList{}},
		"calibrations": {f: "../install/calibrations.csv", l: &meta.CalibrationList{}},
		"cameras":      {f: "../install/cameras.csv", l: &meta.InstalledCameraList{}},
		"channels":     {f: "../install/channels.csv", l: &meta.ChannelList{}},
		"components":   {f: "../install/components.csv", l: &meta.ComponentList{}},
		"connections":  {f: "../install/connections.csv", l: &meta.ConnectionList{}},
		"dataloggers":  {f: "../install/dataloggers.csv", l: &meta.DeployedDataloggerList{}},
		"doases":       {f: "../install/doases.csv", l: &meta.InstalledDoasList{}},
		"firmware":     {f: "../install/firmware.csv", l: &meta.FirmwareHistoryList{}},
		"gains":        {f: "../install/gains.csv", l: &meta.GainList{}},
		"metsensors":   {f: "../install/metsensors.csv", l: &meta.InstalledMetSensorList{}},
		//"polarities":   {f: "../install/polarities.csv", l: &meta.PolarityList{}},
		"preamps":     {f: "../install/preamps.csv", l: &meta.PreampList{}},
		"radomes":     {f: "../install/radomes.csv", l: &meta.InstalledRadomeList{}},
		"receivers":   {f: "../install/receivers.csv", l: &meta.DeployedReceiverList{}},
		"recorders":   {f: "../install/recorders.csv", l: &meta.InstalledRecorderList{}},
		"sensors":     {f: "../install/sensors.csv", l: &meta.InstalledSensorList{}},
		"sessions":    {f: "../install/sessions.csv", l: &meta.SessionList{}},
		"streams":     {f: "../install/streams.csv", l: &meta.StreamList{}},
		"telemetries": {f: "../install/telemetries.csv", l: &meta.TelemetryList{}},
		"timings":     {f: "../install/timings.csv", l: &meta.TimingList{}},

		//"classes":      {f: "../environment/classes.csv", l: &meta.ClassList{}},
		"constituents": {f: "../environment/constituents.csv", l: &meta.ConstituentList{}},
		"features":     {f: "../environment/features.csv", l: &meta.FeatureList{}},
		"gauges":       {f: "../environment/gauges.csv", l: &meta.GaugeList{}},
		"notes":        {f: "../environment/notes.csv", l: &meta.NoteList{}},
		//"placenames":   {f: "../environment/placenames.csv", l: &meta.PlacenameList{}},
		"visibility": {f: "../environment/visibility.csv", l: &meta.VisibilityList{}},

		"citations": {f: "../references/citations.csv", l: &meta.CitationList{}},
	}

	for f, v := range files {

		if err := meta.LoadList(v.f, v.l); err != nil {
			t.Fatalf("unable to load list file %s: %v", v.f, err)
		}

		sort.Sort(v.l)

		for k, fn := range testConsistency {
			t.Run(k+": "+f, fn(v.f, v.l))
		}
	}
}

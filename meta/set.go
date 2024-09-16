package meta

//go:generate bash -c "go run generate/*.go | gofmt -s > set_auto.go; test -s set_auto.go || rm set_auto.go"

import (
	"fmt"
	"io/fs"
	"sort"
)

const (
	AssetFiles = "assets/*.csv"

	DartsFile     = "network/darts.csv"
	MarksFile     = "network/marks.csv"
	MonumentsFile = "network/monuments.csv"
	MountsFile    = "network/mounts.csv"
	NetworksFile  = "network/networks.csv"
	PointsFile    = "network/points.csv"
	SamplesFile   = "network/samples.csv"
	SitesFile     = "network/sites.csv"
	StationsFile  = "network/stations.csv"
	ViewsFile     = "network/views.csv"

	AntennasFile     = "install/antennas.csv"
	CalibrationsFile = "install/calibrations.csv"
	CamerasFile      = "install/cameras.csv"
	ChannelsFile     = "install/channels.csv"
	ComponentsFile   = "install/components.csv"
	ConnectionsFile  = "install/connections.csv"
	DataloggersFile  = "install/dataloggers.csv"
	DoasesFile       = "install/doases.csv"
	FirmwareFile     = "install/firmware.csv"
	GainsFile        = "install/gains.csv"
	MetsensorsFile   = "install/metsensors.csv"
	PolaritiesFile   = "install/polarities.csv"
	PreampsFile      = "install/preamps.csv"
	RadomesFile      = "install/radomes.csv"
	ReceiversFile    = "install/receivers.csv"
	RecordersFile    = "install/recorders.csv"
	SensorsFile      = "install/sensors.csv"
	SessionsFile     = "install/sessions.csv"
	StreamsFile      = "install/streams.csv"
	TelemetriesFile  = "install/telemetries.csv"
	TimingsFile      = "install/timings.csv"

	ClassesFile      = "environment/classes.csv"
	ConstituentsFile = "environment/constituents.csv"
	FeaturesFile     = "environment/features.csv"
	GaugesFile       = "environment/gauges.csv"
	PlacenamesFile   = "environment/placenames.csv"
	VisibilityFile   = "environment/visibility.csv"

	CitationsFile = "references/citations.csv"
	MethodsFile   = "references/methods.csv"
)

// SetPathMap is used to manipulate the filepath inside the Set.
type SetPathMap func(s string) string

// Set allows for extracting and unmarshalling the base delta csv files,
// optional SetPathMap functions can be given to alter the expected default
// file paths prior to reading from the FS set. This is useful for testing
// or using a non-standard file layout.
type Set struct {
	assets AssetList

	darts     DartList
	marks     MarkList
	monuments MonumentList
	mounts    MountList
	networks  NetworkList
	points    PointList
	samples   SampleList
	sites     SiteList
	stations  StationList
	views     ViewList

	installedAntennas   InstalledAntennaList
	calibrations        CalibrationList
	installedCameras    InstalledCameraList
	channels            ChannelList
	components          ComponentList
	connections         ConnectionList
	deployedDataloggers DeployedDataloggerList
	doases              InstalledDoasList
	firmwareHistory     FirmwareHistoryList
	gains               GainList
	installedMetSensors InstalledMetSensorList
	polarities          PolarityList
	preamps             PreampList
	installedRadomes    InstalledRadomeList
	deployedReceivers   DeployedReceiverList
	installedRecorders  InstalledRecorderList
	installedSensors    InstalledSensorList
	sessions            SessionList
	streams             StreamList
	telemetries         TelemetryList
	timings             TimingList

	classes      ClassList
	constituents ConstituentList
	features     FeatureList
	gauges       GaugeList
	placenames   PlacenameList
	visibilities VisibilityList

	citations CitationList
	methods   MethodList
}

func (s *Set) files() map[string]List {
	return map[string]List{
		AssetFiles: &s.assets,

		DartsFile:     &s.darts,
		MarksFile:     &s.marks,
		MonumentsFile: &s.monuments,
		MountsFile:    &s.mounts,
		NetworksFile:  &s.networks,
		PointsFile:    &s.points,
		SamplesFile:   &s.samples,
		SitesFile:     &s.sites,
		StationsFile:  &s.stations,
		ViewsFile:     &s.views,

		AntennasFile:     &s.installedAntennas,
		CalibrationsFile: &s.calibrations,
		CamerasFile:      &s.installedCameras,
		ChannelsFile:     &s.channels,
		ComponentsFile:   &s.components,
		ConnectionsFile:  &s.connections,
		DataloggersFile:  &s.deployedDataloggers,
		DoasesFile:       &s.doases,
		FirmwareFile:     &s.firmwareHistory,
		GainsFile:        &s.gains,
		MetsensorsFile:   &s.installedMetSensors,
		PolaritiesFile:   &s.polarities,
		PreampsFile:      &s.preamps,
		RadomesFile:      &s.installedRadomes,
		ReceiversFile:    &s.deployedReceivers,
		RecordersFile:    &s.installedRecorders,
		SensorsFile:      &s.installedSensors,
		SessionsFile:     &s.sessions,
		StreamsFile:      &s.streams,
		TelemetriesFile:  &s.telemetries,
		TimingsFile:      &s.timings,

		ClassesFile:      &s.classes,
		ConstituentsFile: &s.constituents,
		FeaturesFile:     &s.features,
		GaugesFile:       &s.gauges,
		PlacenamesFile:   &s.placenames,
		VisibilityFile:   &s.visibilities,

		CitationsFile: &s.citations,
		MethodsFile:   &s.methods,
	}
}

// NewSet returns a Set pointer for the given FS and optional SetPathMap
// functions to manipulate the internal csv file paths.
func NewSet(fsys fs.FS, maps ...SetPathMap) (*Set, error) {
	var set Set

	for path, list := range set.files() {
		switch path {
		case AssetFiles:
			for _, p := range maps {
				path = p(path)
			}
			names, err := fs.Glob(fsys, path)
			if err != nil {
				return nil, fmt.Errorf("glob error %s: %w", path, err)
			}
			for _, name := range names {
				var assets AssetList
				data, err := fs.ReadFile(fsys, name)
				if err != nil {
					return nil, fmt.Errorf("read error %s: %w", path, err)
				}
				if err := UnmarshalList(data, &assets); err != nil {
					return nil, fmt.Errorf("unmarshal error %s: %w", path, err)
				}
				set.assets = append(set.assets, assets...)
			}
			sort.Sort(set.assets)
		default:
			for _, p := range maps {
				path = p(path)
			}
			names, err := fs.Glob(fsys, path)
			if err != nil {
				return nil, fmt.Errorf("glob error %s: %w", path, err)
			}
			for _, name := range names {
				data, err := fs.ReadFile(fsys, name)
				if err != nil {
					return nil, fmt.Errorf("read error %s: %w", path, err)
				}
				if err := UnmarshalList(data, list); err != nil {
					return nil, fmt.Errorf("unmarshal error %s: %w", path, err)
				}
			}
			sort.Sort(list)
		}
	}

	return &set, nil
}

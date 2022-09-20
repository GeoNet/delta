package meta

//go:generate bash -c "go run generate/*.go | gofmt -s > set_auto.go; test -s set_auto.go || rm set_auto.go"

import (
	"fmt"
	"io/fs"
	"sort"
)

const (
	AssetFiles = "assets/*.csv"

	MarksFile     = "network/marks.csv"
	MonumentsFile = "network/monuments.csv"
	MountsFile    = "network/mounts.csv"
	NetworksFile  = "network/networks.csv"
	SitesFile     = "network/sites.csv"
	StationsFile  = "network/stations.csv"
	ViewsFile     = "network/views.csv"

	AntennasFile     = "install/antennas.csv"
	CamerasFile      = "install/cameras.csv"
	ConnectionsFile  = "install/connections.csv"
	DataloggersFile  = "install/dataloggers.csv"
	DoasesFile       = "install/doases.csv"
	FirmwareFile     = "install/firmware.csv"
	GainsFile        = "install/gains.csv"
	CalibrationsFile = "install/calibrations.csv"
	MetsensorsFile   = "install/metsensors.csv"
	RadomesFile      = "install/radomes.csv"
	ReceiversFile    = "install/receivers.csv"
	RecordersFile    = "install/recorders.csv"
	SensorsFile      = "install/sensors.csv"
	SessionsFile     = "install/sessions.csv"
	StreamsFile      = "install/streams.csv"

	ConstituentsFile = "environment/constituents.csv"
	FeaturesFile     = "environment/features.csv"
	GaugesFile       = "environment/gauges.csv"
	VisibilityFile   = "environment/visibility.csv"
	PlacenamesFile   = "environment/placenames.csv"

	ChannelsFile   = "install/channels.csv"
	ComponentsFile = "install/components.csv"
	CodenamesFile  = "install/codenames.csv"

	CitationsFile = "references/citations.csv"
)

// SetPathMap is used to manipulate the filepath inside the Set.
type SetPathMap func(s string) string

// Set allows for extracting and unmarshalling the base delta csv files,
// optional SetPathMap functions can be given to alter the expected default
// file paths prior to reading from the FS set. This is useful for testing
// or using a non-standard file layout.
type Set struct {
	assets AssetList

	marks     MarkList
	monuments MonumentList
	mounts    MountList
	networks  NetworkList
	sites     SiteList
	stations  StationList
	views     ViewList

	installedAntennas   InstalledAntennaList
	installedCameras    InstalledCameraList
	connections         ConnectionList
	deployedDataloggers DeployedDataloggerList
	doases              InstalledDoasList
	firmwareHistory     FirmwareHistoryList
	gains               GainList
	calibrations        CalibrationList
	installedMetSensors InstalledMetSensorList
	installedRadomes    InstalledRadomeList

	channels   ChannelList
	components ComponentList

	installedSensors   InstalledSensorList
	installedRecorders InstalledRecorderList
	deployedReceivers  DeployedReceiverList
	sessions           SessionList
	streams            StreamList

	constituents ConstituentList
	features     FeatureList
	gauges       GaugeList
	visibilities VisibilityList
	placenames   PlacenameList
	citations    CitationList
}

func (s *Set) files() map[string]List {
	return map[string]List{
		AssetFiles: &s.assets,

		MarksFile:     &s.marks,
		MonumentsFile: &s.monuments,
		MountsFile:    &s.mounts,
		NetworksFile:  &s.networks,
		SitesFile:     &s.sites,
		StationsFile:  &s.stations,
		ViewsFile:     &s.views,

		AntennasFile:     &s.installedAntennas,
		CamerasFile:      &s.installedCameras,
		ConnectionsFile:  &s.connections,
		DataloggersFile:  &s.deployedDataloggers,
		DoasesFile:       &s.doases,
		FirmwareFile:     &s.firmwareHistory,
		GainsFile:        &s.gains,
		CalibrationsFile: &s.calibrations,
		MetsensorsFile:   &s.installedMetSensors,
		RadomesFile:      &s.installedRadomes,
		ReceiversFile:    &s.deployedReceivers,
		RecordersFile:    &s.installedRecorders,
		SensorsFile:      &s.installedSensors,
		SessionsFile:     &s.sessions,
		StreamsFile:      &s.streams,

		ChannelsFile:   &s.channels,
		ComponentsFile: &s.components,

		ConstituentsFile: &s.constituents,
		FeaturesFile:     &s.features,
		GaugesFile:       &s.gauges,
		VisibilityFile:   &s.visibilities,
		PlacenamesFile:   &s.placenames,
		CitationsFile:    &s.citations,
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

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  To update any changes, run "go generate" in the main project
 *  directory and then commit this file.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

package meta

// Assets is a helper function to return a slice copy of Asset values.
func (s Set) Assets() []Asset {
	assets := make([]Asset, len(s.assets))
	copy(assets, s.assets)
	return assets
}

// Calibrations is a helper function to return a slice copy of Calibration values.
func (s Set) Calibrations() []Calibration {
	calibrations := make([]Calibration, len(s.calibrations))
	copy(calibrations, s.calibrations)
	return calibrations
}

// Channels is a helper function to return a slice copy of Channel values.
func (s Set) Channels() []Channel {
	channels := make([]Channel, len(s.channels))
	copy(channels, s.channels)
	return channels
}

// Citations is a helper function to return a slice copy of Citation values.
func (s Set) Citations() []Citation {
	citations := make([]Citation, len(s.citations))
	copy(citations, s.citations)
	return citations
}

// Components is a helper function to return a slice copy of Component values.
func (s Set) Components() []Component {
	components := make([]Component, len(s.components))
	copy(components, s.components)
	return components
}

// Connections is a helper function to return a slice copy of Connection values.
func (s Set) Connections() []Connection {
	connections := make([]Connection, len(s.connections))
	copy(connections, s.connections)
	return connections
}

// Constituents is a helper function to return a slice copy of Constituent values.
func (s Set) Constituents() []Constituent {
	constituents := make([]Constituent, len(s.constituents))
	copy(constituents, s.constituents)
	return constituents
}

// DeployedDataloggers is a helper function to return a slice copy of DeployedDatalogger values.
func (s Set) DeployedDataloggers() []DeployedDatalogger {
	deployedDataloggers := make([]DeployedDatalogger, len(s.deployedDataloggers))
	copy(deployedDataloggers, s.deployedDataloggers)
	return deployedDataloggers
}

// DeployedReceivers is a helper function to return a slice copy of DeployedReceiver values.
func (s Set) DeployedReceivers() []DeployedReceiver {
	deployedReceivers := make([]DeployedReceiver, len(s.deployedReceivers))
	copy(deployedReceivers, s.deployedReceivers)
	return deployedReceivers
}

// Doases is a helper function to return a slice copy of InstalledDoas values.
func (s Set) Doases() []InstalledDoas {
	doases := make([]InstalledDoas, len(s.doases))
	copy(doases, s.doases)
	return doases
}

// Features is a helper function to return a slice copy of Feature values.
func (s Set) Features() []Feature {
	features := make([]Feature, len(s.features))
	copy(features, s.features)
	return features
}

// FirmwareHistory is a helper function to return a slice copy of FirmwareHistory values.
func (s Set) FirmwareHistory() []FirmwareHistory {
	firmwareHistory := make([]FirmwareHistory, len(s.firmwareHistory))
	copy(firmwareHistory, s.firmwareHistory)
	return firmwareHistory
}

// Gains is a helper function to return a slice copy of Gain values.
func (s Set) Gains() []Gain {
	gains := make([]Gain, len(s.gains))
	copy(gains, s.gains)
	return gains
}

// Gauges is a helper function to return a slice copy of Gauge values.
func (s Set) Gauges() []Gauge {
	gauges := make([]Gauge, len(s.gauges))
	copy(gauges, s.gauges)
	return gauges
}

// InstalledAntennas is a helper function to return a slice copy of InstalledAntenna values.
func (s Set) InstalledAntennas() []InstalledAntenna {
	installedAntennas := make([]InstalledAntenna, len(s.installedAntennas))
	copy(installedAntennas, s.installedAntennas)
	return installedAntennas
}

// InstalledCameras is a helper function to return a slice copy of InstalledCamera values.
func (s Set) InstalledCameras() []InstalledCamera {
	installedCameras := make([]InstalledCamera, len(s.installedCameras))
	copy(installedCameras, s.installedCameras)
	return installedCameras
}

// InstalledMetSensors is a helper function to return a slice copy of InstalledMetSensor values.
func (s Set) InstalledMetSensors() []InstalledMetSensor {
	installedMetSensors := make([]InstalledMetSensor, len(s.installedMetSensors))
	copy(installedMetSensors, s.installedMetSensors)
	return installedMetSensors
}

// InstalledRadomes is a helper function to return a slice copy of InstalledRadome values.
func (s Set) InstalledRadomes() []InstalledRadome {
	installedRadomes := make([]InstalledRadome, len(s.installedRadomes))
	copy(installedRadomes, s.installedRadomes)
	return installedRadomes
}

// InstalledRecorders is a helper function to return a slice copy of InstalledRecorder values.
func (s Set) InstalledRecorders() []InstalledRecorder {
	installedRecorders := make([]InstalledRecorder, len(s.installedRecorders))
	copy(installedRecorders, s.installedRecorders)
	return installedRecorders
}

// InstalledSensors is a helper function to return a slice copy of InstalledSensor values.
func (s Set) InstalledSensors() []InstalledSensor {
	installedSensors := make([]InstalledSensor, len(s.installedSensors))
	copy(installedSensors, s.installedSensors)
	return installedSensors
}

// Marks is a helper function to return a slice copy of Mark values.
func (s Set) Marks() []Mark {
	marks := make([]Mark, len(s.marks))
	copy(marks, s.marks)
	return marks
}

// Monuments is a helper function to return a slice copy of Monument values.
func (s Set) Monuments() []Monument {
	monuments := make([]Monument, len(s.monuments))
	copy(monuments, s.monuments)
	return monuments
}

// Mounts is a helper function to return a slice copy of Mount values.
func (s Set) Mounts() []Mount {
	mounts := make([]Mount, len(s.mounts))
	copy(mounts, s.mounts)
	return mounts
}

// Networks is a helper function to return a slice copy of Network values.
func (s Set) Networks() []Network {
	networks := make([]Network, len(s.networks))
	copy(networks, s.networks)
	return networks
}

// Placenames is a helper function to return a slice copy of Placename values.
func (s Set) Placenames() []Placename {
	placenames := make([]Placename, len(s.placenames))
	copy(placenames, s.placenames)
	return placenames
}

// Sessions is a helper function to return a slice copy of Session values.
func (s Set) Sessions() []Session {
	sessions := make([]Session, len(s.sessions))
	copy(sessions, s.sessions)
	return sessions
}

// Sites is a helper function to return a slice copy of Site values.
func (s Set) Sites() []Site {
	sites := make([]Site, len(s.sites))
	copy(sites, s.sites)
	return sites
}

// Stations is a helper function to return a slice copy of Station values.
func (s Set) Stations() []Station {
	stations := make([]Station, len(s.stations))
	copy(stations, s.stations)
	return stations
}

// Streams is a helper function to return a slice copy of Stream values.
func (s Set) Streams() []Stream {
	streams := make([]Stream, len(s.streams))
	copy(streams, s.streams)
	return streams
}

// Views is a helper function to return a slice copy of View values.
func (s Set) Views() []View {
	views := make([]View, len(s.views))
	copy(views, s.views)
	return views
}

// Visibilities is a helper function to return a slice copy of Visibility values.
func (s Set) Visibilities() []Visibility {
	visibilities := make([]Visibility, len(s.visibilities))
	copy(visibilities, s.visibilities)
	return visibilities
}

// Asset is a helper function to return a Asset value and true if one exists.
func (s Set) Asset(make, model, serial string) (Asset, bool) {
	for _, v := range s.assets {
		if make != v.Make {
			continue
		}
		if model != v.Model {
			continue
		}
		if serial != v.Serial {
			continue
		}
		return v, true
	}
	return Asset{}, false
}

// Mark is a helper function to return a Mark value and true if one exists.
func (s Set) Mark(code string) (Mark, bool) {
	for _, v := range s.marks {
		if code != v.Code {
			continue
		}
		return v, true
	}
	return Mark{}, false
}

// Monument is a helper function to return a Monument value and true if one exists.
func (s Set) Monument(mark string) (Monument, bool) {
	for _, v := range s.monuments {
		if mark != v.Mark {
			continue
		}
		return v, true
	}
	return Monument{}, false
}

// Mount is a helper function to return a Mount value and true if one exists.
func (s Set) Mount(code string) (Mount, bool) {
	for _, v := range s.mounts {
		if code != v.Code {
			continue
		}
		return v, true
	}
	return Mount{}, false
}

// Network is a helper function to return a Network value and true if one exists.
func (s Set) Network(code string) (Network, bool) {
	for _, v := range s.networks {
		if code != v.Code {
			continue
		}
		return v, true
	}
	return Network{}, false
}

// Placename is a helper function to return a Placename value and true if one exists.
func (s Set) Placename(name string) (Placename, bool) {
	for _, v := range s.placenames {
		if name != v.Name {
			continue
		}
		return v, true
	}
	return Placename{}, false
}

// Site is a helper function to return a Site value and true if one exists.
func (s Set) Site(station, location string) (Site, bool) {
	for _, v := range s.sites {
		if station != v.Station {
			continue
		}
		if location != v.Location {
			continue
		}
		return v, true
	}
	return Site{}, false
}

// Station is a helper function to return a Station value and true if one exists.
func (s Set) Station(code string) (Station, bool) {
	for _, v := range s.stations {
		if code != v.Code {
			continue
		}
		return v, true
	}
	return Station{}, false
}

// View is a helper function to return a View value and true if one exists.
func (s Set) View(mount, code string) (View, bool) {
	for _, v := range s.views {
		if mount != v.Mount {
			continue
		}
		if code != v.Code {
			continue
		}
		return v, true
	}
	return View{}, false
}

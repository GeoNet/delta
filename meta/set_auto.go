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

// Assets is a helper function to return a slice of Asset values.
func (s Set) Assets() []Asset {
	return s.assets
}

// Calibrations is a helper function to return a slice of Calibration values.
func (s Set) Calibrations() []Calibration {
	return s.calibrations
}

// Channels is a helper function to return a slice of Channel values.
func (s Set) Channels() []Channel {
	return s.channels
}

// Citations is a helper function to return a slice of Citation values.
func (s Set) Citations() []Citation {
	return s.citations
}

// Classes is a helper function to return a slice of Class values.
func (s Set) Classes() []Class {
	return s.classes
}

// Components is a helper function to return a slice of Component values.
func (s Set) Components() []Component {
	return s.components
}

// Connections is a helper function to return a slice of Connection values.
func (s Set) Connections() []Connection {
	return s.connections
}

// Constituents is a helper function to return a slice of Constituent values.
func (s Set) Constituents() []Constituent {
	return s.constituents
}

// Darts is a helper function to return a slice of Dart values.
func (s Set) Darts() []Dart {
	return s.darts
}

// DeployedDataloggers is a helper function to return a slice of DeployedDatalogger values.
func (s Set) DeployedDataloggers() []DeployedDatalogger {
	return s.deployedDataloggers
}

// DeployedReceivers is a helper function to return a slice of DeployedReceiver values.
func (s Set) DeployedReceivers() []DeployedReceiver {
	return s.deployedReceivers
}

// Doases is a helper function to return a slice of InstalledDoas values.
func (s Set) Doases() []InstalledDoas {
	return s.doases
}

// Features is a helper function to return a slice of Feature values.
func (s Set) Features() []Feature {
	return s.features
}

// FirmwareHistory is a helper function to return a slice of FirmwareHistory values.
func (s Set) FirmwareHistory() []FirmwareHistory {
	return s.firmwareHistory
}

// Gains is a helper function to return a slice of Gain values.
func (s Set) Gains() []Gain {
	return s.gains
}

// Gauges is a helper function to return a slice of Gauge values.
func (s Set) Gauges() []Gauge {
	return s.gauges
}

// InstalledAntennas is a helper function to return a slice of InstalledAntenna values.
func (s Set) InstalledAntennas() []InstalledAntenna {
	return s.installedAntennas
}

// InstalledCameras is a helper function to return a slice of InstalledCamera values.
func (s Set) InstalledCameras() []InstalledCamera {
	return s.installedCameras
}

// InstalledMetSensors is a helper function to return a slice of InstalledMetSensor values.
func (s Set) InstalledMetSensors() []InstalledMetSensor {
	return s.installedMetSensors
}

// InstalledRadomes is a helper function to return a slice of InstalledRadome values.
func (s Set) InstalledRadomes() []InstalledRadome {
	return s.installedRadomes
}

// InstalledRecorders is a helper function to return a slice of InstalledRecorder values.
func (s Set) InstalledRecorders() []InstalledRecorder {
	return s.installedRecorders
}

// InstalledSensors is a helper function to return a slice of InstalledSensor values.
func (s Set) InstalledSensors() []InstalledSensor {
	return s.installedSensors
}

// Marks is a helper function to return a slice of Mark values.
func (s Set) Marks() []Mark {
	return s.marks
}

// Monuments is a helper function to return a slice of Monument values.
func (s Set) Monuments() []Monument {
	return s.monuments
}

// Mounts is a helper function to return a slice of Mount values.
func (s Set) Mounts() []Mount {
	return s.mounts
}

// Networks is a helper function to return a slice of Network values.
func (s Set) Networks() []Network {
	return s.networks
}

// Placenames is a helper function to return a slice of Placename values.
func (s Set) Placenames() []Placename {
	return s.placenames
}

// Points is a helper function to return a slice of Point values.
func (s Set) Points() []Point {
	return s.points
}

// Polarities is a helper function to return a slice of Polarity values.
func (s Set) Polarities() []Polarity {
	return s.polarities
}

// Preamps is a helper function to return a slice of Preamp values.
func (s Set) Preamps() []Preamp {
	return s.preamps
}

// Samples is a helper function to return a slice of Sample values.
func (s Set) Samples() []Sample {
	return s.samples
}

// Sessions is a helper function to return a slice of Session values.
func (s Set) Sessions() []Session {
	return s.sessions
}

// Sites is a helper function to return a slice of Site values.
func (s Set) Sites() []Site {
	return s.sites
}

// Stations is a helper function to return a slice of Station values.
func (s Set) Stations() []Station {
	return s.stations
}

// Streams is a helper function to return a slice of Stream values.
func (s Set) Streams() []Stream {
	return s.streams
}

// Telemetries is a helper function to return a slice of Telemetry values.
func (s Set) Telemetries() []Telemetry {
	return s.telemetries
}

// Timings is a helper function to return a slice of Timing values.
func (s Set) Timings() []Timing {
	return s.timings
}

// Views is a helper function to return a slice of View values.
func (s Set) Views() []View {
	return s.views
}

// Visibilities is a helper function to return a slice of Visibility values.
func (s Set) Visibilities() []Visibility {
	return s.visibilities
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

// Citation is a helper function to return a Citation value and true if one exists.
func (s Set) Citation(key string) (Citation, bool) {
	for _, v := range s.citations {
		if key != v.Key {
			continue
		}
		return v, true
	}
	return Citation{}, false
}

// Class is a helper function to return a Class value and true if one exists.
func (s Set) Class(station string) (Class, bool) {
	for _, v := range s.classes {
		if station != v.Station {
			continue
		}
		return v, true
	}
	return Class{}, false
}

// Constituent is a helper function to return a Constituent value and true if one exists.
func (s Set) Constituent(gauge, location string) (Constituent, bool) {
	for _, v := range s.constituents {
		if gauge != v.Gauge {
			continue
		}
		if location != v.Location {
			continue
		}
		return v, true
	}
	return Constituent{}, false
}

// Dart is a helper function to return a Dart value and true if one exists.
func (s Set) Dart(station string) (Dart, bool) {
	for _, v := range s.darts {
		if station != v.Station {
			continue
		}
		return v, true
	}
	return Dart{}, false
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

// Point is a helper function to return a Point value and true if one exists.
func (s Set) Point(sample, location string) (Point, bool) {
	for _, v := range s.points {
		if sample != v.Sample {
			continue
		}
		if location != v.Location {
			continue
		}
		return v, true
	}
	return Point{}, false
}

// Sample is a helper function to return a Sample value and true if one exists.
func (s Set) Sample(code string) (Sample, bool) {
	for _, v := range s.samples {
		if code != v.Code {
			continue
		}
		return v, true
	}
	return Sample{}, false
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

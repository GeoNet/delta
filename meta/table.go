package meta

// Table holds internal settings suitable for automatic code generation.
type Table struct {
	name    string
	headers Header
	primary []string
	native  []string
	foreign map[string][]string
	remap   map[string]string
	ignore  bool
	start   string
	end     string
}

// Name returns the table name.
func (t Table) Name() string {
	return t.name
}

// Start returns the table start and whether it has been set.
func (t Table) Start() (string, bool) {
	return t.start, t.start != ""
}

// End returns the table start and whether it has been set.
func (t Table) End() (string, bool) {
	return t.end, t.end != ""
}

// Columns returns the header columns.
func (t Table) Columns() []string {
	return t.headers.Columns()
}

// IsPrimary returns whether a column is considered a primary column,
// this is usually an indication that it will be uniquie together with
// any other primary columns.
func (t Table) IsPrimary(col int) bool {
	for _, s := range t.primary {
		if n, ok := t.headers[s]; !ok || n != col {
			continue
		}
		return true
	}
	return false
}

// IsForeign returns whether a column is considered a foreign column,
// this is usually an indication that it will be uniquie together with
// any other primary columns.
func (t Table) IsForeign(col int) (string, bool) {
	for k, v := range t.foreign {
		for _, s := range v {
			if n, ok := t.headers[s]; !ok || n != col {
				continue
			}
			return k, true
		}
	}
	return "", false
}

// IsNative returns whether the column should be displayed without
// quotes, such as numbers.
func (t Table) IsNative(col int) bool {
	for _, s := range t.native {
		if n, ok := t.headers[s]; !ok || n != col {
			continue
		}
		return true
	}
	return false
}

// IsDateTime returns whether the column is a date and may need extra formatting.
func (t Table) IsDateTime(col int) bool {
	for _, s := range []string{t.start, t.end} {
		if n, ok := t.headers[s]; !ok || n != col {
			continue
		}
		return true
	}
	return false
}

// HasDateTime returns whether there are any date columns in the table.
func (t Table) HasDateTime() (string, string, bool) {
	if t.ignore {
		return "", "", false
	}
	if t.start != "" && t.end != "" {
		return t.start, t.end, true
	}
	return "", "", false
}

// Remap provides a mechanism to swap between the columns expected in the
// struct elements and what is used in the CSV columns.
func (t Table) Remap(s string) string {
	if v, ok := t.remap[s]; ok {
		return v
	}
	return s
}

// Scan returns a list of remapped columns.
func (t Table) Scan() []string {
	var list []string
	for _, x := range t.headers.Columns() {
		list = append(list, t.Remap(x))
	}
	return list
}

// Encode returns the standard text listing of a table.
func (t Table) Encode(list ListEncoder) [][]string {
	return list.encode()
}

type Database interface {
	Drop(Table) []string
	Create(Table) []string
	Insert(Table, ListEncoder) []string
}

func Init(d Database, table Table, list ListEncoder) []string {
	var res []string
	for _, v := range d.Drop(table) {
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	for _, v := range d.Create(table) {
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	for _, v := range d.Insert(table, list) {
		if v == "" {
			continue
		}
		res = append(res, v)
	}

	return res
}

// Init builds the set of commands needed to initialise a given database.
func (s *Set) Init(d Database) []string {
	var cmds []string
	cmds = append(cmds, Init(d, NetworkTable, NetworkList(s.Networks()))...)
	cmds = append(cmds, Init(d, StationTable, StationList(s.Stations()))...)
	cmds = append(cmds, Init(d, SiteTable, SiteList(s.Sites()))...)
	cmds = append(cmds, Init(d, MarkTable, MarkList(s.Marks()))...)
	cmds = append(cmds, Init(d, MonumentTable, MonumentList(s.Monuments()))...)
	cmds = append(cmds, Init(d, MountTable, MountList(s.Mounts()))...)
	cmds = append(cmds, Init(d, ViewTable, ViewList(s.Views()))...)
	cmds = append(cmds, Init(d, SampleTable, SampleList(s.Samples()))...)
	cmds = append(cmds, Init(d, PointTable, PointList(s.Points()))...)
	cmds = append(cmds, Init(d, InstalledSensorTable, InstalledSensorList(s.InstalledSensors()))...)
	cmds = append(cmds, Init(d, InstalledAntennaTable, InstalledAntennaList(s.InstalledAntennas()))...)
	cmds = append(cmds, Init(d, AssetTable, AssetList(s.Assets()))...)
	cmds = append(cmds, Init(d, CalibrationTable, CalibrationList(s.Calibrations()))...)
	cmds = append(cmds, Init(d, InstalledCameraTable, InstalledCameraList(s.InstalledCameras()))...)
	cmds = append(cmds, Init(d, ChannelTable, ChannelList(s.Channels()))...)
	cmds = append(cmds, Init(d, CitationTable, CitationList(s.Citations()))...)
	cmds = append(cmds, Init(d, ClassTable, ClassList(s.Classes()))...)
	cmds = append(cmds, Init(d, ComponentTable, ComponentList(s.Components()))...)
	cmds = append(cmds, Init(d, ConnectionTable, ConnectionList(s.Connections()))...)
	cmds = append(cmds, Init(d, ConstituentTable, ConstituentList(s.Constituents()))...)
	cmds = append(cmds, Init(d, DartTable, DartList(s.Darts()))...)
	cmds = append(cmds, Init(d, DeployedDataloggerTable, DeployedDataloggerList(s.DeployedDataloggers()))...)
	cmds = append(cmds, Init(d, InstalledDoasTable, InstalledDoasList(s.Doases()))...)
	cmds = append(cmds, Init(d, FeatureTable, FeatureList(s.Features()))...)
	cmds = append(cmds, Init(d, FirmwareHistoryTable, FirmwareHistoryList(s.FirmwareHistory()))...)
	cmds = append(cmds, Init(d, GainTable, GainList(s.Gains()))...)
	cmds = append(cmds, Init(d, GaugeTable, GaugeList(s.Gauges()))...)
	cmds = append(cmds, Init(d, InstalledMetSensorTable, InstalledMetSensorList(s.InstalledMetSensors()))...)
	cmds = append(cmds, Init(d, PlacenameTable, PlacenameList(s.Placenames()))...)
	cmds = append(cmds, Init(d, PolarityTable, PolarityList(s.Polarities()))...)
	cmds = append(cmds, Init(d, PreampTable, PreampList(s.Preamps()))...)
	cmds = append(cmds, Init(d, InstalledRadomeTable, InstalledRadomeList(s.InstalledRadomes()))...)
	cmds = append(cmds, Init(d, DeployedReceiverTable, DeployedReceiverList(s.DeployedReceivers()))...)
	cmds = append(cmds, Init(d, InstalledRecorderTable, InstalledRecorderList(s.InstalledRecorders()))...)
	cmds = append(cmds, Init(d, SessionTable, SessionList(s.Sessions()))...)
	cmds = append(cmds, Init(d, StreamTable, StreamList(s.Streams()))...)
	cmds = append(cmds, Init(d, TelemetryTable, TelemetryList(s.Telemetries()))...)
	cmds = append(cmds, Init(d, TimingTable, TimingList(s.Timings()))...)
	cmds = append(cmds, Init(d, VisibilityTable, VisibilityList(s.Visibilities()))...)

	return cmds
}

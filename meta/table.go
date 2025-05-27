package meta

import (
	"slices"
	"sort"
)

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

// TableList pairs together a Table and a ListEncoder.
type TableList struct {
	Table Table
	List  ListEncoder
}

// KeyValue contains key value pairs and their labels.
type KeyValue struct {
	label   string
	desc    string
	entries map[string]string
}

// encode implements the ListEncoder interface.
func (kv KeyValue) encode() [][]string {
	var keys []string
	for k := range kv.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	list := [][]string{{kv.label, kv.desc}}
	for _, k := range keys {
		v, ok := kv.entries[k]
		if !ok {
			continue
		}
		list = append(list, []string{k, v})
	}
	return list
}

// KeyValue returns a bespoke TableList for a key value style dataset.
func (s *Set) KeyValue(name, label, desc string, entries map[string]string) TableList {
	return TableList{
		Table: Table{
			name: name,
			headers: Header{
				label: 0,
				desc:  1,
			},
			primary: []string{label},
		},
		List: KeyValue{
			label:   label,
			desc:    desc,
			entries: entries,
		},
	}
}

// TableList returns a pre-built set of Tables and associated Lists.
func (s *Set) TableList(extra ...TableList) []TableList {
	return append(slices.Clone(extra), []TableList{
		{Table: AssetTable, List: AssetList(s.Assets())},
		{Table: NetworkTable, List: NetworkList(s.Networks())},
		{Table: StationTable, List: StationList(s.Stations())},
		{Table: SiteTable, List: SiteList(s.Sites())},
		{Table: MarkTable, List: MarkList(s.Marks())},
		{Table: MonumentTable, List: MonumentList(s.Monuments())},
		{Table: MountTable, List: MountList(s.Mounts())},
		{Table: ViewTable, List: ViewList(s.Views())},
		{Table: SampleTable, List: SampleList(s.Samples())},
		{Table: PointTable, List: PointList(s.Points())},
		{Table: InstalledSensorTable, List: InstalledSensorList(s.InstalledSensors())},
		{Table: InstalledAntennaTable, List: InstalledAntennaList(s.InstalledAntennas())},
		{Table: CalibrationTable, List: CalibrationList(s.Calibrations())},
		{Table: InstalledCameraTable, List: InstalledCameraList(s.InstalledCameras())},
		{Table: ChannelTable, List: ChannelList(s.Channels())},
		{Table: CitationTable, List: CitationList(s.Citations())},
		{Table: ClassTable, List: ClassList(s.Classes())},
		{Table: ComponentTable, List: ComponentList(s.Components())},
		{Table: ConnectionTable, List: ConnectionList(s.Connections())},
		{Table: DartTable, List: DartList(s.Darts())},
		{Table: DeployedDataloggerTable, List: DeployedDataloggerList(s.DeployedDataloggers())},
		{Table: InstalledDoasTable, List: InstalledDoasList(s.Doases())},
		{Table: FeatureTable, List: FeatureList(s.Features())},
		{Table: FirmwareHistoryTable, List: FirmwareHistoryList(s.FirmwareHistory())},
		{Table: GainTable, List: GainList(s.Gains())},
		{Table: GaugeTable, List: GaugeList(s.Gauges())},
		{Table: ConstituentTable, List: ConstituentList(s.Constituents())},
		{Table: InstalledMetSensorTable, List: InstalledMetSensorList(s.InstalledMetSensors())},
		{Table: PlacenameTable, List: PlacenameList(s.Placenames())},
		{Table: PolarityTable, List: PolarityList(s.Polarities())},
		{Table: PreampTable, List: PreampList(s.Preamps())},
		{Table: InstalledRadomeTable, List: InstalledRadomeList(s.InstalledRadomes())},
		{Table: DeployedReceiverTable, List: DeployedReceiverList(s.DeployedReceivers())},
		{Table: InstalledRecorderTable, List: InstalledRecorderList(s.InstalledRecorders())},
		{Table: SessionTable, List: SessionList(s.Sessions())},
		{Table: StreamTable, List: StreamList(s.Streams())},
		{Table: TelemetryTable, List: TelemetryList(s.Telemetries())},
		{Table: TimingTable, List: TimingList(s.Timings())},
		{Table: VisibilityTable, List: VisibilityList(s.Visibilities())},
	}...)
}

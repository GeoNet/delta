package metadb

type MetaDB struct {
	// network details
	networks
	marks
	stations
	sites
	gauges
	constituents

	// instrument details
	sensors
	recorders
	dataloggers

	// instrment configuration
	connections
	streams

	// base directory for raw meta files
	base string
}

func NewMetaDB(base string) *MetaDB {
	return &MetaDB{
		base: base,
	}
}

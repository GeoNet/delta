package metadb

type MetaDB struct {
	// network details
	networks
	stations
	sites
	gauges
	constituents

	// instrument details
	sensors
	recorders
	dataloggers
	calibrations

	// instrument configuration
	connections
	streams
	gains

	// base directory for raw meta files
	base string
}

func NewMetaDB(base string) *MetaDB {
	return &MetaDB{
		base: base,
	}
}

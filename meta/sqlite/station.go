package sqlite

import (
	"fmt"
)

const methodCreate = `
DROP TABLE IF EXISTS method;
CREATE TABLE IF NOT EXISTS method (
  method_id INTEGER PRIMARY KEY NOT NULL,
  method TEXT DEFAULT "Unknown" NOT NULL,
  UNIQUE (method)
);`

var method = Table{
	Create: methodCreate,
	Select: func() string {
		return "SELECT method_id FROM method WHERE method = ?"
	},
	Insert: func() string {
		return "INSERT INTO method (method) VALUES (?) ON CONFLICT(method) DO NOTHING;"
	},
	Fields: []string{"Method"},
}

const stationCreate = `
DROP TABLE IF EXISTS station;
CREATE TABLE IF NOT EXISTS station (
  station_id INTEGER PRIMARY KEY NOT NULL,
  reference_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (reference_id) REFERENCES reference (reference_id),
  UNIQUE (reference_id)
);`

var station = Table{
	Create: stationCreate,
	/*
		Select: func() string {
			return "SELECT station_id FROM station WHERE station = ?"
		},
		Insert: func() string {
			return fmt.Sprintf("INSERT INTO station (datum_id, station, name, latitude, longitude, elevation, depth, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?);", datum.Select())
		},
		Fields: []string{"Datum", "Station", "Name", "Latitude", "Longitude", "Elevation", "Depth", "Start Date", "End Date"},
		Nulls:  []string{"Elevation", "Depth"},
	*/
	Select: func() string {
		return fmt.Sprintf("SELECT station_id FROM station WHERE reference_id = (%s)", reference.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO station (reference_id, start_date, end_date) VALUES ((%s), ?, ?);", reference.Select())
	},
	Fields: []string{"Station", "Start Date", "End Date"},
}

const stationNetworkCreate = `
DROP TABLE IF EXISTS station_network;
CREATE TABLE IF NOT EXISTS station_network (
  station_network_id INTEGER PRIMARY KEY NOT NULL,
  station_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (station_id) REFERENCES station (station_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (station_id, network_id)
);`

var stationNetwork = Table{
	Create: stationNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT station_network_id FROM station_network WHERE station_id = (%s) AND network_id = (%s)",
			station.Select(), network.Select())
	},
	Insert: func() string {
		// not all networks are in the networks table so simply ignore any that fail
		return fmt.Sprintf("INSERT OR IGNORE INTO station_network (station_id, network_id) VALUES ((%s), (%s));",
			station.Select(), network.Select())
	},
	Fields: []string{"Station", "Network"},
}

const siteCreate = `
DROP TABLE IF EXISTS site;
CREATE TABLE IF NOT EXISTS site (
  site_id INTEGER PRIMARY KEY NOT NULL,
  station_id INTEGER NOT NULL,
  datum_id INTEGER NOT NULL,
  location TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  elevation REAL NULL,
  depth REAL NULL,
  survey TEXT DEFAULT "Unknown" NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (station_id) REFERENCES station (station_id),
  FOREIGN KEY (datum_id) REFERENCES datum (datum_id),
  UNIQUE (station_id, location)
);
CREATE TRIGGER IF NOT EXISTS site_too_soon BEFORE INSERT ON site
WHEN NEW.start_date < (SELECT station.start_date FROM station WHERE station.station_id = new.station_id)
BEGIN
  SELECT RAISE(FAIL, "site too soon for station");
END;
CREATE TRIGGER IF NOT EXISTS site_too_late BEFORE INSERT ON site
WHEN NEW.end_date > (SELECT station.end_date FROM station WHERE station.station_id = new.station_id)
BEGIN
  SELECT RAISE(FAIL, "site too late for station");
END;
`

var site = Table{
	Create: siteCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT site_id FROM site WHERE station_id = (%s) AND location = ?", station.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO site (station_id, datum_id, location, latitude, longitude, elevation, depth, survey, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			station.Select(), datum.Select())
	},

	Fields: []string{"Station", "Datum", "Location", "Latitude", "Longitude", "Elevation", "Depth", "Survey", "Start Date", "End Date"},
	Nulls:  []string{"Elevation", "Depth"},
}

const sampleNetworkCreate = `
DROP TABLE IF EXISTS sample_network;
CREATE TABLE IF NOT EXISTS sample_network (
  sample_network_id INTEGER PRIMARY KEY NOT NULL,
  sample_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (sample_id) REFERENCES sample (sample_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (sample_id, network_id)
);`

var sampleNetwork = Table{
	Create: sampleNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT sample_network_id FROM sample_network WHERE sample_id = (%s) AND network_id = (%s)",
			sample.Select(), network.Select())
	},
	Insert: func() string {
		// not all networks are in the networks table so simply ignore any that fail
		return fmt.Sprintf("INSERT OR IGNORE INTO sample_network (sample_id, network_id) VALUES ((%s), (%s));",
			sample.Select(), network.Select())
	},
	Fields: []string{"Station", "Network"},
}

const sampleCreate = `
DROP TABLE IF EXISTS sample;
CREATE TABLE IF NOT EXISTS sample (
  sample_id INTEGER PRIMARY KEY NOT NULL,
  reference_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (reference_id) REFERENCES reference (reference_id),
  UNIQUE (reference_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_sample BEFORE INSERT ON sample
WHEN EXISTS (
  SELECT * FROM sample
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND reference_id =  NEW.reference_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on sample");
END;
`

var sample = Table{
	Create: sampleCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT sample_id FROM sample WHERE reference_id = (%s)", reference.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO sample (reference_id, start_date, end_date) VALUES ((%s), ?, ?);", reference.Select())
	},
	Fields: []string{"Station", "Start Date", "End Date"},
}

const pointCreate = `
DROP TABLE IF EXISTS point;
CREATE TABLE IF NOT EXISTS point (
  point_id INTEGER PRIMARY KEY NOT NULL,
  sample_id INTEGER NOT NULL,
  datum_id INTEGER NOT NULL,
  location TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  elevation REAL DEFAULT 0 NOT NULL,
  depth REAL DEFAULT 0 NOT NULL,
  survey TEXT DEFAULT "Unknown" NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (sample_id) REFERENCES sample (sample_id),
  FOREIGN KEY (datum_id) REFERENCES datum (datum_id),
  UNIQUE (sample_id, location)
);
CREATE TRIGGER IF NOT EXISTS point_too_soon BEFORE INSERT ON point
WHEN NEW.start_date < (SELECT sample.start_date FROM sample WHERE sample.sample_id = new.sample_id)
BEGIN
  SELECT RAISE(FAIL, "point too soon for sample");
END;
CREATE TRIGGER IF NOT EXISTS site_too_late BEFORE INSERT ON point
WHEN NEW.end_date > (SELECT sample.end_date FROM sample WHERE sample.sample_id = new.sample_id)
BEGIN
  SELECT RAISE(FAIL, "point too late for sample");
END;
`

var point = Table{
	Create: pointCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO point (sample_id, datum_id, location, latitude, longitude, elevation, depth, survey, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			sample.Select(), datum.Select())
	},
	Fields: []string{"Sample", "Datum", "Location", "Latitude", "Longitude", "Elevation", "Depth", "Survey", "Start Date", "End Date"},
}

const featureCreate = `
DROP TABLE IF EXISTS feature;
CREATE TABLE IF NOT EXISTS feature (
  feature_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  sublocation TEXT NULL,
  property TEXT NOT NULL,
  description TEXT NULL,
  aspect TEXT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(site_id, sublocation, property, description, aspect, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_feature BEFORE INSERT ON feature
WHEN EXISTS (
  SELECT * FROM feature
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
      AND sublocation =  NEW.sublocation
      AND property =  NEW.property
      AND description =  NEW.description
      AND aspect =  NEW.aspect
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on feature");
END;
`

var feature = Table{
	Create: featureCreate,
	Insert: func() string {
		// currently a feature could reference a site or a point or a sample, solution is consolidation.
		return fmt.Sprintf("INSERT OR IGNORE INTO feature (site_id, sublocation, property, description, aspect, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?);",
			site.Select())
	},
	Fields: []string{"Station", "Location", "Sublocation", "Property", "Description", "Aspect", "Start Date", "End Date"},
}

const classCreate = `
DROP TABLE IF EXISTS class;
CREATE TABLE IF NOT EXISTS class (
  class_id INTEGER PRIMARY KEY NOT NULL,
  station_id INTEGER NOT NULL,
  site_class TEXT NOT NULL,
  vs30 REAL NOT NULL,
  vs30_quality TEXT NOT NULL,
  tsite TEXT NOT NULL,
  tsite_method TEXT NOT NULL,
  tsite_quality TEXT NOT NULL,
  basement_depth REAL NOT NULL,
  depth_quality TEXT NOT NULL,
  link TEXT NULL,
  notes TEXT NULL,
  FOREIGN KEY (station_id) REFERENCES station (station_id),
  UNIQUE(station_id)
);
`

var class = Table{
	Create: classCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT class_id FROM class WHERE station_id = (%s)", station.Select())
	},
	Insert: func() string {
		// not all stations are in the stations file, ignore any conflicts for now
		return fmt.Sprintf("INSERT OR IGNORE INTO class (station_id, site_class, vs30, vs30_quality, tsite, tsite_method, tsite_quality, basement_depth, depth_quality, link, notes) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);", station.Select())
	},
	Fields: []string{"Station", "Site Class", "Vs30", "Vs30 Quality", "Tsite", "Tsite Method", "Tsite Quality", "Basement Depth", "Depth Quality", "Link", "Notes"},
	Nulls:  []string{"Link", "Notes"},
	Unwrap: "Citations",
}

const classCitationCreate = `
DROP TABLE IF EXISTS class_citation;
CREATE TABLE IF NOT EXISTS class_citation (
  class_citation_id INTEGER PRIMARY KEY NOT NULL,
  class_id INTEGER NOT NULL,
  citation_id INTEGER NOT NULL,
  FOREIGN KEY (class_id) REFERENCES class (class_id),
  FOREIGN KEY (citation_id) REFERENCES citation (citation_id),
  UNIQUE (class_id, citation_id)
);`

var classCitation = Table{
	Create: classCitationCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT class_citation_id FROM class_citation WHERE class_id = (%s) AND citation_id = (%s)",
			class.Select(), citation.Select(),
		)
	},
	Insert: func() string {
		// not all stations are in the stations file, ignore any conflicts for now
		return fmt.Sprintf("INSERT OR IGNORE INTO class_citation (class_id, citation_id) VALUES ((%s), (%s));",
			class.Select(), citation.Select(),
		)
	},
	Fields: []string{"Station", "Citations"},
}

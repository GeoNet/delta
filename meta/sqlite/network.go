package sqlite

import (
	"fmt"
)

const networkCreate = `
DROP TABLE IF EXISTS network;
CREATE TABLE IF NOT EXISTS network (
  network_id INTEGER PRIMARY KEY NOT NULL,
  network TEXT NOT NULL,
  external TEXT NOT NULL,
  description TEXT DEFAULT "" NOT NULL,
  restricted BOOLEAN DEFAULT false NOT NULL,
  UNIQUE (network)
);`

var network = Table{
	Create: networkCreate,
	Select: func() string {
		return "SELECT network_id FROM network WHERE network = ?"
	},
	Insert: func() string {
		return "INSERT INTO network (network, external, description, restricted) VALUES (?, ?, ?, ?);"
	},
	Fields: []string{"Network", "External", "Description", "Restricted"},
}

const datumCreate = `
DROP TABLE IF EXISTS datum;
CREATE TABLE IF NOT EXISTS datum (
  datum_id INTEGER PRIMARY KEY NOT NULL,
  datum TEXT NOT NULL,
  UNIQUE (datum)
);`

var datum = Table{
	Create: datumCreate,
	Select: func() string {
		return "SELECT datum_id FROM datum WHERE datum = ?"
	},
	Insert: func() string {
		return "INSERT INTO datum (datum) VALUES (?) ON CONFLICT(datum) DO NOTHING;"
	},
	Fields: []string{"Datum"},
}

const locationCreate = `
DROP TABLE IF EXISTS location;
CREATE TABLE IF NOT EXISTS location (
  location_id INTEGER PRIMARY KEY NOT NULL,
  datum_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  elevation REAL NULL,
  depth REAL NULL,
  FOREIGN KEY (datum_id) REFERENCES datum (datum_id),
  UNIQUE (code)
);`

var location = Table{
	Create: locationCreate,
	Select: func() string {
		return "SELECT location_id FROM location WHERE code = ?"
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO location (datum_id, code, name, latitude, longitude, elevation, depth) VALUES ((%s), ?, ?, ?, ?, ?, ?) ON CONFLICT(code) DO NOTHING;", datum.Select())
	},
	Fields: []string{"Datum", "Code", "Name", "Latitude", "Longitude", "Elevation", "Depth"},
	Nulls:  []string{"Elevation", "Depth"},
	Remap: map[string][]string{
		"Code": {"Station", "Mark", "Mount", "Gauge"},
	},
}

const locationNetworkCreate = `
DROP TABLE IF EXISTS location_network;
CREATE TABLE IF NOT EXISTS location_network (
  location_network_id INTEGER PRIMARY KEY NOT NULL,
  location_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (location_id) REFERENCES location (location_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (location_id, network_id)
);`

var locationNetwork = Table{
	Create: locationNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT location_network_id FROM location_network WHERE location_id = (%s) AND network_id = (%s)",
			location.Select(), network.Select())
	},
	Insert: func() string {
		// not all networks are in the networks table so simply ignore any that fail
		return fmt.Sprintf("INSERT OR IGNORE INTO location_network (location_id, network_id) VALUES ((%s), (%s));",
			location.Select(), network.Select())
	},
	Fields: []string{"Code", "Network"},
	Remap: map[string][]string{
		"Code": {"Station", "Mark", "Mount"},
	},
}

const placenameCreate = `
DROP TABLE IF EXISTS placename;
CREATE TABLE IF NOT EXISTS placename (
/*
** A set of well known place names that can be used to
** describe the relative position of locations.
**
** The level is used as an indication of place scale so
** the at greater distances higher level places will be
** prefered over smaller ones. An example would be a
** suburb of a city is use when the location is close
** to the city, but the city itself should be used when
** the distance increases.
*/
  placename_id INTEGER PRIMARY KEY NOT NULL,
  name TEXT NOT NULL, -- The name of the place being described.
  latitude REAL NOT NULL, -- Geographical Latitude in degrees.
  longitude REAL NOT NULL, -- Geographical Longitude in degrees.
  level INTEGER NOT NULL, -- The place level, or importance.
  UNIQUE(name)
);
`

var placename = Table{
	Create: placenameCreate,
	Insert: func() string {
		return "INSERT INTO placename (name, latitude, longitude, level) VALUES (?, ?, ?, ?);"
	},
	Fields: []string{"Name", "Latitude", "Longitude", "Level"},
}

const citationCreate = `
DROP TABLE IF EXISTS citation;
CREATE TABLE IF NOT EXISTS citation (
  citation_id INTEGER PRIMARY KEY NOT NULL,
  key TEXT NOT NULL,
  author TEXT NOT NULL,
  year REAL NOT NULL,
  title TEXT NOT NULL,
  published TEXT NULL,
  volume TEXT NULL,
  pages TEXT NULL,
  doi TEXT NULL,
  link TEXT NULL,
  retrieved TEXT NULL,
  UNIQUE(key)
);
`

var citation = Table{
	Create: citationCreate,
	Select: func() string {
		return "SELECT citation_id FROM citation WHERE key = ?"
	},
	Insert: func() string {
		return "INSERT INTO citation (key, author, year, title, published, volume, pages, doi, link, retrieved) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	},
	Fields: []string{"Key", "Author", "Year", "Title", "Published", "Volume", "Pages", "DOI", "Link", "Retrieved"},
	Nulls:  []string{"Published", "Volume", "Pages", "DOI", "Link", "Retrieved"},
}

const noteCreate = `
DROP TABLE IF EXISTS note;
CREATE TABLE IF NOT EXISTS note (
  note_id INTEGER PRIMARY KEY NOT NULL,
  location_network_id INTEGER NOT NULL,
  entry TEXT NULL,
  FOREIGN KEY (location_network_id) REFERENCES location_network (location_network_id)
);
`

var note = Table{
	Create: noteCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO note (location_network_id, entry) VALUES ((%s), ?);",
			locationNetwork.Select())
	},
	Fields: []string{"Code", "Network", "Entry"},
}

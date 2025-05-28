package sqlite

import (
	"log"

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

const referenceCreate = `
DROP TABLE IF EXISTS reference;
CREATE TABLE IF NOT EXISTS reference (
  reference_id INTEGER PRIMARY KEY NOT NULL,
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

var reference = Table{
	Create: referenceCreate,
	Select: func() string {
		return "SELECT reference_id FROM reference WHERE code = ?"
	},
	Insert: func() string {
		//log.Println("INSERT INTO reference (datum_id, code, name, latitude, longitude, elevation, depth)")
		x := fmt.Sprintf("INSERT INTO reference (datum_id, code, name, latitude, longitude, elevation, depth) VALUES ((%s), ?, ?, ?, ?, ?, ?) ON CONFLICT(code) DO NOTHING;", datum.Select())
		//x := fmt.Sprintf("INSERT INTO reference (datum_id, code, name, latitude, longitude, elevation, depth) VALUES ((%s), ?, ?, ?, ?, ?, ?);", datum.Select())
		log.Println(x)
		return x
	},
	Fields: []string{"Datum", "Code", "Name", "Latitude", "Longitude", "Elevation", "Depth"},
	Nulls:  []string{"Elevation", "Depth"},
	Remap: map[string][]string{
		"Code": {"Station", "Mark"},
	},
}

/*
const referenceNetworkCreate = `
DROP TABLE IF EXISTS reference_network;
CREATE TABLE IF NOT EXISTS reference_network (
  reference_network_id INTEGER PRIMARY KEY NOT NULL,
  reference_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (reference_id) REFERENCES reference (reference_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (reference_id, network_id)
);`

var referenceNetwork = Table{
	Create: referenceNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT reference_network_id FROM reference_network WHERE reference_id = (%s) AND network_id = (%s)",
			reference.Select(), network.Select())
	},
	Insert: func() string {
		// not all networks are in the networks table so simply ignore any that fail
		return fmt.Sprintf("INSERT OR IGNORE INTO reference_network (reference_id, network_id) VALUES ((%s), (%s));",
			reference.Select(), network.Select())
	},
	Fields: []string{"Reference", "Network"},
}
*/

const placenameCreate = `
DROP TABLE IF EXISTS placename;
CREATE TABLE IF NOT EXISTS placename (
  placename_id INTEGER PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  level INTEGER NOT NULL,
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

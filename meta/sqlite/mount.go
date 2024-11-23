package sqlite

import (
	"fmt"
)

const mountCreate = `
DROP TABLE IF EXISTS mount;
CREATE TABLE IF NOT EXISTS mount (
  mount_id INTEGER PRIMARY KEY NOT NULL,
  datum_id INTEGER NOT NULL,
  mount TEXT NOT NULL,
  name TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  elevation REAL NOT NULL,
  description TEXT DEFAULT "" NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (datum_id) REFERENCES datum (datum_id),
  UNIQUE(mount, start_date, end_date)
);
`

var mount = Table{
	Create: mountCreate,
	Select: func() string {
		return "SELECT mount_id FROM mount WHERE mount = ?"
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO mount (datum_id, mount, name, latitude, longitude, elevation, description, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			datum.Select(),
		)
	},
	Fields: []string{"Datum", "Mount", "Name", "Latitude", "Longitude", "Elevation", "Description", "Start Date", "End Date"},
}

const mountNetworkCreate = `
DROP TABLE IF EXISTS mount_network;
CREATE TABLE IF NOT EXISTS mount_network (
  mount_network_id INTEGER PRIMARY KEY NOT NULL,
  mount_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (mount_id) REFERENCES mount (mount_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (mount_id, network_id)
);`

var mountNetwork = Table{
	Create: mountNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT mount_network_id FROM mount_network WHERE mount_id = (%s) AND network_id = (%s)",
			mount.Select(), network.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO mount_network (mount_id, network_id) VALUES ((%s), (%s));",
			mount.Select(), network.Select())
	},

	Fields: []string{"Mount", "Network"},
}

const viewCreate = `
DROP TABLE IF EXISTS view;
CREATE TABLE IF NOT EXISTS view (
  view_id INTEGER PRIMARY KEY NOT NULL,
  mount_id INTEGER NOT NULL,
  view TEXT NOT NULL,
  label TEXT NULL,
  azimuth REAL NOT NULL,
  dip REAL NOT NULL,
  method TEXT NULL,
  description TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (mount_id) REFERENCES mount (mount_id),
  UNIQUE(mount_id, start_date, view)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_view BEFORE INSERT ON view
WHEN EXISTS (
  SELECT * FROM view
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND mount_id = NEW.mount_id
      AND view =  NEW.view
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on view");
END;
CREATE TRIGGER IF NOT EXISTS view_too_soon BEFORE INSERT ON view
WHEN NEW.start_date < (SELECT mount.start_date FROM mount WHERE mount.mount_id = new.mount_id)
BEGIN
  SELECT RAISE(FAIL, "view too soon for mount");
END;
CREATE TRIGGER IF NOT EXISTS view_too_late BEFORE INSERT ON view
WHEN NEW.end_date > (SELECT mount.end_date FROM mount WHERE mount.mount_id = new.mount_id)
BEGIN
  SELECT RAISE(FAIL, "view too late for mount");
END;
`

var view = Table{
	Create: viewCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT view_id FROM view WHERE mount_id = (%s) AND view = ? AND start_date <= ? AND end_date >= ?",
			mount.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO view (mount_id, view, label, azimuth, dip, method, description, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			mount.Select(),
		)
	},
	Fields: []string{"Mount", "View", "Label", "Azimuth", "Dip", "Method", "Description", "Start Date", "End Date"},
}

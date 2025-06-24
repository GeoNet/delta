package sqlite

import (
	"fmt"
)

const mountCreate = `
DROP TABLE IF EXISTS mount;
CREATE TABLE IF NOT EXISTS mount (
  mount_id INTEGER PRIMARY KEY NOT NULL,
  location_id INTEGER NOT NULL,
  description TEXT DEFAULT "" NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (location_id) REFERENCES location (location_id),
  UNIQUE(location_id, start_date, end_date)
);
`

var mount = Table{
	Create: mountCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT mount_id FROM mount WHERE location_id = (%s)", location.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO mount (location_id, description, start_date, end_date) VALUES ((%s), ?, ?, ?);",
			location.Select(),
		)
	},
	Fields: []string{"Mount", "Description", "Start Date", "End Date"},
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

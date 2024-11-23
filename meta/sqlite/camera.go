package sqlite

import (
	"fmt"
)

const cameraCreate = `
DROP TABLE IF EXISTS camera;
CREATE TABLE IF NOT EXISTS camera (
  camera_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  view_id INTEGER NOT NULL,
  dip REAL DEFAULT 0.0 NOT NULL,
  azimuth REAL DEFAULT 0.0 NOT NULL,
  height REAL DEFAULT 0.0 NOT NULL,
  north REAL DEFAULT 0.0 NOT NULL,
  east REAL DEFAULT 0.0 NOT NULL,
  notes TEXT DEFAULT "" NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (view_id) REFERENCES view (view_id),
  UNIQUE(asset_id, view_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_camera BEFORE INSERT ON camera
WHEN EXISTS (
  SELECT * FROM camera
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND view_id = NEW.view_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on camera");
END;
CREATE TRIGGER IF NOT EXISTS camera_too_soon BEFORE INSERT ON camera
WHEN NEW.start_date < (SELECT view.start_date FROM view WHERE view.view_id = new.view_id)
BEGIN
  SELECT RAISE(FAIL, "camera too soon for view");
END;
CREATE TRIGGER IF NOT EXISTS camera_too_late BEFORE INSERT ON camera
WHEN NEW.end_date > (SELECT view.end_date FROM view WHERE view.view_id = new.view_id)
BEGIN
  SELECT RAISE(FAIL, "camera too late for view");
END;
`

var camera = Table{
	Create: cameraCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT camera_id FROM camera WHERE asset_id = (%s) AND view_id = (%s)",
			asset.Select(), view.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO camera (asset_id, view_id, dip, azimuth, height, north, east, notes, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			asset.Select(), view.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mount", "View", "Start Date", "End Date", "Dip", "Azimuth", "Height", "North", "East", "Notes", "Start Date", "End Date"},
}

var doasCreate = `
DROP TABLE IF EXISTS doas;
CREATE TABLE IF NOT EXISTS doas (
  doas_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  view_id INTEGER NOT NULL,
  dip REAL NOT NULL,
  azimuth REAL NOT NULL,
  height REAL NOT NULL,
  north REAL NOT NULL,
  east REAL NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (view_id) REFERENCES view (view_id),
  UNIQUE(asset_id, view_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_doas BEFORE INSERT ON doas
WHEN EXISTS (
  SELECT * FROM doas
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND view_id = NEW.view_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on doas");
END;
CREATE TRIGGER IF NOT EXISTS doas_too_soon BEFORE INSERT ON doas
WHEN NEW.start_date < (SELECT view.start_date FROM view WHERE view.view_id = new.view_id)
BEGIN
  SELECT RAISE(FAIL, "doas too soon for view");
END;
CREATE TRIGGER IF NOT EXISTS doas_too_late BEFORE INSERT ON doas
WHEN NEW.end_date > (SELECT view.end_date FROM view WHERE view.view_id = new.view_id)
BEGIN
  SELECT RAISE(FAIL, "doas too late for view");
END;
`

var doas = Table{
	Create: doasCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT doas_id FROM doas WHERE asset_id = (%s) AND view_id = (%s)",
			asset.Select(), view.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO doas (asset_id, view_id, dip, azimuth, height, north, east, notes, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			asset.Select(), view.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mount", "View", "Dip", "Azimuth", "Height", "North", "East", "Notes", "Start Date", "End Date"},
}

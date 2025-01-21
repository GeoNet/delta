package sqlite

import (
	"fmt"
)

const antennaCreate = `
DROP TABLE IF EXISTS antenna;
CREATE TABLE IF NOT EXISTS antenna (
  antenna_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  mark_id INTEGER NOT NULL,
  height REAL NOT NULL,
  north REAL NOT NULL,
  east REAL NOT NULL,
  azimuth REAL NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  UNIQUE(asset_id, mark_id, start_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_antenna BEFORE INSERT ON antenna
WHEN EXISTS (
  SELECT * FROM antenna
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND mark_id = NEW.mark_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on antenna");
END;
`

var antenna = Table{
	Create: antennaCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO antenna (asset_id, mark_id, height, north, east, azimuth, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?);",
			asset.Select(), mark.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mark", "Height", "North", "East", "Azimuth", "Start Date", "End Date"},
}

const metsensorCreate = `
DROP TABLE IF EXISTS metsensor;
CREATE TABLE IF NOT EXISTS metsensor (
  metsensor_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  mark_id INTEGER NOT NULL,
  datum_id INTEGER NOT NULL,
  ims_comment TEXT NULL,
  humidity REAL NOT NULL,
  pressure REAL NOT NULL,
  temperature REAL NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  elevation REAL NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  FOREIGN KEY (datum_id) REFERENCES datum (datum_id),
  UNIQUE(asset_id, mark_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_metsensor BEFORE INSERT ON metsensor
WHEN EXISTS (
  SELECT * FROM metsensor
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND mark_id = NEW.mark_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on metsensor");
END;
`

var metsensor = Table{
	Create: metsensorCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO metsensor (asset_id, mark_id, datum_id, ims_comment, humidity, pressure, temperature, latitude, longitude, elevation, start_date, end_date) VALUES ((%s), (%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			asset.Select(), mark.Select(), datum.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mark", "Datum", "IMS Comment", "Humidity", "Pressure", "Temperature", "Latitude", "Longitude", "Elevation", "Start Date", "End Date"},
}

var radomeCreate = `
DROP TABLE IF EXISTS radome;
CREATE TABLE IF NOT EXISTS radome (
  radome_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  mark_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  UNIQUE(asset_id, mark_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_radome BEFORE INSERT ON radome
WHEN EXISTS (
  SELECT * FROM radome
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND mark_id = NEW.mark_id
      AND asset_id = NEW.asset_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on radome");
END;
`

var radome = Table{
	Create: radomeCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO radome (asset_id, mark_id, start_date, end_date) VALUES ((%s), (%s), ?, ?);",
			asset.Select(), mark.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mark", "Start Date", "End Date"},
}

var receiverCreate = `
DROP TABLE IF EXISTS receiver;
CREATE TABLE IF NOT EXISTS receiver (
  receiver_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  mark_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  UNIQUE(asset_id, mark_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_receiver BEFORE INSERT ON receiver
WHEN EXISTS (
  SELECT * FROM receiver
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND mark_id = NEW.mark_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on receiver");
END;
`

var receiver = Table{
	Create: receiverCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO receiver (asset_id, mark_id, start_date, end_date) VALUES ((%s), (%s), ?, ?);",
			asset.Select(), mark.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Mark", "Start Date", "End Date"},
}

const sessionCreate = `
DROP TABLE IF EXISTS session;
CREATE TABLE IF NOT EXISTS session (
  session_id INTEGER PRIMARY KEY NOT NULL,
  mark_id INTEGER NOT NULL,
  operator TEXT NOT NULL,
  agency TEXT NOT NULL,
  model TEXT NOT NULL,
  satellite_system TEXT NOT NULL,
  interval TEXT NOT NULL,
  elevation_mask REAL NOT NULL,
  header_comment TEXT NOT NULL,
  format TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  UNIQUE(interval, mark_id, start_date)
);

CREATE TRIGGER IF NOT EXISTS no_overlap_on_session BEFORE INSERT ON session
WHEN EXISTS (
  SELECT * FROM session
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND mark_id = NEW.mark_id
      AND interval =  NEW.interval
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on session");
END;
`

var session = Table{
	Create: sessionCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO session (mark_id, operator, agency, model, satellite, interval, elevation_mask, header_comment, format, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			mark.Select(),
		)
	},
	Fields: []string{"Make", "Operator", "Agency", "Model", "Satellite", "Interval", "Elevation Mask", "Header Comment", "Format", "Start Date", "End Date"},
}

package sqlite

import (
	"fmt"
)

const gaugeCreate = `
DROP TABLE IF EXISTS gauge;
CREATE TABLE IF NOT EXISTS gauge (
  gauge_id INTEGER PRIMARY KEY NOT NULL,
  location_id INTEGER NOT NULL,
  identification_number TEXT NOT NULL,
  analysis_time_zone REAL NOT NULL,
  analysis_latitude REAL NOT NULL,
  analysis_longitude REAL NOT NULL,
  crex_tag TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (location_id) REFERENCES location (location_id),
  UNIQUE(location_id, start_date, end_date)
);
`

var gauge = Table{
	Create: gaugeCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT gauge_id FROM gauge WHERE location_id = (%s)", location.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO gauge (location_id, identification_number, analysis_time_zone, analysis_latitude, analysis_longitude, crex_tag, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?);", location.Select())
	},
	Fields: []string{"Gauge", "Identification Number", "Analysis Time Zone", "Analysis Latitude", "Analysis Longitude", "Crex Tag", "Start Date", "End Date"},
}

const constituentCreate = `
DROP TABLE IF EXISTS constituent;
CREATE TABLE IF NOT EXISTS constituent (
  constituent_id INTEGER PRIMARY KEY NOT NULL,
  gauge_id INTEGER NOT NULL,
  location TEXT DEFAULT "" NOT NULL,
  number TEXT NOT NULL,
  constituent TEXT NOT NULL,
  amplitude REAL NOT NULL,
  lag TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (gauge_id) REFERENCES gauge (gauge_id),
  UNIQUE(gauge_id, number, start_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_constituent BEFORE INSERT ON constituent
WHEN EXISTS (
  SELECT * FROM constituent
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND gauge_id = NEW.gauge_id
      AND number =  NEW.number
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on constituent");
END;
`

var constituent = Table{
	Create: constituentCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO constituent (gauge_id, location, number, constituent, amplitude, lag, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?);",
			gauge.Select(),
		)
	},
	Fields: []string{"Gauge", "Location", "Number", "Constituent", "Amplitude", "Lag", "Start Date", "End Date"},
}

const dartCreate = `
DROP TABLE IF EXISTS dart;
CREATE TABLE IF NOT EXISTS dart (
  dart_id INTEGER PRIMARY KEY NOT NULL,
  station_id INTEGER NOT NULL,
  pid TEXT NOT NULL,
  wmo_identifier TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (station_id) REFERENCES station (station_id),
  UNIQUE(station_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_dart BEFORE INSERT ON dart
WHEN EXISTS (
  SELECT * FROM dart
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND station_id =  NEW.station_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on dart");
END;
`

var dart = Table{
	Create: dartCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO dart (station_id, pid, wmo_identifier, start_date, end_date) VALUES ((%s), ?, ?, ?, ?);", station.Select())
	},
	Fields: []string{"Station", "Pid", "WMO Identifier", "Start Date", "End Date"},
}

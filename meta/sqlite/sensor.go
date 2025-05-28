package sqlite

import (
	"fmt"
)

const timingCreate = `
DROP TABLE IF EXISTS timing;
CREATE TABLE IF NOT EXISTS timing (
  timing_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  correction TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(site_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_timing BEFORE INSERT ON timing
WHEN EXISTS (
  SELECT * FROM timing
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on timing");
END;
`

var timing = Table{
	Create: timingCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT timing_id FROM timing WHERE site_id = (%s) AND start_date = ? AND end_date = ?",
			site.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO timing (site_id, correction, start_date, end_date) VALUES ((%s), ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Correction", "Start Date", "End Date"},
}

const telemetryCreate = `
DROP TABLE IF EXISTS telemetry;
CREATE TABLE IF NOT EXISTS telemetry (
  telemetry_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  scale_factor REAL NOT NULL ON CONFLICT REPLACE DEFAULT 1.0,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(site_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_telemetry BEFORE INSERT ON telemetry
WHEN EXISTS (
  SELECT * FROM telemetry
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on timing");
END;
`

var telemetry = Table{
	Create: telemetryCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT telemetry_id FROM telemetry WHERE site_id = (%s) AND start_date = ? AND end_date = ?",
			site.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO telemetry (site_id, scale_factor, start_date, end_date) VALUES ((%s), ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Scale Factor", "Start Date", "End Date"},
	Nulls: []string{
		"Scale Factor",
	},
}

const polarityCreate = `
DROP TABLE IF EXISTS polarity;
CREATE TABLE IF NOT EXISTS polarity (
  polarity_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  sublocation TEXT NULL,
  subsource TEXT NULL,
  preferred BOOLEAN NOT NULL ON CONFLICT REPLACE DEFAULT true,
  reversed BOOLEAN NOT NULL ON CONFLICT REPLACE DEFAULT false,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_polarity BEFORE INSERT ON polarity
WHEN EXISTS (
  SELECT * FROM polarity
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
      AND sublocation =  NEW.sublocation
      AND subsource =  NEW.subsource
      AND preferred =  NEW.preferred
      AND preferred =  true
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on preferred polarity");
END;
`

var polarity = Table{
	Create: polarityCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO polarity (site_id, sublocation, subsource, preferred, reversed, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Sublocation", "Subsource", "Primary", "Reversed", "Start Date", "End Date"},
}

const preampCreate = `
DROP TABLE IF EXISTS preamp;
CREATE TABLE IF NOT EXISTS preamp (
  preamp_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  subsource TEXT NULL,
  scale_factor REAL NOT NULL ON CONFLICT REPLACE DEFAULT 1.0,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(site_id, subsource, start_date, end_date)
);

CREATE TRIGGER IF NOT EXISTS no_overlap_on_preamp BEFORE INSERT ON preamp
WHEN EXISTS (
  SELECT * FROM preamp
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
      AND subsource =  NEW.subsource
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on preamp");
END;
`

var preamp = Table{
	Create: preampCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO preamp (site_id, subsource, scale_factor, start_date, end_date) VALUES ((%s), ?, ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Subsource", "Scale Factor", "Start Date", "End Date"},
}

const gainCreate = `
DROP TABLE IF EXISTS gain;
CREATE TABLE IF NOT EXISTS gain (
  gain_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  sublocation TEXT NULL,
  subsource TEXT NULL,
  scale_factor REAL NOT NULL ON CONFLICT REPLACE DEFAULT 1.0,
  scale_bias REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  absolute_bias REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(site_id, sublocation, subsource, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_gain BEFORE INSERT ON gain
WHEN EXISTS (
  SELECT * FROM gain
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id =  NEW.site_id
      AND sublocation = NEW.sublocation
      AND subsource = NEW.subsource
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on gain");
END;
`

var gain = Table{
	Create: gainCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO gain (site_id, sublocation, subsource, scale_factor, scale_bias, absolute_bias, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Sublocation", "Subsource", "Scale Factor", "Scale Bias", "Absolute Bias", "Start Date", "End Date"},
	Nulls: []string{
		"Sublocation", "Subsource", "Scale Factor", "Scale Bias", "Absolute Bias",
	},
}

const dataloggerCreate = `
DROP TABLE IF EXISTS datalogger;
CREATE TABLE IF NOT EXISTS datalogger (
  datalogger_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  place_role_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (place_role_id) REFERENCES place_role (place_role_id),
  UNIQUE(asset_id, place_role_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_datalogger BEFORE INSERT ON datalogger
WHEN EXISTS (
  SELECT * FROM datalogger
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND place_role_id = NEW.place_role_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on datalogger");
END;
`

var datalogger = Table{
	Create: dataloggerCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO datalogger (asset_id, place_role_id, start_date, end_date) VALUES ((%s), (%s), ?, ?);",
			asset.Select(), placeRole.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Place", "Role", "Start Date", "End Date"},
}

const sensorCreate = `
DROP TABLE IF EXISTS sensor;
CREATE TABLE IF NOT EXISTS sensor (
  sensor_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  site_id INTEGER NOT NULL,
  method_id INTEGER NOT NULL,
  azimuth REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  dip REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  depth REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  north REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  east REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  scale_factor REAL NOT NULL ON CONFLICT REPLACE DEFAULT 1.0,
  scale_bias REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  FOREIGN KEY (method_id) REFERENCES method (method_id),
  UNIQUE(asset_id, site_id, start_date, end_date)
);

CREATE TRIGGER IF NOT EXISTS no_overlap_on_sensor BEFORE INSERT ON sensor
WHEN EXISTS (
  SELECT * FROM sensor
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND site_id = NEW.site_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on sensor");
END;
`

var sensor = Table{
	Create: sensorCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO sensor (asset_id, site_id, method_id, azimuth, dip, depth, north, east, scale_factor, scale_bias, start_date, end_date) VALUES ((%s), (%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			asset.Select(), site.Select(), method.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Station", "Location", "Method", "Azimuth", "Dip", "Depth", "North", "East", "Scale Factor", "Scale Bias", "Start Date", "End Date"},
}

const recorderCreate = `
DROP TABLE IF EXISTS recorder;
CREATE TABLE IF NOT EXISTS recorder (
  recorder_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  model_id INTEGER NOT NULL,
  site_id INTEGER NOT NULL,
  method_id INTEGER NOT NULL,
  azimuth REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  dip REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  depth REAL NOT NULL ON CONFLICT REPLACE DEFAULT 0.0,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  FOREIGN KEY (model_id) REFERENCES model (model_id),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  FOREIGN KEY (method_id) REFERENCES method (method_id),
  UNIQUE(asset_id, model_id, site_id, start_date, end_date)
);
`

/*
CREATE TRIGGER IF NOT EXISTS no_overlap_on_recorder BEFORE INSERT ON recorder
WHEN EXISTS (
  SELECT * FROM recorder
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND model_id = NEW.model_id
      AND site_id = NEW.site_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on recorder");
END;
*/

var recorder = Table{
	Create: recorderCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO recorder (asset_id, model_id, site_id, method_id, azimuth, dip, depth, start_date, end_date) VALUES ((%s), (%s), (%s), (%s), ?, ?, ?, ?, ?);",
			asset.Select(), model.Select(), site.Select(), method.Select(),
		)
	},
	Fields: []string{"Make", "Datalogger", "Serial", "Make", "Sensor", "Station", "Location", "Method", "Azimuth", "Dip", "Depth", "Start Date", "End Date"},
}

var recorderModel = Table{
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO model (make_id, model) VALUES ((%s), ?) ON CONFLICT(make_id, model) DO NOTHING;",
			mmake.Select(),
		)
	},
	Fields: []string{"Make", "Sensor"},
}

const streamCreate = `
DROP TABLE IF EXISTS stream;
CREATE TABLE IF NOT EXISTS stream (
  stream_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  band TEXT NOT NULL ON CONFLICT REPLACE DEFAULT "",
  source TEXT NOT NULL ON CONFLICT REPLACE DEFAULT "",
  sampling_rate REAL NOT NULL,
  axial TEXT NOT NULL,
  reversed TEXT NOT NULL,
  triggered TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  UNIQUE(sampling_rate, site_id, source, start_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_stream BEFORE INSERT ON stream
WHEN EXISTS (
  SELECT * FROM stream
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
      AND source =  NEW.source
      AND sampling_rate =  NEW.sampling_rate
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on stream");
END;
`

var stream = Table{
	Create: streamCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO stream (site_id, band, source, sampling_rate, axial, reversed, triggered, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			site.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Band", "Source", "Sampling Rate", "Axial", "Reversed", "Triggered", "Start Date", "End Date"},
}

const connectionCreate = `
DROP TABLE IF EXISTS connection;
CREATE TABLE IF NOT EXISTS connection (
  connection_id INTEGER PRIMARY KEY NOT NULL,
  site_id INTEGER NOT NULL,
  place_role_id INTEGER NOT NULL,
  number TEXT NOT NULL ON CONFLICT REPLACE DEFAULT "",
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (site_id) REFERENCES site (site_id),
  FOREIGN KEY (place_role_id) REFERENCES place_role (place_role_id),
  UNIQUE(site_id, place_role_id, number, start_date, end_date)
);

CREATE TRIGGER IF NOT EXISTS no_overlap_on_connection BEFORE INSERT ON connection
WHEN EXISTS (
  SELECT * FROM connection
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND site_id = NEW.site_id
      AND place_role_id =  NEW.place_role_id
      AND number =  NEW.number
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on connection");
END;
`

var connection = Table{
	Create: connectionCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO connection (site_id, place_role_id, number, start_date, end_date) VALUES ((%s), (%s), ?, ?, ?);",
			site.Select(), placeRole.Select(),
		)
	},
	Fields: []string{"Station", "Location", "Place", "Role", "Number", "Start Date", "End Date"},
}

const placeRoleCreate = `
DROP TABLE IF EXISTS place_role;
CREATE TABLE IF NOT EXISTS place_role (
  place_role_id INTEGER PRIMARY KEY NOT NULL,
  place TEXT NOT NULL,
  role TEXT NOT NULL ON CONFLICT REPLACE DEFAULT "",
  UNIQUE (place, role)
);`

var placeRole = Table{
	Create: placeRoleCreate,
	Select: func() string {
		return "SELECT place_role_id FROM place_role WHERE place = ? AND role = ?"
	},
	Insert: func() string {
		return "INSERT INTO place_role (place, role) VALUES (?, ?) ON CONFLICT(place, role) DO NOTHING;"
	},
	Fields: []string{"Place", "Role"},
}

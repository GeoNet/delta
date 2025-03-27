package sqlite

import (
	"fmt"
)

const makeCreate = `
DROP TABLE IF EXISTS make;
CREATE TABLE IF NOT EXISTS make (
  make_id INTEGER PRIMARY KEY NOT NULL,
  make TEXT NOT NULL,
  UNIQUE (make)
);`

var mmake = Table{
	Create: makeCreate,
	Select: func() string {
		return "SELECT make_id FROM make WHERE make = ?"
	},
	Insert: func() string {
		return "INSERT INTO make (make) VALUES (?) ON CONFLICT (make) DO NOTHING;"
	},

	Fields: []string{"Make"},
}

const modelCreate = `
DROP TABLE IF EXISTS model;
CREATE TABLE IF NOT EXISTS model (
  model_id INTEGER PRIMARY KEY NOT NULL,
  make_id INTEGER NOT NULL,
  model TEXT NOT NULL,
  FOREIGN KEY (make_id) REFERENCES make (make_id),
  UNIQUE (make_id, model)
);`

var model = Table{
	Create: modelCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT model_id FROM model WHERE make_id = (%s) AND model = ?", mmake.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO model (make_id, model) VALUES ((%s), ?) ON CONFLICT (make_id, model) DO NOTHING;", mmake.Select())
	},
	Fields: []string{"Make", "Model"},
}

const assetCreate = `
DROP TABLE IF EXISTS asset;
CREATE TABLE IF NOT EXISTS asset (
  asset_id INTEGER PRIMARY KEY NOT NULL,
  model_id INTEGER NOT NULL,
  serial TEXT NOT NULL,
  number TEXT DEFAULT "" NOT NULL,
  notes TEXT DEFAULT "" NOT NULL,
  FOREIGN KEY (model_id) REFERENCES model (model_id),
  UNIQUE (model_id,serial)
);`

var asset = Table{
	Create: assetCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT asset_id FROM asset WHERE model_id = (%s) AND serial = ?", model.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO asset (model_id, serial, number, notes) VALUES ((%s), ?, ?, ?);", model.Select())
	},
	Fields: []string{"Make", "Model", "Serial", "Number", "Notes"},
}

const firmwareCreate = `
DROP TABLE IF EXISTS firmware;
CREATE TABLE IF NOT EXISTS firmware (
  firmware_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  version TEXT NOT NULL,
  notes TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  UNIQUE (asset_id, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_firmware BEFORE INSERT ON firmware
WHEN EXISTS (
  SELECT * FROM firmware
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on firmware");
END;
`

var firmware = Table{
	Create: firmwareCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT firmware_id FROM firmware WHERE asset_id = (%s) AND start_date = ? AND end_date = ?",
			asset.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO firmware (asset_id, version, notes, start_date, end_date) VALUES ((%s), ?, ?, ?, ?);",
			asset.Select())
	},
	Fields: []string{"Make", "Model", "Serial", "Version", "Notes", "Start Date", "End Date"},
}

const channelCreate = `
DROP TABLE IF EXISTS channel;
CREATE TABLE IF NOT EXISTS channel (
  channel_id INTEGER PRIMARY KEY NOT NULL,
  model_id INTEGER NOT NULL,
  response_id INTEGER NULL,
  channel_type TEXT NOT NULL,
  number REAL DEFAULT 0 NOT NULL,
  sampling_rate REAL NOT NULL,
  FOREIGN KEY (model_id) REFERENCES model (model_id),
  FOREIGN KEY (response_id) REFERENCES response (response_id),
  UNIQUE(model_id, channel_type, number, sampling_rate)
);
`

var channel = Table{
	Create: channelCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO channel (model_id, response_id, channel_type, number, sampling_rate) VALUES ((%s), (%s), ?, ?, ?);",
			model.Select(), response.Select())
	},
	Fields: []string{"Make", "Model", "Response", "Type", "Number", "SamplingRate"},
	Nulls:  []string{"Response"},
}

const componentCreate = `
DROP TABLE IF EXISTS component;
CREATE TABLE IF NOT EXISTS component (
  component_id INTEGER PRIMARY KEY NOT NULL,
  model_id INTEGER NOT NULL,
  response_id INTEGER NOT NULL,
  component_type TEXT NULL,
  number REAL NOT NULL,
  source TEXT NULL,
  subsource TEXT NOT NULL,
  dip REAL NOT NULL,
  azimuth REAL NOT NULL,
  types TEXT NOT NULL,
  sampling_rate REAL NULL,
  FOREIGN KEY (model_id) REFERENCES model (model_id),
  FOREIGN KEY (response_id) REFERENCES response (response_id),
  UNIQUE(model_id, number, source, subsource, sampling_rate)
);
`

var component = Table{
	Create: componentCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO component (model_id, response_id, component_type, number, source, subsource, dip, azimuth, types, sampling_rate) VALUES ((%s), (%s), ?, ?, ?, ?, ?, ?, ?, ?);",
			model.Select(), response.Select())
	},
	Fields: []string{"Make", "Model", "Response", "Type", "Number", "Source", "Subsource", "Dip", "Azimuth", "Types", "Sampling Rate"},
}

const calibrationCreate = `
DROP TABLE IF EXISTS calibration;
CREATE TABLE IF NOT EXISTS calibration (
  calibration_id INTEGER PRIMARY KEY NOT NULL,
  asset_id INTEGER NOT NULL,
  number TEXT NOT NULL,
  scale_factor REAL DEFAULT 1.0 NOT NULL,
  scale_bias REAL DEFAULT 0.0 NOT NULL,
  scale_absolute REAL DEFAULT 0.0 NOT NULL,
  frequency REAL NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (asset_id) REFERENCES asset (asset_id),
  UNIQUE(asset_id, number, start_date, end_date)
);
CREATE TRIGGER IF NOT EXISTS no_overlap_on_calibration BEFORE INSERT ON calibration
WHEN EXISTS (
  SELECT * FROM calibration
      WHERE datetime(start_date) <= datetime(NEW.end_date)
      AND datetime(end_date) > datetime(NEW.start_date)
      AND asset_id = NEW.asset_id
      AND number =  NEW.number
)
BEGIN
  SELECT RAISE(FAIL, "overlapping intervals on calibration");
END;
`

var calibration = Table{
	Create: calibrationCreate,
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO calibration (asset_id, number, scale_factor, scale_bias, scale_absolute, frequency, start_date, end_date) VALUES ((%s), ?, ?, ?, ?, ?, ?, ?);",
			asset.Select(),
		)
	},
	Fields: []string{"Make", "Model", "Serial", "Number", "Scale Factor", "Scale Bias", "Scale Absolute", "Frequency", "Start Date", "End Date"},
}

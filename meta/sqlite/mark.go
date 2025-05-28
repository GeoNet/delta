package sqlite

import (
	"fmt"
)

// should be loaded from a reference file
const bedrockCreate = `
DROP TABLE IF EXISTS bedrock;
CREATE TABLE IF NOT EXISTS bedrock (
  bedrock_id INTEGER PRIMARY KEY NOT NULL,
  bedrock TEXT NOT NULL,
  UNIQUE (bedrock)
);`

var bedrock = Table{
	Create: bedrockCreate,
	Select: func() string {
		return "SELECT bedrock_id FROM bedrock WHERE bedrock = ?"
	},
	Insert: func() string {
		return "INSERT INTO bedrock (bedrock) VALUES (?) ON CONFLICT(bedrock) DO NOTHING;"
	},
	Fields: []string{"Bedrock"},
}

// should be loaded from a reference file
const markTypeCreate = `
DROP TABLE IF EXISTS mark_type;
CREATE TABLE IF NOT EXISTS mark_type (
  mark_type_id INTEGER PRIMARY KEY NOT NULL,
  mark_type TEXT NOT NULL,
  UNIQUE (mark_type)
);`

var markType = Table{
	Create: markTypeCreate,
	Select: func() string {
		return "SELECT mark_type_id FROM mark_type WHERE mark_type = ?"
	},
	Insert: func() string {
		return "INSERT INTO mark_type (mark_type) VALUES (?) ON CONFLICT(mark_type) DO NOTHING;"
	},
	Fields: []string{"Mark Type"},
}

// should be loaded from a reference file
const monumentTypeCreate = `
DROP TABLE IF EXISTS monument_type;
CREATE TABLE IF NOT EXISTS monument_type (
  monument_type_id INTEGER PRIMARY KEY NOT NULL,
  monument_type TEXT NOT NULL,
  UNIQUE (monument_type)
);`

var monumentType = Table{
	Create: monumentTypeCreate,
	Select: func() string {
		return "SELECT monument_type_id FROM monument_type WHERE monument_type = ?"
	},
	Insert: func() string {
		return "INSERT INTO monument_type (monument_type) VALUES (?) ON CONFLICT(monument_type) DO NOTHING;"
	},
	Fields: []string{"Type"},
}

// should be loaded from a reference file
const foundationTypeCreate = `
DROP TABLE IF EXISTS foundation_type;
CREATE TABLE IF NOT EXISTS foundation_type (
  foundation_type_id INTEGER PRIMARY KEY NOT NULL,
  foundation_type TEXT NOT NULL,
  UNIQUE (foundation_type)
);`

var foundationType = Table{
	Create: foundationTypeCreate,
	Select: func() string {
		return "SELECT foundation_type_id FROM foundation_type WHERE foundation_type = ?"
	},
	Insert: func() string {
		return "INSERT INTO foundation_type (foundation_type) VALUES (?) ON CONFLICT(foundation_type) DO NOTHING;"
	},
	Fields: []string{"Foundation Type"},
}

// should be loaded from a reference file
const geologyCreate = `
DROP TABLE IF EXISTS geology;
CREATE TABLE IF NOT EXISTS geology (
  geology_id INTEGER PRIMARY KEY NOT NULL,
  geology TEXT NOT NULL,
  UNIQUE (geology)
);`

var geology = Table{
	Create: geologyCreate,
	Select: func() string {
		return "SELECT geology_id FROM geology WHERE geology = ?"
	},
	Insert: func() string {
		return "INSERT INTO geology (geology) VALUES (?) ON CONFLICT(geology) DO NOTHING;"
	},
	Fields: []string{"Geology"},
}

const markCreate = `
DROP TABLE IF EXISTS mark;
CREATE TABLE IF NOT EXISTS mark (
  mark_id INTEGER PRIMARY KEY NOT NULL,
  reference_id INTEGER NOT NULL,
  igs BOOLEAN NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (reference_id) REFERENCES reference (reference_id),
  UNIQUE (reference_id)
);`

var mark = Table{
	Create: markCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT mark_id FROM mark WHERE reference_id = (%s)", reference.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO mark (reference_id, igs, start_date, end_date) VALUES ((%s), ?, ?, ?);",
			reference.Select(),
		)
	},
	Fields: []string{"Mark", "Igs", "Start Date", "End Date"},
}

const markNetworkCreate = `
DROP TABLE IF EXISTS mark_network;
CREATE TABLE IF NOT EXISTS mark_network (
  mark_network_id INTEGER PRIMARY KEY NOT NULL,
  mark_id INTEGER NOT NULL,
  network_id INTEGER NOT NULL,
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  FOREIGN KEY (network_id) REFERENCES network (network_id),
  UNIQUE (mark_id, network_id)
);`

var markNetwork = Table{
	Create: markNetworkCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT mark_network_id FROM mark_network WHERE mark_id = (%s) AND network_id = (%s)",
			mark.Select(), network.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO mark_network (mark_id, network_id) VALUES ((%s), (%s));",
			mark.Select(), network.Select(),
		)
	},
	Fields: []string{"Mark", "Network"},
}

const monumentCreate = `
DROP TABLE IF EXISTS monument;
CREATE TABLE IF NOT EXISTS monument (
  monument_id INTEGER PRIMARY KEY NOT NULL,
  mark_id INTEGER NOT NULL,
  mark_type_id INTEGER NOT NULL,
  monument_type_id INTEGER NOT NULL,
  foundation_type_id INTEGER NOT NULL,
  bedrock_id INTEGER NOT NULL,
  geology_id INTEGER NOT NULL,
  domes_number TEXT NOT NULL,
  ground_relationship REAL NOT NULL,
  foundation_depth REAL NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  FOREIGN KEY (mark_type_id) REFERENCES mark_type (mark_type_id),
  FOREIGN KEY (monument_type_id) REFERENCES monument_type (monument_type_id),
  FOREIGN KEY (foundation_type_id) REFERENCES foundation_type (foundation_type_id),
  FOREIGN KEY (bedrock_id) REFERENCES bedrock (bedrock_id),
  FOREIGN KEY (geology_id) REFERENCES geology (geology_id),
  UNIQUE (mark_id)
);`

var monument = Table{
	Create: monumentCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT monument_id FROM monument WHERE mark_id = (%s)", mark.Select())
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO monument (mark_id, mark_type_id, monument_type_id, foundation_type_id, bedrock_id, geology_id, domes_number, ground_relationship, foundation_depth, start_date, end_date) VALUES ((%s), (%s), (%s), (%s), (%s), (%s), ?, ?, ?, ?, ?);",
			mark.Select(), markType.Select(), monumentType.Select(), foundationType.Select(), bedrock.Select(), geology.Select(),
		)
	},
	Fields: []string{"Mark", "Mark Type", "Type", "Foundation Type", "Bedrock", "Geology", "Domes Number", "Ground Relationship", "Foundation Depth", "Start Date", "End Date"},
}

const visibilityCreate = `
DROP TABLE IF EXISTS visibility;
CREATE TABLE IF NOT EXISTS visibility (
  visibility_id INTEGER PRIMARY KEY NOT NULL,
  mark_id INTEGER NOT NULL,
  sky_visibility TEXT NOT NULL,
  start_date DATETIME NOT NULL CHECK (start_date IS strftime('%Y-%m-%dT%H:%M:%SZ', start_date)),
  end_date DATETIME NOT NULL CHECK (end_date IS strftime('%Y-%m-%dT%H:%M:%SZ', end_date)),
  FOREIGN KEY (mark_id) REFERENCES mark (mark_id),
  UNIQUE(mark_id, sky_visibility, start_date, end_date)
);
`

var visibility = Table{
	Create: visibilityCreate,
	Select: func() string {
		return fmt.Sprintf("SELECT visibility_id FROM visibility WHERE mark_id = (%s) AND start_date = ? AND end_date = ?",
			mark.Select(),
		)
	},
	Insert: func() string {
		return fmt.Sprintf("INSERT INTO visibility (mark_id, sky_visibility, start_date, end_date) VALUES ((%s), ?, ?, ?);",
			mark.Select(),
		)
	},
	Fields: []string{"Mark", "Sky Visibility", "Start Date", "End Date"},
}

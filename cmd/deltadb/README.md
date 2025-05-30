## deltadb

This application can build a standalone sqlite database which can be used as desired in third party applications.

This has been made possible by using the pure go module `modernc.org/sqlite`, which uses some interesting automatic
conversion of the core c-code to go to build a version that doesn't require the cgo library wrapper.

### options

```
Build a DELTA Sqlite DB

Usage:

  ./deltadb [options]

Options:

  -base string
        base directory of delta files on disk
  -db string
        name of the database file on disk
  -debug
        add extra operational info
  -init
        initialise the database if a file on disk
  -resp string
        base directory of resp files on disk
  -response string
        optional database response table name to use (default "Response")
```

### example

e.g.

```
./deltadb -init -db delta.db -debug
2024/11/03 20:49:20 initialise database
2024/11/03 20:49:22 database initialised in 2.179325139s
```

and then to examine the file

```
sqlite3  delta.db
SQLite version 3.46.1 2024-08-13 09:16:08
Enter ".help" for usage hints.
sqlite> .schema network
CREATE TABLE network (
  network_id INTEGER PRIMARY KEY NOT NULL,
  network TEXT NOT NULL,
  external TEXT NOT NULL,
  description TEXT DEFAULT "" NOT NULL,
  restricted BOOLEAN DEFAULT false NOT NULL,
  UNIQUE (network)
);
sqlite>
```

## schema

[![Schema](delta.svg)](delta.svg)

### building the schema

This uses the application from [sqlite-schema-diagram](https://github.com/o0101/sqlite-schema-diagram) and graphviz.

```
sqlite3 path/to/database.db -init sqlite-schema-diagram.sql "" > schema.dot
dot -Tsvg schema.dot > schema.svg
```


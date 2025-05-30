## deltadb

This application can build a standalone sqlite database which can be used as desired in third party applications.

This has been made possible by using the pure go module `modernc.org/sqlite`, which uses some interesting automatic
conversion of the core c-code to go to build a version that doesn't require the cgo library wrapper.

### options

```
Build and initialise a DELTA Sqlite database

Usage:

  ./deltadb [options]

Options:

  -base string
        base directory of delta files on disk, default uses embedded files
  -path string
        name of the database file on disk, default is to use memory only
  -resp string
        base directory of resp files on disk, default uses embedded files
```

### example

e.g.

```
./deltadb -path delta.db
2024/11/03 20:49:20 initialise database
2024/11/03 20:49:22 database initialised in 2.179325139s
```

and then to examine the file

```
sqlite3 delta.db
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

### building the schema

This uses the application from [sqlite-schema-diagram](https://github.com/o0101/sqlite-schema-diagram) and graphviz.

```
sqlite3 delta.db -init sqlite-schema-diagram.sql "" | dot -Tsvg > delta.svg
```

## design notes

There are a few key elements to how the delta schema is laid out.
The first relates to where equipment is installed or where measurements are read from.
The second relates to the equipment and where and when it was installed, and a third main area
is how equipment was configured and its related capabilities.
Finally, there are entries related to the physical properties of what is being measured,
which may be needed to interpret the recorded data.

The main table for equipment recording and measurement readings is `location`.
An entry includes a unique `code` and a geographic location.
Specific types of observations can be associated with the location entry,
e.g. GNSS marks, or seismic stations.
The benefit of grouping these together is the ability to reference a location `code`
without knowing what it is actually used for.
This makes discovery easier and allows cross field tables, such as `notes`.

It is assumed that any location can be part of one or possibly more networks.
The `network` cannot be assumed from the location directly but all locations
for a given `network` would be readily available.

## schema

[![Schema](/cmd/deltadb/delta.svg)](delta.svg)

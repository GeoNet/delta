package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/GeoNet/delta/meta"
)

type Database struct {
	db     *sql.DB
	schema string
}

func New(db *sql.DB, schema string) Database {
	return Database{
		db:     db,
		schema: schema,
	}
}
func (d Database) Schema() string {
	if d.schema != "" {
		return d.schema + "."
	}
	return ""
}

func (d Database) exec(ctx context.Context, tx *sql.Tx, cmds ...string) error {
	for _, cmd := range cmds {
		if _, err := tx.ExecContext(ctx, cmd); err != nil {
			return fmt.Errorf("cmd %q: %w", cmd, err)
		}
	}
	return nil
}

func (d Database) prepare(ctx context.Context, tx *sql.Tx, cmd string, values ...[]any) error {

	stmt, err := tx.PrepareContext(ctx, cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range values {
		if _, err := stmt.ExecContext(ctx, v...); err != nil {
			return err
		}
	}

	return nil
}

func (d Database) Init(ctx context.Context, tables []meta.TableList) error {

	// Get a Tx for making transaction requests.
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails, not actually
	// worried about any rollback error.
	defer func() { _ = tx.Rollback() }()

	for _, t := range tables {
		if err := d.exec(ctx, tx, d.create(t.Table)...); err != nil {
			return err
		}

		cmd, values, ok := d.insert(t.Table, t.List)
		if !ok {
			continue
		}

		if err := d.prepare(ctx, tx, cmd, values...); err != nil {
			return err
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d Database) create(table meta.Table) []string {
	var drop strings.Builder
	fmt.Fprintf(&drop, "DROP TABLE IF EXISTS %s%s;\n", d.Schema(), table.Name())

	var create strings.Builder

	var primary []string
	for n, x := range table.Columns() {
		if !table.IsPrimary(n) {
			continue
		}
		primary = append(primary, table.Remap(x))
	}

	fmt.Fprintf(&create, "CREATE TABLE IF NOT EXISTS %s%s(\n", d.Schema(), table.Name())
	for n, x := range table.Columns() {
		if n > 0 {
			fmt.Fprintf(&create, ",\n")
		}
		switch {
		case table.IsPrimary(n) && len(primary) == 1:
			fmt.Fprintf(&create, "  %s TEXT PRIMARY KEY", table.Remap(x))
		case table.IsNative(n):
			fmt.Fprintf(&create, "  %s REAL", table.Remap(x))
		case table.IsDateTime(n):
			fmt.Fprintf(&create, "  %s DATETIME CHECK (%s IS strftime('%%Y-%%m-%%dT%%H:%%M:%%SZ', %s))", table.Remap(x), table.Remap(x), table.Remap(x))
		default:
			fmt.Fprintf(&create, "  %s TEXT", table.Remap(x))
		}
	}
	if len(primary) > 1 {
		fmt.Fprintf(&create, ",\n  PRIMARY KEY(%s)", strings.Join(primary, ","))
	}

	foreign := make(map[string][]string)
	for n, x := range table.Columns() {
		if v, ok := table.IsForeign(n); ok {
			foreign[v] = append(foreign[v], table.Remap(x))

		}
	}

	if len(foreign) > 0 {
		for k, v := range foreign {
			fmt.Fprintf(&create, ",\n  FOREIGN KEY(%s) REFERENCES %s (%s)", strings.Join(v, ","), k, strings.Join(v, ","))
		}
	}

	fmt.Fprintln(&create, "\n);")

	var trigger strings.Builder
	if start, end, ok := table.HasDateTime(); ok {
		var primary []string
		for n, x := range table.Columns() {
			if !table.IsPrimary(n) {
				continue
			}
			if t, ok := table.Start(); ok && x == t {
				continue
			}
			if t, ok := table.End(); ok && x == t {
				continue
			}
			primary = append(primary, table.Remap(x))
		}
		fmt.Fprintf(&trigger, "CREATE TRIGGER IF NOT EXISTS NoOverlapOn%s", table.Name())
		fmt.Fprintf(&trigger, " BEFORE INSERT ON %s%s", d.Schema(), table.Name())
		fmt.Fprintf(&trigger, " WHEN EXISTS (\n  SELECT * FROM %s%s\n    WHERE ", d.Schema(), table.Name())
		if len(primary) > 0 {
			for n, v := range primary {
				if n > 0 {
					fmt.Fprintf(&trigger, "\n    AND ")
				}
				fmt.Fprintf(&trigger, "%s == NEW.%s", v, v)
			}
			fmt.Fprintf(&trigger, "\n    AND ")
		}
		fmt.Fprintf(&trigger, "datetime(%s) <= datetime(NEW.%s)\n    AND ", table.Remap(start), table.Remap(end))
		fmt.Fprintf(&trigger, "datetime(%s) >  datetime(NEW.%s)\n)\n", table.Remap(end), table.Remap(start))
		fmt.Fprintf(&trigger, "\nBEGIN\n")
		fmt.Fprintf(&trigger, "SELECT RAISE(FAIL, \"Overlapping Intervals on %s%s\");\n", d.Schema(), table.Name())
		fmt.Fprintf(&trigger, "END;\n")
	}

	return []string{drop.String(), create.String(), trigger.String()}
}

func (d Database) insert(table meta.Table, list meta.ListEncoder) (string, [][]any, bool) {

	lines := table.Encode(list)
	if !(len(lines) > 0) {
		return "", nil, false
	}

	var header []string
	for _, x := range lines[0] {
		header = append(header, table.Remap(x))
	}
	var parts []string
	for n := range header {
		parts = append(parts, fmt.Sprintf("$%d", n+1))
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "INSERT INTO %s%s (%s) VALUES (%s);\n", d.Schema(), table.Name(), strings.Join(header, ","), strings.Join(parts, ","))

	var values [][]any
	for _, line := range lines[1:] {
		var parts []any
		for n, p := range line {
			switch {
			case table.IsNative(n) && p == "":
				parts = append(parts, "0")
			default:
				parts = append(parts, p)
			}
		}
		values = append(values, parts)
	}

	return sb.String(), values, true
}

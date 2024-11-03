package sqlite

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GeoNet/delta/meta"
)

type Database struct {
	schema string
}

func New(schema string) Database {
	return Database{
		schema: schema,
	}
}
func (d Database) Schema() string {
	if d.schema != "" {
		return d.schema + "."
	}
	return ""
}

func (d Database) Drop(table meta.Table) []string {
	var drop strings.Builder
	fmt.Fprintf(&drop, "DROP TABLE IF EXISTS %s%s;\n", d.Schema(), table.Name())
	return []string{drop.String()}
}

func (d Database) Create(table meta.Table) []string {
	var create strings.Builder

	var primary []string
	for n, x := range table.Columns() {
		if !table.IsPrimary(n) {
			continue
		}
		primary = append(primary, table.Remap(x))
	}

	fmt.Fprintf(&create, "CREATE TABLE IF NOT EXISTS %s%s (\n", d.Schema(), table.Name())
	for n, x := range table.Columns() {
		if n > 0 {
			fmt.Fprintf(&create, ",\n")
		}
		switch {
		case table.IsPrimary(n) && len(primary) == 1:
			fmt.Fprintf(&create, "\t%s TEXT PRIMARY KEY", table.Remap(x))
		case table.IsNative(n):
			fmt.Fprintf(&create, "\t%s REAL", table.Remap(x))
		case table.IsDateTime(n):
			fmt.Fprintf(&create, "\t%s DATETIME", table.Remap(x))
		default:
			fmt.Fprintf(&create, "\t%s TEXT", table.Remap(x))
		}
	}
	if len(primary) > 1 {
		fmt.Fprintf(&create, ",\n\n\tPRIMARY KEY(%s)", strings.Join(primary, ","))
	}

	foreign := make(map[string][]string)
	for n, x := range table.Columns() {
		if v, ok := table.IsForeign(n); ok {
			foreign[v] = append(foreign[v], table.Remap(x))

		}
	}

	if len(foreign) > 0 {
		for k, v := range foreign {
			fmt.Fprintf(&create, ",\n\tFOREIGN KEY(%s) REFERENCES %s (%s)\n", strings.Join(v, ","), k, strings.Join(v, ","))
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
		fmt.Fprintf(&trigger, "\nCREATE TRIGGER IF NOT EXISTS NoOverlapOn%s", table.Name())
		fmt.Fprintf(&trigger, " BEFORE INSERT ON %s%s", d.Schema(), table.Name())
		fmt.Fprintf(&trigger, " WHEN EXISTS (\n\tSELECT * FROM %s%s\n\t\tWHERE ", d.Schema(), table.Name())
		if len(primary) > 0 {
			for n, v := range primary {
				if n > 0 {
					fmt.Fprintf(&trigger, "\n\t\tAND ")
				}
				fmt.Fprintf(&trigger, "\"%s\" == NEW.\"%s\"", v, v)
			}
			fmt.Fprintf(&trigger, "\n\t\tAND ")
		}
		fmt.Fprintf(&trigger, "\"%s\" < NEW.\"%s\"\n\t\tAND ", table.Remap(start), table.Remap(end))
		fmt.Fprintf(&trigger, "\"%s\" >  NEW.\"%s\"\n)\n", table.Remap(end), table.Remap(start))
		fmt.Fprintf(&trigger, "\nBEGIN\n")
		fmt.Fprintf(&trigger, "SELECT RAISE(FAIL, \"Overlapping Intervals on %s%s\");\n", d.Schema(), table.Name())
		fmt.Fprintf(&trigger, "END;\n")
	}

	return []string{create.String(), trigger.String()}
}

func (d Database) Insert(table meta.Table, list meta.ListEncoder) []string {
	var sb strings.Builder

	lines := table.Encode(list)
	if !(len(lines) > 0) {
		return nil
	}

	var header []string
	for _, x := range lines[0] {
		header = append(header, table.Remap(x))
	}

	for _, line := range lines[1:] {
		var parts []string
		for n, p := range line {
			switch {
			case table.IsNative(n) && p == "":
				parts = append(parts, "0")
			case table.IsNative(n):
				parts = append(parts, p)
			default:
				parts = append(parts, strconv.Quote(p))
			}
		}
		fmt.Fprintf(&sb, "INSERT INTO %s%s (%s) VALUES (%s);\n", d.Schema(), table.Name(), strings.Join(header, ","), strings.Join(parts, ","))
	}

	return []string{sb.String()}
}

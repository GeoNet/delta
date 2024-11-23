package sqlite

import (
	"strings"

	"github.com/GeoNet/delta/meta"
)

type Table struct {
	// Create holds the SQL code needed to create the table and all associated
	// triggers and constraints.
	Create string
	// Select holds the prepared statement that can be used to select the primary
	// key from the table.
	Select func() string
	// Insert holds the full prepared statement to insert a row into the table.
	Insert func() string
	// Fields gives the delta csv column names to use for inserting a row into the table.
	Fields []string
	// Nulls holds the set of columns that are allowed to be NULL in the table, an empty
	// string in the CSV field will indicate a NULL value should be passed into the row.
	Nulls []string
	// Unwrap can be used to build a linking table when a column has multiple fields
	Unwrap string
}

// Links returns the rows to insert into a Linking table for the given unwrapping column.
func (t Table) Links(list meta.TableList) [][]any {

	lines := list.Table.Encode(list.List)
	if !(len(lines) > 0) {
		return nil
	}

	lookup := make(map[string]int)
	for n, v := range list.Table.Columns() {
		lookup[v] = n
	}

	w, ok := lookup[t.Unwrap]
	if !ok {
		return nil
	}

	var res [][]any
	for _, line := range lines[1:] {
		if !(w < len(line)) {
			continue
		}
		for _, c := range strings.Fields(strings.TrimSpace(line[w])) {
			var parts []any
			for _, f := range t.Fields {
				n, ok := lookup[f]
				if !ok {
					return nil
				}
				if !(n < len(line)) {
					return nil
				}
				switch {
				case n == w:
					parts = append(parts, c)
				default:
					parts = append(parts, line[n])
				}
			}
			res = append(res, parts)
		}
	}
	return res
}

// checkNull returns null if the value is empty and is a member
// of the nulls map, otherwise it returns the value.
func checkNull(nulls map[string]interface{}, key, value string) any {
	if _, ok := nulls[key]; ok && value == "" {
		return nil
	}
	return value
}

// Columns returns the expected rows for the given TableList.
func (t Table) Columns(list meta.TableList) [][]any {

	lines := list.Table.Encode(list.List)
	if !(len(lines) > 0) {
		return nil
	}

	nulls := make(map[string]interface{})
	for _, v := range t.Nulls {
		nulls[v] = true
	}

	lookup := make(map[string]int)
	for n, v := range list.Table.Columns() {
		lookup[v] = n
	}

	var res [][]any
	for _, line := range lines[1:] {
		var parts []any

		for _, f := range t.Fields {
			n, ok := lookup[f]
			if !ok {
				return nil
			}
			if !(n < len(line)) {
				return nil
			}

			parts = append(parts, checkNull(nulls, f, line[n]))
		}

		res = append(res, parts)
	}

	return res
}

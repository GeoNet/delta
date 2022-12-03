package meta

import (
	"sort"
)

// Header is used to manage the encoding and decoding of csv formatted files by using their column headers.
// This allows adding or removing columns in the file without breaking existing applications that may not
// understand or expect extra columns.
//
// An example usage where data is a slice of csv rows.
//
//	var headers Header = map[string]int{
//	  "Label":    labelId,
//	  ...
//	}
//
//	fields := headers.Fields(data[0])
//	for _, v := range data[1:] {
//	  d := fields.Remap(v)
//
//	  label :=  d[labelId]
//
//	  ...
//	}
type Header map[string]int

// Columns will return the csv column headers in the sorted order of the Header values.
func (h Header) Columns() []string {
	var columns []string
	for k := range h {
		columns = append(columns, k)
	}
	sort.Slice(columns, func(i, j int) bool {
		return h[columns[i]] < h[columns[j]]
	})
	return columns
}

// Fields converts a csv header row into a Column for use with field name lookups.
func (h Header) Fields(columns []string) Field {
	lookup := make(map[int]int)
	for i, v := range columns {
		lookup[h[v]] = i
	}
	return lookup
}

// Field is a mapping between the expected field entry and the csv file column number. This allows finding
// the desired column without being explicitly constrained by the column contents or order.
type Field map[int]int

// Remap will convert a csv data row into a entry key map.
func (f Field) Remap(data []string) map[int]string {
	remap := make(map[int]string)
	for i, v := range data {
		remap[f[i]] = v
	}
	return remap
}

package ntrip

import (
	"fmt"
	"strings"
)

const (
	formatFormat int = iota
	formatDetails
	formatLast
)

// Format represents the RTCM message types with update periods in parenthesis in seconds.
type Format struct {
	Format  string
	Details []string
}

func (f *Format) decode(row []string) error {
	if l := len(row); l != formatLast {
		return fmt.Errorf("incorrect \"format\" \"%s\": found %d items, expected %d", strings.Join(row, ","), l, formatLast)
	}

	var details []string
	for _, d := range strings.Split(row[formatDetails], ":") {
		details = append(details, strings.TrimSpace(d))
	}

	*f = Format{
		Format:  strings.TrimSpace(row[formatFormat]),
		Details: details,
	}

	return nil
}

func (f Format) encode() []string {
	var row []string

	row = append(row, f.Format)
	row = append(row, strings.Join(f.Details, ":"))

	return row
}

type Formats []Format

func ReadFormats(path string) (map[string][]string, error) {
	var formats Formats
	if err := ReadFile(path, &formats); err != nil {
		return nil, err
	}

	res := make(map[string][]string)
	for _, f := range formats {
		if _, ok := res[f.Format]; ok {
			return nil, fmt.Errorf("duplicate format name: %s", f.Format)
		}
		res[f.Format] = f.Details
	}

	return res, nil
}

func (f Formats) Header() []string {
	return []string{"#Format", "Details"}
}

func (f Formats) Fields() int {
	return formatLast
}

func (f *Formats) Decode(data [][]string) error {
	for _, v := range data {
		var format Format
		if err := format.decode(v); err != nil {
			return err
		}
		*f = append(*f, format)
	}
	return nil
}

func (f Formats) Encode() [][]string {
	var items [][]string

	for _, format := range f {
		items = append(items, format.encode())
	}

	return items
}

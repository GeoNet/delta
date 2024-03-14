package ntrip

import (
	"fmt"
	"sort"
	"strings"
)

const (
	mountMount int = iota
	mountMark
	mountCountry
	mountFormat
	mountDetails
	mountNavigation
	mountGroups
	mountUser
	mountAddress
	mountLast
)

// Mount represents an ntripcaster mount point.
type Mount struct {
	Mount      string
	Mark       string
	Country    string
	Format     string
	Details    string
	Navigation string
	Groups     []string
	User       string
	Address    string
}

func (m *Mount) decode(row []string) error {
	if l := len(row); l != mountLast {
		return fmt.Errorf("incorrect \"mount\" \"%s\": found %d items, expected %d", strings.Join(row, ","), l, mountLast)
	}

	var groups []string
	for _, g := range strings.Split(row[mountGroups], ":") {
		groups = append(groups, strings.TrimSpace(g))
	}

	sort.Strings(groups)

	*m = Mount{
		Mount:      strings.TrimSpace(row[mountMount]),
		Mark:       strings.TrimSpace(row[mountMark]),
		Country:    strings.TrimSpace(row[mountCountry]),
		Format:     strings.TrimSpace(row[mountFormat]),
		Details:    strings.TrimSpace(row[mountDetails]),
		Navigation: strings.TrimSpace(row[mountNavigation]),
		Groups:     groups,
		User:       strings.TrimSpace(row[mountUser]),
		Address:    strings.TrimSpace(row[mountAddress]),
	}

	return nil
}

func (m Mount) encode() []string {
	var row []string

	row = append(row, m.Mount)
	row = append(row, m.Mark)
	row = append(row, m.Country)
	row = append(row, m.Format)
	row = append(row, m.Details)
	row = append(row, m.Navigation)
	row = append(row, strings.Join(m.Groups, ":"))
	row = append(row, m.User)
	row = append(row, m.Address)

	return row
}

// Mounts represents a list of ntripcaster mount information.
type Mounts []Mount

func ReadMounts(path string) ([]Mount, error) {
	var mounts Mounts
	if err := ReadFile(path, &mounts); err != nil {
		return nil, err
	}

	check := make(map[string]Mount)
	for _, m := range mounts {
		if _, ok := check[m.Mount]; ok {
			return nil, fmt.Errorf("duplicate mount name: %s", m.Mount)
		}
		check[m.Mount] = m
	}

	sort.Slice(mounts, func(i, j int) bool { return mounts[i].Mount < mounts[j].Mount })

	return mounts, nil
}

// Header returns a slice of the expected csv header fields.
func (m Mounts) Header() []string {
	return []string{"#Mount", "Mark", "Country", "Format", "Details", "Navigation", "Groups", "User", "Address"}

}

// Fields returns the number of expected csv fields.
func (m Mounts) Fields() int {
	return mountLast
}

// Encode builds a set of string slices representing the csv mount entries.
func (m Mounts) Encode() [][]string {
	var items [][]string

	for _, format := range m {
		items = append(items, format.encode())
	}

	return items
}

// Decode extracts a set of mount information from rows as expected from csv mount entries.
func (m *Mounts) Decode(data [][]string) error {
	for _, v := range data {
		var mount Mount
		if err := mount.decode(v); err != nil {
			return err
		}
		*m = append(*m, mount)
	}
	return nil
}

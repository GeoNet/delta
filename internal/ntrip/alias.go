package ntrip

import (
	"fmt"
	"sort"
	"strings"
)

const (
	aliasAlias int = iota
	aliasMount
	aliasLast
)

// Mount represents an ntripcaster alias mount point.
type Alias struct {
	Alias string
	Mount string
}

func (a *Alias) decode(row []string) error {
	if l := len(row); l != aliasLast {
		return fmt.Errorf("incorrect \"alias\" \"%s\": found %d items, expected %d", strings.Join(row, ","), l, aliasLast)
	}

	*a = Alias{
		Alias: strings.TrimSpace(row[aliasAlias]),
		Mount: strings.TrimSpace(row[aliasMount]),
	}

	return nil
}

func (a Alias) encode() []string {
	var row []string

	row = append(row, a.Alias)
	row = append(row, a.Mount)

	return row
}

// Mounts represents a list of ntripcaster alias mount information.
type Aliases []Alias

func ReadAliases(path string) ([]Alias, error) {
	var aliases Aliases
	if err := ReadFile(path, &aliases); err != nil {
		return nil, err
	}

	check := make(map[string]Alias)
	for _, a := range aliases {
		if _, ok := check[a.Alias]; ok {
			return nil, fmt.Errorf("duplicate alias name: %s", a.Alias)
		}
		check[a.Alias] = a
	}

	sort.Slice(aliases, func(i, j int) bool { return aliases[i].Alias < aliases[j].Alias })

	return aliases, nil
}

// Header returns a slice of the expected csv header fields.
func (a Aliases) Header() []string {
	return []string{"#Alias", "Mount"}
}

// Fields returns the number of expected csv fields.
func (a Aliases) Fields() int {
	return aliasLast
}

// Encode builds a set of string slices representing the csv alias mount entries.
func (a Aliases) Encode() [][]string {
	var items [][]string

	for _, format := range a {
		items = append(items, format.encode())
	}

	return items
}

// Decode extracts a set of alias mount information from rows as expected from csv mount entries.
func (a *Aliases) Decode(data [][]string) error {
	for _, v := range data {
		var alias Alias
		if err := alias.decode(v); err != nil {
			return err
		}
		*a = append(*a, alias)
	}
	return nil
}

package meta

import (
	"sort"
	"strings"
)

const (
	domainName int = iota
	domainDescription
	domainLast
)

var domainHeaders Header = map[string]int{
	"Domain":      domainName,
	"Description": domainDescription,
}

type Domain struct {
	Name        string
	Description string
}

type DomainList []Domain

func (d DomainList) Len() int           { return len(d) }
func (d DomainList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DomainList) Less(i, j int) bool { return d[i].Name < d[j].Name }

func (d DomainList) encode() [][]string {
	var data [][]string

	data = append(data, domainHeaders.Columns())

	for _, row := range d {
		data = append(data, []string{
			row.Name,
			row.Description,
		})
	}

	return data
}

func (d *DomainList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var domains []Domain

	fields := domainHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		domains = append(domains, Domain{
			Name:        strings.TrimSpace(d[domainName]),
			Description: strings.TrimSpace(d[domainDescription]),
		})
	}

	*d = DomainList(domains)

	return nil
}

func LoadDomains(path string) ([]Domain, error) {
	var d []Domain

	if err := LoadList(path, (*DomainList)(&d)); err != nil {
		return nil, err
	}

	sort.Sort(DomainList(d))

	return d, nil
}

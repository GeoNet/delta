package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	domainStation int = iota
	domainNetwork
	domainDomain
	domainStart
	domainEnd
	domainLast
)

type Domain struct {
	Span

	Network string
	Station string
	Domain  string
}

type DomainList []Domain

func (n DomainList) Len() int      { return len(n) }
func (n DomainList) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n DomainList) Less(i, j int) bool {
	switch {
	case n[i].Station < n[j].Station:
		return true
	case n[i].Station > n[j].Station:
		return false
	case n[i].Network < n[j].Network:
		return true
	case n[i].Network > n[j].Network:
		return false
	case n[i].Domain < n[j].Domain:
		return true
	case n[i].Domain > n[j].Domain:
		return false
	case n[i].Start.Before(n[j].Start):
		return true
	default:
		return false
	}
}

func (n DomainList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Network",
		"Domain",
		"Start Date",
		"End Date",
	}}
	for _, v := range n {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Domain),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (n *DomainList) decode(data [][]string) error {
	var domains []Domain
	if !(len(data) > 1) {
		*n = DomainList(domains)
	}

	for _, d := range data[1:] {
		if len(d) != domainLast {
			return fmt.Errorf("incorrect number of installed domain fields")
		}

		start, err := time.Parse(DateTimeFormat, d[domainStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[domainEnd])
		if err != nil {
			return err
		}

		domains = append(domains, Domain{
			Station: strings.TrimSpace(d[domainStation]),
			Network: strings.TrimSpace(d[domainNetwork]),
			Domain:  strings.TrimSpace(d[domainDomain]),
			Span: Span{
				Start: start,
				End:   end,
			},
		})
	}

	*n = DomainList(domains)

	return nil
}

func LoadDomains(path string) ([]Domain, error) {
	var n []Domain

	if err := LoadList(path, (*DomainList)(&n)); err != nil {
		return nil, err
	}

	sort.Sort(DomainList(n))

	return n, nil
}

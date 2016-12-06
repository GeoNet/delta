package meta

import (
	"fmt"
	"sort"
	"strings"
)

const (
	associationCode = iota
	associationTags
	associationLast
)

type Association struct {
	Code string
	Tags []string
}

type AssociationList []Association

func (m AssociationList) Len() int           { return len(m) }
func (m AssociationList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m AssociationList) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m AssociationList) encode() [][]string {
	data := [][]string{{
		"Code",
		"Tags",
	}}
	for _, v := range m {
		var tags []string
		for _, t := range v.Tags {
			tags = append(tags, strings.TrimSpace(t))
		}
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.Join(tags, "|"),
		})
	}
	return data
}

func (m *AssociationList) decode(data [][]string) error {
	var services []Association
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != serviceLast {
				return fmt.Errorf("incorrect number of installed service fields")
			}

			services = append(services, Association{
				Code: strings.TrimSpace(d[associationCode]),
				Tags: strings.Split(d[associationTags], "|"),
			})
		}

		*m = AssociationList(services)
	}
	return nil
}

func LoadAssociations(path string) ([]Association, error) {
	var m []Association

	if err := LoadList(path, (*AssociationList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(AssociationList(m))

	return m, nil
}

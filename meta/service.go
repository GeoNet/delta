package meta

import (
	"fmt"
	"sort"
	"strings"
)

const (
	serviceTag = iota
	serviceDescription
	serviceLast
)

type Service struct {
	Tag         string
	Description string
}

type ServiceList []Service

func (m ServiceList) Len() int           { return len(m) }
func (m ServiceList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ServiceList) Less(i, j int) bool { return m[i].Tag < m[j].Tag }

func (m ServiceList) encode() [][]string {
	data := [][]string{{
		"Tag",
		"Description",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Tag),
			strings.TrimSpace(v.Description),
		})
	}
	return data
}

func (m *ServiceList) decode(data [][]string) error {
	var services []Service
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != serviceLast {
				return fmt.Errorf("incorrect number of installed service fields")
			}

			services = append(services, Service{
				Tag:         strings.TrimSpace(d[serviceTag]),
				Description: strings.TrimSpace(d[serviceDescription]),
			})
		}

		*m = ServiceList(services)
	}
	return nil
}

func LoadServices(path string) ([]Service, error) {
	var m []Service

	if err := LoadList(path, (*ServiceList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(ServiceList(m))

	return m, nil
}

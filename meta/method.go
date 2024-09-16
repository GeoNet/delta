package meta

import (
	"sort"
	"strings"
)

const (
	methodDomain int = iota
	methodName
	methodDescription
	methodReference
	methodLast
)

var methodHeaders Header = map[string]int{
	"Domain":      methodDomain,
	"Method":      methodName,
	"Description": methodDescription,
	"Reference":   methodReference,
}

type Method struct {
	Domain      string
	Name        string
	Description string
	Reference   string
}

type MethodList []Method

func (m MethodList) Len() int      { return len(m) }
func (m MethodList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m MethodList) Less(i, j int) bool {
	switch {
	case m[i].Domain < m[j].Domain:
		return true
	case m[i].Domain < m[j].Domain:
		return false
	case m[i].Name < m[j].Name:
		return true
	default:
		return false
	}
}

func (m MethodList) encode() [][]string {
	var data [][]string

	data = append(data, methodHeaders.Columns())

	for _, row := range m {
		data = append(data, []string{
			row.Domain,
			row.Name,
			row.Description,
			row.Reference,
		})
	}

	return data
}

func (m *MethodList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var methods []Method

	fields := methodHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		methods = append(methods, Method{
			Domain:      strings.TrimSpace(d[methodDomain]),
			Name:        strings.TrimSpace(d[methodName]),
			Description: strings.TrimSpace(d[methodDescription]),
			Reference:   strings.TrimSpace(d[methodReference]),
		})
	}

	*m = MethodList(methods)

	return nil
}

func LoadMethods(path string) ([]Method, error) {
	var m []Method

	if err := LoadList(path, (*MethodList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MethodList(m))

	return m, nil
}

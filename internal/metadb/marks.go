package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type marks struct {
	list     meta.MarkList
	lookup   map[string]meta.Mark
	networks map[string][]meta.Mark
	once     sync.Once
}

func (m *marks) loadMarks(base string) error {
	var err error

	m.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "network", "marks.csv"), &m.list); err == nil {
			lookup := make(map[string]meta.Mark)
			for _, v := range m.list {
				lookup[v.Code] = v
			}
			m.lookup = lookup

			networks := make(map[string][]meta.Mark)
			for _, v := range m.list {
				networks[v.Network] = append(networks[v.Network], v)
			}
			m.networks = networks
		}
	})

	return err
}

func (m *MetaDB) Marks() ([]meta.Mark, error) {

	if err := m.loadMarks(m.base); err != nil {
		return nil, err
	}

	return m.marks.list, nil
}

func (m *MetaDB) Mark(code string) (*meta.Mark, error) {

	if err := m.loadMarks(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.marks.lookup[code]; ok {
		return &s, nil
	}

	return nil, nil
}

func (m *MetaDB) NetworkMark(code string) ([]meta.Mark, error) {

	if err := m.loadMarks(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.marks.networks[code]; ok {
		return s, nil
	}

	return nil, nil
}

package metadb

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/GeoNet/delta/meta"
)

type streams struct {
	list      meta.StreamList
	locations map[string]map[string][]meta.Stream
	once      sync.Once
}

func (s *streams) loadStreams(base string) error {
	var err error

	s.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "streams.csv"), &s.list); err == nil {
			locations := make(map[string]map[string][]meta.Stream)
			for _, v := range s.list {
				if _, ok := locations[v.Station]; !ok {
					locations[v.Station] = make(map[string][]meta.Stream)
				}
				if _, ok := locations[v.Station][v.Location]; !ok {
					locations[v.Station][v.Location] = []meta.Stream{}
				}
				locations[v.Station][v.Location] = append(locations[v.Station][v.Location], v)
			}
			s.locations = locations
		}
	})

	return err
}

func (m *MetaDB) StationLocationSamplingRateStartStream(sta, loc string, rate float64, start time.Time) (*meta.Stream, error) {
	if err := m.loadStreams(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.streams.locations[sta]; ok {
		if l, ok := s[loc]; ok {
			for _, r := range l {
				switch {
				case r.Start.After(start):
				case r.End.Before(start):
				case r.SamplingRate != rate:
				default:
					return &r, nil
				}
			}
		}
	}

	return nil, nil
}

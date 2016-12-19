package metadb

import (
	"time"

	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) StationLocationStreams(station, location string) ([]meta.Stream, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("stream", "site", station, location)
	if err != nil {
		return nil, err
	}
	var streams []meta.Stream
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		streams = append(streams, s.(meta.Stream))
	}

	return streams, nil
}

func (db *MetaDB) StationLocationSamplingRateStartStream(station, location string, rate float64, start time.Time) (*meta.Stream, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("stream", "interval", station, location, rate)
	if err != nil {
		return nil, err
	}
	var stream *meta.Stream
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		if s.(meta.Stream).End.Before(start) {
			continue
		}
		if s.(meta.Stream).Start.After(start) {
			break
		}
		result := s.(meta.Stream)
		stream = &result
	}

	return stream, nil
}

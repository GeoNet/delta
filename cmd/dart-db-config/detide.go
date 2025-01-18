package main

import (
	"context"
	"time"

	"github.com/GeoNet/delta"
)

func QueryDetide(db *delta.DB, buoy string, on, off time.Time) (*Detide, error) {

	detide, err := QueryGauge(db, buoy, on, off)
	if err != nil {
		return nil, err
	}

	constituents, err := QueryConsituents(db, buoy, on, off)
	if err != nil {
		return nil, err
	}

	detide.Constituents = constituents

	return detide, nil
}

func QueryGauge(db *delta.DB, buoy string, on, off time.Time) (*Detide, error) {

	ctx := context.Background()

	query := `
	SELECT
	  	gauge.analysis_time_zone,
  		gauge.analysis_latitude
	FROM
		gauge
	WHERE
		gauge.gauge = ?
	AND
		datetime(gauge.end_date) > ?
	AND
		datetime(gauge.start_date) <= ?
	LIMIT 1
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, buoy, on, off)

	var detide Detide
	if err := row.Scan(
		&detide.TimeZone,
		&detide.Latitude,
	); err != nil {
		return nil, err
	}

	return &detide, nil
}

// Find constituent values need by the de-tiding algorithm
func QueryConsituents(db *delta.DB, buoy string, on, off time.Time) ([]Constituent, error) {

	ctx := context.Background()

	query := `
	SELECT
	  	constituent.constituent,
  		constituent.number,
  		constituent.amplitude,
  		constituent.lag
	FROM
		constituent
		INNER JOIN gauge ON gauge.gauge_id = constituent.gauge_id
	WHERE
		gauge.gauge = ?  AND
		datetime(constituent.end_date) > ?  AND
		datetime(constituent.start_date) <= ?
	ORDER BY
		CAST(constituent.number AS INTEGER)
	;`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	results, err := stmt.QueryContext(ctx, buoy, on, off)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	constituents := make([]Constituent, 0)
	for results.Next() {
		var constituent Constituent
		if err := results.Scan(
			&constituent.Name,
			&constituent.Number,
			&constituent.Amplitude,
			&constituent.Lag,
		); err != nil {
			return nil, err
		}
		constituents = append(constituents, constituent)
	}

	if err := results.Err(); err != nil {
		return nil, err
	}

	return constituents, nil
}

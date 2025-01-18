package main

import (
	"context"
	"time"

	"github.com/GeoNet/delta"
)

type Site struct {
	Network    string
	Buoy       string
	Location   string
	Latitude   float64
	Longitude  float64
	Depth      float64
	Correction time.Duration
	Start      time.Time
	End        time.Time
}

type SiteQuery struct {
	Network   string
	Buoy      string
	Location  string
	Latitude  float64
	Longitude float64
	Depth     float64
	Start     int64
	End       int64

	CorrectionTime  *string
	CorrectionStart *time.Time
	CorrectionEnd   *time.Time
}

func (s SiteQuery) Correction() (time.Duration, error) {
	// perhaps a time correction is needed
	if s.CorrectionTime != nil {
		d, err := time.ParseDuration(*s.CorrectionTime)
		if err != nil {
			return 0, err
		}
		return d, nil
	}
	return 0, nil
}

func (s SiteQuery) On() time.Time {
	start := time.Unix(s.Start, 0).UTC()
	if s.CorrectionStart != nil && s.CorrectionStart.After(start) {
		start = *s.CorrectionStart
	}
	return start
}

func (s SiteQuery) Off() time.Time {
	end := time.Unix(s.End, 0).UTC()
	if s.CorrectionEnd != nil && s.CorrectionEnd.Before(end) {
		end = *s.CorrectionEnd
	}
	return end
}

func QuerySites(db *delta.DB, network string) ([]Site, error) {

	ctx := context.Background()

	query := `
	SELECT DISTINCT
		network.external,
		station.station,
		site.location,
                site.latitude,
                site.longitude,
		site.depth,
		MAX(
			unixepoch(site.start_date),
			unixepoch(recorder.start_date),
			unixepoch(stream.start_date)
		),
		MIN(
			unixepoch(site.end_date),
			unixepoch(recorder.end_date),
			unixepoch(stream.end_date)
		),
		timing.correction,
		timing.start_date,
		timing.end_date
	FROM
		site
		INNER JOIN station ON station.station_id = site.station_id
		INNER JOIN gauge ON gauge.gauge = station.station AND
			datetime(site.start_date) <= datetime(gauge.end_date) AND
			datetime(site.end_date) > datetime(gauge.start_date)
		INNER JOIN station_network ON station_network.station_id = station.station_id
		INNER JOIN network ON network.network_id = station_network.network_id
		INNER JOIN stream ON stream.site_id = site.site_id
		INNER JOIN recorder ON recorder.site_id = site.site_id
		LEFT JOIN timing ON site.site_id = timing.site_id AND
			datetime(site.start_date) <= datetime(timing.end_date) AND
			datetime(site.end_date) > datetime(timing.start_date)
	WHERE
		network.network = ?
	ORDER BY
		network.external,
		station.station,
		site.location
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	results, err := stmt.QueryContext(ctx, network)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	sites := make([]Site, 0)
	for results.Next() {

		var res SiteQuery
		if err := results.Scan(
			&res.Network,
			&res.Buoy,
			&res.Location,
			&res.Latitude,
			&res.Longitude,
			&res.Depth,
			&res.Start,
			&res.End,
			&res.CorrectionTime,
			&res.CorrectionStart,
			&res.CorrectionEnd,
		); err != nil {
			return nil, err
		}

		correction, err := res.Correction()
		if err != nil {
			return nil, err
		}

		sites = append(sites, Site{
			Network:    res.Network,
			Buoy:       res.Buoy,
			Location:   res.Location,
			Latitude:   res.Latitude,
			Longitude:  res.Longitude,
			Depth:      res.Depth,
			Correction: correction,
			Start:      res.On(),
			End:        res.Off(),
		})
	}

	if err := results.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

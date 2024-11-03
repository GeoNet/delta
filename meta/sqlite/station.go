package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Stations(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Station, error) {

	query := `SELECT Code,Network,Name,Latitude,Longitude,Elevation,Depth,Datum,Start,End FROM Station`
	if len(opts) > 0 {
		query += " WHERE "
	}
	for n, opt := range opts {
		if n > 0 {
			query += " AND "
		}
		query += opt.K(n)
	}
	query += ";"

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var args []any
	for _, opt := range opts {
		args = append(args, opt.V())
	}
	results, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	stations := make([]meta.Station, 0)
	for results.Next() {
		var station meta.Station
		if err := results.Scan(&station.Code, &station.Network, &station.Name, &station.Latitude, &station.Longitude, &station.Elevation, &station.Depth, &station.Datum, &station.Start, &station.End); err != nil {
			return nil, err
		}
		stations = append(stations, station)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}

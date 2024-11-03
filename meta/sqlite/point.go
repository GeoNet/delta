package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Points(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Point, error) {

	query := "SELECT Sample,Location,Latitude,Longitude,Elevation,Depth,Datum,Survey,Start,End FROM Point"
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

	points := make([]meta.Point, 0)
	for results.Next() {
		var point meta.Point
		if err := results.Scan(&point.Sample, &point.Location, &point.Latitude, &point.Longitude, &point.Elevation, &point.Depth, &point.Datum, &point.Survey, &point.Start, &point.End); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

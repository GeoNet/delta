package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Samples(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Sample, error) {

	query := `SELECT Code,Network,Name,Latitude,Longitude,Elevation,Depth,Datum,Start,End FROM Sample`
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

	samples := make([]meta.Sample, 0)
	for results.Next() {
		var sample meta.Sample
		if err := results.Scan(&sample.Code, &sample.Network, &sample.Name, &sample.Latitude, &sample.Longitude, &sample.Elevation, &sample.Depth, &sample.Datum, &sample.Start, &sample.End); err != nil {
			return nil, err
		}
		samples = append(samples, sample)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return samples, nil
}

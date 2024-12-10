package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Sites(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Site, error) {

	query := "SELECT Station,Location,Latitude,Longitude,Elevation,Depth,Datum,Survey,Start,End FROM Site"
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

	sites := make([]meta.Site, 0)
	for results.Next() {
		var site meta.Site
		if err := results.Scan(&site.Station, &site.Location, &site.Latitude, &site.Longitude, &site.Elevation, &site.Depth, &site.Datum, &site.Survey, &site.Start, &site.End); err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

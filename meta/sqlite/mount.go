package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Mounts(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Mount, error) {

	query := `SELECT Code,Network,Name,Latitude,Longitude,Elevation,Datum,Description,Start,End FROM Mount`
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

	mounts := make([]meta.Mount, 0)
	for results.Next() {
		var mount meta.Mount
		if err := results.Scan(&mount.Code, &mount.Network, &mount.Name, &mount.Latitude, &mount.Longitude, &mount.Elevation, &mount.Datum, &mount.Description, &mount.Start, &mount.End); err != nil {
			return nil, err
		}
		mounts = append(mounts, mount)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return mounts, nil
}

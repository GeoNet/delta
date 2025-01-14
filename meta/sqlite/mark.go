package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Marks(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Mark, error) {

	query := `SELECT Code,Network,Igs,Name,Latitude,Longitude,Elevation,Datum,Start,End FROM Mark`
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

	marks := make([]meta.Mark, 0)
	for results.Next() {
		var mark meta.Mark
		var igs string
		if err := results.Scan(&mark.Code, &mark.Network, &igs, &mark.Name, &mark.Latitude, &mark.Longitude, &mark.Elevation, &mark.Datum, &mark.Start, &mark.End); err != nil {
			return nil, err
		}
		if b, ok := ParseBool(igs); ok {
			mark.Igs = b
		}
		marks = append(marks, mark)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return marks, nil
}

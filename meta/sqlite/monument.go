package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Monuments(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Monument, error) {

	query := `SELECT Mark,DomesNumber,MarkType,Type,GroundRelationship,FoundationType,FoundationDepth,Start,End,Bedrock,Geology FROM Monument`
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

	monuments := make([]meta.Monument, 0)
	for results.Next() {
		var monument meta.Monument
		if err := results.Scan(&monument.Mark, &monument.DomesNumber, &monument.MarkType, &monument.Type, &monument.GroundRelationship, &monument.FoundationType, &monument.FoundationDepth, &monument.Start, &monument.End, &monument.Bedrock, &monument.Geology); err != nil {
			return nil, err
		}

		monuments = append(monuments, monument)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return monuments, nil
}

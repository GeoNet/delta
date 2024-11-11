package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Networks(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.Network, error) {

	query := "SELECT Code,External,Description,Restricted FROM Network"
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

	networks := make([]meta.Network, 0)
	for results.Next() {
		var network meta.Network
		if err := results.Scan(&network.Code, &network.External, &network.Description, &network.Restricted); err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return networks, nil
}

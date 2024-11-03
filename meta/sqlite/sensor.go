package sqlite

import (
	"context"
	"database/sql"

	"github.com/GeoNet/delta/meta"
)

func Sensors(ctx context.Context, db *sql.DB, opts ...QueryOpt) ([]meta.InstalledSensor, error) {

	query := "SELECT Make,Model,Serial,Station,Location,Azimuth,Method,Dip,Depth,North,East,Factor,Bias,Start,End FROM Sensor"
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

	sensors := make([]meta.InstalledSensor, 0)
	for results.Next() {
		var depth float64
		var sensor meta.InstalledSensor
		if err := results.Scan(&sensor.Make, &sensor.Model, &sensor.Serial, &sensor.Station, &sensor.Location, &sensor.Azimuth, &sensor.Method, &sensor.Dip, &depth, &sensor.North, &sensor.East, &sensor.Factor, &sensor.Bias, &sensor.Start, &sensor.End); err != nil {
			return nil, err
		}
		sensor.Vertical = -depth
		sensors = append(sensors, sensor)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

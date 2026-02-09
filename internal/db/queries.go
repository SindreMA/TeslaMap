package db

import (
	"context"
	"database/sql"

	"github.com/sindrema/teslamap/internal/model"
)

func ListCars(ctx context.Context, db *sql.DB) ([]model.Car, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, model, vin, trim_badging, exterior_color
		FROM cars ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []model.Car
	for rows.Next() {
		var c model.Car
		if err := rows.Scan(&c.ID, &c.Name, &c.Model, &c.VIN, &c.TrimBadging, &c.ExteriorColor); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	return cars, rows.Err()
}

func GetCar(ctx context.Context, db *sql.DB, carID int) (*model.Car, error) {
	var c model.Car
	err := db.QueryRowContext(ctx, `
		SELECT id, name, model, vin, trim_badging, exterior_color
		FROM cars WHERE id = $1
	`, carID).Scan(&c.ID, &c.Name, &c.Model, &c.VIN, &c.TrimBadging, &c.ExteriorColor)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetLatestPosition(ctx context.Context, db *sql.DB, carID int) (*model.Position, error) {
	var p model.Position
	err := db.QueryRowContext(ctx, `
		WITH last2 AS (
			SELECT date, latitude, longitude, speed, battery_level
			FROM positions
			WHERE car_id = $1
			ORDER BY date DESC
			LIMIT 2
		)
		SELECT
			cur.date, cur.latitude, cur.longitude, cur.speed, cur.battery_level,
			CASE WHEN prev.latitude IS NOT NULL
				THEN degrees(atan2(
					sin(radians(cur.longitude - prev.longitude)) * cos(radians(cur.latitude)),
					cos(radians(prev.latitude)) * sin(radians(cur.latitude)) -
					sin(radians(prev.latitude)) * cos(radians(cur.latitude)) * cos(radians(cur.longitude - prev.longitude))
				))
			END AS heading
		FROM (SELECT * FROM last2 ORDER BY date DESC LIMIT 1) cur
		LEFT JOIN (SELECT * FROM last2 ORDER BY date ASC LIMIT 1) prev ON true
	`, carID).Scan(&p.Date, &p.Latitude, &p.Longitude, &p.Speed, &p.BatteryLevel, &p.Heading)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

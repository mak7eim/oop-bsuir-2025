package sqlite

import (
	"context"
	"database/sql"
)

func (d *DB) SeedDemoData(ctx context.Context) error {
	var count int
	if err := d.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM restaurants").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	restaurants := []struct {
		name, desc, addr string
	}{
		{"Pizza Palace", "Italian pizza and pasta", "10 Main St"},
		{"Sushi Express", "Fresh sushi rolls", "22 Ocean Ave"},
	}
	restaurantIDs := make([]int64, 0, len(restaurants))
	for _, r := range restaurants {
		res, err := tx.ExecContext(ctx,
			"INSERT INTO restaurants (name, description, address) VALUES (?, ?, ?)",
			r.name, r.desc, r.addr)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		restaurantIDs = append(restaurantIDs, id)
	}

	dishes := []struct {
		restIdx              int
		name, desc           string
		price                float64
	}{
		{0, "Margherita", "Classic tomato and mozzarella", 8.99},
		{0, "Pepperoni", "Spicy pepperoni pizza", 10.99},
		{0, "Carbonara", "Creamy pasta with bacon", 9.49},
		{1, "California Roll", "Crab, avocado, cucumber", 7.99},
		{1, "Salmon Nigiri", "Fresh salmon over rice", 6.49},
		{1, "Dragon Roll", "Eel and avocado specialty", 12.99},
	}
	for _, dish := range dishes {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO dishes (restaurant_id, name, description, price, available)
			 VALUES (?, ?, ?, ?, 1)`,
			restaurantIDs[dish.restIdx], dish.name, dish.desc, dish.price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) EnsureAdmin(ctx context.Context, email, passwordHash string) error {
	var id int64
	err := d.conn.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&id)
	if err == sql.ErrNoRows {
		_, err = d.conn.ExecContext(ctx,
			"INSERT INTO users (email, password_hash, role) VALUES (?, ?, 'admin')",
			email, passwordHash)
	}
	return err
}

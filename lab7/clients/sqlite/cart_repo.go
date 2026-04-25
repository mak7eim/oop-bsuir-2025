package sqlite

import (
	"context"
)

type CartRepository struct {
	db *DB
}

func NewCartRepository(db *DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetItems(ctx context.Context, userID int64) (map[int64]int, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		"SELECT dish_id, quantity FROM cart_items WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[int64]int)
	for rows.Next() {
		var dishID int64
		var qty int
		if err := rows.Scan(&dishID, &qty); err != nil {
			return nil, err
		}
		items[dishID] = qty
	}
	return items, rows.Err()
}

func (r *CartRepository) SetItem(ctx context.Context, userID, dishID int64, quantity int) error {
	if quantity <= 0 {
		_, err := r.db.conn.ExecContext(ctx,
			"DELETE FROM cart_items WHERE user_id = ? AND dish_id = ?", userID, dishID)
		return err
	}
	_, err := r.db.conn.ExecContext(ctx,
		`INSERT INTO cart_items (user_id, dish_id, quantity) VALUES (?, ?, ?)
		 ON CONFLICT(user_id, dish_id) DO UPDATE SET quantity = excluded.quantity`,
		userID, dishID, quantity)
	return err
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID, dishID int64) error {
	_, err := r.db.conn.ExecContext(ctx,
		"DELETE FROM cart_items WHERE user_id = ? AND dish_id = ?", userID, dishID)
	return err
}

func (r *CartRepository) Clear(ctx context.Context, userID int64) error {
	_, err := r.db.conn.ExecContext(ctx, "DELETE FROM cart_items WHERE user_id = ?", userID)
	return err
}

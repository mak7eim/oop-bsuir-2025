package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	res, err := tx.ExecContext(ctx,
		`INSERT INTO orders (user_id, restaurant_id, courier_id, status, subtotal, delivery_fee, total, address, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.UserID, order.RestaurantID, order.CourierID, order.Status,
		order.Subtotal, order.DeliveryFee, order.Total, order.Address, now, now)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	orderID, _ := res.LastInsertId()
	order.ID = orderID
	order.CreatedAt = now
	order.UpdatedAt = now

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO order_items (order_id, dish_id, dish_name, quantity, price) VALUES (?, ?, ?, ?, ?)",
			orderID, item.DishID, item.DishName, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	row := r.db.conn.QueryRowContext(ctx,
		`SELECT id, user_id, restaurant_id, courier_id, status, subtotal, delivery_fee, total, address, created_at, updated_at
		 FROM orders WHERE id = ?`, id)
	order, err := scanOrder(row)
	if err != nil {
		return nil, err
	}
	items, err := r.loadItems(ctx, id)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID int64) ([]models.Order, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		`SELECT id, user_id, restaurant_id, courier_id, status, subtotal, delivery_fee, total, address, created_at, updated_at
		 FROM orders WHERE user_id = ? ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) ListAll(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		`SELECT id, user_id, restaurant_id, courier_id, status, subtotal, delivery_fee, total, address, created_at, updated_at
		 FROM orders ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	order.UpdatedAt = time.Now()
	_, err := r.db.conn.ExecContext(ctx,
		`UPDATE orders SET courier_id = ?, status = ?, updated_at = ? WHERE id = ?`,
		order.CourierID, order.Status, order.UpdatedAt, order.ID)
	return err
}

func (r *OrderRepository) scanOrders(ctx context.Context, rows *sql.Rows) ([]models.Order, error) {
	var orders []models.Order
	for rows.Next() {
		order, err := scanOrderFromRows(rows)
		if err != nil {
			return nil, err
		}
		items, err := r.loadItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, *order)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) loadItems(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		"SELECT dish_id, dish_name, quantity, price FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.DishID, &item.DishName, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanOrder(row *sql.Row) (*models.Order, error) {
	var o models.Order
	var courierID sql.NullInt64
	var created, updated string
	err := row.Scan(&o.ID, &o.UserID, &o.RestaurantID, &courierID, &o.Status,
		&o.Subtotal, &o.DeliveryFee, &o.Total, &o.Address, &created, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if courierID.Valid {
		o.CourierID = &courierID.Int64
	}
	o.CreatedAt = parseTime(created)
	o.UpdatedAt = parseTime(updated)
	return &o, nil
}

func scanOrderFromRows(rows *sql.Rows) (*models.Order, error) {
	var o models.Order
	var courierID sql.NullInt64
	var created, updated string
	err := rows.Scan(&o.ID, &o.UserID, &o.RestaurantID, &courierID, &o.Status,
		&o.Subtotal, &o.DeliveryFee, &o.Total, &o.Address, &created, &updated)
	if err != nil {
		return nil, err
	}
	if courierID.Valid {
		o.CourierID = &courierID.Int64
	}
	o.CreatedAt = parseTime(created)
	o.UpdatedAt = parseTime(updated)
	return &o, nil
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	if t.IsZero() {
		t, _ = time.Parse(time.RFC3339, s)
	}
	return t
}

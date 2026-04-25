package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

type CourierRepository struct {
	db *DB
}

func NewCourierRepository(db *DB) *CourierRepository {
	return &CourierRepository{db: db}
}

func (r *CourierRepository) List(ctx context.Context) ([]models.Courier, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		"SELECT id, user_id, name, phone, status FROM couriers ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Courier
	for rows.Next() {
		c, err := scanCourierFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *c)
	}
	return list, rows.Err()
}

func (r *CourierRepository) GetByID(ctx context.Context, id int64) (*models.Courier, error) {
	row := r.db.conn.QueryRowContext(ctx,
		"SELECT id, user_id, name, phone, status FROM couriers WHERE id = ?", id)
	return scanCourier(row)
}

func (r *CourierRepository) FindFree(ctx context.Context) (*models.Courier, error) {
	row := r.db.conn.QueryRowContext(ctx,
		"SELECT id, user_id, name, phone, status FROM couriers WHERE status = 'free' ORDER BY id LIMIT 1")
	c, err := scanCourier(row)
	if errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.ErrCourierBusy
	}
	return c, err
}

func (r *CourierRepository) UpdateStatus(ctx context.Context, id int64, status models.CourierStatus) error {
	res, err := r.db.conn.ExecContext(ctx, "UPDATE couriers SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *models.Courier) error {
	res, err := r.db.conn.ExecContext(ctx,
		"INSERT INTO couriers (user_id, name, phone, status) VALUES (?, ?, ?, ?)",
		courier.UserID, courier.Name, courier.Phone, courier.Status)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	courier.ID = id
	return nil
}

func scanCourier(row *sql.Row) (*models.Courier, error) {
	var c models.Courier
	err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return &c, err
}

func scanCourierFromRows(rows *sql.Rows) (*models.Courier, error) {
	var c models.Courier
	err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Phone, &c.Status)
	return &c, err
}

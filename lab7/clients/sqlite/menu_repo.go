package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

type MenuRepository struct {
	db *DB
}

func NewMenuRepository(db *DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		"SELECT id, name, description, address FROM restaurants ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Restaurant
	for rows.Next() {
		var rest models.Restaurant
		if err := rows.Scan(&rest.ID, &rest.Name, &rest.Description, &rest.Address); err != nil {
			return nil, err
		}
		list = append(list, rest)
	}
	return list, rows.Err()
}

func (r *MenuRepository) GetRestaurant(ctx context.Context, id int64) (*models.Restaurant, error) {
	row := r.db.conn.QueryRowContext(ctx,
		"SELECT id, name, description, address FROM restaurants WHERE id = ?", id)
	var rest models.Restaurant
	err := row.Scan(&rest.ID, &rest.Name, &rest.Description, &rest.Address)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	return &rest, err
}

func (r *MenuRepository) ListDishesByRestaurant(ctx context.Context, restaurantID int64) ([]models.Dish, error) {
	rows, err := r.db.conn.QueryContext(ctx,
		`SELECT id, restaurant_id, name, description, price, available
		 FROM dishes WHERE restaurant_id = ? AND available = 1 ORDER BY id`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Dish
	for rows.Next() {
		var d models.Dish
		var avail int
		if err := rows.Scan(&d.ID, &d.RestaurantID, &d.Name, &d.Description, &d.Price, &avail); err != nil {
			return nil, err
		}
		d.Available = avail == 1
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *MenuRepository) GetDish(ctx context.Context, id int64) (*models.Dish, error) {
	row := r.db.conn.QueryRowContext(ctx,
		`SELECT id, restaurant_id, name, description, price, available FROM dishes WHERE id = ?`, id)
	var d models.Dish
	var avail int
	err := row.Scan(&d.ID, &d.RestaurantID, &d.Name, &d.Description, &d.Price, &avail)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	d.Available = avail == 1
	return &d, err
}

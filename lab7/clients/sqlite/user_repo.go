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

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	res, err := r.db.conn.ExecContext(ctx,
		"INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)",
		user.Email, user.PasswordHash, user.Role)
	if err != nil {
		if isUniqueViolation(err) {
			return apperrors.ErrConflict
		}
		return fmt.Errorf("insert user: %w", err)
	}
	id, _ := res.LastInsertId()
	user.ID = id
	user.CreatedAt = time.Now()
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.conn.QueryRowContext(ctx,
		"SELECT id, email, password_hash, role, created_at FROM users WHERE email = ?", email)
	return scanUser(row)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.conn.QueryRowContext(ctx,
		"SELECT id, email, password_hash, role, created_at FROM users WHERE id = ?", id)
	return scanUser(row)
}

func scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	var created string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	if u.CreatedAt.IsZero() {
		u.CreatedAt, _ = time.Parse(time.RFC3339, created)
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "UNIQUE") || contains(err.Error(), "unique"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

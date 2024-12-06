package repository

import (
	"context"
	"database/sql"
	"errors"
	"user_service/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId string) (*models.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error)
	ListUsers(ctx context.Context) ([]*models.UserRecord, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserById(ctx context.Context, userId string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE id = $1`

	user := &models.UserRecord{}

	err := r.db.QueryRowContext(ctx, query,
		userId,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Verified,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE email = $1`

	user := &models.UserRecord{}

	err := r.db.QueryRowContext(ctx, query,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Verified,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.UserRecord
	for rows.Next() {
		user := &models.UserRecord{}
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Verified,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

package repository

import (
	"auth_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.UserRecord) (*models.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error)
	UpdateUserVerification(ctx context.Context, email string, verified bool) error

	Close()
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Close() {
	r.db.Close()
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.UserRecord) (*models.UserRecord, error) {
	query := `INSERT INTO users (id, email, username, password, verified, role)
	          VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5) 
	          RETURNING id, email, username, password, verified, role, created_at, updated_at`

	var newUser models.UserRecord

	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.Verified,
		user.Role,
	).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Username,
		&newUser.Password,
		&newUser.Verified,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, role, verified, created_at, updated_at FROM users WHERE email = $1`

	var user models.UserRecord

	err := r.db.QueryRowContext(ctx, query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) UpdateUserVerification(ctx context.Context, email string, verified bool) error {
	query := `UPDATE users SET verified = $1, updated_at = $2 WHERE email = $3`
	_, err := r.db.ExecContext(ctx, query, verified, time.Now(), email)
	return err
}

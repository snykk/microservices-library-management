package repository

import (
	"auth_service/internal/models"
	"database/sql"
	"errors"
	"time"
)

type AuthRepository interface {
	CreateUser(user *models.UserRecord) (*models.UserRecord, error)
	GetUserByEmail(email *string) (*models.UserRecord, error)
	UpdateUserVerification(email *string, verified *bool) error

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

func (r *authRepository) CreateUser(user *models.UserRecord) (*models.UserRecord, error) {
	query := `INSERT INTO users (id, email, username, password, verified, role)
	          VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5) 
	          RETURNING id, email, username, password, verified, role, created_at, updated_at`

	var newUser models.UserRecord
	err := r.db.QueryRow(query,
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

func (r *authRepository) GetUserByEmail(email *string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, role, verified, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user models.UserRecord
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUserVerification(email *string, verified *bool) error {
	query := `UPDATE users SET verified = $1, updated_at = $2 WHERE email = $3`
	_, err := r.db.Exec(query, *verified, time.Now(), *email)
	return err
}

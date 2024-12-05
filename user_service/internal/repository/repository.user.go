package repository

import (
	"database/sql"
	"errors"
	"user_service/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUserById(userId *string) (*models.UserRecord, error)
	GetUserByEmail(email *string) (*models.UserRecord, error)
	ListUsers() ([]*models.UserRecord, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserById(userId *string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, *userId)

	user := &models.UserRecord{}
	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Verified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email *string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, *email)

	user := &models.UserRecord{}
	if err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Verified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) ListUsers() ([]*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.UserRecord
	for rows.Next() {
		user := &models.UserRecord{}
		if err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Verified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

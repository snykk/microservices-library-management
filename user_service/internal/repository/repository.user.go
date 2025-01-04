package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"user_service/internal/models"
	"user_service/pkg/utils"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId string) (*models.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error)
	ListUsers(ctx context.Context) ([]*models.UserRecord, error)
	CountUsers(ctx context.Context) (int, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetUserById returns a user from the database by ID
func (r *userRepository) GetUserById(ctx context.Context, userId string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE id = $1`
	log.Printf("[%s] Executing query: %s with userId: %s\n", utils.GetLocation(), query, userId)

	user := &models.UserRecord{}
	err := r.db.QueryRowContext(ctx, query, userId).Scan(
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
			log.Printf("[%s] User with ID %s not found\n", utils.GetLocation(), userId)
			return nil, errors.New("user not found")
		}
		log.Printf("[%s] Error executing query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched user with ID: %s\n", utils.GetLocation(), userId)
	return user, nil
}

// GetUserByEmail returns a user from the database by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users WHERE email = $1`
	log.Printf("[%s] Executing query: %s with email: %s\n", utils.GetLocation(), query, email)

	user := &models.UserRecord{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
			log.Printf("[%s] User with email %s not found\n", utils.GetLocation(), email)
			return nil, errors.New("user not found")
		}
		log.Printf("[%s] Error executing query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched user with email: %s\n", utils.GetLocation(), email)
	return user, nil
}

// ListUsers returns a list of users from the database
func (r *userRepository) ListUsers(ctx context.Context) ([]*models.UserRecord, error) {
	query := `SELECT id, email, username, password, verified, role, created_at, updated_at FROM users`
	log.Printf("[%s] Executing query: %s\n", utils.GetLocation(), query)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[%s] Error executing query: %v\n", utils.GetLocation(), err)
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
			log.Printf("[%s] Error scanning row: %v\n", utils.GetLocation(), err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[%s] Row iteration error: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched %d users\n", utils.GetLocation(), len(users))
	return users, nil
}

// CountUsers returns the total number of users in the database
func (r *userRepository) CountUsers(ctx context.Context) (int, error) {
	log.Printf("Counting total users")
	query := `SELECT COUNT(*) FROM users`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting users: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

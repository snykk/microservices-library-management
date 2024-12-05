package repository

import (
	"context"
	"database/sql"
	"errors"
	"loan_service/internal/models"

	"github.com/jmoiron/sqlx"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *models.LoanRecord) (*models.LoanRecord, error)
	GetLoan(ctx context.Context, id string) (*models.LoanRecord, error)
	UpdateLoanStatus(ctx context.Context, loan *models.LoanRecord) (*models.LoanRecord, error)
	ListLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error)
}

type loanRepository struct {
	db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) LoanRepository {
	return &loanRepository{db: db}
}

// CreateLoan inserts a new loan record into the database
func (r *loanRepository) CreateLoan(ctx context.Context, req *models.LoanRecord) (*models.LoanRecord, error) {
	query := `
		INSERT INTO loans (user_id, book_id, loan_date, status)
		VALUES (:user_id, :book_id, :loan_date, :status)
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	loan := &models.LoanRecord{}
	rows, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(loan); err != nil {
			return nil, err
		}
	}
	return loan, nil
}

// GetLoan fetches a loan record by ID
func (r *loanRepository) GetLoan(ctx context.Context, id string) (*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE id = $1
	`

	loan := &models.LoanRecord{}
	var returnDate sql.NullTime
	if err := r.db.GetContext(ctx, loan, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}

	if returnDate.Valid {
		loan.ReturnDate = &returnDate.Time
	}
	return loan, nil
}

// UpdateLoanStatus updates the status and optionally the return date of a loan
func (r *loanRepository) UpdateLoanStatus(ctx context.Context, req *models.LoanRecord) (*models.LoanRecord, error) {
	query := `
		UPDATE loans
		SET status = :status, return_date = :return_date
		WHERE id = :id
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	loan := &models.LoanRecord{}
	rows, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(loan); err != nil {
			return nil, err
		}
	}
	return loan, nil
}

// ListLoans fetches all loans associated with a specific user
func (r *loanRepository) ListLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE user_id = $1
	`

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, userId); err != nil {
		return nil, err
	}
	return loans, nil
}

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
	GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, error)
	UpdateLoanStatus(ctx context.Context, loan *models.LoanRecord) (*models.LoanRecord, error)
	ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error)
	ListLoans(ctx context.Context) ([]*models.LoanRecord, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, error)
	GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, error)
	ReturnLoan(ctx context.Context, id string) (*models.LoanRecord, error)
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

func (r *loanRepository) GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE book_id = $1 AND user_id = $2 AND status = 'BORROWED'
	`

	loan := &models.LoanRecord{}
	var returnDate sql.NullTime
	if err := r.db.GetContext(ctx, loan, query, bookId, userId); err != nil {
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

// updates the status
func (r *loanRepository) UpdateLoanStatus(ctx context.Context, req *models.LoanRecord) (*models.LoanRecord, error) {
	query := `
		UPDATE loans
		SET status = :status
		WHERE id = :id
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	loan := &models.LoanRecord{}
	err := r.db.GetContext(ctx, loan, query, req)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// ListUserLoans fetches all loans associated with a specific user
func (r *loanRepository) ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error) {
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

func (r *loanRepository) ListLoans(ctx context.Context) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
	`

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepository) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE user_id = $1 AND status = $2  
	`

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, userId, status); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepository) GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE status = $1
	`

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, status); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *loanRepository) ReturnLoan(ctx context.Context, id string) (*models.LoanRecord, error) {
	query := `
		UPDATE loans
		SET status = 'RETURNED', return_date = NOW()
		WHERE id = $1
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	loan := &models.LoanRecord{}
	err := r.db.GetContext(ctx, loan, query, id)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

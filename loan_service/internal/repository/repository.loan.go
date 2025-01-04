package repository

import (
	"context"
	"database/sql"
	"errors"
	"loan_service/internal/models"
	"loan_service/pkg/utils"
	"log"

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
	CountLoans(ctx context.Context) (int, error)
	CountLoansByUserId(ctx context.Context, userId string) (int, error)
	CountLoansByStatus(ctx context.Context, status string) (int, error)
	CountLoansByUserIdAndStatus(ctx context.Context, userId string, status string) (int, error)
}

type loanRepository struct {
	db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) CreateLoan(ctx context.Context, req *models.LoanRecord) (*models.LoanRecord, error) {
	query := `
		INSERT INTO loans (user_id, book_id, loan_date, status)
		VALUES (:user_id, :book_id, :loan_date, :status)
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	log.Printf("[%s] Executing query to create loan: %s with parameters: %+v\n", utils.GetLocation(), query, req)

	loan := &models.LoanRecord{}
	rows, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		log.Printf("[%s] Error executing CreateLoan query: %v\n", utils.GetLocation(), err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(loan); err != nil {
			log.Printf("[%s] Error scanning result of CreateLoan query: %v\n", utils.GetLocation(), err)
			return nil, err
		}
	}
	log.Printf("[%s] Loan created successfully: %+v\n", utils.GetLocation(), loan)
	return loan, nil
}

func (r *loanRepository) GetLoan(ctx context.Context, id string) (*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE id = $1
	`

	log.Printf("[%s] Executing query to get loan with ID: %s\n", utils.GetLocation(), id)

	loan := &models.LoanRecord{}
	var returnDate sql.NullTime
	if err := r.db.GetContext(ctx, loan, query, id); err != nil {
		log.Printf("[%s] Error executing GetLoan query: %v\n", utils.GetLocation(), err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}

	if returnDate.Valid {
		loan.ReturnDate = &returnDate.Time
	}
	log.Printf("[%s] Loan retrieved successfully: %+v\n", utils.GetLocation(), loan)
	return loan, nil
}

func (r *loanRepository) GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE book_id = $1 AND user_id = $2 AND status = 'BORROWED'
	`

	log.Printf("[%s] Executing query to get borrowed loan for book ID: %s and user ID: %s\n", utils.GetLocation(), bookId, userId)

	loan := &models.LoanRecord{}
	var returnDate sql.NullTime
	if err := r.db.GetContext(ctx, loan, query, bookId, userId); err != nil {
		log.Printf("[%s] Error executing GetBorrowedLoanByBookIdAndUserId query: %v\n", utils.GetLocation(), err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("loan not found")
		}
		return nil, err
	}

	if returnDate.Valid {
		loan.ReturnDate = &returnDate.Time
	}
	log.Printf("[%s] Borrowed loan retrieved successfully: %+v\n", utils.GetLocation(), loan)
	return loan, nil
}

func (r *loanRepository) UpdateLoanStatus(ctx context.Context, req *models.LoanRecord) (*models.LoanRecord, error) {
	query := `
		UPDATE loans
		SET status = :status
		WHERE id = :id
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	log.Printf("[%s] Executing query to update loan status for loan ID: %s\n", utils.GetLocation(), req.Id)

	loan := &models.LoanRecord{}
	err := r.db.GetContext(ctx, loan, query, req)
	if err != nil {
		log.Printf("[%s] Error executing UpdateLoanStatus query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Loan status updated successfully: %+v\n", utils.GetLocation(), loan)
	return loan, nil
}

func (r *loanRepository) ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE user_id = $1
	`

	log.Printf("[%s] Executing query to list loans for user ID: %s\n", utils.GetLocation(), userId)

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, userId); err != nil {
		log.Printf("[%s] Error executing ListUserLoans query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Found %d loans for user ID: %s\n", utils.GetLocation(), len(loans), userId)
	return loans, nil
}

func (r *loanRepository) ListLoans(ctx context.Context) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
	`

	log.Printf("[%s] Executing query to list all loans\n", utils.GetLocation())

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query); err != nil {
		log.Printf("[%s] Error executing ListLoans query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Found %d loans\n", utils.GetLocation(), len(loans))
	return loans, nil
}

func (r *loanRepository) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE user_id = $1 AND status = $2  
	`

	log.Printf("[%s] Executing query to list loans for user ID: %s with status: %s\n", utils.GetLocation(), userId, status)

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, userId, status); err != nil {
		log.Printf("[%s] Error executing GetUserLoansByStatus query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Found %d loans for user ID: %s with status: %s\n", utils.GetLocation(), len(loans), userId, status)
	return loans, nil
}

func (r *loanRepository) GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, error) {
	query := `
		SELECT id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE status = $1
	`

	log.Printf("[%s] Executing query to list loans with status: %s\n", utils.GetLocation(), status)

	var loans []*models.LoanRecord
	if err := r.db.SelectContext(ctx, &loans, query, status); err != nil {
		log.Printf("[%s] Error executing GetLoansByStatus query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Found %d loans with status: %s\n", utils.GetLocation(), len(loans), status)
	return loans, nil
}

func (r *loanRepository) ReturnLoan(ctx context.Context, id string) (*models.LoanRecord, error) {
	query := `
		UPDATE loans
		SET status = 'RETURNED', return_date = NOW()
		WHERE id = $1
		RETURNING id, user_id, book_id, loan_date, return_date, status, created_at, updated_at
	`

	log.Printf("[%s] Executing query to return loan with ID: %s\n", utils.GetLocation(), id)

	loan := &models.LoanRecord{}
	err := r.db.GetContext(ctx, loan, query, id)
	if err != nil {
		log.Printf("[%s] Error executing ReturnLoan query: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Loan returned successfully: %+v\n", utils.GetLocation(), loan)
	return loan, nil
}

// CountLoans returns the total number of loans
func (r *loanRepository) CountLoans(ctx context.Context) (int, error) {
	log.Printf("Counting total loans")
	query := `SELECT COUNT(*) FROM loans`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting loans: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

// CountsLoansByUserId returns the total number of loans by user ID
func (r *loanRepository) CountLoansByUserId(ctx context.Context, userId string) (int, error) {
	log.Printf("Counting total loans by user ID: %s\n", userId)
	query := `SELECT COUNT(*) FROM loans WHERE user_id = $1`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, userId).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting loans: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

// CountLoansByStatus returns the total number of loans by status
func (r *loanRepository) CountLoansByStatus(ctx context.Context, status string) (int, error) {
	log.Printf("Counting total loans by status: %s\n", status)
	query := `SELECT COUNT(*) FROM loans WHERE status = $1`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, status).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting loans: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

// CountLoansByUserIdAndStatus returns the total number of loans by user ID and status
func (r *loanRepository) CountLoansByUserIdAndStatus(ctx context.Context, userId string, status string) (int, error) {
	log.Printf("Counting total loans by category ID: %s and by status : %s\n", userId, status)
	query := `SELECT COUNT(*) FROM loans WHERE user_id = $1 AND status = $2`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, userId, status).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting loans: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

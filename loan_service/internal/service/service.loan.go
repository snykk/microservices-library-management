package service

import (
	"context"
	"errors"
	"loan_service/internal/models"
	"loan_service/internal/repository"
	"time"
)

type LoanService interface {
	CreateLoan(ctx context.Context, userId, bookId string) (*models.LoanRecord, error)
	GetLoan(ctx context.Context, id string) (*models.LoanRecord, error)
	UpdateLoanStatus(ctx context.Context, id, status string, returnDate time.Time) (*models.LoanRecord, error)
	ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error)
	ListLoans(ctx context.Context) ([]*models.LoanRecord, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, error)
	GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, error)
}

type loanService struct {
	repo repository.LoanRepository
}

func NewLoanService(repo repository.LoanRepository) LoanService {
	return &loanService{repo: repo}
}

func (s *loanService) CreateLoan(ctx context.Context, userId, bookId string) (*models.LoanRecord, error) {
	loan := &models.LoanRecord{
		UserId:   userId,
		BookId:   bookId,
		LoanDate: time.Now(),
		Status:   "BORROWED",
	}
	return s.repo.CreateLoan(ctx, loan)
}

func (s *loanService) GetLoan(ctx context.Context, id string) (*models.LoanRecord, error) {
	return s.repo.GetLoan(ctx, id)
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id, status string, returnDate time.Time) (*models.LoanRecord, error) {
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, errors.New("loan not found")
	}

	if loan.Status == status {
		return nil, errors.New("loan already " + status)
	}

	loan.Status = status
	loan.ReturnDate = &returnDate

	return s.repo.UpdateLoanStatus(ctx, loan)
}

func (s *loanService) ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error) {
	return s.repo.ListUserLoans(ctx, userId)
}

func (s *loanService) ListLoans(ctx context.Context) ([]*models.LoanRecord, error) {
	return s.repo.ListLoans(ctx)
}

func (s *loanService) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, error) {
	return s.repo.GetUserLoansByStatus(ctx, userId, status)
}

func (s *loanService) GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, error) {
	return s.repo.GetLoansByStatus(ctx, status)
}

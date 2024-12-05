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
	ListLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error)
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

	loan.Status = status
	loan.ReturnDate = &returnDate

	return s.repo.UpdateLoanStatus(ctx, loan)
}

func (s *loanService) ListLoans(ctx context.Context, userId string) ([]*models.LoanRecord, error) {
	return s.repo.ListLoans(ctx, userId)
}

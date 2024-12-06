package service

import (
	"context"
	"errors"
	"fmt"
	"loan_service/internal/models"
	"loan_service/internal/repository"
	"time"

	"google.golang.org/grpc/codes"
)

type LoanService interface {
	CreateLoan(ctx context.Context, userId, bookId string) (*models.LoanRecord, codes.Code, error)
	GetLoan(ctx context.Context, id string) (*models.LoanRecord, codes.Code, error)
	UpdateLoanStatus(ctx context.Context, id, userId, role, status string, returnDate time.Time) (*models.LoanRecord, codes.Code, error)
	ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, codes.Code, error)
	ListLoans(ctx context.Context) ([]*models.LoanRecord, codes.Code, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, codes.Code, error)
	GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, codes.Code, error)
}

type loanService struct {
	repo repository.LoanRepository
}

func NewLoanService(repo repository.LoanRepository) LoanService {
	return &loanService{repo: repo}
}

func (s *loanService) CreateLoan(ctx context.Context, userId, bookId string) (*models.LoanRecord, codes.Code, error) {
	loan := &models.LoanRecord{
		UserId:   userId,
		BookId:   bookId,
		LoanDate: time.Now(),
		Status:   "BORROWED",
	}

	createdLoan, err := s.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, codes.Internal, errors.New("failet to create new loan")
	}
	return createdLoan, codes.OK, nil
}

func (s *loanService) GetLoan(ctx context.Context, id string) (*models.LoanRecord, codes.Code, error) {
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get loan with id %s", id)
	}
	return loan, codes.OK, nil
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id, userId, role, status string, returnDate time.Time) (*models.LoanRecord, codes.Code, error) {
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, codes.NotFound, errors.New("loan not found")
	}

	if loan.UserId != userId && role != "admin" {
		return nil, codes.PermissionDenied, errors.New("you dont have access to update this loan")
	}

	if loan.Status == status {
		return nil, codes.Internal, fmt.Errorf("loan already %s", status)
	}

	loan.Status = status
	loan.ReturnDate = &returnDate

	updatedLoan, err := s.repo.UpdateLoanStatus(ctx, loan)
	if err != nil {
		return nil, codes.Internal, errors.New("failed to update loan status")
	}
	return updatedLoan, codes.OK, nil
}

func (s *loanService) ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, codes.Code, error) {
	loans, err := s.repo.ListUserLoans(ctx, userId)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get list loan with user id %s", userId)
	}
	return loans, codes.OK, nil
}

func (s *loanService) ListLoans(ctx context.Context) ([]*models.LoanRecord, codes.Code, error) {
	loans, err := s.repo.ListLoans(ctx)
	if err != nil {
		return nil, codes.Internal, errors.New("failed to get list loan")
	}
	return loans, codes.OK, err
}

func (s *loanService) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, codes.Code, error) {
	loans, err := s.repo.GetUserLoansByStatus(ctx, userId, status)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get list loan of user %s with status %s", userId, status)
	}
	return loans, codes.OK, nil
}

func (s *loanService) GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, codes.Code, error) {
	loans, err := s.repo.GetLoansByStatus(ctx, status)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get list loan with status %s", status)
	}
	return loans, codes.OK, nil
}

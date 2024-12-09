package service

import (
	"context"
	"errors"
	"fmt"
	"loan_service/internal/clients"
	"loan_service/internal/models"
	"loan_service/internal/repository"
	"loan_service/pkg/rabbitmq"
	"time"

	"google.golang.org/grpc/codes"
)

type LoanService interface {
	CreateLoan(ctx context.Context, userId, email, bookId string) (*models.LoanRecord, codes.Code, error)
	ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (*models.LoanRecord, codes.Code, error)
	GetLoan(ctx context.Context, id string) (*models.LoanRecord, codes.Code, error)
	GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, codes.Code, error)
	UpdateLoanStatus(ctx context.Context, id, status string, returnDate time.Time) (*models.LoanRecord, codes.Code, error)
	ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, codes.Code, error)
	ListLoans(ctx context.Context) ([]*models.LoanRecord, codes.Code, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, codes.Code, error)
	GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, codes.Code, error)
}

type loanService struct {
	bookClient clients.BookClient
	repo       repository.LoanRepository
	publisher  *rabbitmq.Publisher
}

func NewLoanService(repo repository.LoanRepository, bookClient clients.BookClient, publsisher *rabbitmq.Publisher) LoanService {
	return &loanService{
		bookClient: bookClient,
		repo:       repo,
		publisher:  publsisher,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, userId, email, bookId string) (*models.LoanRecord, codes.Code, error) {
	userLoan, _ := s.repo.GetBorrowedLoanByBookIdAndUserId(ctx, bookId, userId)
	if userLoan != nil && userLoan.Status == "BORROWED" {
		return nil, codes.Canceled, errors.New("user must return the borrowed book before creating a new loan")
	}

	book, _ := s.bookClient.GetBook(ctx, bookId)
	if book == nil {
		return nil, codes.NotFound, fmt.Errorf("book '%s' not found", bookId)
	}
	if book.Stock == 0 {
		return nil, codes.Unavailable, errors.New("book is not available")
	}

	err := s.bookClient.DecrementBookStock(ctx, book.Id)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed when update stock book '%s'", book.Id)
	}

	loan := &models.LoanRecord{
		UserId:   userId,
		BookId:   book.Id,
		LoanDate: time.Now(),
		Status:   "BORROWED",
	}

	createdLoan, err := s.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, codes.Internal, errors.New("failed to create new loan")
	}

	err = s.publisher.Publish("email_exchange", "loan_notification", models.LoanNotificationMessage{
		Email: email,
		Book:  book.Title,
		Due:   time.Now().AddDate(0, 0, 7),
	})
	if err != nil {
		return nil, codes.Internal, errors.New("failed to publish loan notification to queue")
	}
	return createdLoan, codes.OK, nil
}

func (s *loanService) ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (*models.LoanRecord, codes.Code, error) {
	// check loan
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil { // supposed loan not found
		return nil, codes.NotFound, fmt.Errorf("loan '%s' not found", id)
	}
	if loan.Status == "RETURNED" {
		return nil, codes.Canceled, fmt.Errorf("loan '%s' already returned", loan.Id)
	}
	if loan.UserId != userId {
		return nil, codes.PermissionDenied, errors.New("you don't have access to this resource")
	}

	// check book
	book, err := s.bookClient.GetBook(ctx, loan.BookId)
	if err != nil {
		return nil, codes.Internal, errors.New("failed to get loan book")
	}

	returnedLoan, err := s.repo.ReturnLoan(ctx, id)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to return loan with id %s", id)
	}

	err = s.bookClient.IncrementBookStock(ctx, book.Id)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed when update stock book '%s'", book.Id)
	}

	err = s.publisher.Publish("email_exchange", "return_notification", models.ReturnNotificationMessage{
		Email: email,
		Book:  book.Title,
	})
	if err != nil {
		return nil, codes.Internal, errors.New("failed to publish loan notification to queue")
	}

	return returnedLoan, codes.OK, nil
}

func (s *loanService) GetLoan(ctx context.Context, id string) (*models.LoanRecord, codes.Code, error) {
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get loan with id %s", id)
	}
	return loan, codes.OK, nil
}

func (s *loanService) GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, codes.Code, error) {
	loan, err := s.repo.GetBorrowedLoanByBookIdAndUserId(ctx, bookId, userId)
	if err != nil {
		return nil, codes.Internal, fmt.Errorf("failed to get loan with id %s and userId %s", bookId, userId)
	}
	return loan, codes.OK, nil
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id, status string, returnDate time.Time) (*models.LoanRecord, codes.Code, error) {
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		return nil, codes.NotFound, errors.New("loan not found")
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

package service

import (
	"context"
	"errors"
	"fmt"
	"loan_service/internal/clients"
	"loan_service/internal/constants"
	"loan_service/internal/models"
	"loan_service/internal/repository"
	"loan_service/pkg/rabbitmq"
	"loan_service/pkg/utils"
	"log"
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

func NewLoanService(repo repository.LoanRepository, bookClient clients.BookClient, publisher *rabbitmq.Publisher) LoanService {
	return &loanService{
		bookClient: bookClient,
		repo:       repo,
		publisher:  publisher,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, userId, email, bookId string) (*models.LoanRecord, codes.Code, error) {
	// Get requestID from context
	requestID, ok := ctx.Value(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	log.Printf("[%s] Creating new loan for user %s and book %s\n", utils.GetLocation(), userId, bookId)

	// Check if the user already borrowed this book
	userLoan, _ := s.repo.GetBorrowedLoanByBookIdAndUserId(ctx, bookId, userId)
	if userLoan != nil && userLoan.Status == "BORROWED" {
		log.Printf("[%s] User %s already has a borrowed book with ID %s\n", utils.GetLocation(), userId, bookId)
		return nil, codes.Canceled, errors.New("user must return the borrowed book before creating a new loan")
	}

	// Fetch the book details
	book, _ := s.bookClient.GetBook(ctx, bookId)
	if book == nil {
		log.Printf("[%s] Book with ID %s not found\n", utils.GetLocation(), bookId)
		return nil, codes.NotFound, fmt.Errorf("book '%s' not found", bookId)
	}
	if book.Stock == 0 {
		log.Printf("[%s] Book '%s' is not available\n", utils.GetLocation(), bookId)
		return nil, codes.Unavailable, errors.New("book is not available")
	}

	// Decrease book stock
	err := s.bookClient.DecrementBookStock(ctx, book.Id)
	if err != nil {
		log.Printf("[%s] Failed to decrement stock for book %s: %v\n", utils.GetLocation(), book.Id, err)
		return nil, codes.Internal, fmt.Errorf("failed when updating stock for book '%s'", book.Id)
	}

	// Create the loan record
	loan := &models.LoanRecord{
		UserId:   userId,
		BookId:   book.Id,
		LoanDate: time.Now(),
		Status:   "BORROWED",
	}

	createdLoan, err := s.repo.CreateLoan(ctx, loan)
	if err != nil {
		log.Printf("[%s] Failed to create loan for user %s and book %s: %v\n", utils.GetLocation(), userId, bookId, err)
		return nil, codes.Internal, errors.New("failed to create new loan")
	}

	// Publish loan notification
	err = s.publisher.Publish(constants.EmailExchange, constants.LoanNotificationQueue, models.LoanNotificationMessage{
		RequestID: requestID,
		Email:     email,
		Book:      book.Title,
		Due:       time.Now().AddDate(0, 0, 7),
	})
	if err != nil {
		log.Printf("[%s] Failed to publish loan notification for user %s: %v\n", utils.GetLocation(), userId, err)
		return nil, codes.Internal, errors.New("failed to publish loan notification to queue")
	}

	log.Printf("[%s] Loan created successfully for user %s and book %s\n", utils.GetLocation(), userId, bookId)
	return createdLoan, codes.OK, nil
}

func (s *loanService) ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (*models.LoanRecord, codes.Code, error) {
	// Get requestID from context
	requestID, ok := ctx.Value(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	log.Printf("[%s] Returning loan with ID %s for user %s\n", utils.GetLocation(), id, userId)

	// Fetch the loan record
	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("[%s] Loan with ID %s not found for user %s: %v\n", utils.GetLocation(), id, userId, err)
		return nil, codes.NotFound, fmt.Errorf("loan '%s' not found", id)
	}
	if loan.Status == "RETURNED" {
		log.Printf("[%s] Loan '%s' already returned\n", utils.GetLocation(), id)
		return nil, codes.Canceled, fmt.Errorf("loan '%s' already returned", loan.Id)
	}
	if loan.UserId != userId {
		log.Printf("[%s] User %s does not have access to loan '%s'\n", utils.GetLocation(), userId, id)
		return nil, codes.PermissionDenied, errors.New("you don't have access to this resource")
	}

	// Fetch book details
	book, err := s.bookClient.GetBook(ctx, loan.BookId)
	if err != nil {
		log.Printf("[%s] Failed to get book details for loan '%s': %v\n", utils.GetLocation(), id, err)
		return nil, codes.Internal, errors.New("failed to get loan book")
	}

	// Return the loan
	returnedLoan, err := s.repo.ReturnLoan(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to return loan with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, codes.Internal, fmt.Errorf("failed to return loan with id %s", id)
	}

	// Increment book stock
	err = s.bookClient.IncrementBookStock(ctx, book.Id)
	if err != nil {
		log.Printf("[%s] Failed to increment stock for book %s: %v\n", utils.GetLocation(), book.Id, err)
		return nil, codes.Internal, fmt.Errorf("failed when updating stock for book '%s'", book.Id)
	}

	// Publish return notification
	err = s.publisher.Publish(constants.EmailExchange, constants.ReturnNotificationQueue, models.ReturnNotificationMessage{
		RequestID: requestID,
		Email:     email,
		Book:      book.Title,
	})
	if err != nil {
		log.Printf("[%s] Failed to publish return notification for user %s: %v\n", utils.GetLocation(), userId, err)
		return nil, codes.Internal, errors.New("failed to publish loan notification to queue")
	}

	log.Printf("[%s] Loan with ID %s successfully returned by user %s\n", utils.GetLocation(), id, userId)
	return returnedLoan, codes.OK, nil
}

func (s *loanService) GetLoan(ctx context.Context, id string) (*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching loan with ID: %s\n", utils.GetLocation(), id)

	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to get loan with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, codes.Internal, fmt.Errorf("failed to get loan with id %s", id)
	}

	log.Printf("[%s] Loan with ID %s fetched successfully\n", utils.GetLocation(), id)
	return loan, codes.OK, nil
}

func (s *loanService) GetBorrowedLoanByBookIdAndUserId(ctx context.Context, bookId, userId string) (*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching borrowed loan for book %s by user %s\n", utils.GetLocation(), bookId, userId)

	loan, err := s.repo.GetBorrowedLoanByBookIdAndUserId(ctx, bookId, userId)
	if err != nil {
		log.Printf("[%s] Failed to get borrowed loan for book %s by user %s: %v\n", utils.GetLocation(), bookId, userId, err)
		return nil, codes.Internal, fmt.Errorf("failed to get borrowed loan with bookId %s and userId %s", bookId, userId)
	}

	log.Printf("[%s] Borrowed loan for book %s by user %s fetched successfully\n", utils.GetLocation(), bookId, userId)
	return loan, codes.OK, nil
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id, status string, returnDate time.Time) (*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Updating status of loan with ID %s to %s\n", utils.GetLocation(), id, status)

	loan, err := s.repo.GetLoan(ctx, id)
	if err != nil {
		log.Printf("[%s] Loan with ID %s not found: %v\n", utils.GetLocation(), id, err)
		return nil, codes.NotFound, fmt.Errorf("loan '%s' not found", id)
	}

	loan.Status = status
	loan.ReturnDate = &returnDate

	updatedLoan, err := s.repo.UpdateLoanStatus(ctx, loan)
	if err != nil {
		log.Printf("[%s] Failed to update loan status for loan with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, codes.Internal, errors.New("failed to update loan status")
	}

	log.Printf("[%s] Loan with ID %s status updated to %s\n", utils.GetLocation(), id, status)
	return updatedLoan, codes.OK, nil
}

func (s *loanService) ListUserLoans(ctx context.Context, userId string) ([]*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching all loans for user %s\n", utils.GetLocation(), userId)

	loans, err := s.repo.ListUserLoans(ctx, userId)
	if err != nil {
		log.Printf("[%s] Failed to fetch loans for user %s: %v\n", utils.GetLocation(), userId, err)
		return nil, codes.Internal, errors.New("failed to fetch loans")
	}

	log.Printf("[%s] Found %d loans for user %s\n", utils.GetLocation(), len(loans), userId)
	return loans, codes.OK, nil
}

func (s *loanService) ListLoans(ctx context.Context) ([]*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching all loans\n", utils.GetLocation())

	loans, err := s.repo.ListLoans(ctx)
	if err != nil {
		log.Printf("[%s] Failed to fetch all loans: %v\n", utils.GetLocation(), err)
		return nil, codes.Internal, errors.New("failed to fetch all loans")
	}

	log.Printf("[%s] Found %d loans\n", utils.GetLocation(), len(loans))
	return loans, codes.OK, nil
}

func (s *loanService) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching loans for user %s with status %s\n", utils.GetLocation(), userId, status)

	loans, err := s.repo.GetUserLoansByStatus(ctx, userId, status)
	if err != nil {
		log.Printf("[%s] Failed to fetch loans for user %s with status %s: %v\n", utils.GetLocation(), userId, status, err)
		return nil, codes.Internal, errors.New("failed to fetch loans")
	}

	log.Printf("[%s] Found %d loans for user %s with status %s\n", utils.GetLocation(), len(loans), userId, status)
	return loans, codes.OK, nil
}

func (s *loanService) GetLoansByStatus(ctx context.Context, status string) ([]*models.LoanRecord, codes.Code, error) {
	log.Printf("[%s] Fetching loans with status %s\n", utils.GetLocation(), status)

	loans, err := s.repo.GetLoansByStatus(ctx, status)
	if err != nil {
		log.Printf("[%s] Failed to fetch loans with status %s: %v\n", utils.GetLocation(), status, err)
		return nil, codes.Internal, errors.New("failed to fetch loans")
	}

	log.Printf("[%s] Found %d loans with status %s\n", utils.GetLocation(), len(loans), status)
	return loans, codes.OK, nil
}

package grpc_server

import (
	"context"
	"fmt"
	"loan_service/internal/clients"
	"loan_service/internal/service"
	protoLoan "loan_service/proto/loan_service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loanGRPCServer struct {
	loanService service.LoanService
	bookClient  clients.BookClient
	protoLoan.UnimplementedLoanServiceServer
}

func NewLoanGRPCServer(loanService service.LoanService, bookClient clients.BookClient) protoLoan.LoanServiceServer {
	return &loanGRPCServer{
		loanService: loanService,
		bookClient:  bookClient,
	}
}

func (s *loanGRPCServer) CreateLoan(ctx context.Context, req *protoLoan.CreateLoanRequest) (*protoLoan.LoanResponse, error) {
	// check book existence
	book, _ := s.bookClient.GetBook(ctx, req.BookId)
	if book == nil {
		return nil, status.Error(codes.NotFound, "book not found")
	}

	userLoan, _, _ := s.loanService.GetLoanByBookIdAndUserId(ctx, req.BookId, req.UserId)
	if userLoan != nil && userLoan.Status == "BORROWED" {
		return nil, status.Error(codes.Canceled, "user must return the borrowed book before creating a new loan")
	}

	if book.Stock == 0 {
		return nil, status.Error(codes.Unavailable, "book is not available")
	}

	err := s.bookClient.DecrementBookStock(ctx, book.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed when update stock book")
	}

	loan, code, err := s.loanService.CreateLoan(ctx, req.UserId, req.Email, book)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return &protoLoan.LoanResponse{
		Loan: &protoLoan.Loan{
			Id:        loan.Id,
			UserId:    loan.UserId,
			BookId:    loan.BookId,
			LoanDate:  loan.LoanDate.Unix(),
			Status:    loan.Status,
			CreatedAt: loan.CreatedAt.Unix(),
			UpdatedAt: loan.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *loanGRPCServer) GetLoan(ctx context.Context, req *protoLoan.GetLoanRequest) (*protoLoan.LoanResponse, error) {
	loan, code, err := s.loanService.GetLoan(ctx, req.Id)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var returnDate int64
	if loan.ReturnDate != nil {
		returnDate = loan.ReturnDate.Unix()
	} else {
		returnDate = 0
	}

	return &protoLoan.LoanResponse{
		Loan: &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate, // Menggunakan returnDate yang sudah diproses
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *loanGRPCServer) UpdateLoanStatus(ctx context.Context, req *protoLoan.UpdateLoanStatusRequest) (*protoLoan.LoanResponse, error) {
	loan, code, err := s.loanService.UpdateLoanStatus(ctx, req.Id, req.Status, time.Unix(req.ReturnDate, 0))
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var returnDate int64
	if loan.ReturnDate != nil {
		returnDate = loan.ReturnDate.Unix()
	} else {
		returnDate = 0
	}

	return &protoLoan.LoanResponse{
		Loan: &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate,
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *loanGRPCServer) ListUserLoans(ctx context.Context, req *protoLoan.ListUserLoansRequest) (*protoLoan.ListLoansResponse, error) {
	loans, code, err := s.loanService.ListUserLoans(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var protoLoans []*protoLoan.Loan
	for _, loan := range loans {
		var returnDate int64
		if loan.ReturnDate != nil {
			returnDate = loan.ReturnDate.Unix()
		} else {
			returnDate = 0
		}

		protoLoans = append(protoLoans, &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate,
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		})
	}

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) ListLoans(ctx context.Context, req *protoLoan.ListLoansRequest) (*protoLoan.ListLoansResponse, error) {
	loans, code, err := s.loanService.ListLoans(ctx)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var protoLoans []*protoLoan.Loan
	for _, loan := range loans {
		var returnDate int64
		if loan.ReturnDate != nil {
			returnDate = loan.ReturnDate.Unix()
		} else {
			returnDate = 0
		}

		protoLoans = append(protoLoans, &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate,
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		})
	}

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) GetUserLoansByStatus(ctx context.Context, req *protoLoan.GetUserLoansByStatusRequest) (*protoLoan.ListLoansResponse, error) {
	// Call service layer to get uer loans by status
	loans, code, err := s.loanService.GetUserLoansByStatus(ctx, req.UserId, req.Status)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var protoLoans []*protoLoan.Loan
	for _, loan := range loans {
		var returnDate int64
		if loan.ReturnDate != nil {
			returnDate = loan.ReturnDate.Unix()
		} else {
			returnDate = 0
		}

		protoLoans = append(protoLoans, &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate,
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		})
	}

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) GetLoansByStatus(ctx context.Context, req *protoLoan.GetLoansByStatusRequest) (*protoLoan.ListLoansResponse, error) {
	// Call service layer to get loans by status
	loans, code, err := s.loanService.GetLoansByStatus(ctx, req.Status)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	var protoLoans []*protoLoan.Loan
	for _, loan := range loans {
		var returnDate int64
		if loan.ReturnDate != nil {
			returnDate = loan.ReturnDate.Unix()
		} else {
			returnDate = 0
		}

		protoLoans = append(protoLoans, &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: returnDate,
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		})
	}

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) ReturnLoan(ctx context.Context, req *protoLoan.ReturnLoanRequest) (*protoLoan.LoanResponse, error) {
	loan, code, _ := s.loanService.GetLoan(ctx, req.Id)
	if loan == nil {
		return nil, status.Error(code, fmt.Sprintf("loan not found with id %s", req.Id))
	}

	if loan.Status == "RETURNED" {
		return nil, status.Error(codes.Canceled, "loan already returned")
	}

	if loan.UserId != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "you don't have access to this resource")
	}

	book, err := s.bookClient.GetBook(ctx, loan.BookId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	loan, code, err = s.loanService.ReturnLoan(ctx, req.Id, req.UserId, req.Email, book.Title, time.Unix(req.ReturnDate, 0))
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	err = s.bookClient.IncrementBookStock(ctx, loan.BookId)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	return &protoLoan.LoanResponse{
		Loan: &protoLoan.Loan{
			Id:         loan.Id,
			UserId:     loan.UserId,
			BookId:     loan.BookId,
			LoanDate:   loan.LoanDate.Unix(),
			ReturnDate: loan.ReturnDate.Unix(),
			Status:     loan.Status,
			CreatedAt:  loan.CreatedAt.Unix(),
			UpdatedAt:  loan.UpdatedAt.Unix(),
		},
	}, nil
}

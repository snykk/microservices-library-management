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

	loan, code, err := s.loanService.CreateLoan(ctx, req.UserId, req.BookId)
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
	loan, code, err := s.loanService.UpdateLoanStatus(ctx, req.Id, req.UserId, req.Role, req.Status, time.Unix(req.ReturnDate, 0))
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

	fmt.Println("setelah listloan")

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

	fmt.Println("success")

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) ListLoans(ctx context.Context, req *protoLoan.ListLoansRequest) (*protoLoan.ListLoansResponse, error) {
	loans, code, err := s.loanService.ListLoans(ctx)
	if err != nil {
		return nil, status.Error(code, err.Error())
	}

	fmt.Println("setelah listloan")

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

	fmt.Println("success")

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

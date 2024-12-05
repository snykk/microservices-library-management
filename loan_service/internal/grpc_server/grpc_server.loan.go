package grpc_server

import (
	"context"
	"fmt"
	"loan_service/internal/service"
	protoLoan "loan_service/proto/loan_service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loanGRPCServer struct {
	loanService service.LoanService
	protoLoan.UnimplementedLoanServiceServer
}

func NewLoanGRPCServer(loanService service.LoanService) protoLoan.LoanServiceServer {
	return &loanGRPCServer{
		loanService: loanService,
	}
}

func (s *loanGRPCServer) CreateLoan(ctx context.Context, req *protoLoan.CreateLoanRequest) (*protoLoan.LoanResponse, error) {
	loan, err := s.loanService.CreateLoan(ctx, req.UserId, req.BookId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create loan")
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
	loan, err := s.loanService.GetLoan(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "loan not found")
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
	loan, err := s.loanService.UpdateLoanStatus(ctx, req.Id, req.Status, time.Unix(req.ReturnDate, 0))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update loan status")
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
	loans, err := s.loanService.ListUserLoans(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list loans")
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
	loans, err := s.loanService.ListLoans(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list loans")
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

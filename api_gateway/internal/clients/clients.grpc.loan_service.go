package clients

import (
	"api_gateway/internal/datatransfers"
	protoLoan "api_gateway/proto/loan_service"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoanClient interface {
	CreateLoan(ctx context.Context, userId string, dto datatransfers.LoanRequest) (datatransfers.LoanResponse, error)
	GetLoan(ctx context.Context, id string) (datatransfers.LoanResponse, error)
	UpdateLoanStatus(ctx context.Context, loanId, userId, role, status string, returnDate time.Time) (datatransfers.LoanResponse, error)
	ListUserLoans(ctx context.Context, userId string) ([]datatransfers.LoanResponse, error)
	ListLoans(ctx context.Context) ([]datatransfers.LoanResponse, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]datatransfers.LoanResponse, error)
	GetLoansByStatus(ctx context.Context, status string) ([]datatransfers.LoanResponse, error)
}

type loanClient struct {
	client protoLoan.LoanServiceClient
}

func NewLoanClient() (LoanClient, error) {
	conn, err := grpc.NewClient("loan-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoLoan.NewLoanServiceClient(conn)
	return &loanClient{
		client: client,
	}, nil
}

func (l *loanClient) CreateLoan(ctx context.Context, userId string, dto datatransfers.LoanRequest) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.CreateLoanRequest{
		UserId: userId,
		BookId: dto.BookId,
	}

	resp, err := l.client.CreateLoan(ctx, &reqProto)
	if err != nil {
		return datatransfers.LoanResponse{}, err
	}

	return datatransfers.LoanResponse{
		Id:         resp.Loan.Id,
		UserId:     resp.Loan.UserId,
		BookId:     resp.Loan.BookId,
		LoanDate:   time.Unix(resp.Loan.LoanDate, 0),
		ReturnDate: nil,
		Status:     resp.Loan.Status,
		CreatedAt:  time.Unix(resp.Loan.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Loan.UpdatedAt, 0),
	}, nil
}

func (l *loanClient) GetLoan(ctx context.Context, id string) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.GetLoanRequest{
		Id: id,
	}

	resp, err := l.client.GetLoan(ctx, &reqProto)
	if err != nil {
		return datatransfers.LoanResponse{}, err
	}

	loanResponse := datatransfers.LoanResponse{
		Id:        resp.Loan.Id,
		UserId:    resp.Loan.UserId,
		BookId:    resp.Loan.BookId,
		LoanDate:  time.Unix(resp.Loan.LoanDate, 0),
		Status:    resp.Loan.Status,
		CreatedAt: time.Unix(resp.Loan.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Loan.UpdatedAt, 0),
	}

	if resp.Loan.ReturnDate != 0 {
		returnDate := time.Unix(resp.Loan.ReturnDate, 0)
		loanResponse.ReturnDate = &returnDate
	}

	return loanResponse, nil
}

func (l *loanClient) UpdateLoanStatus(ctx context.Context, loanId, userId, role, status string, returnDate time.Time) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.UpdateLoanStatusRequest{
		Id:         loanId,
		UserId:     userId,
		Role:       role,
		Status:     status,
		ReturnDate: returnDate.Unix(),
	}

	resp, err := l.client.UpdateLoanStatus(ctx, &reqProto)
	if err != nil {
		return datatransfers.LoanResponse{}, err
	}

	loanResponse := datatransfers.LoanResponse{
		Id:        resp.Loan.Id,
		UserId:    resp.Loan.UserId,
		BookId:    resp.Loan.BookId,
		LoanDate:  time.Unix(resp.Loan.LoanDate, 0),
		Status:    resp.Loan.Status,
		CreatedAt: time.Unix(resp.Loan.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Loan.UpdatedAt, 0),
	}

	if resp.Loan.ReturnDate != 0 {
		returnDate := time.Unix(resp.Loan.ReturnDate, 0)
		loanResponse.ReturnDate = &returnDate
	}

	return loanResponse, nil
}

func (l *loanClient) ListUserLoans(ctx context.Context, userId string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.ListUserLoansRequest{
		UserId: userId,
	}

	resp, err := l.client.ListUserLoans(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var loans []datatransfers.LoanResponse
	for _, loan := range resp.Loans {
		loanResponse := datatransfers.LoanResponse{
			Id:        loan.Id,
			UserId:    loan.UserId,
			BookId:    loan.BookId,
			LoanDate:  time.Unix(loan.LoanDate, 0),
			Status:    loan.Status,
			CreatedAt: time.Unix(loan.CreatedAt, 0),
			UpdatedAt: time.Unix(loan.UpdatedAt, 0),
		}

		if loan.ReturnDate != 0 {
			returnDate := time.Unix(loan.ReturnDate, 0)
			loanResponse.ReturnDate = &returnDate
		}

		loans = append(loans, loanResponse)
	}

	return loans, nil
}

func (l *loanClient) ListLoans(ctx context.Context) ([]datatransfers.LoanResponse, error) {
	resp, err := l.client.ListLoans(ctx, &protoLoan.ListLoansRequest{})
	if err != nil {
		return nil, err
	}

	var loans []datatransfers.LoanResponse
	for _, loan := range resp.Loans {
		loanResponse := datatransfers.LoanResponse{
			Id:        loan.Id,
			UserId:    loan.UserId,
			BookId:    loan.BookId,
			LoanDate:  time.Unix(loan.LoanDate, 0),
			Status:    loan.Status,
			CreatedAt: time.Unix(loan.CreatedAt, 0),
			UpdatedAt: time.Unix(loan.UpdatedAt, 0),
		}

		if loan.ReturnDate != 0 {
			returnDate := time.Unix(loan.ReturnDate, 0)
			loanResponse.ReturnDate = &returnDate
		}

		loans = append(loans, loanResponse)
	}

	return loans, nil
}

func (l *loanClient) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.GetUserLoansByStatusRequest{
		UserId: userId,
		Status: status,
	}

	resp, err := l.client.GetUserLoansByStatus(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var loans []datatransfers.LoanResponse
	for _, loan := range resp.Loans {
		loanResponse := datatransfers.LoanResponse{
			Id:        loan.Id,
			UserId:    loan.UserId,
			BookId:    loan.BookId,
			LoanDate:  time.Unix(loan.LoanDate, 0),
			Status:    loan.Status,
			CreatedAt: time.Unix(loan.CreatedAt, 0),
			UpdatedAt: time.Unix(loan.UpdatedAt, 0),
		}

		if loan.ReturnDate != 0 {
			returnDate := time.Unix(loan.ReturnDate, 0)
			loanResponse.ReturnDate = &returnDate
		}

		loans = append(loans, loanResponse)
	}

	return loans, nil
}

func (l *loanClient) GetLoansByStatus(ctx context.Context, status string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.GetLoansByStatusRequest{
		Status: status,
	}

	resp, err := l.client.GetLoansByStatus(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var loans []datatransfers.LoanResponse
	for _, loan := range resp.Loans {
		loanResponse := datatransfers.LoanResponse{
			Id:        loan.Id,
			UserId:    loan.UserId,
			BookId:    loan.BookId,
			LoanDate:  time.Unix(loan.LoanDate, 0),
			Status:    loan.Status,
			CreatedAt: time.Unix(loan.CreatedAt, 0),
			UpdatedAt: time.Unix(loan.UpdatedAt, 0),
		}

		if loan.ReturnDate != 0 {
			returnDate := time.Unix(loan.ReturnDate, 0)
			loanResponse.ReturnDate = &returnDate
		}

		loans = append(loans, loanResponse)
	}

	return loans, nil
}

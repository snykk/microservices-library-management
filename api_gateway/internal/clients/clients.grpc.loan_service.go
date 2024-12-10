package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoLoan "api_gateway/proto/loan_service"
	"context"
	"time"

	"api_gateway/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoanClient interface {
	CreateLoan(ctx context.Context, userId, email string, dto datatransfers.LoanRequest) (datatransfers.LoanResponse, error)
	GetLoan(ctx context.Context, id string) (datatransfers.LoanResponse, error)
	UpdateLoanStatus(ctx context.Context, loanId, status string, returnDate time.Time) (datatransfers.LoanResponse, error)
	ListUserLoans(ctx context.Context, userId string) ([]datatransfers.LoanResponse, error)
	ListLoans(ctx context.Context) ([]datatransfers.LoanResponse, error)
	GetUserLoansByStatus(ctx context.Context, userId, status string) ([]datatransfers.LoanResponse, error)
	GetLoansByStatus(ctx context.Context, status string) ([]datatransfers.LoanResponse, error)
	ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (datatransfers.LoanResponse, error)
}

type loanClient struct {
	client protoLoan.LoanServiceClient
}

func NewLoanClient() (LoanClient, error) {
	conn, err := grpc.NewClient("loan-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Failed to create LoanClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)
		return nil, err
	}

	client := protoLoan.NewLoanServiceClient(conn)

	logger.Log.Info("Successfully created LoanClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)

	return &loanClient{
		client: client,
	}, nil
}

func (l *loanClient) CreateLoan(ctx context.Context, userId, email string, dto datatransfers.LoanRequest) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.CreateLoanRequest{
		UserId: userId,
		BookId: dto.BookId,
		Email:  email,
	}

	logger.Log.Info("Sending CreateLoan request to Loan Service",
		zap.String("user_id", userId),
		zap.String("book_id", dto.BookId),
		zap.String("email", email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.CreateLoan(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("CreateLoan request failed",
			zap.String("user_id", userId),
			zap.String("book_id", dto.BookId),
			zap.String("email", email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.LoanResponse{}, err
	}

	logger.Log.Info("CreateLoan request succeeded",
		zap.String("loan_id", resp.Loan.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

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

	logger.Log.Info("Sending GetLoan request to Loan Service",
		zap.String("loan_id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.GetLoan(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetLoan request failed",
			zap.String("loan_id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("GetLoan request succeeded",
		zap.String("loan_id", resp.Loan.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loanResponse, nil
}

func (l *loanClient) UpdateLoanStatus(ctx context.Context, loanId, status string, returnDate time.Time) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.UpdateLoanStatusRequest{
		Id:         loanId,
		Status:     status,
		ReturnDate: returnDate.Unix(),
	}

	logger.Log.Info("Sending UpdateLoanStatus request to Loan Service",
		zap.String("loan_id", loanId),
		zap.String("status", status),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.UpdateLoanStatus(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("UpdateLoanStatus request failed",
			zap.String("loan_id", loanId),
			zap.String("status", status),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("UpdateLoanStatus request succeeded",
		zap.String("loan_id", resp.Loan.Id),
		zap.String("status", resp.Loan.Status),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loanResponse, nil
}

func (l *loanClient) ListUserLoans(ctx context.Context, userId string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.ListUserLoansRequest{
		UserId: userId,
	}

	logger.Log.Info("Sending ListUserLoans request to Loan Service",
		zap.String("user_id", userId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.ListUserLoans(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ListUserLoans request failed",
			zap.String("user_id", userId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("ListUserLoans request succeeded",
		zap.String("user_id", userId),
		zap.Int("loans_count", len(loans)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loans, nil
}

func (l *loanClient) ListLoans(ctx context.Context) ([]datatransfers.LoanResponse, error) {
	logger.Log.Info("Sending ListLoans request to Loan Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.ListLoans(ctx, &protoLoan.ListLoansRequest{})
	if err != nil {
		logger.Log.Error("ListLoans request failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("ListLoans request succeeded",
		zap.Int("loans_count", len(loans)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loans, nil
}

func (l *loanClient) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.GetUserLoansByStatusRequest{
		UserId: userId,
		Status: status,
	}

	logger.Log.Info("Sending GetUserLoansByStatus request to Loan Service",
		zap.String("user_id", userId),
		zap.String("status", status),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.GetUserLoansByStatus(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetUserLoansByStatus request failed",
			zap.String("user_id", userId),
			zap.String("status", status),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("GetUserLoansByStatus request succeeded",
		zap.String("user_id", userId),
		zap.String("status", status),
		zap.Int("loans_count", len(loans)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loans, nil
}

func (l *loanClient) GetLoansByStatus(ctx context.Context, status string) ([]datatransfers.LoanResponse, error) {
	reqProto := protoLoan.GetLoansByStatusRequest{
		Status: status,
	}

	logger.Log.Info("Sending GetLoansByStatus request to Loan Service",
		zap.String("status", status),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.GetLoansByStatus(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetLoansByStatus request failed",
			zap.String("status", status),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("GetLoansByStatus request succeeded",
		zap.String("status", status),
		zap.Int("loans_count", len(loans)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loans, nil
}

func (l *loanClient) ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (datatransfers.LoanResponse, error) {
	reqProto := protoLoan.ReturnLoanRequest{
		Id:         id,
		Email:      email,
		UserId:     userId,
		ReturnDate: returnDate.Unix(),
	}

	logger.Log.Info("Sending ReturnLoan request to Loan Service",
		zap.String("loan_id", id),
		zap.String("user_id", userId),
		zap.String("email", email),
		zap.String("return_date", returnDate.String()),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := l.client.ReturnLoan(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ReturnLoan request failed",
			zap.String("loan_id", id),
			zap.String("user_id", userId),
			zap.String("email", email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("ReturnLoan request succeeded",
		zap.String("loan_id", resp.Loan.Id),
		zap.String("status", resp.Loan.Status),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return loanResponse, nil
}

package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoLoan "api_gateway/proto/loan_service"
	"context"
	"log"
	"time"

	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

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
	logger *logger.Logger
}

func NewLoanClient(logger *logger.Logger) (LoanClient, error) {
	conn, err := grpc.NewClient("loan-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create LoanClient:", err)
		return nil, err
	}
	client := protoLoan.NewLoanServiceClient(conn)

	log.Println("Successfully created LoanClient")

	return &loanClient{
		client: client,
		logger: logger,
	}, nil
}

func (l *loanClient) CreateLoan(ctx context.Context, userId, email string, dto datatransfers.LoanRequest) (datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.CreateLoanRequest{
		UserId: userId,
		BookId: dto.BookId,
		Email:  email,
	}

	extra := map[string]interface{}{
		"user_id": userId,
		"book_id": dto.BookId,
		"email":   email,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending CreateLoan request to Loan Service", extra, nil)

	resp, err := l.client.CreateLoan(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "CreateLoan request failed", extra, err)
		return datatransfers.LoanResponse{}, err
	}

	extra["loan_id"] = resp.Loan.Id

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "CreateLoan request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.GetLoanRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"loan_id": id,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetLoan request to Loan Service", extra, nil)

	resp, err := l.client.GetLoan(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetLoan request failed", extra, err)

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

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetLoan request succeeded", extra, nil)

	return loanResponse, nil
}

func (l *loanClient) UpdateLoanStatus(ctx context.Context, loanId, status string, returnDate time.Time) (datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.UpdateLoanStatusRequest{
		Id:         loanId,
		Status:     status,
		ReturnDate: returnDate.Unix(),
	}

	extra := map[string]interface{}{
		"loan_id":          loanId,
		"loan_status":      status,
		"loan_return_date": returnDate,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending UpdateLoanStatus request to Loan Service", extra, nil)

	resp, err := l.client.UpdateLoanStatus(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "UpdateLoanStatus request failed", extra, err)
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

	extra["status"] = resp.Loan.Status
	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "UpdateLoanStatus request succeeded", extra, nil)

	return loanResponse, nil
}

func (l *loanClient) ListUserLoans(ctx context.Context, userId string) ([]datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.ListUserLoansRequest{
		UserId: userId,
	}

	extra := map[string]interface{}{
		"user_id": userId,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListUserLoans request to Loan Service", extra, nil)

	resp, err := l.client.ListUserLoans(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListUserLoans request failed", extra, err)
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

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListUserLoans request succeeded", extra, nil)

	return loans, nil
}

func (l *loanClient) ListLoans(ctx context.Context) ([]datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListLoans request to Loan Service", nil, nil)

	resp, err := l.client.ListLoans(ctx, &protoLoan.ListLoansRequest{})
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListLoans request failed", nil, err)

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

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListLoans request to Loan Service", map[string]interface{}{"loans_count": len(loans)}, nil)

	return loans, nil
}

func (l *loanClient) GetUserLoansByStatus(ctx context.Context, userId, status string) ([]datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.GetUserLoansByStatusRequest{
		UserId: userId,
		Status: status,
	}

	extra := map[string]interface{}{
		"user_id":     userId,
		"loan_status": status,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetUserLoansByStatus request to Loan Service", extra, nil)

	resp, err := l.client.GetUserLoansByStatus(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetUserLoansByStatus request failed", extra, err)
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

	extra["loans_count"] = len(loans)
	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetUserLoansByStatus request succeeded", extra, nil)

	return loans, nil
}

func (l *loanClient) GetLoansByStatus(ctx context.Context, status string) ([]datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.GetLoansByStatusRequest{
		Status: status,
	}

	extra := map[string]interface{}{
		"loan_status": status,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetLoansByStatus request to Loan Service", extra, nil)

	resp, err := l.client.GetLoansByStatus(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetLoansByStatus request failed", extra, err)
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

	extra["loans_count"] = len(loans)
	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetLoansByStatus request succeeded", extra, nil)

	return loans, nil
}

func (l *loanClient) ReturnLoan(ctx context.Context, id, userId, email string, returnDate time.Time) (datatransfers.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoLoan.ReturnLoanRequest{
		Id:         id,
		Email:      email,
		UserId:     userId,
		ReturnDate: returnDate.Unix(),
	}

	extra := map[string]interface{}{
		"loan_id":          id,
		"email":            email,
		"user_id":          userId,
		"loan_return_date": returnDate,
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ReturnLoan request to Loan Service", extra, nil)

	resp, err := l.client.ReturnLoan(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ReturnLoan request failed", extra, err)
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

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ReturnLoan request succeeded", extra, nil)

	return loanResponse, nil
}

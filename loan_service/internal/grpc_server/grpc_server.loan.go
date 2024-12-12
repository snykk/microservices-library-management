package grpc_server

import (
	"context"
	"loan_service/internal/constants"
	"loan_service/internal/service"
	"loan_service/pkg/logger"
	"loan_service/pkg/utils"
	protoLoan "loan_service/proto/loan_service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loanGRPCServer struct {
	loanService service.LoanService
	logger      *logger.Logger
	protoLoan.UnimplementedLoanServiceServer
}

func NewLoanGRPCServer(loanService service.LoanService, logger *logger.Logger) protoLoan.LoanServiceServer {
	return &loanGRPCServer{
		loanService: loanService,
		logger:      logger,
	}
}

func (s *loanGRPCServer) CreateLoan(ctx context.Context, req *protoLoan.CreateLoanRequest) (*protoLoan.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received CreateLoan request", map[string]interface{}{"user_id": req.UserId, "book_id": req.BookId, "email": req.Email}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid CreateLoan request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loan, code, err := s.loanService.CreateLoan(context.WithValue(ctx, constants.ContextRequestIDKey, requestID), req.UserId, req.Email, req.BookId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create loan", nil, err)
		return nil, status.Error(code, err.Error())
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan created successfully", map[string]interface{}{"loan_id": loan.Id}, nil)

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

func (s *loanGRPCServer) ReturnLoan(ctx context.Context, req *protoLoan.ReturnLoanRequest) (*protoLoan.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ReturnLoan request", map[string]interface{}{"loan_id": req.Id, "user_id": req.UserId, "email": req.Email, "return_date": req.ReturnDate}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ReturnLoan request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loan, code, err := s.loanService.ReturnLoan(context.WithValue(ctx, constants.ContextRequestIDKey, requestID), req.Id, req.UserId, req.Email, time.Unix(req.ReturnDate, 0))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to return loan", nil, err)
		return nil, status.Error(code, err.Error())
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan returned successfully", map[string]interface{}{"loan_id": loan.Id}, nil)

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

func (s *loanGRPCServer) GetLoan(ctx context.Context, req *protoLoan.GetLoanRequest) (*protoLoan.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetLoan request", map[string]interface{}{"loan_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetLoan request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loan, code, err := s.loanService.GetLoan(ctx, req.Id)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve loan", nil, err)
		return nil, status.Error(code, err.Error())
	}

	var returnDate int64
	if loan.ReturnDate != nil {
		returnDate = loan.ReturnDate.Unix()
	} else {
		returnDate = 0
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan retrieved successfully", map[string]interface{}{"loan_id": loan.Id}, nil)

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

func (s *loanGRPCServer) UpdateLoanStatus(ctx context.Context, req *protoLoan.UpdateLoanStatusRequest) (*protoLoan.LoanResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received UpdateLoanStatus request", map[string]interface{}{"loan_id": req.Id, "status": req.Status}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid UpdateLoanStatus request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loan, code, err := s.loanService.UpdateLoanStatus(ctx, req.Id, req.Status, time.Unix(req.ReturnDate, 0))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to update loan status", nil, err)
		return nil, status.Error(code, err.Error())
	}

	var returnDate int64
	if loan.ReturnDate != nil {
		returnDate = loan.ReturnDate.Unix()
	} else {
		returnDate = 0
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan status updated successfully", map[string]interface{}{"loan_id": loan.Id}, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListUserLoans request", map[string]interface{}{"user_id": req.UserId}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListUserLoans request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loans, code, err := s.loanService.ListUserLoans(ctx, req.UserId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve user loans", nil, err)
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

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User loans retrieved successfully", nil, nil)

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) ListLoans(ctx context.Context, req *protoLoan.ListLoansRequest) (*protoLoan.ListLoansResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListLoans request", nil, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListLoans request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	loans, code, err := s.loanService.ListLoans(ctx)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve loan list", nil, err)
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

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan list retrieved successfully", nil, nil)

	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) GetUserLoansByStatus(ctx context.Context, req *protoLoan.GetUserLoansByStatusRequest) (*protoLoan.ListLoansResponse, error) {
	// Mendapatkan request ID dari konteks untuk logging
	requestID := utils.GetRequestIDFromContext(ctx)

	// Logging permintaan untuk mendapatkan pinjaman berdasarkan status pengguna
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetUserLoansByStatus request", map[string]interface{}{"user_id": req.UserId, "status": req.Status}, nil)

	// Validasi permintaan dari klien
	if err := req.Validate(); err != nil {
		// Logging error pada saat validasi
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetUserLoansByStatus request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Mengambil pinjaman berdasarkan status pengguna
	loans, code, err := s.loanService.GetUserLoansByStatus(ctx, req.UserId, req.Status)
	if err != nil {
		// Logging error pada saat mengambil data pinjaman
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve user loans by status", map[string]interface{}{"user_id": req.UserId, "status": req.Status}, err)
		return nil, status.Error(code, err.Error())
	}

	// Membentuk daftar pinjaman untuk respons
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

	// Logging bahwa data pinjaman berdasarkan status pengguna berhasil didapatkan
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User loans by status retrieved successfully", map[string]interface{}{"user_id": req.UserId, "status": req.Status}, nil)

	// Mengembalikan respons dengan daftar pinjaman
	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

func (s *loanGRPCServer) GetLoansByStatus(ctx context.Context, req *protoLoan.GetLoansByStatusRequest) (*protoLoan.ListLoansResponse, error) {
	// Mendapatkan request ID dari konteks untuk logging
	requestID := utils.GetRequestIDFromContext(ctx)

	// Logging permintaan untuk mendapatkan pinjaman berdasarkan status
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetLoansByStatus request", map[string]interface{}{"status": req.Status}, nil)

	// Validasi permintaan dari klien
	if err := req.Validate(); err != nil {
		// Logging error pada saat validasi
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetLoansByStatus request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Mengambil pinjaman berdasarkan status
	loans, code, err := s.loanService.GetLoansByStatus(ctx, req.Status)
	if err != nil {
		// Logging error pada saat mengambil data pinjaman
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve loans by status", map[string]interface{}{"status": req.Status}, err)
		return nil, status.Error(code, err.Error())
	}

	// Membentuk daftar pinjaman untuk respons
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

	// Logging bahwa data pinjaman berdasarkan status berhasil didapatkan
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loans by status retrieved successfully", map[string]interface{}{"status": req.Status}, nil)

	// Mengembalikan respons dengan daftar pinjaman
	return &protoLoan.ListLoansResponse{
		Loans: protoLoans,
	}, nil
}

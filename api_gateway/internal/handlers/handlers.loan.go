package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoanHandler struct {
	client clients.LoanClient
}

func NewLoanHandler(client clients.LoanClient) LoanHandler {
	return LoanHandler{
		client: client,
	}
}

func (l *LoanHandler) CreateLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	var req datatransfers.LoanRequest
	if err := c.BodyParser(&req); err != nil {
		// Log error on failed request body parsing
		logger.Log.Error("Failed to parse create loan request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := l.client.CreateLoan(c.Context(), userID, userEmail, req)
	if err != nil {
		// Log error from loan creation failure
		logger.Log.Error("Failed to create loan",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create loan", err))
	}

	// Log successful loan creation
	logger.Log.Info("Loan created successfully",
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
		zap.String("user_email", userEmail),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Loan created successfully", resp))
}

func (l *LoanHandler) ReturnLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")
	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	resp, err := l.client.ReturnLoan(c.Context(), loanId, userID, userEmail, time.Now())
	if err != nil {
		// Log error from returning loan failure
		logger.Log.Error("Failed to return loan",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to return loan", err))
	}

	// Log successful loan return
	logger.Log.Info("Loan status updated successfully",
		zap.String("request_id", requestID),
		zap.String("loan_id", loanId),
		zap.String("user_id", userID),
		zap.String("user_email", userEmail),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) GetLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")

	resp, err := l.client.GetLoan(c.Context(), loanId)
	if err != nil {
		// Log error from fetching loan
		logger.Log.Error("Failed to get loan",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get loan", err))
	}

	// Log successful loan fetch
	logger.Log.Info("Loan data fetched successfully",
		zap.String("request_id", requestID),
		zap.String("loan_id", loanId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Loan data with id '%s' fetched successfully", loanId), resp))
}

func (l *LoanHandler) UpdateLoanStatusHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")
	var req datatransfers.LoanStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		// Log error on failed request body parsing
		logger.Log.Error("Failed to parse update loan status request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	req.ReturnDate = time.Now()
	resp, err := l.client.UpdateLoanStatus(c.Context(), loanId, req.Status, req.ReturnDate)
	if err != nil {
		// Log error from loan status update failure
		logger.Log.Error("Failed to update loan status",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update loan status", err))
	}

	// Log successful loan status update
	logger.Log.Info("Loan status updated successfully",
		zap.String("request_id", requestID),
		zap.String("loan_id", loanId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) ListUserLoansHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userId := c.Locals("userID").(string)
	status := c.Query("status")

	var (
		resp []datatransfers.LoanResponse
		err  error
	)

	if status != "" {
		resp, err = l.client.GetUserLoansByStatus(c.Context(), userId, status)
	} else {
		resp, err = l.client.ListUserLoans(c.Context(), userId)
	}

	if err != nil {
		// Log error from fetching loans
		logger.Log.Error("Failed to list user loans",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	// Log successful loan list fetch
	logger.Log.Info("List user loans fetched successfully",
		zap.String("request_id", requestID),
		zap.String("user_id", userId),
		zap.String("status", status),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List user loans fetched successfully", resp))
}

func (l *LoanHandler) ListLoansHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	status := c.Query("status")

	var (
		resp []datatransfers.LoanResponse
		err  error
	)

	if status != "" {
		resp, err = l.client.GetLoansByStatus(c.Context(), status)
	} else {
		resp, err = l.client.ListLoans(c.Context())
	}

	if err != nil {
		// Log error from fetching loans
		logger.Log.Error("Failed to list loans",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	// Log successful loan list fetch
	logger.Log.Info("List loans fetched successfully",
		zap.String("request_id", requestID),
		zap.String("status", status),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List loans fetched successfully", resp))
}

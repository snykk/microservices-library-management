package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoanHandler struct {
	client clients.LoanClient
	logger *logger.Logger
}

func NewLoanHandler(client clients.LoanClient, logger *logger.Logger) LoanHandler {
	return LoanHandler{
		client: client,
		logger: logger,
	}
}

func (l *LoanHandler) CreateLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	extra := map[string]interface{}{
		"method":     c.Method(),
		"url":        c.OriginalURL(),
		"user_id":    userID,
		"user_email": userEmail,
	}

	var req datatransfers.LoanRequest
	if err := c.BodyParser(&req); err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse create loan request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["loan_book_id"] = req.BookId

	resp, err := l.client.CreateLoan(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userID, userEmail, req)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create loan", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create loan", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan created successfully", extra, nil)

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Loan created successfully", resp))
}

func (l *LoanHandler) ReturnLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")
	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	extra := map[string]interface{}{
		"method":     c.Method(),
		"url":        c.OriginalURL(),
		"loan_id":    loanId,
		"user_id":    userID,
		"user_email": userEmail,
	}

	resp, err := l.client.ReturnLoan(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), loanId, userID, userEmail, time.Now())
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to return loan", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to return loan", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan status updated successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) GetLoanHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")

	extra := map[string]interface{}{
		"method":  c.Method(),
		"url":     c.OriginalURL(),
		"loan_id": loanId,
	}

	resp, err := l.client.GetLoan(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), loanId)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get loan", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get loan", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan data fetched successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Loan data with id '%s' fetched successfully", loanId), resp))
}

func (l *LoanHandler) UpdateLoanStatusHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	loanId := c.Params("id")

	extra := map[string]interface{}{
		"method":  c.Method(),
		"url":     c.OriginalURL(),
		"loan_id": loanId,
	}

	var req datatransfers.LoanStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse update loan status request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["loan_status"] = req.Status
	extra["loan_return_date"] = req.ReturnDate

	req.ReturnDate = time.Now()
	resp, err := l.client.UpdateLoanStatus(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), loanId, req.Status, req.ReturnDate)
	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to update loan status", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update loan status", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Loan status updated successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) ListUserLoansHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userId := c.Locals("userID").(string)
	status := c.Query("status")

	extra := map[string]interface{}{
		"method":  c.Method(),
		"url":     c.OriginalURL(),
		"user_id": userId,
		"status":  status,
	}

	var resp []datatransfers.LoanResponse
	var err error

	if status != "" {
		resp, err = l.client.GetUserLoansByStatus(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userId, status)
	} else {
		resp, err = l.client.ListUserLoans(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userId)
	}

	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to list user loans", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "List user loans fetched successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List user loans fetched successfully", resp))
}

func (l *LoanHandler) ListLoansHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	status := c.Query("status")

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
		"status": status,
	}

	var resp []datatransfers.LoanResponse
	var err error

	if status != "" {
		resp, err = l.client.GetLoansByStatus(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), status)
	} else {
		resp, err = l.client.ListLoans(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID))
	}

	if err != nil {
		l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to list loans", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	l.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "List loans fetched successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List loans fetched successfully", resp))
}

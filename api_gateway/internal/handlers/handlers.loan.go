package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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
	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	var req datatransfers.LoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := l.client.CreateLoan(c.Context(), userID, userEmail, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create loan", err))
	}

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Loan created successfully", resp))
}

func (l *LoanHandler) ReturnLoanHandler(c *fiber.Ctx) error {
	loanId := c.Params("id")
	userID := c.Locals("userID").(string)
	userEmail := c.Locals("email").(string)

	resp, err := l.client.ReturnLoan(c.Context(), loanId, userID, userEmail, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to return loan", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) GetLoanHandler(c *fiber.Ctx) error {
	loanId := c.Params("id")

	resp, err := l.client.GetLoan(c.Context(), loanId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get loan", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Loan data with id '%s' fetched successfully", loanId), resp))
}

func (l *LoanHandler) UpdateLoanStatusHandler(c *fiber.Ctx) error {
	loanId := c.Params("id")
	var req datatransfers.LoanStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	req.ReturnDate = time.Now()
	resp, err := l.client.UpdateLoanStatus(c.Context(), loanId, req.Status, req.ReturnDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update loan status", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Loan status updated successfully", resp))
}

func (l *LoanHandler) ListUserLoansHandler(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List user loans fetched successfully", resp))
}

func (l *LoanHandler) ListLoansHandler(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list loans", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("List loans fetched successfully", resp))
}

package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	client clients.AuthClient
}

func NewAuthHandler(client clients.AuthClient) AuthHandler {
	return AuthHandler{
		client: client,
	}
}

func (authH *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	var req datatransfers.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := authH.client.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to register", err))
	}

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Registration successful", resp))
}

func (authH *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	var req datatransfers.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := authH.client.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ResponseError("Failed to login", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Login successful", resp))
}

func (authH *AuthHandler) SendOtpHandler(c *fiber.Ctx) error {
	var req datatransfers.SendOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := authH.client.SendOtp(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to verify email", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Otp code will be sent to  %s", req.Email), resp))
}

func (authH *AuthHandler) VerifyEmailHandler(c *fiber.Ctx) error {
	var req datatransfers.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := authH.client.VerifyEmail(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to verify email", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Email verification successful", resp))
}

func (authH *AuthHandler) ValidateTokenHandler(c *fiber.Ctx) error {
	var req datatransfers.ValidateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := authH.client.ValidateToken(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to validate token", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Token validation successful", resp))
}

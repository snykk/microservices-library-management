package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse register request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to register
	resp, err := authH.client.Register(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to register", err))
	}

	logger.Log.Info("User registration successful",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Registration successful", resp))
}

func (authH *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse login request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to login
	resp, err := authH.client.Login(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to login",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ResponseError("Failed to login", err))
	}

	logger.Log.Info("User login successful",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Login successful", resp))
}

func (authH *AuthHandler) SendOtpHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.SendOtpRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse send OTP request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to send OTP
	resp, err := authH.client.SendOtp(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to send OTP",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to verify email", err))
	}

	logger.Log.Info("OTP sent successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Otp code will be sent to  %s", req.Email), resp))
}

func (authH *AuthHandler) VerifyEmailHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse verify email request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to verify email
	resp, err := authH.client.VerifyEmail(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to verify email",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to verify email", err))
	}

	logger.Log.Info("Email verification successful",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Email verification successful", resp))
}

func (authH *AuthHandler) ValidateTokenHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.ValidateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse validate token request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to validate token
	resp, err := authH.client.ValidateToken(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to validate token",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to validate token", err))
	}

	logger.Log.Info("Token validation successful",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Token validation successful", resp))
}

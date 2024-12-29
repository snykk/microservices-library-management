package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	client clients.AuthClient
	logger *logger.Logger
}

func NewAuthHandler(client clients.AuthClient, logger *logger.Logger) AuthHandler {
	return AuthHandler{
		client: client,
		logger: logger,
	}
}

// func (authH *AuthHandler) logMessage(caller, requestID, level, message string, extra map[string]interface{}, err error) {
// 	logMsg := models.LogMessage{
// 		Timestamp:      time.Now(),
// 		Service:        constants.LogServiceApiGateway,
// 		Level:          level,
// 		XCorrelationID: requestID,
// 		Caller:         caller,
// 		Message:        message,
// 		Extra:          extra,
// 	}

// 	if err != nil {
// 		logMsg.Error = err.Error()
// 	}

// 	// go func() {
// 	// 	defer func() {
// 	// 		if r := recover(); r != nil {
// 	// 			log.Printf("[%s] Recovered from panic in logMessage goroutine: %v\n", utils.GetLocation(), r)
// 	// 		}
// 	// 	}()
// 	// 	if pubErr := authH.rabbitMQPublisher.Publish(constants.LogExchange, constants.LogQueue, logMsg); pubErr != nil {
// 	// 		log.Printf("[%s] Failed to publish log to RabbitMQ: %v\n", caller, pubErr)
// 	// 	}
// 	// }()

// 	// kirim log ke channel
// 	authH.logger.LogChannel <- logMsg
// }

func (authH *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse register request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["email"] = req.Email
	extra["password"] = req.Password

	resp, err := authH.client.Register(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to register user", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to register", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User registration successful", extra, nil)
	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Registration successful", resp))
}

func (authH *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse login request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["email"] = req.Email

	resp, err := authH.client.Login(context.WithValue(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to login", extra, err)
		return c.Status(fiber.StatusUnauthorized).JSON(datatransfers.ResponseError("Failed to login", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User login successful", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Login successful", resp))
}

func (authH *AuthHandler) SendOtpHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.SendOtpRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse send OTP request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["email"] = req.Email

	resp, err := authH.client.SendOtp(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to send OTP", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to send OTP", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "OTP sent successfully", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Otp code will be sent to  %s", req.Email), resp))
}

func (authH *AuthHandler) VerifyEmailHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse verify email request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["email"] = req.Email
	extra["otp"] = req.OTP

	resp, err := authH.client.VerifyEmail(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to verify email", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to verify email", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Email verification successful", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Email verification successful", resp))
}

func (authH *AuthHandler) ValidateTokenHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.ValidateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse validate token request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	resp, err := authH.client.ValidateToken(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to validate token", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to validate token", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Token validation successful", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Token validation successful", resp))
}

func (authH *AuthHandler) RefreshTokenHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userID := c.Locals("userID").(string)

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse refresh token request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["refresh_token"] = req.RefreshToken

	resp, err := authH.client.RefreshToken(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userID, req)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to refresh token", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to refresh token", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Token refresh successful", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Token refresh successful", resp))
}

func (authH *AuthHandler) LogoutHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	userID := c.Locals("userID").(string)

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	extra["user_id"] = userID

	resp, err := authH.client.Logout(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), userID)
	if err != nil {
		authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to logout", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to logout", err))
	}

	authH.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Logout successful", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Logout successful", resp))
}

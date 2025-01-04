package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthorHandler struct {
	client     clients.AuthorClient
	bookClient clients.BookClient
	logger     *logger.Logger
}

func NewAuthorHandler(client clients.AuthorClient, bookClient clients.BookClient, logger *logger.Logger) AuthorHandler {
	return AuthorHandler{
		client:     client,
		bookClient: bookClient,
		logger:     logger,
	}
}

func (a *AuthorHandler) CreateAuthorHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	// Parse the request body
	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse create author request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["author_name"] = req.Name
	extra["author_biography"] = req.Biography

	// Call client to create author
	resp, err := a.client.CreateAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create author", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to create author", err))
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author created successfully", extra, nil)

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Author created successfully", resp))
}

func (a *AuthorHandler) GetAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")
	includeBooks := c.Query("includeBooks", "false") == "true"

	extra := map[string]interface{}{
		"method":    c.Method(),
		"url":       c.OriginalURL(),
		"author_id": authorId,
	}

	// Call client to get author by id
	resp, err := a.client.GetAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), authorId)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get author by ID", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get author", err))
	}

	// If includeBooks query param is true, get books for the author
	if includeBooks {
		books, totalItems, _, err := a.bookClient.GetBooksByAuthorId(
			context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), resp.Id,
			1,
			5,
		)
		if err != nil {
			a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get books by author ID", extra, err)
			return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get books by author ID", err))
		}

		resp.SampleBooks = &books
		resp.TotalBooks = &totalItems
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author data fetched successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Author data with id '%s' fetched successfully", authorId), resp))
}

func (a *AuthorHandler) GetAllAuthorsHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeBooks := c.Query("includeBooks", "false") == "true"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":       c.Method(),
		"url":          c.OriginalURL(),
		"includeBooks": includeBooks,
		"page":         page,
		"page_size":    pageSize,
	}

	// Call client to get all authors
	authors, totalItems, totalPages, err := a.client.ListAuthors(
		context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID),
		page,
		pageSize,
	)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get author list", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get author list", err))
	}

	// If includeBooks query param is true, get books for each author
	if includeBooks {
		for i := range authors {
			books, totalItems, _, err := a.bookClient.GetBooksByAuthorId(
				context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), authors[i].Id,
				1,
				5,
			)
			if err != nil {
				a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get books by author ID", extra, err)
				return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get books by author ID", err))
			}

			authors[i].SampleBooks = &books
			authors[i].TotalBooks = &totalItems
		}
	}

	extra["authors_count"] = len(authors)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "All authors data fetched successfully", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Author data fetched successfully", map[string]interface{}{
		"authors": authors,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
}

func (a *AuthorHandler) UpdateAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")

	extra := map[string]interface{}{
		"method":    c.Method(),
		"url":       c.OriginalURL(),
		"author_id": authorId,
	}

	// Parse the request body
	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse update author request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["author_name"] = req.Name
	extra["author_biography"] = req.Biography

	// Call client to update author by id
	resp, err := a.client.UpdateAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), authorId, req)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to update author", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to update author", err))
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author updated successfully", extra, nil)

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Author updated successfully", resp))
}

func (a *AuthorHandler) DeleteAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")

	extra := map[string]interface{}{
		"method":    c.Method(),
		"url":       c.OriginalURL(),
		"author_id": authorId,
	}

	// Call client to delete author by id
	err := a.client.DeleteAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), authorId)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to delete author", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to delete author", err))
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author deleted successfully", extra, nil)

	return c.SendStatus(fiber.StatusNoContent)
}

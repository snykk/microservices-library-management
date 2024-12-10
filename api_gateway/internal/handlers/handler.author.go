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

type AuthorHandler struct {
	client     clients.AuthorClient
	bookClient clients.BookClient
}

func NewAuthorHandler(client clients.AuthorClient, bookClient clients.BookClient) AuthorHandler {
	return AuthorHandler{
		client:     client,
		bookClient: bookClient,
	}
}

func (a *AuthorHandler) CreateAuthorHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	// Parse the request body
	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse create author request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to create author
	resp, err := a.client.CreateAuthor(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to create author",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create author", err))
	}

	logger.Log.Info("Author created successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Author created successfully", resp))
}

func (a *AuthorHandler) GetAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")
	includeBooks := c.Query("includeBooks", "false") == "true"

	// Call client to get author by id
	resp, err := a.client.GetAuthor(c.Context(), authorId)
	if err != nil {
		logger.Log.Error("Failed to get author by ID",
			zap.String("request_id", requestID),
			zap.String("author_id", authorId),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get author", err))
	}

	// If includeBooks query param is true, get books for the author
	if includeBooks {
		books, err := a.bookClient.GetBooksByAuthorId(c.Context(), resp.Id)
		if err == nil {
			resp.Books = &books
		}
	}

	logger.Log.Info("Author data fetched successfully",
		zap.String("request_id", requestID),
		zap.String("author_id", authorId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Author data with id '%s' fetched successfully", authorId), resp))
}

func (a *AuthorHandler) GetAllAuthorsHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeBooks := c.Query("includeBooks", "false") == "true"

	// Call client to get all authors
	resp, err := a.client.ListAuthors(c.Context())
	if err != nil {
		logger.Log.Error("Failed to list authors",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list authors", err))
	}

	// If includeBooks query param is true, get books for each author
	if includeBooks {
		for i := range resp {
			books, err := a.bookClient.GetBooksByAuthorId(c.Context(), resp[i].Id)
			if err == nil {
				resp[i].Books = &books
			}
		}
	}

	logger.Log.Info("All authors data fetched successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Author data fetched successfully", resp))
}

func (a *AuthorHandler) UpdateAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")

	// Parse the request body
	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse update author request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("author_id", authorId),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	// Call client to update author by id
	resp, err := a.client.UpdateAuthor(c.Context(), authorId, req)
	if err != nil {
		logger.Log.Error("Failed to update author",
			zap.String("request_id", requestID),
			zap.String("author_id", authorId),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update author", err))
	}

	logger.Log.Info("Author updated successfully",
		zap.String("request_id", requestID),
		zap.String("author_id", authorId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Author updated successfully", resp))
}

func (a *AuthorHandler) DeleteAuthorByIdHandler(c *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("id")

	// Call client to delete author by id
	err := a.client.DeleteAuthor(c.Context(), authorId)
	if err != nil {
		logger.Log.Error("Failed to delete author",
			zap.String("request_id", requestID),
			zap.String("author_id", authorId),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete author", err))
	}

	logger.Log.Info("Author deleted successfully",
		zap.String("request_id", requestID),
		zap.String("author_id", authorId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.SendStatus(fiber.StatusNoContent)
}

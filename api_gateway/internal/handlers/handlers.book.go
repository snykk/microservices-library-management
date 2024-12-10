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

type BookHandler struct {
	client         clients.BookClient
	authorClient   clients.AuthorClient
	categoryClient clients.CategoryClient
}

func NewBookHandler(client clients.BookClient, authorClient clients.AuthorClient, categoryClient clients.CategoryClient) BookHandler {
	return BookHandler{
		client:         client,
		authorClient:   authorClient,
		categoryClient: categoryClient,
	}
}

func (b *BookHandler) CreateBookHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	var req datatransfers.BookRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse create book request",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := b.client.CreateBook(c.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to create book",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create book", err))
	}

	logger.Log.Info("Book created successfully",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Book created successfully", resp))
}

func (b *BookHandler) GetBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")
	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	resp, err := b.client.GetBook(c.Context(), bookId)
	if err != nil {
		logger.Log.Error("Failed to get book by id",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get book", err))
	}

	if includeAuthor || includeCategory {
		if includeAuthor {
			author, err := b.authorClient.GetAuthor(c.Context(), *resp.AuthorId)
			if err == nil {
				resp.Author = &author
				resp.AuthorId = nil
			}
		}

		if includeCategory {
			category, err := b.categoryClient.GetCategory(c.Context(), *resp.CategoryId)
			if err == nil {
				resp.Category = &category
				resp.CategoryId = nil
			}
		}
	}

	logger.Log.Info("Fetched book data by id",
		zap.String("request_id", requestID),
		zap.String("book_id", bookId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Book data with id '%s' fetched successfully", bookId), resp))
}

func (b *BookHandler) GetBooksByAuthorIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("authorId")

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	// Fetch books by author
	resp, err := b.client.GetBooksByAuthorId(c.Context(), authorId)
	if err != nil {
		logger.Log.Error("Failed to get books by author",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get books by author", err))
	}

	// If includeAuthor or includeCategory are true, fetch related data
	if includeAuthor || includeCategory {
		for i := range resp {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(c.Context(), *resp[i].AuthorId)
				if err == nil {
					resp[i].Author = &author
					resp[i].AuthorId = nil
				}
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(c.Context(), *resp[i].CategoryId)
				if err == nil {
					resp[i].Category = &category
					resp[i].CategoryId = nil
				}
			}
		}
	}

	logger.Log.Info("Fetched books by author",
		zap.String("request_id", requestID),
		zap.String("author_id", authorId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Books by author '%s' fetched successfully", authorId), resp))
}

func (b *BookHandler) GetBooksByCategoryIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := c.Params("categoryId")

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	// Fetch books by category
	resp, err := b.client.GetBooksByCategoryId(c.Context(), categoryId)
	if err != nil {
		logger.Log.Error("Failed to get books by category",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get books by category", err))
	}

	// If includeAuthor or includeCategory are true, fetch related data
	if includeAuthor || includeCategory {
		for i := range resp {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(c.Context(), *resp[i].AuthorId)
				if err == nil {
					resp[i].Author = &author
					resp[i].AuthorId = nil
				}
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(c.Context(), *resp[i].CategoryId)
				if err == nil {
					resp[i].Category = &category
					resp[i].CategoryId = nil
				}
			}
		}
	}

	logger.Log.Info("Fetched books by category",
		zap.String("request_id", requestID),
		zap.String("category_id", categoryId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Books under category '%s' fetched successfully", categoryId), resp))
}

func (b *BookHandler) GetAllBooksHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	resp, err := b.client.ListBooks(c.Context())
	if err != nil {
		logger.Log.Error("Failed to list books",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list books", err))
	}

	if includeAuthor || includeCategory {
		for i := range resp {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(c.Context(), *resp[i].AuthorId)
				if err == nil {
					resp[i].Author = &author
					resp[i].AuthorId = nil
				}
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(c.Context(), *resp[i].CategoryId)
				if err == nil {
					resp[i].Category = &category
					resp[i].CategoryId = nil
				}
			}
		}
	}

	logger.Log.Info("Fetched all books",
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Book data fetched successfully", resp))
}

func (b *BookHandler) UpdateBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")

	var req datatransfers.BookRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse update book request",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := b.client.UpdateBook(c.Context(), bookId, req)
	if err != nil {
		logger.Log.Error("Failed to update book",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update book", err))
	}

	logger.Log.Info("Book updated successfully",
		zap.String("request_id", requestID),
		zap.String("book_id", bookId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Book updated successfully", resp))
}

func (b *BookHandler) DeleteBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")

	err := b.client.DeleteBook(c.Context(), bookId)
	if err != nil {
		logger.Log.Error("Failed to delete book",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete book", err))
	}

	logger.Log.Info("Book deleted successfully",
		zap.String("request_id", requestID),
		zap.String("book_id", bookId),
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()),
	)
	return c.SendStatus(fiber.StatusNoContent)
}

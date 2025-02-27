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

type BookHandler struct {
	client         clients.BookClient
	authorClient   clients.AuthorClient
	categoryClient clients.CategoryClient
	logger         *logger.Logger
}

func NewBookHandler(client clients.BookClient, authorClient clients.AuthorClient, categoryClient clients.CategoryClient, logger *logger.Logger) BookHandler {
	return BookHandler{
		client:         client,
		authorClient:   authorClient,
		categoryClient: categoryClient,
		logger:         logger,
	}
}

func (b *BookHandler) CreateBookHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": c.Method(),
		"url":    c.OriginalURL(),
	}

	var req datatransfers.BookRequest
	if err := c.BodyParser(&req); err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to parse create book request", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["book_title"] = req.Title
	extra["book_author_id"] = req.AuthorId
	extra["book_category_id"] = req.CategoryId
	extra["book_stock"] = req.Stock

	resp, err := b.client.CreateBook(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to create book", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to create book", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book created successfully", extra, nil)
	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Book created successfully", resp))
}

func (b *BookHandler) GetBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")
	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	extra := map[string]interface{}{
		"method":          c.Method(),
		"url":             c.OriginalURL(),
		"book_id":         bookId,
		"includeAuthor":   includeAuthor,
		"includeCategory": includeCategory,
	}

	resp, err := b.client.GetBook(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), bookId)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get book by id", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get book", err))
	}

	if includeAuthor || includeCategory {
		if includeAuthor {
			author, err := b.authorClient.GetAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *resp.AuthorId)
			if err == nil {
				resp.Author = &author
				resp.AuthorId = nil
			}
		}

		if includeCategory {
			category, err := b.categoryClient.GetCategory(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *resp.CategoryId)
			if err == nil {
				resp.Category = &category
				resp.CategoryId = nil
			}
		}
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched book data by id", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Book data with id '%s' fetched successfully", bookId), resp))
}

func (b *BookHandler) GetBooksByAuthorIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	authorId := c.Params("authorId")
	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":          c.Method(),
		"url":             c.OriginalURL(),
		"author_id":       authorId,
		"includeAuthor":   includeAuthor,
		"includeCategory": includeCategory,
		"page":            page,
		"page_size":       pageSize,
	}

	// Fetch books by author
	books, totalItems, totalPages, err := b.client.GetBooksByAuthorId(
		context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID),
		authorId,
		page,
		pageSize,
	)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get books by author", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get books by author", err))
	}

	// If includeAuthor or includeCategory are true, fetch related data
	if includeAuthor || includeCategory {
		for i := range books {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].AuthorId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get author by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get author", err))
				}

				books[i].Author = &author
				books[i].AuthorId = nil
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].CategoryId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get category by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get category", err))
				}

				books[i].Category = &category
				books[i].CategoryId = nil
			}
		}
	}

	extra["books_count"] = len(books)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched books by author", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Books by author '%s' fetched successfully", authorId), map[string]interface{}{
		"books": books,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
}

func (b *BookHandler) GetBooksByCategoryIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := c.Params("categoryId")
	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":          c.Method(),
		"url":             c.OriginalURL(),
		"category_id":     categoryId,
		"includeAuthor":   includeAuthor,
		"includeCategory": includeCategory,
		"page":            page,
		"page_size":       pageSize,
	}

	// Fetch books by category
	books, totalItems, totalPages, err := b.client.GetBooksByCategoryId(
		context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID),
		categoryId,
		page,
		pageSize,
	)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get books by category", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get books by category", err))
	}

	// If includeAuthor or includeCategory are true, fetch related data
	if includeAuthor || includeCategory {
		for i := range books {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].AuthorId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get author by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get author", err))
				}

				books[i].Author = &author
				books[i].AuthorId = nil
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].CategoryId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get category by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get category", err))
				}

				books[i].Category = &category
				books[i].CategoryId = nil
			}
		}
	}

	extra["books_count"] = len(books)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched books by category", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Books under category '%s' fetched successfully", categoryId), map[string]interface{}{
		"books": books,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
}

func (b *BookHandler) GetAllBooksHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":          c.Method(),
		"url":             c.OriginalURL(),
		"includeAuthor":   includeAuthor,
		"includeCategory": includeCategory,
		"page":            page,
		"page_size":       pageSize,
	}

	books, totalItems, totalPages, err := b.client.ListBooks(
		context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID),
		page,
		pageSize,
	)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to list books", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to list books", err))
	}

	// If includeAuthor or includeCategory are true, fetch related data
	if includeAuthor || includeCategory {
		for i := range books {
			if includeAuthor {
				author, err := b.authorClient.GetAuthor(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].AuthorId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get author by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get author", err))
				}

				books[i].Author = &author
				books[i].AuthorId = nil
			}

			if includeCategory {
				category, err := b.categoryClient.GetCategory(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), *books[i].CategoryId)
				if err != nil {
					b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to get category by ID", extra, err)
					return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get category", err))
				}

				books[i].Category = &category
				books[i].CategoryId = nil
			}
		}
	}

	extra["books_count"] = len(books)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched all books", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Book data fetched successfully", map[string]interface{}{
		"books": books,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
}

func (b *BookHandler) UpdateBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")

	extra := map[string]interface{}{
		"method":  c.Method(),
		"url":     c.OriginalURL(),
		"book_id": bookId,
	}

	var req datatransfers.BookUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to parse update book request", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["book_title"] = req.Title
	extra["book_author_id"] = req.AuthorId
	extra["book_category_id"] = req.CategoryId
	extra["book_stock"] = req.Stock

	resp, err := b.client.UpdateBook(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), bookId, req)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to update book", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to update book", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book updated successfully", extra, nil)
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Book updated successfully", resp))
}

func (b *BookHandler) DeleteBookByIdHandler(c *fiber.Ctx) error {
	requestID, ok := c.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	bookId := c.Params("id")

	extra := map[string]interface{}{
		"method":  c.Method(),
		"url":     c.OriginalURL(),
		"book_id": bookId,
	}

	// Parse the request body
	var req datatransfers.BookDeleteRequest
	if err := c.BodyParser(&req); err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse update book request body", extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return c.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["book_version"] = req.Version

	err := b.client.DeleteBook(context.WithValue(c.Context(), constants.ContextRequestIDKey, requestID), bookId, req.Version)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, "error", "Failed to delete book", extra, err)
		return c.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to delete book", err))
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book deleted successfully", extra, nil)
	return c.Status(fiber.StatusNoContent).JSON(datatransfers.ResponseSuccess("Book deleted successfully", nil))
}

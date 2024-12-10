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

type CategoryHandler struct {
	client     clients.CategoryClient
	bookClient clients.BookClient
}

func NewCategoryHandler(client clients.CategoryClient, bookClient clients.BookClient) CategoryHandler {
	return CategoryHandler{
		client:     client,
		bookClient: bookClient,
	}
}

func (c *CategoryHandler) CreateCategoryHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse create category request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := c.client.CreateCategory(ctx.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to create category",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create category", err))
	}

	logger.Log.Info("Category created successfully",
		zap.String("request_id", requestID),
		zap.String("method", ctx.Method()),
		zap.String("url", ctx.OriginalURL()),
	)
	return ctx.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Category created successfully", resp))
}

func (c *CategoryHandler) GetCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")
	includeBooks := ctx.Query("includeBooks", "false") == "true"

	resp, err := c.client.GetCategory(ctx.Context(), categoryId)
	if err != nil {
		logger.Log.Error("Failed to get category",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get category", err))
	}

	if includeBooks {
		books, err := c.bookClient.GetBooksByCategoryId(ctx.Context(), resp.Id)
		if err == nil {
			resp.Books = &books
		}
	}

	logger.Log.Info("Fetched category data successfully",
		zap.String("request_id", requestID),
		zap.String("category_id", categoryId),
		zap.String("method", ctx.Method()),
		zap.String("url", ctx.OriginalURL()),
	)
	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Category data with id '%s' fetched successfully", categoryId), resp))
}

func (c *CategoryHandler) GetAllCategoriesHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeBooks := ctx.Query("includeBooks", "false") == "true"

	resp, err := c.client.ListCategories(ctx.Context())
	if err != nil {
		logger.Log.Error("Failed to list categories",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list categories", err))
	}

	if includeBooks {
		for i := range resp {
			books, err := c.bookClient.GetBooksByCategoryId(ctx.Context(), resp[i].Id)
			if err == nil {
				resp[i].Books = &books
			}
		}
	}

	logger.Log.Info("Fetched all categories successfully",
		zap.String("request_id", requestID),
		zap.String("method", ctx.Method()),
		zap.String("url", ctx.OriginalURL()),
	)
	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Category data fetched successfully", resp))
}

func (c *CategoryHandler) UpdateCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")

	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse update category request body",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := c.client.UpdateCategory(ctx.Context(), categoryId, req)
	if err != nil {
		logger.Log.Error("Failed to update category",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update category", err))
	}

	logger.Log.Info("Category updated successfully",
		zap.String("request_id", requestID),
		zap.String("category_id", categoryId),
		zap.String("method", ctx.Method()),
		zap.String("url", ctx.OriginalURL()),
	)
	return ctx.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Category updated successfully", resp))
}

func (c *CategoryHandler) DeleteCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals("requestID").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")

	err := c.client.DeleteCategory(ctx.Context(), categoryId)
	if err != nil {
		logger.Log.Error("Failed to delete category",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("method", ctx.Method()),
			zap.String("url", ctx.OriginalURL()),
		)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete category", err))
	}

	logger.Log.Info("Category deleted successfully",
		zap.String("request_id", requestID),
		zap.String("category_id", categoryId),
		zap.String("method", ctx.Method()),
		zap.String("url", ctx.OriginalURL()),
	)
	return ctx.SendStatus(fiber.StatusNoContent)
}

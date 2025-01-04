package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	"context"
	"fmt"
	"strconv"

	"api_gateway/internal/constants"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	client     clients.CategoryClient
	bookClient clients.BookClient
	logger     *logger.Logger
}

func NewCategoryHandler(client clients.CategoryClient, bookClient clients.BookClient, logger *logger.Logger) CategoryHandler {
	return CategoryHandler{
		client:     client,
		bookClient: bookClient,
		logger:     logger,
	}
}

func (c *CategoryHandler) CreateCategoryHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	extra := map[string]interface{}{
		"method": ctx.Method(),
		"url":    ctx.OriginalURL(),
	}

	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse create category request body", extra, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["category_name"] = req.Name

	// Call client to create category
	resp, err := c.client.CreateCategory(context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), req)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create category", extra, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to create category", err))
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category created successfully", extra, nil)
	return ctx.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Category created successfully", resp))
}

func (c *CategoryHandler) GetCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")
	includeBooks := ctx.Query("includeBooks", "false") == "true"

	extra := map[string]interface{}{
		"method":      ctx.Method(),
		"url":         ctx.OriginalURL(),
		"category_id": categoryId,
	}

	// Call client to get category by id
	resp, err := c.client.GetCategory(context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), categoryId)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get category", extra, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to get category", err))
	}

	// If includeBooks query param is true, get books for the category
	if includeBooks {
		books, totalItems, _, err := c.bookClient.GetBooksByCategoryId(
			context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), resp.Id,
			1,
			10,
		)
		if err != nil {
			c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to include books", extra, err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to include books", err))
		}

		resp.SampleBooks = &books
		resp.TotalBooks = &totalItems
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category data fetched successfully", extra, nil)
	return ctx.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess(fmt.Sprintf("Category data with id '%s' fetched successfully", categoryId), resp))
}

func (c *CategoryHandler) GetAllCategoriesHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	includeBooks := ctx.Query("includeBooks", "false") == "true"
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize", "10"))

	extra := map[string]interface{}{
		"method":       ctx.Method(),
		"url":          ctx.OriginalURL(),
		"includeBooks": includeBooks,
		"page":         page,
		"page_size":    pageSize,
	}

	// Call client to get all categories
	categories, totalItems, totalPages, err := c.client.ListCategories(
		context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID),
		1,
		5,
	)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to list categories", extra, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to list categories", err))
	}

	// If includeBooks query param is true, get books for each category
	if includeBooks {
		for i := range categories {
			books, totalItems, _, err := c.bookClient.GetBooksByCategoryId(
				context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), categories[i].Id,
				1,
				10,
			)
			if err != nil {
				c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to include books", extra, err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to include books", err))
			}

			categories[i].SampleBooks = &books
			categories[i].TotalBooks = &totalItems
		}
	}

	extra["categories_count"] = len(categories)
	extra["total_items"] = totalItems
	extra["total_pages"] = totalPages
	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Fetched all categories successfully", extra, nil)
	return ctx.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Category data fetched successfully", map[string]interface{}{
		"categories": categories,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"page_size":   pageSize,
			"totalItems":  totalItems,
			"totalPages":  totalPages,
		},
	}))
}

func (c *CategoryHandler) UpdateCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")
	extra := map[string]interface{}{
		"method":      ctx.Method(),
		"url":         ctx.OriginalURL(),
		"category_id": categoryId,
	}

	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to parse update category request body", extra, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError("Invalid request body", err))
	}

	if errorsMap, err := utils.ValidatePayloads(req); err != nil {
		extra["errors"] = errorsMap
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, constants.ErrValidationMessage, extra, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(datatransfers.ResponseError(constants.ErrValidationMessage, errorsMap))
	}

	extra["category_name"] = req.Name

	// Call client to update category
	resp, err := c.client.UpdateCategory(context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), categoryId, req)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to update category", extra, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to update category", err))
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category updated successfully", extra, nil)
	return ctx.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Category updated successfully", resp))
}

func (c *CategoryHandler) DeleteCategoryByIdHandler(ctx *fiber.Ctx) error {
	// Retrieve requestID from context
	requestID, ok := ctx.Locals(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	categoryId := ctx.Params("id")
	extra := map[string]interface{}{
		"method":      ctx.Method(),
		"url":         ctx.OriginalURL(),
		"category_id": categoryId,
	}

	// Call client to delete category
	err := c.client.DeleteCategory(context.WithValue(ctx.Context(), constants.ContextRequestIDKey, requestID), categoryId)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to delete category", extra, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(datatransfers.ResponseError("Failed to delete category", err))
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category deleted successfully", extra, nil)
	return ctx.SendStatus(fiber.StatusNoContent)
}

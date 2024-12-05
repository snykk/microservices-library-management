package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := c.client.CreateCategory(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create category", err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Category created successfully", resp))
}

func (c *CategoryHandler) GetCategoryByIdHandler(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	includeBooks := ctx.Query("includeBooks", "false") == "true"

	resp, err := c.client.GetCategory(ctx.Context(), categoryId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get category", err))
	}

	if includeBooks {
		books, err := c.bookClient.GetBooksByCategoryId(ctx.Context(), resp.Id)
		if err == nil {
			resp.Books = &books
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Category data with id '%s' fetched successfully", categoryId), resp))
}

func (c *CategoryHandler) GetAllCategoriesHandler(ctx *fiber.Ctx) error {
	includeBooks := ctx.Query("includeBooks", "false") == "true"

	resp, err := c.client.ListCategories(ctx.Context())
	if err != nil {
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

	return ctx.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Category data fetched successfully", resp))
}

func (c *CategoryHandler) UpdateCategoryByIdHandler(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	var req datatransfers.CategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := c.client.UpdateCategory(ctx.Context(), categoryId, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update category", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Category updated successfully", resp))
}

func (c *CategoryHandler) DeleteCategoryByIdHandler(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("id")

	err := c.client.DeleteCategory(ctx.Context(), categoryId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete category", err))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

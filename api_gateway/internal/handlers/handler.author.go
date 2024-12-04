package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthorHandler struct {
	client clients.AuthorClient
}

func NewAuthorHandler(client clients.AuthorClient) AuthorHandler {
	return AuthorHandler{
		client: client,
	}
}

func (a *AuthorHandler) CreateAuthorHandler(c *fiber.Ctx) error {
	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := a.client.CreateAuthor(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create author", err))
	}

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Author created successfully", resp))
}

func (a *AuthorHandler) GetAuthorByIdHandler(c *fiber.Ctx) error {
	authorId := c.Params("id")

	resp, err := a.client.GetAuthor(c.Context(), authorId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get author", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Author data with id '%s' fetched successfully", authorId), resp))
}

func (a *AuthorHandler) GetAllAuthorsHandler(c *fiber.Ctx) error {
	resp, err := a.client.ListAuthors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list authors", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Author data fetched successfully", resp))
}

func (a *AuthorHandler) UpdateAuthorByIdHandler(c *fiber.Ctx) error {
	authorId := c.Params("id")

	var req datatransfers.AuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := a.client.UpdateAuthor(c.Context(), authorId, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update author", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Author updated successfully", resp))
}

func (a *AuthorHandler) DeleteAuthorByIdHandler(c *fiber.Ctx) error {
	authorId := c.Params("id")

	err := a.client.DeleteAuthor(c.Context(), authorId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete author", err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

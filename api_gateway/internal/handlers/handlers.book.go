package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	client clients.BookClient
}

func NewBookHandler(client clients.BookClient) BookHandler {
	return BookHandler{
		client: client,
	}
}

func (b *BookHandler) CreateBookHandler(c *fiber.Ctx) error {
	var req datatransfers.BookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := b.client.CreateBook(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to create book", err))
	}

	return c.Status(fiber.StatusCreated).JSON(datatransfers.ResponseSuccess("Book created successfully", resp))
}

func (b *BookHandler) GetBookByIdHandler(c *fiber.Ctx) error {
	bookId := c.Params("id")

	resp, err := b.client.GetBook(c.Context(), bookId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to get book", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Book data with id '%s' fetched successfully", bookId), resp))
}

func (b *BookHandler) GetAllBooksHandler(c *fiber.Ctx) error {
	resp, err := b.client.ListBooks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to list books", err))
	}

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess("Book data fetched successfully", resp))
}

func (b *BookHandler) UpdateBookByIdHandler(c *fiber.Ctx) error {
	bookId := c.Params("id")

	var req datatransfers.BookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ResponseError("Invalid request body", err))
	}

	resp, err := b.client.UpdateBook(c.Context(), bookId, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to update book", err))
	}

	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("Book updated successfully", resp))
}

func (b *BookHandler) DeleteBookByIdHandler(c *fiber.Ctx) error {
	bookId := c.Params("id")

	err := b.client.DeleteBook(c.Context(), bookId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ResponseError("Failed to delete book", err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

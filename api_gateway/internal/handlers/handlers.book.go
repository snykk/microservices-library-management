package handlers

import (
	"api_gateway/internal/clients"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	resp, err := b.client.GetBook(c.Context(), bookId)
	if err != nil {
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

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Book data with id '%s' fetched successfully", bookId), resp))
}

func (b *BookHandler) GetBooksByAuthorIdHandler(c *fiber.Ctx) error {
	authorId := c.Params("authorId")

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	// Fetch books by author
	resp, err := b.client.GetBooksByAuthorId(c.Context(), authorId)
	if err != nil {
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

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Books by author '%s' fetched successfully", authorId), resp))
}

func (b *BookHandler) GetBooksByCategoryIdHandler(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	// Fetch books by category
	resp, err := b.client.GetBooksByCategoryId(c.Context(), categoryId)
	if err != nil {
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

	return c.Status(fiber.StatusOK).JSON(utils.ResponseSuccess(fmt.Sprintf("Books under category '%s' fetched successfully", categoryId), resp))
}

func (b *BookHandler) GetAllBooksHandler(c *fiber.Ctx) error {
	includeAuthor := c.Query("includeAuthor", "false") == "true"
	includeCategory := c.Query("includeCategory", "false") == "true"

	resp, err := b.client.ListBooks(c.Context())
	if err != nil {
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

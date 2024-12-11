package service

import (
	"book_service/internal/clients"
	"book_service/internal/models"
	"book_service/internal/repository"
	"book_service/pkg/utils"
	"context"
	"fmt"
	"log"
)

type BookService interface {
	CreateBook(ctx context.Context, req *models.BookRequest) (*models.BookRecord, error)
	GetBook(ctx context.Context, id string) (*models.BookRecord, error)
	GetBookByAuthorId(ctx context.Context, authorId string) ([]*models.BookRecord, error)
	GetBookByCategoryId(ctx context.Context, categoryId string) ([]*models.BookRecord, error)
	ListBooks(ctx context.Context) ([]*models.BookRecord, error)
	UpdateBook(ctx context.Context, id string, req *models.BookRequest) (*models.BookRecord, error)
	DeleteBook(ctx context.Context, id string) error
	UpdateBookStock(ctx context.Context, id string, newStock int) error
	IncrementBookStock(ctx context.Context, id string) error
	DecrementBookStock(ctx context.Context, id string) error

	ValidateAuthorExistence(ctx context.Context, authorId string) error
	ValidateCategoryExistence(ctx context.Context, categoryId string) error
}

type bookService struct {
	repo           repository.BookRepository
	authorClient   clients.AuthorClient
	categoryClient clients.CategoryClient
}

func NewBookService(repo repository.BookRepository, authorClient clients.AuthorClient, categoryClient clients.CategoryClient) BookService {
	return &bookService{
		repo:           repo,
		authorClient:   authorClient,
		categoryClient: categoryClient,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req *models.BookRequest) (*models.BookRecord, error) {
	log.Printf("[%s] Creating new book with title: %s\n", utils.GetLocation(), req.Title)

	book := &models.BookRecord{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      req.Stock,
	}

	createdBook, err := s.repo.CreateBook(ctx, book)
	if err != nil {
		log.Printf("[%s] Failed to create book: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Book %s created successfully with ID %s\n", utils.GetLocation(), req.Title, createdBook.Id)
	return createdBook, nil
}

func (s *bookService) GetBook(ctx context.Context, id string) (*models.BookRecord, error) {
	log.Printf("[%s] Fetching book with ID: %s\n", utils.GetLocation(), id)

	book, err := s.repo.GetBook(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to get book with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Book with ID %s fetched successfully\n", utils.GetLocation(), id)
	return book, nil
}

func (s *bookService) GetBookByAuthorId(ctx context.Context, authorId string) ([]*models.BookRecord, error) {
	log.Printf("[%s] Fetching books by author ID: %s\n", utils.GetLocation(), authorId)

	books, err := s.repo.GetBookByAuthorId(ctx, authorId)
	if err != nil {
		log.Printf("[%s] Failed to fetch books by author ID %s: %v\n", utils.GetLocation(), authorId, err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched %d books by author ID %s\n", utils.GetLocation(), len(books), authorId)
	return books, nil
}

func (s *bookService) GetBookByCategoryId(ctx context.Context, categoryId string) ([]*models.BookRecord, error) {
	log.Printf("[%s] Fetching books by category ID: %s\n", utils.GetLocation(), categoryId)

	books, err := s.repo.GetBookByCategoryId(ctx, categoryId)
	if err != nil {
		log.Printf("[%s] Failed to fetch books by category ID %s: %v\n", utils.GetLocation(), categoryId, err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched %d books by category ID %s\n", utils.GetLocation(), len(books), categoryId)
	return books, nil
}

func (s *bookService) ListBooks(ctx context.Context) ([]*models.BookRecord, error) {
	log.Printf("[%s] Fetching list of books\n", utils.GetLocation())

	books, err := s.repo.ListBooks(ctx)
	if err != nil {
		log.Printf("[%s] Failed to list books: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched %d books\n", utils.GetLocation(), len(books))
	return books, nil
}

func (s *bookService) UpdateBook(ctx context.Context, id string, req *models.BookRequest) (*models.BookRecord, error) {
	log.Printf("[%s] Updating book with ID: %s\n", utils.GetLocation(), id)

	book := &models.BookRecord{
		Id:         id,
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	}

	updatedBook, err := s.repo.UpdateBook(ctx, book)
	if err != nil {
		log.Printf("[%s] Failed to update book with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Book with ID %s updated successfully\n", utils.GetLocation(), id)
	return updatedBook, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	log.Printf("[%s] Deleting book with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.DeleteBook(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to delete book with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Book with ID %s deleted successfully\n", utils.GetLocation(), id)
	return nil
}

func (s *bookService) UpdateBookStock(ctx context.Context, id string, newStock int) error {
	log.Printf("[%s] Updating stock for book with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.UpdateBookStock(ctx, id, newStock)
	if err != nil {
		log.Printf("[%s] Failed to update stock for book with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Stock for book with ID %s updated to %d successfully\n", utils.GetLocation(), id, newStock)
	return nil
}

func (s *bookService) IncrementBookStock(ctx context.Context, id string) error {
	log.Printf("[%s] Incrementing stock for book with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.IncrementBookStock(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to increment stock for book with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Stock for book with ID %s incremented successfully\n", utils.GetLocation(), id)
	return nil
}

func (s *bookService) DecrementBookStock(ctx context.Context, id string) error {
	log.Printf("[%s] Decrementing stock for book with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.DecrementBookStock(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to decrement stock for book with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Stock for book with ID %s decremented successfully\n", utils.GetLocation(), id)
	return nil
}

func (s *bookService) ValidateAuthorExistence(ctx context.Context, authorId string) error {
	log.Printf("[%s] Validating author existence with ID: %s\n", utils.GetLocation(), authorId)

	author, _ := s.authorClient.GetAuthor(ctx, authorId)
	if author == nil {
		log.Printf("[%s] Author with ID %s not found\n", utils.GetLocation(), authorId)
		return fmt.Errorf("author with ID '%s' not found", authorId)
	}

	log.Printf("[%s] Author with ID %s exists\n", utils.GetLocation(), authorId)
	return nil
}

func (s *bookService) ValidateCategoryExistence(ctx context.Context, categoryId string) error {
	log.Printf("[%s] Validating category existence with ID: %s\n", utils.GetLocation(), categoryId)

	category, _ := s.categoryClient.GetCategory(ctx, categoryId)
	if category == nil {
		log.Printf("[%s] Category with ID %s not found\n", utils.GetLocation(), categoryId)
		return fmt.Errorf("category with ID '%s' not found", categoryId)
	}

	log.Printf("[%s] Category with ID %s exists\n", utils.GetLocation(), categoryId)
	return nil
}

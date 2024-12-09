package service

import (
	"book_service/internal/clients"
	"book_service/internal/models"
	"book_service/internal/repository"
	"context"
	"fmt"
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
	return s.repo.CreateBook(ctx, &models.BookRecord{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      req.Stock,
	})
}

func (s *bookService) GetBook(ctx context.Context, id string) (*models.BookRecord, error) {
	return s.repo.GetBook(ctx, id)
}

func (s *bookService) GetBookByAuthorId(ctx context.Context, authorId string) ([]*models.BookRecord, error) {
	return s.repo.GetBookByAuthorId(ctx, authorId)
}

func (s *bookService) GetBookByCategoryId(ctx context.Context, categoryId string) ([]*models.BookRecord, error) {
	return s.repo.GetBookByCategoryId(ctx, categoryId)
}

func (s *bookService) ListBooks(ctx context.Context) ([]*models.BookRecord, error) {
	return s.repo.ListBooks(ctx)
}

func (s *bookService) UpdateBook(ctx context.Context, id string, req *models.BookRequest) (*models.BookRecord, error) {
	return s.repo.UpdateBook(ctx, &models.BookRecord{
		Id:         id,
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
}

func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	return s.repo.DeleteBook(ctx, id)
}

func (s *bookService) UpdateBookStock(ctx context.Context, id string, newStock int) error {
	return s.repo.UpdateBookStock(ctx, id, newStock)
}

func (s *bookService) IncrementBookStock(ctx context.Context, id string) error {
	return s.repo.IncrementBookStock(ctx, id)
}

func (s *bookService) DecrementBookStock(ctx context.Context, id string) error {
	return s.repo.DecrementBookStock(ctx, id)
}

func (s *bookService) ValidateAuthorExistence(ctx context.Context, authorId string) error {
	author, _ := s.authorClient.GetAuthor(ctx, authorId)
	if author == nil {
		return fmt.Errorf("author with id '%s' not found", authorId)
	}
	return nil

}

func (s *bookService) ValidateCategoryExistence(ctx context.Context, categoryId string) error {
	category, _ := s.categoryClient.GetCategory(ctx, categoryId)
	if category == nil {
		return fmt.Errorf("category with id '%s' not found", categoryId)
	}
	return nil
}

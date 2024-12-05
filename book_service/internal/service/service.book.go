package service

import (
	"book_service/internal/models"
	"book_service/internal/repository"
	"context"
)

type BookService interface {
	CreateBook(ctx context.Context, req *models.BookRequest) (*models.BookRecord, error)
	GetBook(ctx context.Context, id *string) (*models.BookRecord, error)
	GetBookByAuthorId(ctx context.Context, authorId *string) ([]*models.BookRecord, error)
	GetBookByCategoryId(ctx context.Context, categoryId *string) ([]*models.BookRecord, error)
	ListBooks(ctx context.Context) ([]*models.BookRecord, error)
	UpdateBook(ctx context.Context, id *string, req *models.BookRequest) (*models.BookRecord, error)
	DeleteBook(ctx context.Context, id *string) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req *models.BookRequest) (*models.BookRecord, error) {
	return s.repo.CreateBook(&models.BookRecord{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      req.Stock,
	})
}

func (s *bookService) GetBook(ctx context.Context, id *string) (*models.BookRecord, error) {
	return s.repo.GetBook(id)
}

func (s *bookService) GetBookByAuthorId(ctx context.Context, authorId *string) ([]*models.BookRecord, error) {
	return s.repo.GetBookByAuthorId(authorId)
}

func (s *bookService) GetBookByCategoryId(ctx context.Context, categoryId *string) ([]*models.BookRecord, error) {
	return s.repo.GetBookByCategoryId(categoryId)
}

func (s *bookService) ListBooks(ctx context.Context) ([]*models.BookRecord, error) {
	return s.repo.ListBooks()
}

func (s *bookService) UpdateBook(ctx context.Context, id *string, req *models.BookRequest) (*models.BookRecord, error) {
	return s.repo.UpdateBook(&models.BookRecord{
		Id:         *id,
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
}

func (s *bookService) DeleteBook(ctx context.Context, id *string) error {
	return s.repo.DeleteBook(id)
}

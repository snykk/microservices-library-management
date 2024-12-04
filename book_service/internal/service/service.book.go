package service

import (
	"book_service/internal/models"
	"book_service/internal/repository"
	"context"
)

type BookService interface {
	CreateBook(ctx context.Context, req *models.BookRequest) (*models.BookRecord, error)
	GetBook(ctx context.Context, id *string) (*models.BookRecord, error)
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
	book := models.BookRecord{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      req.Stock,
	}

	createdBook, err := s.repo.CreateBook(&book)
	if err != nil {
		return nil, err
	}
	return createdBook, nil
}

func (s *bookService) GetBook(ctx context.Context, id *string) (*models.BookRecord, error) {
	book, err := s.repo.GetBook(id)
	if err != nil {
		return nil, ErrGetBook
	}
	return book, nil
}

func (s *bookService) ListBooks(ctx context.Context) ([]*models.BookRecord, error) {
	books, err := s.repo.ListBooks()
	if err != nil {
		return nil, ErrGetListBook
	}
	return books, nil
}

func (s *bookService) UpdateBook(ctx context.Context, id *string, req *models.BookRequest) (*models.BookRecord, error) {
	book := models.BookRecord{
		Id:         *id,
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	}

	updatedBook, err := s.repo.UpdateBook(&book)
	if err != nil {
		return nil, ErrUpdateBook
	}
	return updatedBook, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id *string) error {
	err := s.repo.DeleteBook(id)
	if err != nil {
		return ErrDeleteBook
	}
	return nil
}

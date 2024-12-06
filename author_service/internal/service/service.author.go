package service

import (
	"author_service/internal/models"
	"author_service/internal/repository"
	"context"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, req *models.AuthorRequest) (*models.AuthorRecord, error)
	GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error)
	ListAuthors(ctx context.Context) ([]*models.AuthorRecord, error)
	UpdateAuthor(ctx context.Context, id string, req *models.AuthorRequest) (*models.AuthorRecord, error)
	DeleteAuthor(ctx context.Context, id string) error
}

type authorService struct {
	repo repository.AuthorRepository
}

func NewAuthorService(repo repository.AuthorRepository) AuthorService {
	return &authorService{
		repo: repo,
	}
}

func (s *authorService) CreateAuthor(ctx context.Context, req *models.AuthorRequest) (*models.AuthorRecord, error) {
	return s.repo.CreateAuthor(ctx, &models.AuthorRecord{
		Name:      req.Name,
		Biography: req.Biography,
	})
}

func (s *authorService) GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error) {
	return s.repo.GetAuthor(ctx, id)
}

func (s *authorService) ListAuthors(ctx context.Context) ([]*models.AuthorRecord, error) {
	return s.repo.ListAuthors(ctx)
}

func (s *authorService) UpdateAuthor(ctx context.Context, id string, req *models.AuthorRequest) (*models.AuthorRecord, error) {
	return s.repo.UpdateAuthor(ctx, &models.AuthorRecord{
		Id:        id,
		Name:      req.Name,
		Biography: req.Biography,
	})
}

func (s *authorService) DeleteAuthor(ctx context.Context, id string) error {
	return s.repo.DeleteAuthor(ctx, id)
}

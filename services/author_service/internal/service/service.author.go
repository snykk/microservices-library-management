package service

import (
	"author_service/internal/models"
	"author_service/internal/repository"
	"author_service/pkg/utils"
	"context"
	"log"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, req *models.AuthorRequest) (*models.AuthorRecord, error)
	GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error)
	ListAuthors(ctx context.Context, page int, pageSize int) (authors []*models.AuthorRecord, totalItems int, err error)
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

// CreateAuthor creates a new author and returns the created author.
func (s *authorService) CreateAuthor(ctx context.Context, req *models.AuthorRequest) (*models.AuthorRecord, error) {
	log.Printf("[%s] Creating new author with name: %s\n", utils.GetLocation(), req.Name)
	author := &models.AuthorRecord{
		Name:      req.Name,
		Biography: req.Biography,
	}

	createdAuthor, err := s.repo.CreateAuthor(ctx, author)
	if err != nil {
		log.Printf("[%s] Failed to create author: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Author %s created successfully with ID %s\n", utils.GetLocation(), req.Name, createdAuthor.Id)
	return createdAuthor, nil
}

// GetAuthor fetches an author by ID and returns the author.
func (s *authorService) GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error) {
	log.Printf("[%s] Fetching author with ID: %s\n", utils.GetLocation(), id)

	author, err := s.repo.GetAuthor(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to get author with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Author with ID %s fetched successfully\n", utils.GetLocation(), id)
	return author, nil
}

// ListAuthors fetches a list of authors with pagination and returns the authors.
func (s *authorService) ListAuthors(ctx context.Context, page int, pageSize int) (authors []*models.AuthorRecord, totalItems int, err error) {
	log.Printf("[%s] Fetching list of authors with pagination (Page: %d, PageSize: %d)\n", utils.GetLocation(), page, pageSize)

	authors, err = s.repo.ListAuthors(ctx, page, pageSize)
	if err != nil {
		log.Printf("[%s] Failed to list authors: %v\n", utils.GetLocation(), err)
		return nil, 0, err
	}

	totalItems, err = s.repo.CountAuthors(ctx)
	if err != nil {
		log.Printf("[%s] Failed to count authors: %v\n", utils.GetLocation(), err)
		return nil, 0, err
	}

	log.Printf("[%s] Successfully fetched %d authors\n", utils.GetLocation(), len(authors))
	return authors, totalItems, nil
}

// UpdateAuthor updates an author and returns the updated author.
func (s *authorService) UpdateAuthor(ctx context.Context, id string, req *models.AuthorRequest) (*models.AuthorRecord, error) {
	log.Printf("[%s] Updating author with ID: %s\n", utils.GetLocation(), id)

	author := &models.AuthorRecord{
		Id:        id,
		Name:      req.Name,
		Biography: req.Biography,
	}

	updatedAuthor, err := s.repo.UpdateAuthor(ctx, author)
	if err != nil {
		log.Printf("[%s] Failed to update author with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Author with ID %s updated successfully\n", utils.GetLocation(), id)
	return updatedAuthor, nil
}

// DeleteAuthor deletes an author by ID.
func (s *authorService) DeleteAuthor(ctx context.Context, id string) error {
	log.Printf("[%s] Deleting author with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.DeleteAuthor(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to delete author with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Author with ID %s deleted successfully\n", utils.GetLocation(), id)
	return nil
}

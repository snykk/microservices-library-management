package service

import (
	"category_service/internal/models"
	"category_service/internal/repository"
	"context"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.CategoryRecord, error)
	GetCategory(ctx context.Context, id *string) (*models.CategoryRecord, error)
	ListCategories(ctx context.Context) ([]*models.CategoryRecord, error)
	UpdateCategory(ctx context.Context, id *string, req *models.CategoryRequest) (*models.CategoryRecord, error)
	DeleteCategory(ctx context.Context, id *string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.CategoryRecord, error) {
	category := models.CategoryRecord{Name: req.Name}
	return s.repo.CreateCategory(&category)
}

func (s *categoryService) GetCategory(ctx context.Context, id *string) (*models.CategoryRecord, error) {
	return s.repo.GetCategory(id)
}

func (s *categoryService) ListCategories(ctx context.Context) ([]*models.CategoryRecord, error) {
	return s.repo.ListCategories()
}

func (s *categoryService) UpdateCategory(ctx context.Context, id *string, req *models.CategoryRequest) (*models.CategoryRecord, error) {
	return s.repo.UpdateCategory(&models.CategoryRecord{
		Id:   *id,
		Name: req.Name,
	})
}

func (s *categoryService) DeleteCategory(ctx context.Context, id *string) error {
	return s.repo.DeleteCategory(id)
}

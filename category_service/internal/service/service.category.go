package service

import (
	"category_service/internal/models"
	"category_service/internal/repository"
	"category_service/pkg/utils"
	"context"
	"log"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.CategoryRecord, error)
	GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error)
	ListCategories(ctx context.Context, page int, pageSize int) (categories []*models.CategoryRecord, totalItems int, err error)
	UpdateCategory(ctx context.Context, id string, req *models.CategoryRequest) (*models.CategoryRecord, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.CategoryRecord, error) {
	log.Printf("[%s] Creating new category with name: %s\n", utils.GetLocation(), req.Name)
	category := &models.CategoryRecord{Name: req.Name}

	createdCategory, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		log.Printf("[%s] Failed to create category: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Category %s created successfully with ID %s\n", utils.GetLocation(), req.Name, createdCategory.Id)
	return createdCategory, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error) {
	log.Printf("[%s] Fetching category with ID: %s\n", utils.GetLocation(), id)

	category, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to get category with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Category with ID %s fetched successfully\n", utils.GetLocation(), id)
	return category, nil
}

func (s *categoryService) ListCategories(ctx context.Context, page int, pageSize int) (categories []*models.CategoryRecord, totalItems int, err error) {
	log.Printf("[%s] Fetching list of authors with pagination (Page: %d, PageSize: %d)\n", utils.GetLocation(), page, pageSize)

	categories, err = s.repo.ListCategories(ctx)
	if err != nil {
		log.Printf("[%s] Failed to list categories: %v\n", utils.GetLocation(), err)
		return nil, 0, err
	}

	totalItems, err = s.repo.CountCategories(ctx)
	if err != nil {
		log.Printf("[%s] Failed to count authors: %v\n", utils.GetLocation(), err)
		return nil, 0, err
	}

	log.Printf("[%s] Successfully fetched %d categories\n", utils.GetLocation(), len(categories))
	return categories, totalItems, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *models.CategoryRequest) (*models.CategoryRecord, error) {
	log.Printf("[%s] Updating category with ID: %s\n", utils.GetLocation(), id)

	category := &models.CategoryRecord{
		Id:   id,
		Name: req.Name,
	}

	updatedCategory, err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		log.Printf("[%s] Failed to update category with ID %s: %v\n", utils.GetLocation(), id, err)
		return nil, err
	}

	log.Printf("[%s] Category with ID %s updated successfully\n", utils.GetLocation(), id)
	return updatedCategory, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	log.Printf("[%s] Deleting category with ID: %s\n", utils.GetLocation(), id)

	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		log.Printf("[%s] Failed to delete category with ID %s: %v\n", utils.GetLocation(), id, err)
		return err
	}

	log.Printf("[%s] Category with ID %s deleted successfully\n", utils.GetLocation(), id)
	return nil
}

package repository

import (
	"category_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error)
	GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error)
	ListCategories(ctx context.Context) ([]*models.CategoryRecord, error)
	UpdateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error) {
	log.Printf("Executing CreateCategory with name: %s\n", req.Name)
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at`
	category := &models.CategoryRecord{}
	err := r.db.QueryRowContext(ctx, query, req.Name).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		log.Printf("Error in CreateCategory: %v\n", err)
		return nil, err
	}
	log.Printf("Successfully created category with ID: %s\n", category.Id)
	return category, nil
}

func (r *categoryRepository) GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error) {
	log.Printf("Executing GetCategory for ID: %s\n", id)
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	category := &models.CategoryRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Printf("Category with ID %s not found\n", id)
		return nil, errors.New("category not found")
	}
	if err != nil {
		log.Printf("Error in GetCategory for ID %s: %v\n", id, err)
		return nil, err
	}
	log.Printf("Successfully retrieved category with ID: %s\n", category.Id)
	return category, nil
}

func (r *categoryRepository) ListCategories(ctx context.Context) ([]*models.CategoryRecord, error) {
	log.Printf("Executing ListCategories")
	query := `SELECT id, name, created_at, updated_at FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error in ListCategories: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var categories []*models.CategoryRecord
	for rows.Next() {
		category := &models.CategoryRecord{}
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("Error scanning row in ListCategories: %v\n", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows in ListCategories: %v\n", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d categories\n", len(categories))
	return categories, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error) {
	log.Printf("Executing UpdateCategory for ID: %s\n", req.Id)
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name, created_at, updated_at`
	category := &models.CategoryRecord{}
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		log.Printf("Error in UpdateCategory for ID %s: %v\n", req.Id, err)
		return nil, err
	}
	log.Printf("Successfully updated category with ID: %s\n", category.Id)
	return category, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	log.Printf("Executing DeleteCategory for ID: %s\n", id)
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error in DeleteCategory for ID %s: %v\n", id, err)
		return err
	}
	log.Printf("Successfully deleted category with ID: %s\n", id)
	return nil
}

package repository

import (
	"category_service/internal/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error)
	GetCategory(id *string) (*models.CategoryRecord, error)
	ListCategories() ([]*models.CategoryRecord, error)
	UpdateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error)
	DeleteCategory(id *string) error
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at, updated_at`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, req.Name).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	return category, err
}

func (r *categoryRepository) GetCategory(id *string) (*models.CategoryRecord, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, *id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	return category, err
}

func (r *categoryRepository) ListCategories() ([]*models.CategoryRecord, error) {
	query := `SELECT id, name, created_at, updated_at FROM categories`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.CategoryRecord
	for rows.Next() {
		category := &models.CategoryRecord{}
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) UpdateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error) {
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name, created_at, updated_at`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, req.Name, req.Id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	return category, err
}

func (r *categoryRepository) DeleteCategory(id *string) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, *id)
	return err
}

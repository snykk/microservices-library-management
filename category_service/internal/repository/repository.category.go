package repository

import (
	"category_service/internal/models"
	"database/sql"
	"errors"
)

type CategoryRepository interface {
	CreateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error)
	GetCategory(id *string) (*models.CategoryRecord, error)
	ListCategories() ([]*models.CategoryRecord, error)
	UpdateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error)
	DeleteCategory(id *string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, req.Name).Scan(&category.Id, &category.Name)
	return category, err
}

func (r *categoryRepository) GetCategory(id *string) (*models.CategoryRecord, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, *id).Scan(&category.Id, &category.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	return category, err
}

func (r *categoryRepository) ListCategories() ([]*models.CategoryRecord, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.CategoryRecord
	for rows.Next() {
		category := &models.CategoryRecord{}
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) UpdateCategory(req *models.CategoryRecord) (*models.CategoryRecord, error) {
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name`
	category := &models.CategoryRecord{}
	err := r.db.QueryRow(query, req.Name, req.Id).Scan(&category.Id, &category.Name)
	return category, err
}

func (r *categoryRepository) DeleteCategory(id *string) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, *id)
	return err
}

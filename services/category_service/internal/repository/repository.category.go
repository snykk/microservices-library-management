package repository

import (
	"category_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error)
	GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error)
	ListCategories(ctx context.Context) ([]*models.CategoryRecord, error)
	UpdateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error)
	DeleteCategory(ctx context.Context, id string, version int) error
	CountCategories(ctx context.Context) (int, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// CreateCategory inserts a new category into the database and returns the created category.
func (r *categoryRepository) CreateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error) {
	log.Printf("Executing CreateCategory with name: %s\n", req.Name)
	query := `INSERT INTO 
				categories (name) 
			  VALUES 
			  	($1) 
			  RETURNING 
			  	id, name, version, created_at, updated_at`
	category := &models.CategoryRecord{}
	err := r.db.QueryRowContext(ctx, query,
		req.Name,
	).Scan(
		&category.Id,
		&category.Name,
		&category.Version,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error in CreateCategory: %v\n", err)
		return nil, err
	}
	log.Printf("Successfully created category with ID: %s\n", category.Id)
	return category, nil
}

// GetCategory fetches a category by ID and returns the category.
func (r *categoryRepository) GetCategory(ctx context.Context, id string) (*models.CategoryRecord, error) {
	log.Printf("Executing GetCategory for ID: %s\n", id)
	query := `SELECT 
				id, name, version, created_at, updated_at 
			  FROM 
			  	categories 
			  WHERE 
			  	id = $1`
	category := &models.CategoryRecord{}
	err := r.db.QueryRowContext(ctx, query,
		id,
	).Scan(
		&category.Id,
		&category.Name,
		&category.Version,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
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

// ListCategories fetches all categories and returns a list of categories.
func (r *categoryRepository) ListCategories(ctx context.Context) ([]*models.CategoryRecord, error) {
	log.Printf("Executing ListCategories")
	query := `SELECT 
				id, name, version, created_at, updated_at 
			  FROM 
			  	categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error in ListCategories: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var categories []*models.CategoryRecord
	for rows.Next() {
		category := &models.CategoryRecord{}
		if err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Version,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
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

// UpdateCategory updates a category and returns the updated category.
func (r *categoryRepository) UpdateCategory(ctx context.Context, req *models.CategoryRecord) (*models.CategoryRecord, error) {
	log.Printf("Executing UpdateCategory for ID: %s\n", req.Id)

	const maxRetries = 3
	resultChan := make(chan *models.CategoryRecord, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		for attempt := range maxRetries {
			query := `UPDATE 
						categories 
					  SET 
					  	name = $1 
					  WHERE 
					  	id = $2  AND version = $3
					  RETURNING 
					  	id, name, version, created_at, updated_at`
			category := &models.CategoryRecord{}
			err := r.db.QueryRowContext(ctx, query,
				req.Name,
				req.Id,
				req.Version, // Optimistic locking
			).Scan(
				&category.Id,
				&category.Name,
				&category.Version,
				&category.CreatedAt,
				&category.UpdatedAt,
			)
			if err == nil {
				log.Printf("Updated category with ID %s successfully on attempt %d\n", req.Id, attempt+1)
				resultChan <- category
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Optimistic locking failed for book ID %s, retrying... (attempt %d)\n", req.Id, attempt+1)
				time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond) // Exponential backoff

				// Fetch latest category record
				updatedCategory, err := r.GetCategory(ctx, req.Id)
				if err != nil {
					log.Printf("Error fetching latest category record: %v\n", err)
					errorChan <- fmt.Errorf("error updating category with ID %s: %v", req.Id, err) // always makes error be general
					return
				}
				req.Version = updatedCategory.Version
				continue
			}

			errorChan <- err
			return
		}

		errorChan <- fmt.Errorf("error updating category with ID %s: max retries exceeded", req.Id)
	}()

	select {
	case result := <-resultChan:
		log.Printf("Successfully updated category with ID: %s\n", result.Id)
		return result, nil
	case err := <-errorChan:
		log.Printf("Error in UpdateCategory for ID %s: %v\n", req.Id, err)
		return nil, err
	case <-ctx.Done():
		log.Printf("Request timed out while updating category with ID: %s\n", req.Id)
		return nil, ctx.Err()
	}
}

// DeleteCategory deletes a category by ID.
func (r *categoryRepository) DeleteCategory(ctx context.Context, id string, version int) error {
	log.Printf("Executing DeleteCategory for ID: %s\n", id)

	query := `DELETE FROM 
				categories 
			  WHERE 
			  	id = $1 AND version = $2`
	result, err := r.db.ExecContext(ctx, query, id, version)

	if err != nil {
		log.Printf("Error in DeleteCategory for ID %s: %v\n", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for delete operation: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Optimistic locking failed, no rows deleted for category ID %s with version %d\n", id, version)
		return errors.New("delete failed due to wrong id or concurrent modification")
	}

	log.Printf("Deleted category successfully with ID: %s and version: %d\n", id, version)
	return nil
}

// CountCategories counts the total number of categories in the database.
func (r *categoryRepository) CountCategories(ctx context.Context) (int, error) {
	log.Printf("Counting total categories")
	query := `SELECT COUNT(*) FROM categories`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting categories: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

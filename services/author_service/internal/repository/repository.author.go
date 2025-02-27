package repository

import (
	"author_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// AuthorRepository defines the methods that our repository will implement.
type AuthorRepository interface {
	CreateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error)
	GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error)
	ListAuthors(ctx context.Context, page int, pageSize int) ([]*models.AuthorRecord, error)
	UpdateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error)
	DeleteAuthor(ctx context.Context, id string, version int) error
	CountAuthors(ctx context.Context) (int, error)
}

// authorRepository implements the AuthorRepository interface
type authorRepository struct {
	db *sqlx.DB
}

// NewAuthorRepository creates a new instance of authorRepository
func NewAuthorRepository(db *sqlx.DB) AuthorRepository {
	return &authorRepository{
		db: db,
	}
}

// CreateAuthor inserts a new author into the database and returns the created author.
func (r *authorRepository) CreateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error) {
	query := `INSERT INTO 
				authors (name, biography) 
			  VALUES 
			  	($1, $2) 
			  RETURNING 
			  	id, name, biography, created_at, updated_at`

	author := &models.AuthorRecord{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Name,
		req.Biography,
	).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating author: %v\n", err)
		return nil, err
	}

	log.Printf("Author created with ID: %s\n", author.Id)
	return author, nil
}

// GetAuthor retrieves an author by their ID.
func (r *authorRepository) GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error) {
	query := `SELECT 
				id, name, biography, version, created_at, updated_at 
			  FROM 
			  	authors 
			  WHERE 
			  	id = $1`

	author := &models.AuthorRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.Version,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Author with ID %s not found.\n", id)
			return nil, errors.New("author not found")
		}
		log.Printf("Error fetching author by ID %s: %v\n", id, err)
		return nil, err
	}

	log.Printf("Fetched author with ID: %s\n", id)
	return author, nil
}

// ListAuthors retrieves all authors from the database.
func (r *authorRepository) ListAuthors(ctx context.Context, page int, pageSize int) ([]*models.AuthorRecord, error) {
	log.Printf("Listing authors with pagination (Page: %d, PageSize: %d)", page, pageSize)

	offset := (page - 1) * pageSize
	query := `SELECT 
				id, name, biography, version, created_at, updated_at 
			  FROM 
			  	authors 
			  LIMIT 
			  	$1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		log.Printf("Error listing authors: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var authors []*models.AuthorRecord
	for rows.Next() {
		author := &models.AuthorRecord{}
		if err := rows.Scan(
			&author.Id,
			&author.Name,
			&author.Biography,
			&author.Version,
			&author.CreatedAt,
			&author.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning author: %v\n", err)
			return nil, err
		}
		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over authors: %v\n", err)
		return nil, err
	}

	log.Printf("Fetched %d authors\n", len(authors))
	return authors, nil
}

// UpdateAuthor updates an existing author's data.
func (r *authorRepository) UpdateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error) {
	const maxRetries = 3
	resultChan := make(chan *models.AuthorRecord, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		for attempt := range maxRetries {
			query := `UPDATE 
						authors 
					  SET 
					  	name = $1, biography = $2
                      WHERE 
					  	id = $3 AND version = $4 
                      RETURNING 
					  	id, name, biography, created_at, updated_at, version`

			author := &models.AuthorRecord{}
			err := r.db.QueryRowContext(ctx, query, req.Name, req.Biography, req.Id, req.Version).Scan(
				&author.Id,
				&author.Name,
				&author.Biography,
				&author.CreatedAt,
				&author.UpdatedAt,
				&author.Version,
			)

			if err == nil {
				log.Printf("Updated author with ID %s successfully on attempt %d\n", req.Id, attempt+1)
				resultChan <- author
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Optimistic locking failed for author ID %s, retrying... (attempt %d)\n", req.Id, attempt+1)
				time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond) // Exponential backoff

				// Fetch latest author record
				updatedAuthor, err := r.GetAuthor(ctx, req.Id)
				if err != nil {
					log.Printf("Error fetching latest author record: %v\n", err)
					errorChan <- fmt.Errorf("error updating author with ID %s: %v", req.Id, err) // always makes error be general
					return
				}

				// Update request with newest version
				req.Version = updatedAuthor.Version
				continue
			}

			log.Printf("Error updating author with ID %s: %v\n", req.Id, err)
			errorChan <- err
			return
		}

		errorChan <- errors.New("update failed after max retries")
	}()

	select {
	case author := <-resultChan:
		return author, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	}
}

// DeleteAuthor deletes an author from the database by their ID.
func (r *authorRepository) DeleteAuthor(ctx context.Context, id string, version int) error {
	query := `DELETE FROM authors WHERE id = $1 AND version = $2`

	result, err := r.db.ExecContext(ctx, query, id, version)
	if err != nil {
		log.Printf("Error deleting author with ID %s: %v\n", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for delete operation: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Optimistic locking failed, no rows deleted for author ID %s with version %d\n", id, version)
		return errors.New("delete failed due to wrong id or concurrent modification")
	}

	log.Printf("Deleted author with ID: %s and version: %d\n", id, version)
	return nil
}

// CountAuthors counts the total number of authors in the database.
func (r *authorRepository) CountAuthors(ctx context.Context) (int, error) {
	log.Printf("Counting total authors")
	query := `SELECT COUNT(*) FROM authors`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting authors: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

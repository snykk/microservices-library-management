package repository

import (
	"author_service/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type AuthorRepository interface {
	CreateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error)
	GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error)
	ListAuthors(ctx context.Context) ([]*models.AuthorRecord, error)
	UpdateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error)
	DeleteAuthor(ctx context.Context, id string) error
}

type authorRepository struct {
	db *sqlx.DB
}

func NewAuthorRepository(db *sqlx.DB) AuthorRepository {
	return &authorRepository{
		db: db,
	}
}

func (r *authorRepository) CreateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error) {
	query := `INSERT INTO authors (name, biography) VALUES ($1, $2) RETURNING id, name, biography, created_at, updated_at`

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
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) GetAuthor(ctx context.Context, id string) (*models.AuthorRecord, error) {
	query := `SELECT id, name, biography, created_at, updated_at FROM authors WHERE id = $1`

	author := &models.AuthorRecord{}

	if err := r.db.QueryRowContext(ctx, query,
		id,
	).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.CreatedAt,
		&author.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) ListAuthors(ctx context.Context) ([]*models.AuthorRecord, error) {
	query := `SELECT id, name, biography, created_at, updated_at FROM authors`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
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
			&author.CreatedAt,
			&author.UpdatedAt,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *authorRepository) UpdateAuthor(ctx context.Context, req *models.AuthorRecord) (*models.AuthorRecord, error) {
	query := `UPDATE authors SET name = $1, biography = $2 WHERE id = $3 RETURNING id, name, biography, created_at, updated_at`
	row := r.db.QueryRowContext(ctx, query, req.Name, req.Biography, req.Id)

	author := &models.AuthorRecord{}
	if err := row.Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.CreatedAt,
		&author.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) DeleteAuthor(ctx context.Context, id string) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

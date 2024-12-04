package repository

import (
	"author_service/internal/models"
	"database/sql"
	"errors"
)

type AuthorRepository interface {
	CreateAuthor(req *models.AuthorRecord) (*models.AuthorRecord, error)
	GetAuthor(id *string) (*models.AuthorRecord, error)
	ListAuthors() ([]*models.AuthorRecord, error)
	UpdateAuthor(req *models.AuthorRecord) (*models.AuthorRecord, error)
	DeleteAuthor(id *string) error
}

type authorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) AuthorRepository {
	return &authorRepository{
		db: db,
	}
}

func (r *authorRepository) CreateAuthor(req *models.AuthorRecord) (*models.AuthorRecord, error) {
	query := `INSERT INTO authors (name, biography, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name, biography, created_at, updated_at`

	author := &models.AuthorRecord{}
	err := r.db.QueryRow(
		query,
		req.Name,
		req.Biography,
		req.CreatedAt,
		req.UpdatedAt,
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

func (r *authorRepository) GetAuthor(id *string) (*models.AuthorRecord, error) {
	query := `SELECT id, name, biography FROM authors WHERE id = $1`
	row := r.db.QueryRow(query, *id)

	author := &models.AuthorRecord{}
	if err := row.Scan(&author.Id, &author.Name, &author.Biography); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) ListAuthors() ([]*models.AuthorRecord, error) {
	query := `SELECT id, name, biography FROM authors`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*models.AuthorRecord
	for rows.Next() {
		author := &models.AuthorRecord{}
		if err := rows.Scan(&author.Id, &author.Name, &author.Biography); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *authorRepository) UpdateAuthor(req *models.AuthorRecord) (*models.AuthorRecord, error) {
	query := `UPDATE authors SET name = $1, biography = $2 WHERE id = $3 RETURNING id, name, biography`
	row := r.db.QueryRow(query, req.Name, req.Biography, req.Id)

	author := &models.AuthorRecord{}
	if err := row.Scan(&author.Id, &author.Name, &author.Biography); err != nil {
		return nil, err
	}

	return author, nil
}

func (r *authorRepository) DeleteAuthor(id *string) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := r.db.Exec(query, *id)
	return err
}

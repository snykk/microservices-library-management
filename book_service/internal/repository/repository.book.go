package repository

import (
	"book_service/internal/models"
	"database/sql"
	"errors"
)

type BookRepository interface {
	CreateBook(req *models.BookRecord) (*models.BookRecord, error)
	GetBook(id *string) (*models.BookRecord, error)
	ListBooks() ([]*models.BookRecord, error)
	UpdateBook(req *models.BookRecord) (*models.BookRecord, error)
	DeleteBook(id *string) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) CreateBook(req *models.BookRecord) (*models.BookRecord, error) {
	query := `INSERT INTO books (title, author_id, category_id, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author_id, category_id, stock, created_at, updated_at`

	book := &models.BookRecord{}
	err := r.db.QueryRow(
		query,
		req.Title,
		req.AuthorId,
		req.CategoryId,
		req.Stock,
		req.CreatedAt,
		req.UpdatedAt,
	).Scan(
		&book.Id,
		&book.Title,
		&book.AuthorId,
		&book.CategoryId,
		&book.Stock,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) GetBook(id *string) (*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock FROM books WHERE id = $1`
	row := r.db.QueryRow(query, *id)

	book := &models.BookRecord{}
	if err := row.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) ListBooks() ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock FROM books`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord
	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) UpdateBook(req *models.BookRecord) (*models.BookRecord, error) {
	query := `UPDATE books SET title = $1, author_id = $2, category_id = $3, stock = $4 WHERE id = $5 RETURNING id, title, author_id, category_id, stock`
	row := r.db.QueryRow(query, req.Title, req.AuthorId, req.CategoryId, req.Stock, req.Id)

	book := &models.BookRecord{}
	if err := row.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock); err != nil {
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) DeleteBook(id *string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(query, *id)
	return err
}

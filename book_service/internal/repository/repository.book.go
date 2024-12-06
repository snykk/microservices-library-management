package repository

import (
	"book_service/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	CreateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	GetBook(ctx context.Context, id string) (*models.BookRecord, error)
	GetBookByAuthorId(ctx context.Context, authorId string) ([]*models.BookRecord, error)
	GetBookByCategoryId(ctx context.Context, categoryId string) ([]*models.BookRecord, error)
	ListBooks(ctx context.Context) ([]*models.BookRecord, error)
	UpdateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	DeleteBook(ctx context.Context, id string) error
	UpdateBookStock(ctx context.Context, bookId string, newStock int) error
	IncrementBookStock(ctx context.Context, bookId string) error
	DecrementBookStock(ctx context.Context, bookId string) error
}

type bookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) CreateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error) {
	query := `
		INSERT INTO books (title, author_id, category_id, stock)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, author_id, category_id, stock, created_at, updated_at
	`

	book := &models.BookRecord{}

	err := r.db.QueryRowContext(ctx, query,
		req.Title,
		req.AuthorId,
		req.CategoryId,
		req.Stock,
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

func (r *bookRepository) GetBook(ctx context.Context, id string) (*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE id = $1`

	book := &models.BookRecord{}

	err := r.db.QueryRowContext(ctx, query,
		id,
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
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) GetBookByAuthorId(ctx context.Context, authorId string) ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE author_id = $1`
	rows, err := r.db.QueryContext(ctx, query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord

	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.AuthorId,
			&book.CategoryId,
			&book.Stock,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) GetBookByCategoryId(ctx context.Context, categoryId string) ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord

	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.AuthorId,
			&book.CategoryId,
			&book.Stock,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) ListBooks(ctx context.Context) ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord

	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.AuthorId,
			&book.CategoryId,
			&book.Stock,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) UpdateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error) {
	query := `
		UPDATE books 
		SET title = $1, author_id = $2, category_id = $3, stock = $4 
		WHERE id = $5 
		RETURNING id, title, author_id, category_id, stock, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, req.Title, req.AuthorId, req.CategoryId, req.Stock, req.Id)

	book := &models.BookRecord{}

	if err := row.Scan(
		&book.Id,
		&book.Title,
		&book.AuthorId,
		&book.CategoryId,
		&book.Stock,
		&book.CreatedAt,
		&book.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) DeleteBook(ctx context.Context, id string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *bookRepository) UpdateBookStock(ctx context.Context, bookId string, newStock int) error {
	query := `UPDATE books SET stock = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newStock, bookId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no book found with the given Id")
	}

	return nil
}

func (r *bookRepository) IncrementBookStock(ctx context.Context, bookId string) error {
	query := `UPDATE books SET stock = stock + 1 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, bookId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no book found with the given Id")
	}

	return nil
}

func (r *bookRepository) DecrementBookStock(ctx context.Context, bookId string) error {
	query := `UPDATE books SET stock = stock - 1 WHERE id = $1 AND stock > 0`
	result, err := r.db.ExecContext(ctx, query, bookId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no book found with the given Id or insufficient stock")
	}

	return nil
}

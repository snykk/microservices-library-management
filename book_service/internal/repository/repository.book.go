package repository

import (
	"book_service/internal/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	CreateBook(req *models.BookRecord) (*models.BookRecord, error)
	GetBook(id *string) (*models.BookRecord, error)
	GetBookByAuthorId(authorId *string) ([]*models.BookRecord, error)
	GetBookByCategoryId(categoryId *string) ([]*models.BookRecord, error)
	ListBooks() ([]*models.BookRecord, error)
	UpdateBook(req *models.BookRecord) (*models.BookRecord, error)
	DeleteBook(id *string) error
	UpdateBookStock(bookId string, newStock int) error
	IncrementBookStock(bookId string) error
	DecrementBookStock(bookId string) error
}

type bookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) CreateBook(req *models.BookRecord) (*models.BookRecord, error) {
	query := `INSERT INTO books (title, author_id, category_id, stock) VALUES ($1, $2, $3, $4) RETURNING id, title, author_id, category_id, stock, created_at, updated_at`

	book := &models.BookRecord{}
	err := r.db.QueryRow(
		query,
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

func (r *bookRepository) GetBook(id *string) (*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE id = $1`
	row := r.db.QueryRow(query, *id)

	book := &models.BookRecord{}
	if err := row.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) GetBookByAuthorId(authorId *string) ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE author_id = $1`
	rows, err := r.db.Query(query, *authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord
	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) GetBookByCategoryId(categoryId *string) ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE category_id = $1`
	rows, err := r.db.Query(query, *categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord
	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) ListBooks() ([]*models.BookRecord, error) {
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.BookRecord
	for rows.Next() {
		book := &models.BookRecord{}
		if err := rows.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock, &book.CreatedAt, &book.UpdatedAt); err != nil {
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
	query := `UPDATE books SET title = $1, author_id = $2, category_id = $3, stock = $4 WHERE id = $5 RETURNING id, title, author_id, category_id, stock, created_at, updated_at`
	row := r.db.QueryRow(query, req.Title, req.AuthorId, req.CategoryId, req.Stock, req.Id)

	book := &models.BookRecord{}
	if err := row.Scan(&book.Id, &book.Title, &book.AuthorId, &book.CategoryId, &book.Stock, &book.CreatedAt, &book.UpdatedAt); err != nil {
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) DeleteBook(id *string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(query, *id)
	return err
}

func (r *bookRepository) UpdateBookStock(bookId string, newStock int) error {
	query := `UPDATE books SET stock = $1 WHERE id = $2`
	result, err := r.db.Exec(query, newStock, bookId)
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

func (r *bookRepository) IncrementBookStock(bookId string) error {
	query := `UPDATE books SET stock = stock + 1, updated_at = NOW() WHERE id = $1`
	result, err := r.db.Exec(query, bookId)
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

func (r *bookRepository) DecrementBookStock(bookId string) error {
	query := `UPDATE books SET stock = stock - 1, updated_at = NOW() WHERE id = $1 AND stock > 0`
	result, err := r.db.Exec(query, bookId)
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

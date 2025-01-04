package repository

import (
	"book_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	CreateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	GetBook(ctx context.Context, id string) (*models.BookRecord, error)
	GetBookByAuthorId(ctx context.Context, authorId string, page int, pageSize int) ([]*models.BookRecord, error)
	GetBookByCategoryId(ctx context.Context, categoryId string, page int, pageSize int) ([]*models.BookRecord, error)
	ListBooks(ctx context.Context, page int, pageSize int) ([]*models.BookRecord, error)
	UpdateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	DeleteBook(ctx context.Context, id string) error
	UpdateBookStock(ctx context.Context, bookId string, newStock int) error
	IncrementBookStock(ctx context.Context, bookId string) error
	DecrementBookStock(ctx context.Context, bookId string) error
	CountBooks(ctx context.Context) (int, error)
	CountBooksByCategoryId(ctx context.Context, categoryId string) (int, error)
	CountBooksByAuthorId(ctx context.Context, authorid string) (int, error)
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
	log.Printf("Creating book: %+v\n", req)
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
		log.Printf("Error creating book: %v\n", err)
		return nil, err
	}

	log.Printf("Book created successfully: %+v\n", book)
	return book, nil
}

func (r *bookRepository) GetBook(ctx context.Context, id string) (*models.BookRecord, error) {
	log.Printf("Fetching book with ID: %s\n", id)
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE id = $1`
	book := &models.BookRecord{}
	if err := r.db.QueryRowContext(ctx, query,
		id,
	).Scan(
		&book.Id,
		&book.Title,
		&book.AuthorId,
		&book.CategoryId,
		&book.Stock,
		&book.CreatedAt,
		&book.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Book not found with ID: %s\n", id)
			return nil, errors.New("book not found")
		}
		log.Printf("Error fetching book: %v\n", err)
		return nil, err
	}

	log.Printf("Book fetched successfully: %+v\n", book)
	return book, nil
}

func (r *bookRepository) GetBookByAuthorId(ctx context.Context, authorId string, page int, pageSize int) ([]*models.BookRecord, error) {
	log.Printf("Fetching books by author ID: %s with pagination (Page: %d, PageSize: %d)\n", authorId, page, pageSize)

	offset := (page - 1) * pageSize
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE author_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, authorId, pageSize, offset)
	if err != nil {
		log.Printf("Error fetching books by author ID: %v\n", err)
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
			log.Printf("Error scanning book row: %v\n", err)
			return nil, err
		}
		books = append(books, book)
	}

	log.Printf("Books fetched successfully: %d books found for author ID: %s\n", len(books), authorId)
	return books, nil
}

func (r *bookRepository) GetBookByCategoryId(ctx context.Context, categoryId string, page int, pageSize int) ([]*models.BookRecord, error) {
	log.Printf("Fetching books by category ID: %s with pagination (Page: %d, PageSize: %d)\n", categoryId, page, pageSize)

	offset := (page - 1) * pageSize
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books WHERE category_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, categoryId, pageSize, offset)
	if err != nil {
		log.Printf("Error fetching books by category ID: %v\n", err)
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
			log.Printf("Error scanning book row: %v\n", err)
			return nil, err
		}
		books = append(books, book)
	}

	log.Printf("Books fetched successfully: %d books found for category ID: %s\n", len(books), categoryId)
	return books, nil
}

func (r *bookRepository) ListBooks(ctx context.Context, page int, pageSize int) ([]*models.BookRecord, error) {
	log.Printf("Listing books with pagination (Page: %d, PageSize: %d)", page, pageSize)

	offset := (page - 1) * pageSize
	query := `SELECT id, title, author_id, category_id, stock, created_at, updated_at FROM books LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		log.Printf("Error listing books: %v\n", err)
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
			log.Printf("Error scanning book row: %v\n", err)
			return nil, err
		}
		books = append(books, book)
	}

	log.Printf("Books listed successfully: %d books found\n", len(books))
	return books, nil
}

func (r *bookRepository) UpdateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error) {
	log.Printf("Updating book with ID: %s\n", req.Id)
	query := `
		UPDATE books 
		SET title = $1, author_id = $2, category_id = $3, stock = $4 
		WHERE id = $5 
		RETURNING id, title, author_id, category_id, stock, created_at, updated_at
	`

	book := &models.BookRecord{}
	if err := r.db.QueryRowContext(ctx, query,
		req.Title,
		req.AuthorId,
		req.CategoryId,
		req.Stock,
		req.Id,
	).Scan(
		&book.Id,
		&book.Title,
		&book.AuthorId,
		&book.CategoryId,
		&book.Stock,
		&book.CreatedAt,
		&book.UpdatedAt,
	); err != nil {
		log.Printf("Error updating book: %v\n", err)
		return nil, err
	}

	log.Printf("Book updated successfully: %+v\n", book)
	return book, nil
}

func (r *bookRepository) DeleteBook(ctx context.Context, id string) error {
	log.Printf("Deleting book with ID: %s\n", id)
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting book: %v\n", err)
		return err
	}

	log.Printf("Book deleted successfully with ID: %s\n", id)
	return nil
}

func (r *bookRepository) UpdateBookStock(ctx context.Context, bookId string, newStock int) error {
	log.Printf("Updating stock for book ID: %s to %d\n", bookId, newStock)
	query := `UPDATE books SET stock = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newStock, bookId)
	if err != nil {
		log.Printf("Error updating book stock: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No book found with ID: %s\n", bookId)
		return errors.New("no book found with the given ID")
	}

	log.Printf("Stock updated successfully for book ID: %s\n", bookId)
	return nil
}

func (r *bookRepository) IncrementBookStock(ctx context.Context, bookId string) error {
	log.Printf("Incrementing stock for book ID: %s\n", bookId)
	query := `UPDATE books SET stock = stock + 1 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, bookId)
	if err != nil {
		log.Printf("Error incrementing book stock: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No book found with ID: %s\n", bookId)
		return errors.New("no book found with the given ID")
	}

	log.Printf("Stock incremented successfully for book ID: %s\n", bookId)
	return nil
}

func (r *bookRepository) DecrementBookStock(ctx context.Context, bookId string) error {
	log.Printf("Decrementing stock for book ID: %s\n", bookId)
	query := `UPDATE books SET stock = stock - 1 WHERE id = $1 AND stock > 0`
	result, err := r.db.ExecContext(ctx, query, bookId)
	if err != nil {
		log.Printf("Error decrementing book stock: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching affected rows: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No book found with ID: %s or insufficient stock\n", bookId)
		return errors.New("no book found with the given ID or insufficient stock")
	}

	log.Printf("Stock decremented successfully for book ID: %s\n", bookId)
	return nil
}

// CountBooks counts the total number of books in the database.
func (r *bookRepository) CountBooks(ctx context.Context) (int, error) {
	log.Printf("Counting total books")
	query := `SELECT COUNT(*) FROM books`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting books: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

// CountBooksByCategoryId counts the total number of books in the database by category ID.
func (r *bookRepository) CountBooksByCategoryId(ctx context.Context, categoryId string) (int, error) {
	log.Printf("Counting total books by category ID: %s\n", categoryId)
	query := `SELECT COUNT(*) FROM books WHERE category_id = $1`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, categoryId).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting books: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

// CountBooksByAuthorId counts the total number of books in the database by author ID.
func (r *bookRepository) CountBooksByAuthorId(ctx context.Context, authorid string) (int, error) {
	log.Printf("Counting total books by author ID: %s\n", authorid)
	query := `SELECT COUNT(*) FROM books WHERE author_id = $1`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, authorid).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting books: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

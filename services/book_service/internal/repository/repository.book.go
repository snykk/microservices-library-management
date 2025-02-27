package repository

import (
	"book_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	CreateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	GetBook(ctx context.Context, id string) (*models.BookRecord, error)
	GetBookByAuthorId(ctx context.Context, authorId string, page int, pageSize int) ([]*models.BookRecord, error)
	GetBookByCategoryId(ctx context.Context, categoryId string, page int, pageSize int) ([]*models.BookRecord, error)
	ListBooks(ctx context.Context, page int, pageSize int) ([]*models.BookRecord, error)
	UpdateBook(ctx context.Context, req *models.BookRecord) (*models.BookRecord, error)
	DeleteBook(ctx context.Context, id string, version int) error
	UpdateBookStock(ctx context.Context, bookId string, newStock, version int) error
	IncrementBookStock(ctx context.Context, bookId string, version int) error
	DecrementBookStock(ctx context.Context, bookId string, version int) error
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
		INSERT INTO 
			books (title, author_id, category_id, stock)
		VALUES 
			($1, $2, $3, $4)
		RETURNING 
			id, title, author_id, category_id, stock, created_at, updated_at
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
	query := `SELECT 
				id, title, author_id, category_id, stock, version, created_at, updated_at 
			  FROM 
			  	books 
			  WHERE 
			  	id = $1`
	book := &models.BookRecord{}
	if err := r.db.QueryRowContext(ctx, query,
		id,
	).Scan(
		&book.Id,
		&book.Title,
		&book.AuthorId,
		&book.CategoryId,
		&book.Stock,
		&book.Version,
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
	query := `SELECT 
				id, title, author_id, category_id, stock, version, created_at, updated_at 
			  FROM 
			  	books 
			  WHERE 
			  	author_id = $1 
			  LIMIT 
			  	$2 OFFSET $3`
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
			&book.Version,
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
	query := `SELECT 
				id, title, author_id, category_id, stock, version, created_at, updated_at 
			  FROM 
			  	books 
			  WHERE 
			  	category_id = $1 
			  LIMIT 
			  	$2 OFFSET $3`
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
			&book.Version,
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
	query := `SELECT 
				id, title, author_id, category_id, stock, version, created_at, updated_at 
			  FROM 
			  	books 
			  LIMIT 
			  	$1 OFFSET $2`
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
			&book.Version,
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

	const maxRetries = 3
	resultCh := make(chan *models.BookRecord, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(resultCh)
		defer close(errCh)

		for attempt := range maxRetries {
			query := `
				UPDATE 
					books 
				SET 
					title = $1, author_id = $2, category_id = $3, stock = $4
				WHERE 
					id = $5 AND version = $6
				RETURNING 
					id, title, author_id, category_id, stock, version, created_at, updated_at
			`

			book := &models.BookRecord{}
			err := r.db.QueryRowContext(ctx, query,
				req.Title,
				req.AuthorId,
				req.CategoryId,
				req.Stock,
				req.Id,
				req.Version, // Optimistic locking check
			).Scan(
				&book.Id,
				&book.Title,
				&book.AuthorId,
				&book.CategoryId,
				&book.Stock,
				&book.Version,
				&book.CreatedAt,
				&book.UpdatedAt,
			)

			if err == nil {
				log.Printf("Updated book with ID %s successfully on attempt %d\n", req.Id, attempt+1)
				resultCh <- book
				return
			}

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Optimistic locking failed for book ID %s, retrying... (attempt %d)\n", req.Id, attempt+1)
				time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond) // Exponential backoff

				// Fetch latest book record
				updatedBook, err := r.GetBook(ctx, req.Id)
				if err != nil {
					log.Printf("Error fetching latest book record: %v\n", err)
					errCh <- fmt.Errorf("error updating book with ID %s: %v", req.Id, err) // always makes error be general
					return
				}
				req.Version = updatedBook.Version
				continue
			}

			errCh <- err
			return
		}

		errCh <- errors.New("update failed after max retries")
	}()

	select {
	case res := <-resultCh:
		log.Printf("Book updated successfully with ID: %s\n", res.Id)
		return res, nil
	case err := <-errCh:
		log.Printf("Error updating book with ID %s: %v\n", req.Id, err)
		return nil, err
	case <-ctx.Done():
		log.Printf("Request timed out while updating book with ID: %s\n", req.Id)
		return nil, ctx.Err()
	}
}

func (r *bookRepository) DeleteBook(ctx context.Context, id string, version int) error {
	log.Printf("Deleting book with ID: %s and version: %d\n", id, version)

	query := `DELETE 
		 	  FROM 
			  	books 
			  WHERE 
			  	id = $1 AND version = $2`
	result, err := r.db.ExecContext(ctx, query, id, version)

	if err != nil {
		log.Printf("Error deleting book with ID %s: %v\n", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected for delete operation: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Optimistic locking failed, no rows deleted for book ID %s with version %d\n", id, version)
		return errors.New("delete failed due to wrong id or concurrent modification")
	}

	log.Printf("Deleted book successfully with ID: %s and version: %d\n", id, version)
	return nil
}

func (r *bookRepository) UpdateBookStock(ctx context.Context, bookId string, newStock int, version int) error {
	log.Printf("Updating stock for book ID: %s\n", bookId)

	const maxRetries = 3
	doneCh := make(chan bool, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(doneCh)
		defer close(errCh)

		for attempt := range maxRetries {
			query := `
				UPDATE 
					books 
				SET 
					stock = $1
				WHERE 
					id = $2 AND version = $3
				RETURNING 
					version
			`

			var newVersion int
			err := r.db.QueryRowContext(ctx, query, newStock, bookId, version).Scan(&newVersion)

			if err == nil {
				log.Printf("Updated stock for book ID %s successfully on attempt %d\n", bookId, attempt+1)
				doneCh <- true
				return
			}

			// Optimistic locking check
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Optimistic locking failed for book ID %s, retrying... (attempt %d)\n", bookId, attempt+1)
				time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond) // Exponential backoff

				// Fetch latest book version
				updatedBook, fetchErr := r.GetBook(ctx, bookId)
				if fetchErr != nil {
					log.Printf("Error fetching latest book record: %v\n", fetchErr)
					errCh <- fmt.Errorf("error updating stock for book ID %s", bookId)
					return
				}
				version = updatedBook.Version
				continue
			}

			errCh <- err
			return
		}

		errCh <- errors.New("update stock failed after max retries")
	}()

	select {
	case <-doneCh:
		log.Printf("Stock updated successfully for book ID: %s\n", bookId)
		return nil
	case err := <-errCh:
		log.Printf("Error updating stock for book ID %s: %v\n", bookId, err)
		return err
	case <-ctx.Done():
		log.Printf("Request timed out while updating stock for book ID: %s\n", bookId)
		return ctx.Err()
	}
}

func (r *bookRepository) IncrementBookStock(ctx context.Context, bookId string, version int) error {
	log.Printf("Incrementing stock for book ID: %s\n", bookId)
	updatedBook, err := r.GetBook(ctx, bookId)

	if err != nil {
		log.Printf("Error fetching book: %v\n", err)
		return fmt.Errorf("error fetching book: %v", err)
	}

	newStock := updatedBook.Stock + 1
	err = r.UpdateBookStock(ctx, bookId, newStock, version)

	if err != nil {
		log.Printf("Error updating stock: %v\n", err)
		return fmt.Errorf("error updating stock: %v", err)
	}

	log.Printf("Stock incremented successfully for book ID: %s\n", bookId)
	return nil
}

func (r *bookRepository) DecrementBookStock(ctx context.Context, bookId string, version int) error {
	log.Printf("Decrementing stock for book ID: %s\n", bookId)
	updatedBook, err := r.GetBook(ctx, bookId)
	if err != nil {
		log.Printf("Error fetching book: %v\n", err)
		return fmt.Errorf("error fetching book: %v", err)
	}

	// Ensure stock doesn't go below zero
	if updatedBook.Stock <= 0 {
		log.Printf("Stock cannot be negative for book ID %s\n", bookId)
		return fmt.Errorf("stock cannot be negative for book ID %s", bookId)
	}

	newStock := updatedBook.Stock - 1
	err = r.UpdateBookStock(ctx, bookId, newStock, version)
	if err != nil {
		log.Printf("Error updating stock: %v\n", err)
		return fmt.Errorf("error updating stock: %v", err)
	}

	log.Printf("Stock decremented successfully for book ID: %s\n", bookId)
	return nil
}

// CountBooks counts the total number of books in the database.
func (r *bookRepository) CountBooks(ctx context.Context) (int, error) {
	log.Printf("Counting total books")
	query := `SELECT 
				COUNT(*) 
			  FROM 
			  	books`
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
	query := `SELECT 
				COUNT(*) 
			  FROM 
			  	books 
			  WHERE 
			  	category_id = $1`
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
	query := `SELECT 
				COUNT(*) 
			  FROM 
			  	books 
			  WHERE 
			  	author_id = $1`
	var totalItems int
	err := r.db.QueryRowContext(ctx, query, authorid).Scan(&totalItems)
	if err != nil {
		log.Printf("Error counting books: %v\n", err)
		return 0, err
	}
	return totalItems, nil
}

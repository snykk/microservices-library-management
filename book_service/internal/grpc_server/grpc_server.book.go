package grpc_server

import (
	"book_service/internal/clients"
	"book_service/internal/models"
	"book_service/internal/service"
	protoBook "book_service/proto/book_service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bookGRPCServer struct {
	bookService    service.BookService
	authorClient   clients.AuthorClient
	categoryClient clients.CategoryClient
	protoBook.UnimplementedBookServiceServer
}

func NewBookGRPCServer(bookService service.BookService, authorClient clients.AuthorClient, categoryClient clients.CategoryClient) protoBook.BookServiceServer {
	return &bookGRPCServer{
		bookService:    bookService,
		authorClient:   authorClient,
		categoryClient: categoryClient,
	}
}

func (s *bookGRPCServer) CreateBook(ctx context.Context, req *protoBook.CreateBookRequest) (*protoBook.CreateBookResponse, error) {
	// check author existence
	author, _ := s.authorClient.GetAuthor(ctx, req.AuthorId)
	if author == nil {
		return nil, status.Error(codes.NotFound, "author not found")
	}

	// check category existence
	category, _ := s.categoryClient.GetCategory(ctx, req.CategoryId)
	if category == nil {
		return nil, status.Error(codes.NotFound, "category not found")
	}

	createdBook, err := s.bookService.CreateBook(ctx, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create new book")
	}

	return &protoBook.CreateBookResponse{
		Book: &protoBook.Book{
			Id:         createdBook.Id,
			Title:      createdBook.Title,
			AuthorId:   createdBook.AuthorId,
			CategoryId: createdBook.CategoryId,
			Stock:      int32(createdBook.Stock),
			CreatedAt:  createdBook.CreatedAt.Unix(),
			UpdatedAt:  createdBook.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) GetBooksByAuthor(ctx context.Context, req *protoBook.GetBooksByAuthorRequest) (*protoBook.ListBooksResponse, error) {
	books, err := s.bookService.GetBookByAuthorId(ctx, req.AuthorId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve books by author id")
	}

	var protoBooks []*protoBook.Book
	for _, book := range books {
		protoBooks = append(protoBooks, &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	return &protoBook.ListBooksResponse{
		Books: protoBooks,
	}, nil
}

func (s *bookGRPCServer) GetBooksByCategory(ctx context.Context, req *protoBook.GetBooksByCategoryRequest) (*protoBook.ListBooksResponse, error) {
	books, err := s.bookService.GetBookByCategoryId(ctx, req.CategoryId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve books by category id")
	}

	var protoBooks []*protoBook.Book
	for _, book := range books {
		protoBooks = append(protoBooks, &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	return &protoBook.ListBooksResponse{
		Books: protoBooks,
	}, nil
}

func (s *bookGRPCServer) GetBook(ctx context.Context, req *protoBook.GetBookRequest) (*protoBook.GetBookResponse, error) {
	book, err := s.bookService.GetBook(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve book data with id '%s'", req.Id))
	}

	return &protoBook.GetBookResponse{
		Book: &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) ListBooks(ctx context.Context, req *protoBook.ListBooksRequest) (*protoBook.ListBooksResponse, error) {
	books, err := s.bookService.ListBooks(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve book list")
	}

	var protoBooks []*protoBook.Book
	for _, book := range books {
		protoBooks = append(protoBooks, &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	return &protoBook.ListBooksResponse{
		Books: protoBooks,
	}, nil
}

func (s *bookGRPCServer) UpdateBook(ctx context.Context, req *protoBook.UpdateBookRequest) (*protoBook.UpdateBookResponse, error) {
	// check author existence
	author, _ := s.authorClient.GetAuthor(ctx, req.AuthorId)
	if author == nil {
		return nil, status.Error(codes.NotFound, "author not found")
	}

	// check category existence
	category, _ := s.categoryClient.GetCategory(ctx, req.CategoryId)
	if category == nil {
		return nil, status.Error(codes.NotFound, "category not found")
	}

	updatedBook, err := s.bookService.UpdateBook(ctx, req.Id, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update book data with id '%s'", req.Id))
	}

	return &protoBook.UpdateBookResponse{
		Book: &protoBook.Book{
			Id:         updatedBook.Id,
			Title:      updatedBook.Title,
			AuthorId:   updatedBook.AuthorId,
			CategoryId: updatedBook.CategoryId,
			Stock:      int32(updatedBook.Stock),
			CreatedAt:  updatedBook.CreatedAt.Unix(),
			UpdatedAt:  updatedBook.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) DeleteBook(ctx context.Context, req *protoBook.DeleteBookRequest) (*protoBook.DeleteBookResponse, error) {
	err := s.bookService.DeleteBook(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete book dta with id '%s'", req.Id))
	}

	return &protoBook.DeleteBookResponse{Message: fmt.Sprintf("success delete book with id %s", req.Id)}, nil
}

func (s *bookGRPCServer) UpdateBookStock(ctx context.Context, req *protoBook.UpdateBookStockRequest) (*protoBook.UpdateBookStockResponse, error) {
	err := s.bookService.UpdateBookStock(ctx, req.Id, int(req.Stock))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update book data with id '%s'", req.Id))
	}

	return &protoBook.UpdateBookStockResponse{
		Message: fmt.Sprintf("stock book with id %s updates successfully", req.Id),
	}, nil
}

func (s *bookGRPCServer) IncrementBookStock(ctx context.Context, req *protoBook.IncrementBookStockRequest) (*protoBook.IncrementBookStockResponse, error) {
	err := s.bookService.IncrementBookStock(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to increment book data with id '%s'", req.Id))
	}

	return &protoBook.IncrementBookStockResponse{
		Message: fmt.Sprintf("stock book with id %s increments successfully", req.Id),
	}, nil
}

func (s *bookGRPCServer) DecrementBookStock(ctx context.Context, req *protoBook.DecrementBookStockRequest) (*protoBook.DecrementBookStockResponse, error) {
	err := s.bookService.DecrementBookStock(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to decrement book data with id '%s'", req.Id))
	}

	return &protoBook.DecrementBookStockResponse{
		Message: fmt.Sprintf("stock book with id %s increments successfully", req.Id),
	}, nil
}

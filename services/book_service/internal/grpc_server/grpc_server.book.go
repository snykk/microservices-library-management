package grpc_server

import (
	"book_service/internal/constants"
	"book_service/internal/models"
	"book_service/internal/service"
	"book_service/pkg/logger"
	protoBook "book_service/proto/book_service"
	"context"
	"fmt"

	"book_service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bookGRPCServer struct {
	bookService service.BookService
	logger      *logger.Logger
	protoBook.UnimplementedBookServiceServer
}

func NewBookGRPCServer(bookService service.BookService, logger *logger.Logger) protoBook.BookServiceServer {
	return &bookGRPCServer{
		bookService: bookService,
		logger:      logger,
	}
}

func (s *bookGRPCServer) CreateBook(ctx context.Context, req *protoBook.CreateBookRequest) (*protoBook.CreateBookResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received CreateBook request", map[string]interface{}{"title": req.Title, "author_id": req.AuthorId, "category_id": req.CategoryId}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid CreateBook request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// check author existence
	err := s.bookService.ValidateAuthorExistence(ctx, req.AuthorId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Author not found", nil, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// check category existence
	err = s.bookService.ValidateCategoryExistence(ctx, req.CategoryId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Category not found", nil, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	createdBook, err := s.bookService.CreateBook(ctx, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create book", nil, err)
		return nil, status.Error(codes.Internal, "failed to create new book")
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book created successfully", map[string]interface{}{"book_id": createdBook.Id}, nil)

	return &protoBook.CreateBookResponse{
		Book: &protoBook.Book{
			Id:         createdBook.Id,
			Title:      createdBook.Title,
			AuthorId:   createdBook.AuthorId,
			CategoryId: createdBook.CategoryId,
			Stock:      int32(createdBook.Stock),
			Version:    int32(createdBook.Version),
			CreatedAt:  createdBook.CreatedAt.Unix(),
			UpdatedAt:  createdBook.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) GetBooksByAuthor(ctx context.Context, req *protoBook.GetBooksByAuthorRequest) (*protoBook.ListBooksResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetBooksByAuthor request", map[string]interface{}{"author_id": req.AuthorId}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetBooksByAuthor request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	books, totalItems, err := s.bookService.GetBookByAuthorId(ctx, req.AuthorId, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve books by author id", nil, err)
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
			Version:    int32(book.Version),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Books retrieved successfully by author", nil, nil)

	return &protoBook.ListBooksResponse{
		Books:      protoBooks,
		TotalItems: int32(totalItems),
		TotalPages: int32(utils.CalculateTotalPages(totalItems, int(req.PageSize))),
	}, nil
}

func (s *bookGRPCServer) GetBooksByCategory(ctx context.Context, req *protoBook.GetBooksByCategoryRequest) (*protoBook.ListBooksResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetBooksByCategory request", map[string]interface{}{"category_id": req.CategoryId}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetBooksByCategory request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	books, totalItems, err := s.bookService.GetBookByCategoryId(ctx, req.CategoryId, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve books by category id", nil, err)
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
			Version:    int32(book.Version),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Books retrieved successfully by category", nil, nil)

	return &protoBook.ListBooksResponse{
		Books:      protoBooks,
		TotalItems: int32(totalItems),
		TotalPages: int32(utils.CalculateTotalPages(totalItems, int(req.PageSize))),
	}, nil
}

func (s *bookGRPCServer) GetBook(ctx context.Context, req *protoBook.GetBookRequest) (*protoBook.GetBookResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetBook request", map[string]interface{}{"book_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetBook request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	book, err := s.bookService.GetBook(ctx, req.Id)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to retrieve book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book retrieved successfully", map[string]interface{}{"book_id": book.Id}, nil)

	return &protoBook.GetBookResponse{
		Book: &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
			Version:    int32(book.Version),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) ListBooks(ctx context.Context, req *protoBook.ListBooksRequest) (*protoBook.ListBooksResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListBooks request", nil, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListBooks request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Retrieve books list
	books, totalItems, err := s.bookService.ListBooks(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve books list", nil, err)
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
			Version:    int32(book.Version),
			CreatedAt:  book.CreatedAt.Unix(),
			UpdatedAt:  book.UpdatedAt.Unix(),
		})
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Books list retrieved successfully", nil, nil)

	return &protoBook.ListBooksResponse{
		Books:      protoBooks,
		TotalItems: int32(totalItems),
		TotalPages: int32(utils.CalculateTotalPages(totalItems, int(req.PageSize))),
	}, nil
}

func (s *bookGRPCServer) UpdateBook(ctx context.Context, req *protoBook.UpdateBookRequest) (*protoBook.UpdateBookResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received UpdateBook request", map[string]interface{}{"book_id": req.Id, "title": req.Title, "author_id": req.AuthorId, "category_id": req.CategoryId}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid UpdateBook request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// check author existence
	err := s.bookService.ValidateAuthorExistence(ctx, req.AuthorId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Author not found", nil, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// check category existence
	err = s.bookService.ValidateCategoryExistence(ctx, req.CategoryId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Category not found", nil, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	updatedBook, err := s.bookService.UpdateBook(ctx, req.Id, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
		Version:    int(req.Version),
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to update book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book updated successfully", map[string]interface{}{"book_id": updatedBook.Id}, nil)

	return &protoBook.UpdateBookResponse{
		Book: &protoBook.Book{
			Id:         updatedBook.Id,
			Title:      updatedBook.Title,
			AuthorId:   updatedBook.AuthorId,
			CategoryId: updatedBook.CategoryId,
			Stock:      int32(updatedBook.Stock),
			Version:    int32(updatedBook.Version),
			CreatedAt:  updatedBook.CreatedAt.Unix(),
			UpdatedAt:  updatedBook.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *bookGRPCServer) DeleteBook(ctx context.Context, req *protoBook.DeleteBookRequest) (*protoBook.DeleteBookResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received DeleteBook request", map[string]interface{}{"book_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid DeleteBook request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.bookService.DeleteBook(ctx, req.Id, int(req.Version))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to delete book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book deleted successfully", map[string]interface{}{"book_id": req.Id}, nil)

	return &protoBook.DeleteBookResponse{
		Message: fmt.Sprintf("Book with ID '%s' deleted successfully", req.Id),
	}, nil
}

func (s *bookGRPCServer) UpdateBookStock(ctx context.Context, req *protoBook.UpdateBookStockRequest) (*protoBook.UpdateBookStockResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received UpdateBookStock request", map[string]interface{}{"book_id": req.Id, "new_stock": req.Stock}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid UpdateBookStock request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.bookService.UpdateBookStock(ctx, req.Id, int(req.Stock), int(req.Version))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to update stock for book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update stock for book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book stock updated successfully", map[string]interface{}{"book_id": req.Id, "new_stock": req.Stock}, nil)

	return &protoBook.UpdateBookStockResponse{
		Message: fmt.Sprintf("Stock for book with id '%s' updated successfully", req.Id),
	}, nil
}

func (s *bookGRPCServer) IncrementBookStock(ctx context.Context, req *protoBook.IncrementBookStockRequest) (*protoBook.IncrementBookStockResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received IncrementBookStock request", map[string]interface{}{"book_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid IncrementBookStock request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.bookService.IncrementBookStock(ctx, req.Id, int(req.Version))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to increment stock for book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to increment stock for book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book stock incremented successfully", map[string]interface{}{"book_id": req.Id}, nil)

	return &protoBook.IncrementBookStockResponse{
		Message: fmt.Sprintf("Stock for book with id '%s' incremented successfully", req.Id),
	}, nil
}

func (s *bookGRPCServer) DecrementBookStock(ctx context.Context, req *protoBook.DecrementBookStockRequest) (*protoBook.DecrementBookStockResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received DecrementBookStock request", map[string]interface{}{"book_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid DecrementBookStock request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.bookService.DecrementBookStock(ctx, req.Id, int(req.Version))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to decrement stock for book with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to decrement stock for book with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Book stock decremented successfully", map[string]interface{}{"book_id": req.Id}, nil)

	return &protoBook.DecrementBookStockResponse{
		Message: fmt.Sprintf("Stock for book with id '%s' decremented successfully", req.Id),
	}, nil
}

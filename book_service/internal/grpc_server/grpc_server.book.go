package grpc_server

import (
	"book_service/internal/exception"
	"book_service/internal/models"
	"book_service/internal/service"
	protoBook "book_service/proto"
	"context"
	"fmt"
)

type bookGRPCServer struct {
	bookService service.BookService
	protoBook.UnimplementedBookServiceServer
}

func NewBookGRPCServer(bookService service.BookService) protoBook.BookServiceServer {
	return &bookGRPCServer{
		bookService: bookService,
	}
}

func (s *bookGRPCServer) CreateBook(ctx context.Context, req *protoBook.CreateBookRequest) (*protoBook.CreateBookResponse, error) {
	book, err := s.bookService.CreateBook(ctx, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
	if err != nil {
		return nil, exception.GRPCErrorFormatter(err)
	}

	return &protoBook.CreateBookResponse{
		Book: &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
		},
	}, nil
}

func (s *bookGRPCServer) GetBook(ctx context.Context, req *protoBook.GetBookRequest) (*protoBook.GetBookResponse, error) {
	book, err := s.bookService.GetBook(ctx, &req.Id)
	if err != nil {
		return nil, exception.GRPCErrorFormatter(err)
	}

	return &protoBook.GetBookResponse{
		Book: &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
		},
	}, nil
}

func (s *bookGRPCServer) ListBooks(ctx context.Context, req *protoBook.ListBooksRequest) (*protoBook.ListBooksResponse, error) {
	books, err := s.bookService.ListBooks(ctx)
	if err != nil {
		return nil, exception.GRPCErrorFormatter(err)
	}

	var protoBooks []*protoBook.Book
	for _, book := range books {
		protoBooks = append(protoBooks, &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
		})
	}

	return &protoBook.ListBooksResponse{
		Books: protoBooks,
	}, nil
}

func (s *bookGRPCServer) UpdateBook(ctx context.Context, req *protoBook.UpdateBookRequest) (*protoBook.UpdateBookResponse, error) {
	book, err := s.bookService.UpdateBook(ctx, &req.Id, &models.BookRequest{
		Title:      req.Title,
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Stock:      int(req.Stock),
	})
	if err != nil {
		return nil, exception.GRPCErrorFormatter(err)
	}

	return &protoBook.UpdateBookResponse{
		Book: &protoBook.Book{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int32(book.Stock),
		},
	}, nil
}

func (s *bookGRPCServer) DeleteBook(ctx context.Context, req *protoBook.DeleteBookRequest) (*protoBook.DeleteBookResponse, error) {
	err := s.bookService.DeleteBook(ctx, &req.Id)
	if err != nil {
		return nil, exception.GRPCErrorFormatter(err)
	}

	return &protoBook.DeleteBookResponse{Message: fmt.Sprintf("success delete book with id %s", req.Id)}, nil
}

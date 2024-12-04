package clients

import (
	"api_gateway/internal/datatransfers"
	protoBook "api_gateway/proto/book_service"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookClient interface {
	CreateBook(ctx context.Context, dto datatransfers.BookRequest) (datatransfers.BookResponse, error)
	GetBook(ctx context.Context, id string) (datatransfers.BookResponse, error)
	ListBooks(ctx context.Context) ([]datatransfers.BookResponse, error)
	UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookRequest) (datatransfers.BookResponse, error)
	DeleteBook(ctx context.Context, id string) error
}

type bookClient struct {
	client protoBook.BookServiceClient
}

func NewBookClient() (BookClient, error) {
	conn, err := grpc.Dial("book-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoBook.NewBookServiceClient(conn)
	return &bookClient{
		client: client,
	}, nil
}

func (b *bookClient) CreateBook(ctx context.Context, dto datatransfers.BookRequest) (datatransfers.BookResponse, error) {
	reqProto := protoBook.CreateBookRequest{
		Title:      dto.Title,
		AuthorId:   dto.AuthorId,
		CategoryId: dto.CategoryId,
		Stock:      int32(dto.Stock),
	}

	resp, err := b.client.CreateBook(ctx, &reqProto)
	if err != nil {
		return datatransfers.BookResponse{}, err
	}

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   resp.Book.AuthorId,
		CategoryId: resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
	}, nil
}

func (b *bookClient) GetBook(ctx context.Context, id string) (datatransfers.BookResponse, error) {
	reqProto := protoBook.GetBookRequest{
		Id: id,
	}

	resp, err := b.client.GetBook(ctx, &reqProto)
	if err != nil {
		return datatransfers.BookResponse{}, err
	}

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   resp.Book.AuthorId,
		CategoryId: resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
	}, nil
}

func (b *bookClient) ListBooks(ctx context.Context) ([]datatransfers.BookResponse, error) {
	reqProto := protoBook.ListBooksRequest{}

	resp, err := b.client.ListBooks(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   book.AuthorId,
			CategoryId: book.CategoryId,
			Stock:      int(book.Stock),
		})
	}

	return books, nil
}

func (b *bookClient) UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookRequest) (datatransfers.BookResponse, error) {
	reqProto := protoBook.UpdateBookRequest{
		Id:         bookId,
		Title:      dto.Title,
		AuthorId:   dto.AuthorId,
		CategoryId: dto.CategoryId,
		Stock:      int32(dto.Stock),
	}

	resp, err := b.client.UpdateBook(ctx, &reqProto)
	if err != nil {
		return datatransfers.BookResponse{}, err
	}

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   resp.Book.AuthorId,
		CategoryId: resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
	}, nil
}

func (b *bookClient) DeleteBook(ctx context.Context, id string) error {
	reqProto := protoBook.DeleteBookRequest{
		Id: id,
	}

	_, err := b.client.DeleteBook(ctx, &reqProto)
	if err != nil {
		return err
	}

	return nil
}

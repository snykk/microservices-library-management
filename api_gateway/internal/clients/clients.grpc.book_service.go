package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoBook "api_gateway/proto/book_service"
	"context"
	"time"

	"api_gateway/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookClient interface {
	CreateBook(ctx context.Context, dto datatransfers.BookRequest) (datatransfers.BookResponse, error)
	GetBook(ctx context.Context, id string) (datatransfers.BookResponse, error)
	GetBooksByAuthorId(ctx context.Context, authorId string) ([]datatransfers.BookResponse, error)
	GetBooksByCategoryId(ctx context.Context, categoryId string) ([]datatransfers.BookResponse, error)
	ListBooks(ctx context.Context) ([]datatransfers.BookResponse, error)
	UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookRequest) (datatransfers.BookResponse, error)
	DeleteBook(ctx context.Context, id string) error
}

type bookClient struct {
	client protoBook.BookServiceClient
}

func NewBookClient() (BookClient, error) {
	conn, err := grpc.NewClient("book-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Failed to create BookClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)
		return nil, err
	}

	client := protoBook.NewBookServiceClient(conn)

	logger.Log.Info("Successfully created BookClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)

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

	logger.Log.Info("Sending CreateBook request to Book Service",
		zap.String("title", dto.Title),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.CreateBook(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("CreateBook request failed",
			zap.String("title", dto.Title),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.BookResponse{}, err
	}

	logger.Log.Info("CreateBook request succeeded",
		zap.String("title", resp.Book.Title),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   &resp.Book.AuthorId,
		CategoryId: &resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
		CreatedAt:  time.Unix(resp.Book.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Book.UpdatedAt, 0),
	}, nil
}

func (b *bookClient) GetBook(ctx context.Context, id string) (datatransfers.BookResponse, error) {
	reqProto := protoBook.GetBookRequest{
		Id: id,
	}

	logger.Log.Info("Sending GetBook request to Book Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.GetBook(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetBook request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.BookResponse{}, err
	}

	logger.Log.Info("GetBook request succeeded",
		zap.String("id", resp.Book.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   &resp.Book.AuthorId,
		CategoryId: &resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
		CreatedAt:  time.Unix(resp.Book.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Book.UpdatedAt, 0),
	}, nil
}

func (b *bookClient) ListBooks(ctx context.Context) ([]datatransfers.BookResponse, error) {
	reqProto := protoBook.ListBooksRequest{}

	logger.Log.Info("Sending ListBooks request to Book Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.ListBooks(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ListBooks request failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return nil, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	logger.Log.Info("ListBooks request succeeded",
		zap.Int("books_count", len(books)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

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

	logger.Log.Info("Sending UpdateBook request to Book Service",
		zap.String("id", bookId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.UpdateBook(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("UpdateBook request failed",
			zap.String("id", bookId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.BookResponse{}, err
	}

	logger.Log.Info("UpdateBook request succeeded",
		zap.String("id", resp.Book.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.BookResponse{
		Id:         resp.Book.Id,
		Title:      resp.Book.Title,
		AuthorId:   &resp.Book.AuthorId,
		CategoryId: &resp.Book.CategoryId,
		Stock:      int(resp.Book.Stock),
		CreatedAt:  time.Unix(resp.Book.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Book.UpdatedAt, 0),
	}, nil
}

func (b *bookClient) DeleteBook(ctx context.Context, id string) error {
	reqProto := protoBook.DeleteBookRequest{
		Id: id,
	}

	logger.Log.Info("Sending DeleteBook request to Book Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	_, err := b.client.DeleteBook(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("DeleteBook request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return err
	}

	logger.Log.Info("DeleteBook request succeeded",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return nil
}

func (b *bookClient) GetBooksByCategoryId(ctx context.Context, categoryId string) ([]datatransfers.BookResponse, error) {
	reqProto := protoBook.GetBooksByCategoryRequest{
		CategoryId: categoryId,
	}

	logger.Log.Info("Sending GetBooksByCategory request to Book Service",
		zap.String("category_id", categoryId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.GetBooksByCategory(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetBooksByCategory request failed",
			zap.String("category_id", categoryId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return nil, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	logger.Log.Info("GetBooksByCategory request succeeded",
		zap.Int("books_count", len(books)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return books, nil
}

func (b *bookClient) GetBooksByAuthorId(ctx context.Context, authorId string) ([]datatransfers.BookResponse, error) {
	reqProto := protoBook.GetBooksByAuthorRequest{
		AuthorId: authorId,
	}

	logger.Log.Info("Sending GetBooksByAuthor request to Book Service",
		zap.String("author_id", authorId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := b.client.GetBooksByAuthor(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetBooksByAuthor request failed",
			zap.String("author_id", authorId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return nil, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	logger.Log.Info("GetBooksByAuthor request succeeded",
		zap.Int("books_count", len(books)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return books, nil
}

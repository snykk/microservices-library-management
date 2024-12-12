package clients

import (
	"api_gateway/configs"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	protoBook "api_gateway/proto/book_service"
	"context"
	"log"
	"time"

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
	logger *logger.Logger
}

func NewBookClient(logger *logger.Logger) (BookClient, error) {
	conn, err := grpc.NewClient(configs.AppConfig.BookServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create BookClient:", err)
		return nil, err
	}
	client := protoBook.NewBookServiceClient(conn)

	log.Println("Successfully created BookClient")

	return &bookClient{
		client: client,
		logger: logger,
	}, nil
}

func (b *bookClient) CreateBook(ctx context.Context, dto datatransfers.BookRequest) (datatransfers.BookResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.CreateBookRequest{
		Title:      dto.Title,
		AuthorId:   dto.AuthorId,
		CategoryId: dto.CategoryId,
		Stock:      int32(dto.Stock),
	}

	extra := map[string]interface{}{
		"title":       dto.Title,
		"author_id":   dto.AuthorId,
		"category_id": dto.CategoryId,
		"stock":       dto.Stock,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending CreateBook request to Book Service", extra, nil)

	resp, err := b.client.CreateBook(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "CreateBook request failed", extra, err)
		return datatransfers.BookResponse{}, err
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "CreateBook request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.GetBookRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"id": id,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetBook request to Book Service", extra, nil)

	resp, err := b.client.GetBook(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetBook request failed", extra, err)
		return datatransfers.BookResponse{}, err
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetBook request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.ListBooksRequest{}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListBooks request to Book Service", nil, nil)

	resp, err := b.client.ListBooks(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListBooks request failed", nil, err)
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

	extra := map[string]interface{}{
		"books_count": len(books),
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListBooks request succeeded", extra, nil)

	return books, nil
}

func (b *bookClient) UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookRequest) (datatransfers.BookResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.UpdateBookRequest{
		Id:         bookId,
		Title:      dto.Title,
		AuthorId:   dto.AuthorId,
		CategoryId: dto.CategoryId,
		Stock:      int32(dto.Stock),
	}

	extra := map[string]interface{}{
		"id":          bookId,
		"title":       dto.Title,
		"author_id":   dto.AuthorId,
		"category_id": dto.CategoryId,
		"stock":       dto.Stock,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending UpdateBook request to Book Service", extra, nil)

	resp, err := b.client.UpdateBook(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "UpdateBook request failed", extra, err)
		return datatransfers.BookResponse{}, err
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "UpdateBook request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.DeleteBookRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"id": id,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending DeleteBook request to Book Service", extra, nil)

	_, err := b.client.DeleteBook(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "DeleteBook request failed", extra, err)
		return err
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "DeleteBook request succeeded", extra, nil)

	return nil
}

func (b *bookClient) GetBooksByCategoryId(ctx context.Context, categoryId string) ([]datatransfers.BookResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.GetBooksByCategoryRequest{
		CategoryId: categoryId,
	}

	extra := map[string]interface{}{
		"category_id": categoryId,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetBooksByCategory request to Book Service", extra, nil)

	resp, err := b.client.GetBooksByCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetBooksByCategory request failed", extra, err)
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

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetBooksByCategory request succeeded", extra, nil)

	return books, nil
}

func (b *bookClient) GetBooksByAuthorId(ctx context.Context, authorId string) ([]datatransfers.BookResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.GetBooksByAuthorRequest{
		AuthorId: authorId,
	}

	extra := map[string]interface{}{
		"author_id": authorId,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetBooksByAuthor request to Book Service", extra, nil)

	resp, err := b.client.GetBooksByAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetBooksByAuthor request failed", extra, err)
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

	extra["books_count"] = len(books)

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetBooksByAuthor request succeeded", extra, nil)

	return books, nil
}

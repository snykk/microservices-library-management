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
	GetBooksByAuthorId(ctx context.Context, authorId string, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error)
	GetBooksByCategoryId(ctx context.Context, categoryId string, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error)
	ListBooks(ctx context.Context, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error)
	UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookUpdateRequest) (datatransfers.BookResponse, error)
	DeleteBook(ctx context.Context, id string, version int) error
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
		Version:    int(resp.Book.Version),
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
		Version:    int(resp.Book.Version),
		CreatedAt:  time.Unix(resp.Book.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Book.UpdatedAt, 0),
	}, nil
}

func (b *bookClient) ListBooks(ctx context.Context, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.ListBooksRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	extra := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListBooks request to Book Service", extra, nil)

	resp, err := b.client.ListBooks(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListBooks request failed", extra, err)
		return nil, 0, 0, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			Version:    int(book.Version),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	extra["books_count"] = len(books)
	extra["total_pages"] = resp.TotalPages
	extra["total_items"] = resp.TotalItems
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListBooks request succeeded", extra, nil)

	return books, int(resp.TotalItems), int(resp.TotalPages), nil
}

func (b *bookClient) UpdateBook(ctx context.Context, bookId string, dto datatransfers.BookUpdateRequest) (datatransfers.BookResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.UpdateBookRequest{
		Id:         bookId,
		Title:      dto.Title,
		AuthorId:   dto.AuthorId,
		CategoryId: dto.CategoryId,
		Stock:      int32(dto.Stock),
		Version:    int32(dto.Version),
	}

	extra := map[string]interface{}{
		"id":          bookId,
		"title":       dto.Title,
		"author_id":   dto.AuthorId,
		"category_id": dto.CategoryId,
		"stock":       dto.Stock,
		"version":     dto.Version,
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
		Version:    int(resp.Book.Version),
		CreatedAt:  time.Unix(resp.Book.CreatedAt, 0),
		UpdatedAt:  time.Unix(resp.Book.UpdatedAt, 0),
	}, nil
}

func (b *bookClient) DeleteBook(ctx context.Context, id string, version int) error {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.DeleteBookRequest{
		Id:      id,
		Version: int32(version),
	}

	extra := map[string]interface{}{
		"id":      id,
		"version": version,
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

func (b *bookClient) GetBooksByCategoryId(ctx context.Context, categoryId string, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.GetBooksByCategoryRequest{
		CategoryId: categoryId,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	}

	extra := map[string]interface{}{
		"category_id": categoryId,
		"page":        page,
		"page_size":   pageSize,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetBooksByCategory request to Book Service", extra, nil)

	resp, err := b.client.GetBooksByCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetBooksByCategory request failed", extra, err)
		return nil, 0, 0, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			Version:    int(book.Version),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	extra["books_count"] = len(books)
	extra["total_items"] = resp.TotalItems
	extra["total_pages"] = resp.TotalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetBooksByCategory request succeeded", extra, nil)

	return books, int(resp.TotalItems), int(resp.TotalPages), nil
}

func (b *bookClient) GetBooksByAuthorId(ctx context.Context, authorId string, page int, pageSize int) ([]datatransfers.BookResponse, int, int, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	reqProto := protoBook.GetBooksByAuthorRequest{
		AuthorId: authorId,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	extra := map[string]interface{}{
		"author_id": authorId,
		"page":      page,
		"page_size": pageSize,
	}

	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetBooksByAuthor request to Book Service", extra, nil)

	resp, err := b.client.GetBooksByAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetBooksByAuthor request failed", extra, err)
		return nil, 0, 0, err
	}

	var books []datatransfers.BookResponse
	for _, book := range resp.Books {
		books = append(books, datatransfers.BookResponse{
			Id:         book.Id,
			Title:      book.Title,
			AuthorId:   &book.AuthorId,
			CategoryId: &book.CategoryId,
			Stock:      int(book.Stock),
			Version:    int(book.Version),
			CreatedAt:  time.Unix(book.CreatedAt, 0),
			UpdatedAt:  time.Unix(book.UpdatedAt, 0),
		})
	}

	extra["books_count"] = len(books)
	extra["total_items"] = resp.TotalItems
	extra["total_pages"] = resp.TotalPages
	b.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetBooksByAuthor request succeeded", extra, nil)

	return books, int(resp.TotalItems), int(resp.TotalPages), nil
}

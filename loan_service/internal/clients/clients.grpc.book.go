package clients

import (
	"context"
	protoBook "loan_service/proto/book_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookClient interface {
	GetBook(ctx context.Context, id string) (*BookResponse, error)
	UpdateBookStock(ctx context.Context, id string, newStock int) error
	IncrementBookStock(ctx context.Context, id string) error
	DecrementBookStock(ctx context.Context, id string) error
}

type bookClient struct {
	client protoBook.BookServiceClient
}

func NewBookClient() (BookClient, error) {
	conn, err := grpc.NewClient("book-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoBook.NewBookServiceClient(conn)
	return &bookClient{
		client: client,
	}, nil
}

type BookResponse struct { // simplify struct to optimize memory
	Id    string
	Title string
	Stock int
}

func (a *bookClient) GetBook(ctx context.Context, id string) (*BookResponse, error) {
	reqProto := protoBook.GetBookRequest{
		Id: id,
	}

	resp, err := a.client.GetBook(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		Id:    resp.Book.Id,
		Title: resp.Book.Title,
		Stock: int(resp.Book.Stock),
	}, nil
}

func (a *bookClient) UpdateBookStock(ctx context.Context, id string, newStock int) error {
	reqProto := protoBook.UpdateBookStockRequest{
		Id:    id,
		Stock: int32(newStock),
	}
	_, err := a.client.UpdateBookStock(ctx, &reqProto)

	return err
}

func (a *bookClient) IncrementBookStock(ctx context.Context, id string) error {
	reqProto := protoBook.IncrementBookStockRequest{
		Id: id,
	}

	_, err := a.client.IncrementBookStock(ctx, &reqProto)

	return err
}

func (a *bookClient) DecrementBookStock(ctx context.Context, id string) error {
	reqProto := protoBook.DecrementBookStockRequest{
		Id: id,
	}

	_, err := a.client.DecrementBookStock(ctx, &reqProto)

	return err
}

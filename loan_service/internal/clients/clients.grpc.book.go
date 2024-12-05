package clients

import (
	"context"
	protoBook "loan_service/proto/book_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookClient interface {
	GetBook(ctx context.Context, id string) (*BookResponse, error)
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

func (a *bookClient) GetBook(ctx context.Context, id string) (*BookResponse, error) {
	reqProto := protoBook.GetBookRequest{
		Id: id,
	}

	resp, err := a.client.GetBook(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		Id: resp.Book.Id,
	}, nil
}

type BookResponse struct { // simplify struct to optimize memory
	Id string
}

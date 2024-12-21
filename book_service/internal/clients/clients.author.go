package clients

import (
	"book_service/configs"
	protoAuthor "book_service/proto/author_service"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthorClient interface {
	GetAuthor(ctx context.Context, id string) (*AuthorResponse, error)
}

type authorClient struct {
	client protoAuthor.AuthorServiceClient
}

func NewAuthorClient() (AuthorClient, error) {
	conn, err := grpc.NewClient(configs.AppConfig.AuthorServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoAuthor.NewAuthorServiceClient(conn)
	return &authorClient{
		client: client,
	}, nil
}

func (a *authorClient) GetAuthor(ctx context.Context, id string) (*AuthorResponse, error) {
	reqProto := protoAuthor.GetAuthorRequest{
		Id: id,
	}

	resp, err := a.client.GetAuthor(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	return &AuthorResponse{
		Id: resp.Author.Id,
	}, nil
}

type AuthorResponse struct { // simplify struct to optimize memory
	Id string
}

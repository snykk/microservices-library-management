package clients

import (
	"api_gateway/internal/datatransfers"
	protoAuthor "api_gateway/proto/author_service"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthorClient interface {
	CreateAuthor(ctx context.Context, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error)
	GetAuthor(ctx context.Context, id string) (datatransfers.AuthorResponse, error)
	ListAuthors(ctx context.Context) ([]datatransfers.AuthorResponse, error)
	UpdateAuthor(ctx context.Context, authorId string, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error)
	DeleteAuthor(ctx context.Context, id string) error
}

type authorClient struct {
	client protoAuthor.AuthorServiceClient
}

func NewAuthorClient() (AuthorClient, error) {
	conn, err := grpc.NewClient("author-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoAuthor.NewAuthorServiceClient(conn)
	return &authorClient{
		client: client,
	}, nil
}

func (a *authorClient) CreateAuthor(ctx context.Context, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.CreateAuthorRequest{
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	resp, err := a.client.CreateAuthor(ctx, &reqProto)
	if err != nil {
		return datatransfers.AuthorResponse{}, err
	}

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) GetAuthor(ctx context.Context, id string) (datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.GetAuthorRequest{
		Id: id,
	}

	resp, err := a.client.GetAuthor(ctx, &reqProto)
	if err != nil {
		return datatransfers.AuthorResponse{}, err
	}

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) ListAuthors(ctx context.Context) ([]datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.ListAuthorsRequest{}

	resp, err := a.client.ListAuthors(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var authors []datatransfers.AuthorResponse
	for _, author := range resp.Authors {
		authors = append(authors, datatransfers.AuthorResponse{
			Id:        author.Id,
			Name:      author.Name,
			Biography: author.Biography,
			CreatedAt: time.Unix(author.CreatedAt, 0),
			UpdatedAt: time.Unix(author.UpdatedAt, 0),
		})
	}

	return authors, nil
}

func (a *authorClient) UpdateAuthor(ctx context.Context, authorId string, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.UpdateAuthorRequest{
		Id:        authorId,
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	resp, err := a.client.UpdateAuthor(ctx, &reqProto)
	if err != nil {
		return datatransfers.AuthorResponse{}, err
	}

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) DeleteAuthor(ctx context.Context, id string) error {
	reqProto := protoAuthor.DeleteAuthorRequest{
		Id: id,
	}

	_, err := a.client.DeleteAuthor(ctx, &reqProto)
	if err != nil {
		return err
	}

	return nil
}

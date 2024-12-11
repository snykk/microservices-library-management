package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	protoAuthor "api_gateway/proto/author_service"
	"context"
	"log"
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
	logger *logger.Logger
}

func NewAuthorClient(logger *logger.Logger) (AuthorClient, error) {
	conn, err := grpc.NewClient("author-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create AuthorClient:", err)
		return nil, err
	}
	client := protoAuthor.NewAuthorServiceClient(conn)

	log.Println("Successfully created AuthorClient")

	return &authorClient{
		client: client,
		logger: logger,
	}, nil
}

func (a *authorClient) CreateAuthor(ctx context.Context, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuthor.CreateAuthorRequest{
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	extra := map[string]interface{}{
		"name": dto.Name,
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending CreateAuthor request to Author Service", extra, nil)

	resp, err := a.client.CreateAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "CreateAuthor request failed", extra, err)
		return datatransfers.AuthorResponse{}, err
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "CreateAuthor request succeeded", extra, nil)

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) GetAuthor(ctx context.Context, id string) (datatransfers.AuthorResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuthor.GetAuthorRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"id": id,
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetAuthor request to Author Service", extra, nil)

	resp, err := a.client.GetAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetAuthor request failed", extra, err)
		return datatransfers.AuthorResponse{}, err
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetAuthor request succeeded", extra, nil)

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) ListAuthors(ctx context.Context) ([]datatransfers.AuthorResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuthor.ListAuthorsRequest{}
	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListAuthors request to Author Service", nil, nil)

	resp, err := a.client.ListAuthors(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListAuthors request failed", nil, err)
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
	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListAuthors request succeeded", map[string]interface{}{"author_count": len(authors)}, nil)

	return authors, nil
}

func (a *authorClient) UpdateAuthor(ctx context.Context, authorId string, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuthor.UpdateAuthorRequest{
		Id:        authorId,
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	extra := map[string]interface{}{
		"author_id": authorId,
		"name":      dto.Name,
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending UpdateAuthor request to Author Service", extra, nil)

	resp, err := a.client.UpdateAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "UpdateAuthor request failed", extra, err)
		return datatransfers.AuthorResponse{}, err
	}

	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "UpdateAuthor request succeeded", extra, nil)

	return datatransfers.AuthorResponse{
		Id:        resp.Author.Id,
		Name:      resp.Author.Name,
		Biography: resp.Author.Biography,
		CreatedAt: time.Unix(resp.Author.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Author.UpdatedAt, 0),
	}, nil
}

func (a *authorClient) DeleteAuthor(ctx context.Context, id string) error {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuthor.DeleteAuthorRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"id": id,
	}
	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending DeleteAuthor request to Author Service", extra, nil)

	_, err := a.client.DeleteAuthor(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "DeleteAuthor request failed", extra, err)
		return err
	}
	a.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "DeleteAuthor request succeeded", extra, nil)

	return nil
}

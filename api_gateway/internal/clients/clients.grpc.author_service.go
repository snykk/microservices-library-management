package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	protoAuthor "api_gateway/proto/author_service"
	"context"
	"time"

	"go.uber.org/zap"
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
		logger.Log.Error("Failed to create AuthorClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)
		return nil, err
	}

	client := protoAuthor.NewAuthorServiceClient(conn)
	logger.Log.Info("Successfully created AuthorClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)

	return &authorClient{
		client: client,
	}, nil
}

func (a *authorClient) CreateAuthor(ctx context.Context, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.CreateAuthorRequest{
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	logger.Log.Info("Sending CreateAuthor request to Author Service",
		zap.String("name", dto.Name),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := a.client.CreateAuthor(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("CreateAuthor request failed",
			zap.String("name", dto.Name),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.AuthorResponse{}, err
	}

	logger.Log.Info("CreateAuthor request succeeded",
		zap.String("name", resp.Author.Name),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

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

	logger.Log.Info("Sending GetAuthor request to Author Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := a.client.GetAuthor(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetAuthor request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.AuthorResponse{}, err
	}

	logger.Log.Info("GetAuthor request succeeded",
		zap.String("id", resp.Author.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

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

	logger.Log.Info("Sending ListAuthors request to Author Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := a.client.ListAuthors(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ListAuthors request failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
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

	logger.Log.Info("ListAuthors request succeeded",
		zap.Int("author_count", len(authors)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return authors, nil
}

func (a *authorClient) UpdateAuthor(ctx context.Context, authorId string, dto datatransfers.AuthorRequest) (datatransfers.AuthorResponse, error) {
	reqProto := protoAuthor.UpdateAuthorRequest{
		Id:        authorId,
		Name:      dto.Name,
		Biography: dto.Biography,
	}

	logger.Log.Info("Sending UpdateAuthor request to Author Service",
		zap.String("author_id", authorId),
		zap.String("name", dto.Name),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := a.client.UpdateAuthor(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("UpdateAuthor request failed",
			zap.String("author_id", authorId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.AuthorResponse{}, err
	}

	logger.Log.Info("UpdateAuthor request succeeded",
		zap.String("author_id", resp.Author.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

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

	logger.Log.Info("Sending DeleteAuthor request to Author Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	_, err := a.client.DeleteAuthor(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("DeleteAuthor request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return err
	}

	logger.Log.Info("DeleteAuthor request succeeded",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return nil
}

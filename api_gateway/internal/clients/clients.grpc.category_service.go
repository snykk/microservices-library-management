package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoCategory "api_gateway/proto/category_service"
	"context"
	"time"

	"api_gateway/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryClient interface {
	CreateCategory(ctx context.Context, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error)
	GetCategory(ctx context.Context, id string) (datatransfers.CategoryResponse, error)
	ListCategories(ctx context.Context) ([]datatransfers.CategoryResponse, error)
	UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryClient struct {
	client protoCategory.CategoryServiceClient
}

func NewCategoryClient() (CategoryClient, error) {
	conn, err := grpc.NewClient("category-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Failed to create CategoryClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)
		return nil, err
	}

	client := protoCategory.NewCategoryServiceClient(conn)

	logger.Log.Info("Successfully created CategoryClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)

	return &categoryClient{
		client: client,
	}, nil
}

func (c *categoryClient) CreateCategory(ctx context.Context, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.CreateCategoryRequest{
		Name: dto.Name,
	}

	logger.Log.Info("Sending CreateCategory request to Category Service",
		zap.String("name", dto.Name),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := c.client.CreateCategory(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("CreateCategory request failed",
			zap.String("name", dto.Name),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.CategoryResponse{}, err
	}

	logger.Log.Info("CreateCategory request succeeded",
		zap.String("id", resp.Category.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) GetCategory(ctx context.Context, id string) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.GetCategoryRequest{
		Id: id,
	}

	logger.Log.Info("Sending GetCategory request to Category Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := c.client.GetCategory(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetCategory request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.CategoryResponse{}, err
	}

	logger.Log.Info("GetCategory request succeeded",
		zap.String("id", resp.Category.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) ListCategories(ctx context.Context) ([]datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.ListCategoriesRequest{}

	logger.Log.Info("Sending ListCategories request to Category Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := c.client.ListCategories(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ListCategories request failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return nil, err
	}

	var categories []datatransfers.CategoryResponse
	for _, category := range resp.Categories {
		categories = append(categories, datatransfers.CategoryResponse{
			Id:        category.Id,
			Name:      category.Name,
			CreatedAt: time.Unix(category.CreatedAt, 0),
			UpdatedAt: time.Unix(category.UpdatedAt, 0),
		})
	}

	logger.Log.Info("ListCategories request succeeded",
		zap.Int("categories_count", len(categories)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return categories, nil
}

func (c *categoryClient) UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.UpdateCategoryRequest{
		Id:   categoryId,
		Name: dto.Name,
	}

	logger.Log.Info("Sending UpdateCategory request to Category Service",
		zap.String("id", categoryId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := c.client.UpdateCategory(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("UpdateCategory request failed",
			zap.String("id", categoryId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.CategoryResponse{}, err
	}

	logger.Log.Info("UpdateCategory request succeeded",
		zap.String("id", resp.Category.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) DeleteCategory(ctx context.Context, id string) error {
	reqProto := protoCategory.DeleteCategoryRequest{
		Id: id,
	}

	logger.Log.Info("Sending DeleteCategory request to Category Service",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	_, err := c.client.DeleteCategory(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("DeleteCategory request failed",
			zap.String("id", id),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return err
	}

	logger.Log.Info("DeleteCategory request succeeded",
		zap.String("id", id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return nil
}

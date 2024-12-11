package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoCategory "api_gateway/proto/category_service"
	"context"
	"log"
	"time"

	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

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
	logger *logger.Logger
}

func NewCategoryClient(logger *logger.Logger) (CategoryClient, error) {
	conn, err := grpc.NewClient("category-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create CategoryClient:", err)
		return nil, err
	}
	client := protoCategory.NewCategoryServiceClient(conn)

	log.Println("Successfully created CategoryClient")

	return &categoryClient{
		client: client,
		logger: logger,
	}, nil
}

func (c *categoryClient) CreateCategory(ctx context.Context, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.CreateCategoryRequest{
		Name: dto.Name,
	}

	extra := map[string]interface{}{
		"category_name": dto.Name,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending CreateCategory request to Category Service", extra, nil)

	resp, err := c.client.CreateCategory(ctx, &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "CreateCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "CreateCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) GetCategory(ctx context.Context, id string) (datatransfers.CategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.GetCategoryRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"category_id": id,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetCategory request to Category Service", extra, nil)

	resp, err := c.client.GetCategory(ctx, &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) ListCategories(ctx context.Context) ([]datatransfers.CategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.ListCategoriesRequest{}

	logger.Log.Info("Sending ListCategories request to Category Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListCategories request to Category Service", nil, nil)

	resp, err := c.client.ListCategories(ctx, &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Sending ListCategories request to Category Service", nil, err)
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

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListCategories request succeeded", map[string]interface{}{"categories_count": len(categories)}, nil)

	return categories, nil
}

func (c *categoryClient) UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.UpdateCategoryRequest{
		Id:   categoryId,
		Name: dto.Name,
	}

	extra := map[string]interface{}{
		"category_id":   categoryId,
		"category_name": dto.Name,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending UpdateCategory request to Category Service", extra, nil)

	resp, err := c.client.UpdateCategory(ctx, &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "UpdateCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "UpdateCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) DeleteCategory(ctx context.Context, id string) error {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.DeleteCategoryRequest{
		Id: id,
	}

	extra := map[string]interface{}{
		"category_id": id,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending DeleteCategory request to Category Service", extra, nil)

	_, err := c.client.DeleteCategory(ctx, &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "DeleteCategory request failed", extra, err)
		return err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "DeleteCategory request succeeded", extra, nil)

	return nil
}

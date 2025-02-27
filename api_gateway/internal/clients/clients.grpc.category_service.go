package clients

import (
	"api_gateway/configs"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoCategory "api_gateway/proto/category_service"
	"context"
	"log"
	"time"

	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryClient interface {
	CreateCategory(ctx context.Context, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error)
	GetCategory(ctx context.Context, id string) (datatransfers.CategoryResponse, error)
	ListCategories(ctx context.Context, page int, pageSize int) ([]datatransfers.CategoryResponse, int, int, error)
	UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryUpdateRequest) (datatransfers.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string, version int) error
}

type categoryClient struct {
	client protoCategory.CategoryServiceClient
	logger *logger.Logger
}

func NewCategoryClient(logger *logger.Logger) (CategoryClient, error) {
	conn, err := grpc.NewClient(configs.AppConfig.CategoryServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	resp, err := c.client.CreateCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "CreateCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "CreateCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		Version:   int(resp.Category.Version),
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

	resp, err := c.client.GetCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		Version:   int(resp.Category.Version),
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) ListCategories(ctx context.Context, page int, pageSize int) ([]datatransfers.CategoryResponse, int, int, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.ListCategoriesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	extra := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListCategories request to Category Service", extra, nil)

	resp, err := c.client.ListCategories(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Sending ListCategories request to Category Service", extra, err)
		return nil, 0, 0, err
	}

	var categories []datatransfers.CategoryResponse
	for _, category := range resp.Categories {
		categories = append(categories, datatransfers.CategoryResponse{
			Id:        category.Id,
			Name:      category.Name,
			Version:   int(category.Version),
			CreatedAt: time.Unix(category.CreatedAt, 0),
			UpdatedAt: time.Unix(category.UpdatedAt, 0),
		})
	}

	extra["categories_count"] = len(categories)
	extra["total_items"] = resp.TotalItems
	extra["total_pages"] = resp.TotalPages
	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListCategories request succeeded", extra, nil)

	return categories, int(resp.TotalItems), int(resp.TotalPages), nil
}

func (c *categoryClient) UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryUpdateRequest) (datatransfers.CategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.UpdateCategoryRequest{
		Id:      categoryId,
		Name:    dto.Name,
		Version: int32(dto.Version),
	}

	extra := map[string]interface{}{
		"category_id":      categoryId,
		"category_name":    dto.Name,
		"category_version": dto.Version,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending UpdateCategory request to Category Service", extra, nil)

	resp, err := c.client.UpdateCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "UpdateCategory request failed", extra, err)
		return datatransfers.CategoryResponse{}, err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "UpdateCategory request succeeded", extra, nil)

	return datatransfers.CategoryResponse{
		Id:        resp.Category.Id,
		Name:      resp.Category.Name,
		Version:   int(resp.Category.Version),
		CreatedAt: time.Unix(resp.Category.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.Category.UpdatedAt, 0),
	}, nil
}

func (c *categoryClient) DeleteCategory(ctx context.Context, id string, version int) error {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoCategory.DeleteCategoryRequest{
		Id:      id,
		Version: int32(version),
	}

	extra := map[string]interface{}{
		"category_id":      id,
		"category_version": version,
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending DeleteCategory request to Category Service", extra, nil)

	_, err := c.client.DeleteCategory(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "DeleteCategory request failed", extra, err)
		return err
	}

	c.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "DeleteCategory request succeeded", extra, nil)

	return nil
}

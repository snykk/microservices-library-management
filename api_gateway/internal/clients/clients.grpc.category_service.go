package clients

import (
	"api_gateway/internal/datatransfers"
	protoCategory "api_gateway/proto/category_service"
	"context"

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
		return nil, err
	}

	client := protoCategory.NewCategoryServiceClient(conn)
	return &categoryClient{
		client: client,
	}, nil
}

func (c *categoryClient) CreateCategory(ctx context.Context, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.CreateCategoryRequest{
		Name: dto.Name,
	}

	resp, err := c.client.CreateCategory(ctx, &reqProto)
	if err != nil {
		return datatransfers.CategoryResponse{}, err
	}

	return datatransfers.CategoryResponse{
		Id:   resp.Category.Id,
		Name: resp.Category.Name,
	}, nil
}

func (c *categoryClient) GetCategory(ctx context.Context, id string) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.GetCategoryRequest{
		Id: id,
	}

	resp, err := c.client.GetCategory(ctx, &reqProto)
	if err != nil {
		return datatransfers.CategoryResponse{}, err
	}

	return datatransfers.CategoryResponse{
		Id:   resp.Category.Id,
		Name: resp.Category.Name,
	}, nil
}

func (c *categoryClient) ListCategories(ctx context.Context) ([]datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.ListCategoriesRequest{}

	resp, err := c.client.ListCategories(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	var categories []datatransfers.CategoryResponse
	for _, category := range resp.Categories {
		categories = append(categories, datatransfers.CategoryResponse{
			Id:   category.Id,
			Name: category.Name,
		})
	}

	return categories, nil
}

func (c *categoryClient) UpdateCategory(ctx context.Context, categoryId string, dto datatransfers.CategoryRequest) (datatransfers.CategoryResponse, error) {
	reqProto := protoCategory.UpdateCategoryRequest{
		Id:   categoryId,
		Name: dto.Name,
	}

	resp, err := c.client.UpdateCategory(ctx, &reqProto)
	if err != nil {
		return datatransfers.CategoryResponse{}, err
	}

	return datatransfers.CategoryResponse{
		Id:   resp.Category.Id,
		Name: resp.Category.Name,
	}, nil
}

func (c *categoryClient) DeleteCategory(ctx context.Context, id string) error {
	reqProto := protoCategory.DeleteCategoryRequest{
		Id: id,
	}

	_, err := c.client.DeleteCategory(ctx, &reqProto)
	if err != nil {
		return err
	}

	return nil
}

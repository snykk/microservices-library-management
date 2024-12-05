package clients

import (
	protoCategory "book_service/proto/category_service"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryClient interface {
	GetCategory(ctx context.Context, id string) (*CategoryResponse, error)
}

type categoryClient struct {
	client protoCategory.CategoryServiceClient
}

func NewCategoryClient() (CategoryClient, error) {
	conn, err := grpc.NewClient("Category-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoCategory.NewCategoryServiceClient(conn)
	return &categoryClient{
		client: client,
	}, nil
}

func (c *categoryClient) GetCategory(ctx context.Context, id string) (*CategoryResponse, error) {
	reqProto := protoCategory.GetCategoryRequest{
		Id: id,
	}

	resp, err := c.client.GetCategory(ctx, &reqProto)
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id: resp.Category.Id,
	}, nil
}

type CategoryResponse struct { // simplify struct to optimize memory
	Id string
}

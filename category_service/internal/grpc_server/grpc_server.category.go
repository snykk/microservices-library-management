package grpc_server

import (
	"category_service/internal/models"
	"category_service/internal/service"
	protoCategory "category_service/proto/category_service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryGRPCServer struct {
	categoryService service.CategoryService
	protoCategory.UnimplementedCategoryServiceServer
}

func NewCategoryGRPCServer(categoryService service.CategoryService) protoCategory.CategoryServiceServer {
	return &categoryGRPCServer{
		categoryService: categoryService,
	}
}

func (s *categoryGRPCServer) CreateCategory(ctx context.Context, req *protoCategory.CreateCategoryRequest) (*protoCategory.CreateCategoryResponse, error) {
	createdCategory, err := s.categoryService.CreateCategory(ctx, &models.CategoryRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create new category")
	}

	return &protoCategory.CreateCategoryResponse{
		Category: &protoCategory.Category{
			Id:        createdCategory.Id,
			Name:      createdCategory.Name,
			CreatedAt: createdCategory.CreatedAt.Unix(),
			UpdatedAt: createdCategory.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *categoryGRPCServer) GetCategory(ctx context.Context, req *protoCategory.GetCategoryRequest) (*protoCategory.GetCategoryResponse, error) {
	category, err := s.categoryService.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve category data with id '%s'", req.Id))
	}

	return &protoCategory.GetCategoryResponse{
		Category: &protoCategory.Category{
			Id:        category.Id,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Unix(),
			UpdatedAt: category.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *categoryGRPCServer) ListCategories(ctx context.Context, req *protoCategory.ListCategoriesRequest) (*protoCategory.ListCategoriesResponse, error) {
	categories, err := s.categoryService.ListCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve category list")
	}

	var protoCategories []*protoCategory.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, &protoCategory.Category{
			Id:        category.Id,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Unix(),
			UpdatedAt: category.UpdatedAt.Unix(),
		})
	}

	return &protoCategory.ListCategoriesResponse{
		Categories: protoCategories,
	}, nil
}

func (s *categoryGRPCServer) UpdateCategory(ctx context.Context, req *protoCategory.UpdateCategoryRequest) (*protoCategory.UpdateCategoryResponse, error) {
	updatedCategory, err := s.categoryService.UpdateCategory(ctx, req.Id, &models.CategoryRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update category data with id '%s'", req.Id))
	}

	return &protoCategory.UpdateCategoryResponse{
		Category: &protoCategory.Category{
			Id:        updatedCategory.Id,
			Name:      updatedCategory.Name,
			CreatedAt: updatedCategory.CreatedAt.Unix(),
			UpdatedAt: updatedCategory.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *categoryGRPCServer) DeleteCategory(ctx context.Context, req *protoCategory.DeleteCategoryRequest) (*protoCategory.DeleteCategoryResponse, error) {
	err := s.categoryService.DeleteCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete category data with id '%s'", req.Id))
	}

	return &protoCategory.DeleteCategoryResponse{
		Message: fmt.Sprintf("success delete category with id %s", req.Id),
	}, nil
}

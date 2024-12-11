package grpc_server

import (
	"category_service/internal/constants"
	"category_service/internal/models"
	"category_service/internal/service"
	"category_service/pkg/logger"
	"category_service/pkg/utils"
	protoCategory "category_service/proto/category_service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryGRPCServer struct {
	categoryService service.CategoryService
	logger          *logger.Logger
	protoCategory.UnimplementedCategoryServiceServer
}

func NewCategoryGRPCServer(categoryService service.CategoryService, logger *logger.Logger) protoCategory.CategoryServiceServer {
	return &categoryGRPCServer{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (s *categoryGRPCServer) CreateCategory(ctx context.Context, req *protoCategory.CreateCategoryRequest) (*protoCategory.CreateCategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received CreateCategory request", map[string]interface{}{"name": req.Name}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid CreateCategory request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	createdCategory, err := s.categoryService.CreateCategory(ctx, &models.CategoryRequest{Name: req.Name})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create category", nil, err)
		return nil, status.Error(codes.Internal, "failed to create new category")
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category created successfully", map[string]interface{}{"category_id": createdCategory.Id}, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetCategory request", map[string]interface{}{"category_id": req.Id}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetCategory request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	category, err := s.categoryService.GetCategory(ctx, req.Id)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to retrieve category with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve category with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category retrieved successfully", map[string]interface{}{"category_id": category.Id}, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListCategories request", nil, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListCategories request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	categories, err := s.categoryService.ListCategories(ctx)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve categories list", nil, err)
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

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Categories list retrieved successfully", nil, nil)

	return &protoCategory.ListCategoriesResponse{
		Categories: protoCategories,
	}, nil
}

func (s *categoryGRPCServer) UpdateCategory(ctx context.Context, req *protoCategory.UpdateCategoryRequest) (*protoCategory.UpdateCategoryResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received UpdateCategory request", map[string]interface{}{"category_id": req.Id, "name": req.Name}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid UpdateCategory request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	updatedCategory, err := s.categoryService.UpdateCategory(ctx, req.Id, &models.CategoryRequest{Name: req.Name})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to update category with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update category with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category updated successfully", map[string]interface{}{"category_id": updatedCategory.Id}, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received DeleteCategory request", map[string]interface{}{"category_id": req.Id}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid DeleteCategory request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.categoryService.DeleteCategory(ctx, req.Id)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to delete category with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete category with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Category deleted successfully", map[string]interface{}{"category_id": req.Id}, nil)

	return &protoCategory.DeleteCategoryResponse{
		Message: fmt.Sprintf("success delete category with id %s", req.Id),
	}, nil
}

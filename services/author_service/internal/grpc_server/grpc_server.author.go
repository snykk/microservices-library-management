package grpc_server

import (
	"author_service/internal/models"
	"author_service/internal/service"
	"author_service/pkg/logger"
	protoAuthor "author_service/proto/author_service"
	"context"
	"fmt"

	"author_service/internal/constants"
	"author_service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorGRPCServer struct {
	authorService service.AuthorService
	logger        *logger.Logger
	protoAuthor.UnimplementedAuthorServiceServer
}

func NewAuthorGRPCServer(authorService service.AuthorService, logger *logger.Logger) protoAuthor.AuthorServiceServer {
	return &authorGRPCServer{
		authorService: authorService,
		logger:        logger,
	}
}

func (s *authorGRPCServer) CreateAuthor(ctx context.Context, req *protoAuthor.CreateAuthorRequest) (*protoAuthor.CreateAuthorResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received CreateAuthor request", map[string]interface{}{"name": req.Name, "biography": req.Biography}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid CreateAuthor request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	createdAuthor, err := s.authorService.CreateAuthor(ctx, &models.AuthorCreateRequest{
		Name:      req.Name,
		Biography: req.Biography,
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to create author", nil, err)
		return nil, status.Error(codes.Internal, "failed to create new author")
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author created successfully", map[string]interface{}{"created_author": createdAuthor}, nil)

	return &protoAuthor.CreateAuthorResponse{
		Author: &protoAuthor.Author{
			Id:        createdAuthor.Id,
			Name:      createdAuthor.Name,
			Biography: createdAuthor.Biography,
			CreatedAt: createdAuthor.CreatedAt.Unix(),
			UpdatedAt: createdAuthor.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authorGRPCServer) GetAuthor(ctx context.Context, req *protoAuthor.GetAuthorRequest) (*protoAuthor.GetAuthorResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetAuthor request", map[string]interface{}{"author_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetAuthor request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	author, err := s.authorService.GetAuthor(ctx, req.Id)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to retrieve author with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve author with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author retrieved successfully", map[string]interface{}{"author_id": author.Id}, nil)

	return &protoAuthor.GetAuthorResponse{
		Author: &protoAuthor.Author{
			Id:        author.Id,
			Name:      author.Name,
			Biography: author.Biography,
			Version:   int32(author.Version),
			CreatedAt: author.CreatedAt.Unix(),
			UpdatedAt: author.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authorGRPCServer) ListAuthors(ctx context.Context, req *protoAuthor.ListAuthorsRequest) (*protoAuthor.ListAuthorsResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListAuthors request", nil, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListAuthors request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	// Retrieve authors list
	authors, totalItems, err := s.authorService.ListAuthors(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve authors list", nil, err)
		return nil, status.Error(codes.Internal, "failed to retrieve author list")
	}

	var protoAuthors []*protoAuthor.Author
	for _, author := range authors {
		protoAuthors = append(protoAuthors, &protoAuthor.Author{
			Id:        author.Id,
			Name:      author.Name,
			Biography: author.Biography,
			Version:   int32(author.Version),
			CreatedAt: author.CreatedAt.Unix(),
			UpdatedAt: author.UpdatedAt.Unix(),
		})
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Authors list retrieved successfully", nil, nil)

	return &protoAuthor.ListAuthorsResponse{
		Authors:    protoAuthors,
		TotalItems: int32(totalItems),
		TotalPages: int32(utils.CalculateTotalPages(totalItems, int(req.PageSize))),
	}, nil
}

func (s *authorGRPCServer) UpdateAuthor(ctx context.Context, req *protoAuthor.UpdateAuthorRequest) (*protoAuthor.UpdateAuthorResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received UpdateAuthor request", map[string]interface{}{"author_id": req.Id, "name": req.Name, "biography": req.Biography}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid UpdateAuthor request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	updatedAuthor, err := s.authorService.UpdateAuthor(ctx, req.Id, &models.AuthorUpdateRequest{
		Name:      req.Name,
		Biography: req.Biography,
		Version:   int(req.Version),
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to update author with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update author with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author updated successfully", map[string]interface{}{"author_id": updatedAuthor.Id}, nil)

	return &protoAuthor.UpdateAuthorResponse{
		Author: &protoAuthor.Author{
			Id:        updatedAuthor.Id,
			Name:      updatedAuthor.Name,
			Biography: updatedAuthor.Biography,
			Version:   req.Version,
			CreatedAt: updatedAuthor.CreatedAt.Unix(),
			UpdatedAt: updatedAuthor.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authorGRPCServer) DeleteAuthor(ctx context.Context, req *protoAuthor.DeleteAuthorRequest) (*protoAuthor.DeleteAuthorResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received DeleteAuthor request", map[string]interface{}{"author_id": req.Id}, nil)

	// Validate request from client
	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid DeleteAuthor request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	err := s.authorService.DeleteAuthor(ctx, req.Id, int(req.Version))
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to delete author with id '%s'", req.Id), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete author with id '%s'", req.Id))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Author deleted successfully", map[string]interface{}{"author_id": req.Id}, nil)

	return &protoAuthor.DeleteAuthorResponse{Message: fmt.Sprintf("success delete author with id %s", req.Id)}, nil
}

package grpc_server

import (
	"author_service/internal/models"
	"author_service/internal/service"
	protoAuthor "author_service/proto/author_service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorGRPCServer struct {
	authorService service.AuthorService
	protoAuthor.UnimplementedAuthorServiceServer
}

func NewAuthorGRPCServer(authorService service.AuthorService) protoAuthor.AuthorServiceServer {
	return &authorGRPCServer{
		authorService: authorService,
	}
}

func (s *authorGRPCServer) CreateAuthor(ctx context.Context, req *protoAuthor.CreateAuthorRequest) (*protoAuthor.CreateAuthorResponse, error) {
	createdAuthor, err := s.authorService.CreateAuthor(ctx, &models.AuthorRequest{
		Name:      req.Name,
		Biography: req.Biography,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create new author")
	}

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
	author, err := s.authorService.GetAuthor(ctx, &req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve author with id '%s'", req.Id))
	}

	return &protoAuthor.GetAuthorResponse{
		Author: &protoAuthor.Author{
			Id:        author.Id,
			Name:      author.Name,
			Biography: author.Biography,
			CreatedAt: author.CreatedAt.Unix(),
			UpdatedAt: author.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authorGRPCServer) ListAuthors(ctx context.Context, req *protoAuthor.ListAuthorsRequest) (*protoAuthor.ListAuthorsResponse, error) {
	authors, err := s.authorService.ListAuthors(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve author list")
	}

	var protoAuthors []*protoAuthor.Author
	for _, author := range authors {
		protoAuthors = append(protoAuthors, &protoAuthor.Author{
			Id:        author.Id,
			Name:      author.Name,
			Biography: author.Biography,
			CreatedAt: author.CreatedAt.Unix(),
			UpdatedAt: author.UpdatedAt.Unix(),
		})
	}

	return &protoAuthor.ListAuthorsResponse{
		Authors: protoAuthors,
	}, nil
}

func (s *authorGRPCServer) UpdateAuthor(ctx context.Context, req *protoAuthor.UpdateAuthorRequest) (*protoAuthor.UpdateAuthorResponse, error) {
	updatedAuthor, err := s.authorService.UpdateAuthor(ctx, &req.Id, &models.AuthorRequest{
		Name:      req.Name,
		Biography: req.Biography,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update author with id '%s'", req.Id))
	}

	return &protoAuthor.UpdateAuthorResponse{
		Author: &protoAuthor.Author{
			Id:        updatedAuthor.Id,
			Name:      updatedAuthor.Name,
			Biography: updatedAuthor.Biography,
			CreatedAt: updatedAuthor.CreatedAt.Unix(),
			UpdatedAt: updatedAuthor.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authorGRPCServer) DeleteAuthor(ctx context.Context, req *protoAuthor.DeleteAuthorRequest) (*protoAuthor.DeleteAuthorResponse, error) {
	err := s.authorService.DeleteAuthor(ctx, &req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete author with id '%s'", req.Id))
	}

	return &protoAuthor.DeleteAuthorResponse{Message: fmt.Sprintf("success delete author with id %s", req.Id)}, nil
}

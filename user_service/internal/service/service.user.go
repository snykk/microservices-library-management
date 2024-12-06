package service

import (
	"context"
	"user_service/internal/models"
	"user_service/internal/repository"
)

type UserService interface {
	GetUserById(ctx context.Context, userId string) (*models.UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error)
	ListUsers(ctx context.Context) ([]*models.UserRecord, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetUserById(ctx context.Context, userId string) (*models.UserRecord, error) {
	return s.repo.GetUserById(ctx, userId)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.UserRecord, error) {
	return s.repo.ListUsers(ctx)
}

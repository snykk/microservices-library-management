package service

import (
	"context"
	"log"
	"user_service/internal/models"
	"user_service/internal/repository"
	"user_service/pkg/utils"
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
	log.Printf("[%s] Fetching user by ID: %s\n", utils.GetLocation(), userId)

	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		log.Printf("[%s] Failed to fetch user by ID %s: %v\n", utils.GetLocation(), userId, err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched user by ID %s\n", utils.GetLocation(), userId)
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	log.Printf("[%s] Fetching user by email: %s\n", utils.GetLocation(), email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("[%s] Failed to fetch user by email %s: %v\n", utils.GetLocation(), email, err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched user by email %s\n", utils.GetLocation(), email)
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.UserRecord, error) {
	log.Printf("[%s] Fetching list of users\n", utils.GetLocation())

	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		log.Printf("[%s] Failed to fetch list of users: %v\n", utils.GetLocation(), err)
		return nil, err
	}

	log.Printf("[%s] Successfully fetched %d users\n", utils.GetLocation(), len(users))
	return users, nil
}

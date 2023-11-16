package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/entity"
	"github.com/khhini/riset-gcp-api-gateway-auth/pkg/password"
)

type Service struct {
	userRepository entity.Repository
}

func NewService(userRepo entity.Repository) entity.Service {
	return &Service{
		userRepository: userRepo,
	}
}

func (s *Service) InsertUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	user.CreatedOn = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	password.GenerateHash(user.Password, &user.Password)

	insertedID, err := s.userRepository.InsertUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user; %w", err)
	}

	return insertedID, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID string, updatedUser *entity.User) error {
	err := s.userRepository.UpdateUser(ctx, userID, updatedUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepository.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

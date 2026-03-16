package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type hasher interface {
	HashPassword(password string) (string, error)
}

type Service struct {
	repository repository
	hasher     hasher
}

func NewService(repository repository, hasher hasher) *Service {
	return &Service{
		repository: repository,
		hasher:     hasher,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, password string) error {
	hashedPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return err
	}

	userId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	user := User{
		Email:        email,
		PasswordHash: hashedPassword,
		ID:           userId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) VerifyEmailIsTaken(ctx context.Context, email string) (bool, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrorUserNotFound) {
		return false, err
	}
	return user != nil, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repository.GetUserByEmail(ctx, email)
}

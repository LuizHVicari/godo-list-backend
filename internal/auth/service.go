package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/user"
)

type userService interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	VerifyEmailIsTaken(ctx context.Context, email string) (bool, error)
}

type repository interface {
	CreateSession(ctx context.Context, session Session, sessionTtlSeconds int) error
	GetSessionByID(ctx context.Context, sessionId uuid.UUID) (*Session, error)
	DeleteSession(ctx context.Context, sessionId uuid.UUID) error
	RefreshSession(ctx context.Context, sessionId uuid.UUID, sessionTtlSeconds int) error
}

type hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) (bool, error)
}

type Service struct {
	userService userService
	repository  repository
	hasher      hasher
}

// Avoids generating a fake hash for both every use case attempt or every service initialization
// this aims to make server initialization faster, making it more compatible with serverless environments
// it must be updated whenever the hashing parameters are updated to ensure it has the same cost as a real hash
const fakePasswordHash = "$argon2id$v=19$m=64,t=3,p=2$RmxReUx4Z09YUkltTGVzUQ$9yEqNkVTgnuafT06+VEeYNarL5ETGDszIbycABuNZNU"
const sessionTtlSeconds = 86400

func NewService(userService userService, repository repository, hasher hasher) *Service {

	return &Service{
		userService: userService,
		repository:  repository,
		hasher:      hasher,
	}
}

func (s *Service) SignUp(ctx context.Context, email, password string) error {

	isEmailTaken, err := s.userService.VerifyEmailIsTaken(ctx, email)
	if err != nil {
		return err
	}
	if isEmailTaken {
		// prevents user enumeration by not revealing that the email is already taken
		_, _ = s.hasher.HashPassword(password)
		return nil
	}
	err = s.userService.CreateUser(ctx, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, email, password string) (*Session, error) {

	userEntity, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, user.ErrorUserNotFound) {
		return nil, err
	}
	if err != nil && errors.Is(err, user.ErrorUserNotFound) {
		// prevents user enumeration by not revealing that the email is not found
		s.hasher.ComparePassword(fakePasswordHash, password)
		return nil, user.ErrorInvalidCredentials
	}

	isPasswordValid, err := s.hasher.ComparePassword(userEntity.PasswordHash, password)
	if err != nil {
		return nil, err
	}
	if !isPasswordValid {
		return nil, user.ErrorInvalidCredentials
	}

	session, err := s.createSession(userEntity.ID)
	if err != nil {
		return nil, err
	}

	err = s.repository.CreateSession(ctx, *session, sessionTtlSeconds)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) SignOut(ctx context.Context, sessionId uuid.UUID) error {
	return s.repository.DeleteSession(ctx, sessionId)
}

func (s *Service) VerifySessionValid(ctx context.Context, sessionId uuid.UUID) (*Session, error) {
	session, err := s.repository.GetSessionByID(ctx, sessionId)
	if errors.Is(err, ErrorSessionNotFound) {
		return nil, ErrorSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	err = s.repository.RefreshSession(ctx, sessionId, sessionTtlSeconds)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) createSession(userId uuid.UUID) (*Session, error) {
	sessionId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	session := Session{
		ID:        sessionId,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &session, nil
}

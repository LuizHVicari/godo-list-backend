package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/platform/db"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	})
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	userModel, err := r.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorUserNotFound
	}
	if err != nil {
		return nil, err
	}

	userEntity := modelToEntity(userModel)

	return &userEntity, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	userModel, err := r.queries.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorUserNotFound
	}
	if err != nil {
		return nil, err
	}

	userEntity := modelToEntity(userModel)

	return &userEntity, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user User) error {
	return r.queries.UpdateUserByID(ctx, db.UpdateUserByIDParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		UpdatedAt:    user.UpdatedAt,
	})
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUserByID(ctx, id)
}

func modelToEntity(user db.AuthUser) User {
	return User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

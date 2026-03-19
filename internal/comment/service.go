package comment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	CreateComment(ctx context.Context, c Comment) error
	GetCommentByID(ctx context.Context, id uuid.UUID) (*Comment, error)
	UpdateComment(ctx context.Context, c Comment) error
	DeleteComment(ctx context.Context, id uuid.UUID) error
	ListCommentsByItemID(ctx context.Context, itemID uuid.UUID, filter ListCommentsFilter) (*ListCommentsResult, error)
	IsItemInOwnedProject(ctx context.Context, itemID, ownerID uuid.UUID) (bool, error)
}

type CreateCommentParams struct {
	ItemID   uuid.UUID
	AuthorID uuid.UUID
	Content  string
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateComment(ctx context.Context, params CreateCommentParams) (*Comment, error) {
	owned, err := s.repo.IsItemInOwnedProject(ctx, params.ItemID, params.AuthorID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorItemNotFound
	}

	now := time.Now()
	c := NewComment(uuid.New(), params.ItemID, params.AuthorID, params.Content, now, now)
	if err := s.repo.CreateComment(ctx, *c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCommentByID(ctx context.Context, id, itemID, ownerID uuid.UUID) (*Comment, error) {
	owned, err := s.repo.IsItemInOwnedProject(ctx, itemID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorItemNotFound
	}

	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.ItemID != itemID {
		return nil, ErrorCommentNotFound
	}
	return c, nil
}

func (s *Service) UpdateComment(ctx context.Context, id, authorID uuid.UUID, content string) error {
	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	if c.AuthorID != authorID {
		return ErrorForbidden
	}

	c.Content = content
	c.UpdatedAt = time.Now()
	return s.repo.UpdateComment(ctx, *c)
}

func (s *Service) DeleteComment(ctx context.Context, id, authorID uuid.UUID) error {
	c, err := s.repo.GetCommentByID(ctx, id)
	if err != nil {
		return err
	}
	if c.AuthorID != authorID {
		return ErrorForbidden
	}

	return s.repo.DeleteComment(ctx, id)
}

func (s *Service) ListCommentsByItemID(ctx context.Context, itemID, ownerID uuid.UUID, filter ListCommentsFilter) (*ListCommentsResult, error) {
	owned, err := s.repo.IsItemInOwnedProject(ctx, itemID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorItemNotFound
	}

	return s.repo.ListCommentsByItemID(ctx, itemID, filter)
}

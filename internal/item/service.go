package item

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	CreateItem(ctx context.Context, item Item) error
	GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error)
	UpdateItem(ctx context.Context, item Item) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	ListItemsByStepID(ctx context.Context, stepID uuid.UUID, filter ListItemsFilter) (*ListItemsResult, error)
	GetLastItemPositionByStepID(ctx context.Context, stepID uuid.UUID) (int32, error)
	IsItemPositionTaken(ctx context.Context, stepID uuid.UUID, position int32, excludeID *uuid.UUID) (bool, error)
	IsStepInOwnedProject(ctx context.Context, stepID, ownerID uuid.UUID) (bool, error)
	RepositionItems(ctx context.Context, params RepositionItemsParams) error
}

type CreateItemParams struct {
	StepID      uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	Description *string
	Priority    ItemPriority
	Position    *int32
}

type UpdateItemParams struct {
	Name        string
	Description *string
	Priority    ItemPriority
	Position    int32
	IsCompleted bool
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateItem(ctx context.Context, params CreateItemParams) error {
	owned, err := s.repo.IsStepInOwnedProject(ctx, params.StepID, params.OwnerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorItemNotFound
	}

	if params.Priority == "" {
		params.Priority = ItemPriorityNone
	}

	position := params.Position
	if position == nil {
		last, err := s.repo.GetLastItemPositionByStepID(ctx, params.StepID)
		if err != nil {
			return err
		}
		next := last + 1
		position = &next
	}

	taken, err := s.repo.IsItemPositionTaken(ctx, params.StepID, *position, nil)
	if err != nil {
		return err
	}
	if taken {
		return ErrorItemPositionTaken
	}

	now := time.Now()
	item := NewItem(uuid.New(), params.Name, params.Description, params.Priority, *position, params.StepID, now, now)
	return s.repo.CreateItem(ctx, *item)
}

func (s *Service) GetItemByID(ctx context.Context, id, stepID, ownerID uuid.UUID) (*Item, error) {
	owned, err := s.repo.IsStepInOwnedProject(ctx, stepID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorItemNotFound
	}

	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.StepID != stepID {
		return nil, ErrorItemNotFound
	}
	return item, nil
}

func (s *Service) UpdateItem(ctx context.Context, id, stepID, ownerID uuid.UUID, params UpdateItemParams) error {
	owned, err := s.repo.IsStepInOwnedProject(ctx, stepID, ownerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorItemNotFound
	}

	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return err
	}
	if item.StepID != stepID {
		return ErrorItemNotFound
	}

	if params.Position != item.Position {
		taken, err := s.repo.IsItemPositionTaken(ctx, stepID, params.Position, &id)
		if err != nil {
			return err
		}
		if taken {
			return ErrorItemPositionTaken
		}
	}

	item.Name = params.Name
	item.Description = params.Description
	item.Priority = params.Priority
	item.Position = params.Position
	item.IsCompleted = params.IsCompleted
	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(ctx, *item)
}

func (s *Service) DeleteItem(ctx context.Context, id, stepID, ownerID uuid.UUID) error {
	owned, err := s.repo.IsStepInOwnedProject(ctx, stepID, ownerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorItemNotFound
	}

	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return err
	}
	if item.StepID != stepID {
		return ErrorItemNotFound
	}
	return s.repo.DeleteItem(ctx, id)
}

func (s *Service) ListItemsByStepID(ctx context.Context, stepID, ownerID uuid.UUID, filter ListItemsFilter) (*ListItemsResult, error) {
	owned, err := s.repo.IsStepInOwnedProject(ctx, stepID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorItemNotFound
	}
	return s.repo.ListItemsByStepID(ctx, stepID, filter)
}

func (s *Service) RepositionItems(ctx context.Context, params RepositionItemsParams) error {
	owned, err := s.repo.IsStepInOwnedProject(ctx, params.StepID, params.OwnerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorItemNotFound
	}
	return s.repo.RepositionItems(ctx, params)
}

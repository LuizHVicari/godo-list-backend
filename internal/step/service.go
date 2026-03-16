package step

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	CreateStep(ctx context.Context, step Step) error
	GetStepByID(ctx context.Context, id uuid.UUID) (*Step, error)
	UpdateStep(ctx context.Context, step Step) error
	DeleteStep(ctx context.Context, id uuid.UUID) error
	ListStepsByProjectID(ctx context.Context, projectID uuid.UUID, filter ListStepsFilter) (*ListStepsResult, error)
	GetLastStepPositionByProjectID(ctx context.Context, projectID uuid.UUID) (int32, error)
	IsStepPositionTaken(ctx context.Context, projectID uuid.UUID, position int32, excludeID *uuid.UUID) (bool, error)
	IsProjectOwnedByUser(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error)
	RepositionSteps(ctx context.Context, params RepositionStepsParams) error
}

type CreateStepParams struct {
	ProjectID uuid.UUID
	OwnerID   uuid.UUID
	Name      string
	Position  *int32
}

type UpdateStepParams struct {
	Name        string
	Position    int32
	IsCompleted bool
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateStep(ctx context.Context, params CreateStepParams) (*Step, error) {
	owned, err := s.repo.IsProjectOwnedByUser(ctx, params.ProjectID, params.OwnerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorStepNotFound
	}

	position := params.Position
	if position == nil {
		last, err := s.repo.GetLastStepPositionByProjectID(ctx, params.ProjectID)
		if err != nil {
			return nil, err
		}
		next := last + 1
		position = &next
	}

	taken, err := s.repo.IsStepPositionTaken(ctx, params.ProjectID, *position, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrorStepPositionTaken
	}

	now := time.Now()
	step := NewStep(uuid.New(), params.ProjectID, params.Name, *position, now, now)
	if err := s.repo.CreateStep(ctx, *step); err != nil {
		return nil, err
	}
	return step, nil
}

func (s *Service) GetStepByID(ctx context.Context, id, projectID, ownerID uuid.UUID) (*Step, error) {
	owned, err := s.repo.IsProjectOwnedByUser(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorStepNotFound
	}

	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if step.ProjectID != projectID {
		return nil, ErrorStepNotFound
	}
	return step, nil
}

func (s *Service) UpdateStep(ctx context.Context, id, ownerID uuid.UUID, params UpdateStepParams) error {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return err
	}

	owned, err := s.repo.IsProjectOwnedByUser(ctx, step.ProjectID, ownerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorStepNotFound
	}

	if params.Position != step.Position {
		taken, err := s.repo.IsStepPositionTaken(ctx, step.ProjectID, params.Position, &id)
		if err != nil {
			return err
		}
		if taken {
			return ErrorStepPositionTaken
		}
	}

	step.Name = params.Name
	step.Position = params.Position
	step.IsCompleted = params.IsCompleted
	step.UpdatedAt = time.Now()
	return s.repo.UpdateStep(ctx, *step)
}

func (s *Service) DeleteStep(ctx context.Context, id, ownerID uuid.UUID) error {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return err
	}

	owned, err := s.repo.IsProjectOwnedByUser(ctx, step.ProjectID, ownerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorStepNotFound
	}

	return s.repo.DeleteStep(ctx, id)
}

func (s *Service) ListStepsByProjectID(ctx context.Context, projectID, ownerID uuid.UUID, filter ListStepsFilter) (*ListStepsResult, error) {
	owned, err := s.repo.IsProjectOwnedByUser(ctx, projectID, ownerID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrorStepNotFound
	}
	return s.repo.ListStepsByProjectID(ctx, projectID, filter)
}

func (s *Service) RepositionSteps(ctx context.Context, params RepositionStepsParams) error {
	owned, err := s.repo.IsProjectOwnedByUser(ctx, params.ProjectID, params.OwnerID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrorStepNotFound
	}
	return s.repo.RepositionSteps(ctx, params)
}

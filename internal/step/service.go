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
	RepositionSteps(ctx context.Context, params RepositionStepsParams) error
}

type CreateStepParams struct {
	ProjectID uuid.UUID
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

func (s *Service) CreateStep(ctx context.Context, params CreateStepParams) error {
	position := params.Position
	if position == nil {
		last, err := s.repo.GetLastStepPositionByProjectID(ctx, params.ProjectID)
		if err != nil {
			return err
		}
		next := last + 1
		position = &next
	}

	now := time.Now()
	step := NewStep(uuid.New(), params.ProjectID, params.Name, *position, now, now)
	return s.repo.CreateStep(ctx, *step)
}

func (s *Service) GetStepByID(ctx context.Context, id, projectID uuid.UUID) (*Step, error) {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if step.ProjectID != projectID {
		return nil, ErrorStepNotFound
	}
	return step, nil
}

func (s *Service) UpdateStep(ctx context.Context, id, projectID uuid.UUID, params UpdateStepParams) error {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return err
	}
	if step.ProjectID != projectID {
		return ErrorStepNotFound
	}
	step.Name = params.Name
	step.Position = params.Position
	step.IsCompleted = params.IsCompleted
	step.UpdatedAt = time.Now()
	return s.repo.UpdateStep(ctx, *step)
}

func (s *Service) DeleteStep(ctx context.Context, id, projectID uuid.UUID) error {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return err
	}
	if step.ProjectID != projectID {
		return ErrorStepNotFound
	}
	return s.repo.DeleteStep(ctx, id)
}

func (s *Service) ListStepsByProjectID(ctx context.Context, projectID uuid.UUID, filter ListStepsFilter) (*ListStepsResult, error) {
	return s.repo.ListStepsByProjectID(ctx, projectID, filter)
}

func (s *Service) RepositionSteps(ctx context.Context, params RepositionStepsParams) error {
	return s.repo.RepositionSteps(ctx, params)
}

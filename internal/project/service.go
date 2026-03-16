package project

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	CreateProject(ctx context.Context, project Project) error
	GetProjectById(ctx context.Context, id uuid.UUID) (*Project, error)
	UpdateProject(ctx context.Context, project Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
	ListProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID, filter ListProjectsFilter) (*ListProjectsResult, error)
}

type CreateProjectParams struct {
	Name        string
	Description *string
	OwnerID     uuid.UUID
}

type UpdateProjectParams struct {
	Name        string
	Description *string
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateProject(ctx context.Context, params CreateProjectParams) (*Project, error) {
	now := time.Now()
	project := NewProject(uuid.New(), params.Name, params.Description, params.OwnerID, now, now)
	if err := s.repo.CreateProject(ctx, *project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *Service) GetProjectById(ctx context.Context, id, ownerID uuid.UUID) (*Project, error) {
	project, err := s.repo.GetProjectById(ctx, id)
	if err != nil {
		return nil, err
	}
	if project.OwnerID != ownerID {
		return nil, ErrorProjectNotFound
	}
	return project, nil
}

func (s *Service) UpdateProject(ctx context.Context, id, ownerID uuid.UUID, params UpdateProjectParams) error {
	project, err := s.repo.GetProjectById(ctx, id)
	if err != nil {
		return err
	}
	if project.OwnerID != ownerID {
		return ErrorProjectNotFound
	}
	project.Name = params.Name
	project.Description = params.Description
	project.UpdatedAt = time.Now()
	return s.repo.UpdateProject(ctx, *project)
}

func (s *Service) DeleteProject(ctx context.Context, id, ownerID uuid.UUID) error {
	project, err := s.repo.GetProjectById(ctx, id)
	if err != nil {
		return err
	}
	if project.OwnerID != ownerID {
		return ErrorProjectNotFound
	}
	return s.repo.DeleteProject(ctx, id)
}

func (s *Service) ListProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID, filter ListProjectsFilter) (*ListProjectsResult, error) {
	return s.repo.ListProjectsByOwnerID(ctx, ownerID, filter)
}

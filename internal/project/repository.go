package project

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

func (r *Repository) CreateProject(ctx context.Context, project Project) error {
	return r.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:   project.ID,
		Name: project.Name,
		Description: sql.NullString{
			String: ptrVal(project.Description),
			Valid:  project.Description != nil,
		},
		OwnerID:   project.OwnerID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	})
}

func (r *Repository) GetProjectById(ctx context.Context, id uuid.UUID) (*Project, error) {
	projectModel, err := r.queries.GetProjectByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	projectEntity := r.modelToEntity(projectModel)

	return projectEntity, nil
}

func (r *Repository) UpdateProject(ctx context.Context, project Project) error {
	return r.queries.UpdateProjectByID(ctx, db.UpdateProjectByIDParams{
		ID:   project.ID,
		Name: project.Name,
		Description: sql.NullString{
			String: ptrVal(project.Description),
			Valid:  project.Description != nil,
		},
		UpdatedAt: project.UpdatedAt,
	})
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteProjectByID(ctx, id)
}

type ListProjectsFilter struct {
	Name      *string
	Sort      *string // "name" | "created_at"
	Direction *string // "asc" | "desc"
	Limit     *int32
	Offset    *int32
}

type ListProjectsResult struct {
	Total    int64
	Projects []*Project
}

func (r *Repository) ListProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID, filter ListProjectsFilter) (*ListProjectsResult, error) {
	if filter.Sort != nil && *filter.Sort != "name" && *filter.Sort != "created_at" {
		return nil, ErrorInvalidFilterParams
	}
	if filter.Direction != nil && *filter.Direction != "asc" && *filter.Direction != "desc" {
		return nil, ErrorInvalidFilterParams
	}

	total, err := r.queries.CountProjectsByOwnerID(ctx, db.CountProjectsByOwnerIDParams{
		OwnerID: ownerID,
		Name:    sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
	})
	if err != nil {
		return nil, err
	}

	projectModels, err := r.queries.ListProjectsByOwnerID(ctx, db.ListProjectsByOwnerIDParams{
		OwnerID:   ownerID,
		Name:      sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
		Sort:      sql.NullString{String: ptrVal(filter.Sort), Valid: filter.Sort != nil},
		Direction: sql.NullString{String: ptrVal(filter.Direction), Valid: filter.Direction != nil},
		Limit:     sql.NullInt32{Int32: ptrVal(filter.Limit), Valid: filter.Limit != nil},
		Offset:    sql.NullInt32{Int32: ptrVal(filter.Offset), Valid: filter.Offset != nil},
	})
	if err != nil {
		return nil, err
	}

	var projects []*Project
	for _, projectModel := range projectModels {
		projects = append(projects, r.modelToEntity(projectModel))
	}

	return &ListProjectsResult{Total: total, Projects: projects}, nil
}

func (r *Repository) modelToEntity(projectModel db.TodoProject) *Project {
	var description *string
	if projectModel.Description.Valid {
		description = &projectModel.Description.String
	}

	return &Project{
		ID:          projectModel.ID,
		Name:        projectModel.Name,
		Description: description,
		OwnerID:     projectModel.OwnerID,
		CreatedAt:   projectModel.CreatedAt,
		UpdatedAt:   projectModel.UpdatedAt,
	}
}


func ptrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

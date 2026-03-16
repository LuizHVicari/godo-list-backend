package step

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/luizhvicari/backend/internal/platform/db"
)

type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewRepository(database *sql.DB, queries *db.Queries) *Repository {
	return &Repository{db: database, queries: queries}
}

type ListStepsFilter struct {
	Name      *string
	Sort      *string // "name" | "position" | "created_at" | "updated_at"
	Direction *string // "asc" | "desc"
	Limit     *int32
	Offset    *int32
}

type ListStepsResult struct {
	Total int64
	Steps []*Step
}

func (r *Repository) CreateStep(ctx context.Context, step Step) error {
	return r.queries.CreateStep(ctx, db.CreateStepParams{
		ID:          step.ID,
		ProjectID:   step.ProjectID,
		Name:        step.Name,
		Position:    step.Position,
		IsCompleted: step.IsCompleted,
		CreatedAt:   step.CreatedAt,
		UpdatedAt:   step.UpdatedAt,
	})
}

func (r *Repository) GetStepByID(ctx context.Context, id uuid.UUID) (*Step, error) {
	stepModel, err := r.queries.GetStepByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorStepNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(stepModel), nil
}

func (r *Repository) UpdateStep(ctx context.Context, step Step) error {
	return r.queries.UpdateStepByID(ctx, db.UpdateStepByIDParams{
		ID:          step.ID,
		Name:        step.Name,
		Position:    step.Position,
		IsCompleted: step.IsCompleted,
		UpdatedAt:   step.UpdatedAt,
	})
}

func (r *Repository) DeleteStep(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteStepByID(ctx, id)
}

type StepReposition struct {
	ID       uuid.UUID
	Position int32
}

type RepositionStepsParams struct {
	ProjectID uuid.UUID
	OwnerID   uuid.UUID
	Steps     []StepReposition
}

func (r *Repository) RepositionSteps(ctx context.Context, params RepositionStepsParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	now := time.Now()

	for _, s := range params.Steps {
		result, err := qtx.UpdateStepPositionByID(ctx, db.UpdateStepPositionByIDParams{
			ID:        s.ID,
			ProjectID: params.ProjectID,
			Position:  s.Position,
			UpdatedAt: now,
		})
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrorStepNotBelongsToProject
		}
	}

	return tx.Commit()
}

func (r *Repository) GetLastStepPositionByProjectID(ctx context.Context, projectID uuid.UUID) (int32, error) {
	return r.queries.GetLastStepPositionByProjectID(ctx, projectID)
}

func (r *Repository) IsProjectOwnedByUser(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error) {
	return r.queries.IsProjectOwnedByUser(ctx, db.IsProjectOwnedByUserParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
}

func (r *Repository) IsStepPositionTaken(ctx context.Context, projectID uuid.UUID, position int32, excludeID *uuid.UUID) (bool, error) {
	var nullID uuid.NullUUID
	if excludeID != nil {
		nullID = uuid.NullUUID{UUID: *excludeID, Valid: true}
	}
	return r.queries.IsStepPositionTaken(ctx, db.IsStepPositionTakenParams{
		ProjectID: projectID,
		Position:  position,
		ExcludeID: nullID,
	})
}

func (r *Repository) ListStepsByProjectID(ctx context.Context, projectID uuid.UUID, filter ListStepsFilter) (*ListStepsResult, error) {
	validSorts := map[string]bool{"name": true, "position": true, "created_at": true, "updated_at": true}
	if filter.Sort != nil && !validSorts[*filter.Sort] {
		return nil, ErrorInvalidFilterParams
	}
	if filter.Direction != nil && *filter.Direction != "asc" && *filter.Direction != "desc" {
		return nil, ErrorInvalidFilterParams
	}

	total, err := r.queries.CountStepsByProjectID(ctx, db.CountStepsByProjectIDParams{
		ProjectID: projectID,
		Name:      sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
	})
	if err != nil {
		return nil, err
	}

	stepModels, err := r.queries.ListStepsByProjectID(ctx, db.ListStepsByProjectIDParams{
		ProjectID: projectID,
		Name:      sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
		Sort:      sql.NullString{String: ptrVal(filter.Sort), Valid: filter.Sort != nil},
		Direction: sql.NullString{String: ptrVal(filter.Direction), Valid: filter.Direction != nil},
		Limit:     sql.NullInt32{Int32: ptrVal(filter.Limit), Valid: filter.Limit != nil},
		Offset:    sql.NullInt32{Int32: ptrVal(filter.Offset), Valid: filter.Offset != nil},
	})
	if err != nil {
		return nil, err
	}

	var steps []*Step
	for _, m := range stepModels {
		steps = append(steps, r.modelToEntity(m))
	}

	return &ListStepsResult{Total: total, Steps: steps}, nil
}

func (r *Repository) modelToEntity(m db.TodoStep) *Step {
	return &Step{
		ID:          m.ID,
		ProjectID:   m.ProjectID,
		Name:        m.Name,
		Position:    m.Position,
		IsCompleted: m.IsCompleted,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ptrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

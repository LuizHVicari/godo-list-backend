package item

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

type ListItemsFilter struct {
	Name      *string
	Sort      *string // "name" | "position" | "priority" | "created_at" | "updated_at"
	Direction *string // "asc" | "desc"
	Limit     *int32
	Offset    *int32
}

type ListItemsResult struct {
	Total int64
	Items []*Item
}

func (r *Repository) CreateItem(ctx context.Context, item Item) error {
	return r.queries.CreateItem(ctx, db.CreateItemParams{
		ID:          item.ID,
		StepID:      item.StepID,
		Name:        item.Name,
		Description: sql.NullString{String: ptrVal(item.Description), Valid: item.Description != nil},
		Priority:    db.TodoItemPriority(item.Priority),
		Position:    item.Position,
		IsCompleted: item.IsCompleted,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (r *Repository) GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	itemModel, err := r.queries.GetItemByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorItemNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(itemModel), nil
}

func (r *Repository) UpdateItem(ctx context.Context, item Item) error {
	return r.queries.UpdateItemByID(ctx, db.UpdateItemByIDParams{
		ID:          item.ID,
		Name:        item.Name,
		Description: sql.NullString{String: ptrVal(item.Description), Valid: item.Description != nil},
		Priority:    db.TodoItemPriority(item.Priority),
		Position:    item.Position,
		IsCompleted: item.IsCompleted,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (r *Repository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteItemByID(ctx, id)
}

type ItemReposition struct {
	ID       uuid.UUID
	Position int32
}

type RepositionItemsParams struct {
	StepID  uuid.UUID
	OwnerID uuid.UUID
	Items   []ItemReposition
}

func (r *Repository) RepositionItems(ctx context.Context, params RepositionItemsParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	now := time.Now()

	for _, i := range params.Items {
		result, err := qtx.UpdateItemPositionByID(ctx, db.UpdateItemPositionByIDParams{
			ID:        i.ID,
			StepID:    params.StepID,
			Position:  i.Position,
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
			return ErrorItemNotBelongsToStep
		}
	}

	return tx.Commit()
}

func (r *Repository) GetLastItemPositionByStepID(ctx context.Context, stepID uuid.UUID) (int32, error) {
	return r.queries.GetLastItemPositionByStepID(ctx, stepID)
}

func (r *Repository) IsStepInOwnedProject(ctx context.Context, stepID, ownerID uuid.UUID) (bool, error) {
	return r.queries.IsStepInOwnedProject(ctx, db.IsStepInOwnedProjectParams{
		ID:      stepID,
		OwnerID: ownerID,
	})
}

func (r *Repository) IsItemPositionTaken(ctx context.Context, stepID uuid.UUID, position int32, excludeID *uuid.UUID) (bool, error) {
	var nullID uuid.NullUUID
	if excludeID != nil {
		nullID = uuid.NullUUID{UUID: *excludeID, Valid: true}
	}
	return r.queries.IsItemPositionTaken(ctx, db.IsItemPositionTakenParams{
		StepID:    stepID,
		Position:  position,
		ExcludeID: nullID,
	})
}

func (r *Repository) ListItemsByStepID(ctx context.Context, stepID uuid.UUID, filter ListItemsFilter) (*ListItemsResult, error) {
	validSorts := map[string]bool{"name": true, "position": true, "priority": true, "created_at": true, "updated_at": true}
	if filter.Sort != nil && !validSorts[*filter.Sort] {
		return nil, ErrorInvalidFilterParams
	}
	if filter.Direction != nil && *filter.Direction != "asc" && *filter.Direction != "desc" {
		return nil, ErrorInvalidFilterParams
	}

	total, err := r.queries.CountItemsByStepID(ctx, db.CountItemsByStepIDParams{
		StepID: stepID,
		Name:   sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
	})
	if err != nil {
		return nil, err
	}

	itemModels, err := r.queries.ListItemsByStepID(ctx, db.ListItemsByStepIDParams{
		StepID:    stepID,
		Name:      sql.NullString{String: ptrVal(filter.Name), Valid: filter.Name != nil},
		Sort:      sql.NullString{String: ptrVal(filter.Sort), Valid: filter.Sort != nil},
		Direction: sql.NullString{String: ptrVal(filter.Direction), Valid: filter.Direction != nil},
		Limit:     sql.NullInt32{Int32: ptrVal(filter.Limit), Valid: filter.Limit != nil},
		Offset:    sql.NullInt32{Int32: ptrVal(filter.Offset), Valid: filter.Offset != nil},
	})
	if err != nil {
		return nil, err
	}

	var items []*Item
	for _, m := range itemModels {
		items = append(items, r.modelToEntity(m))
	}

	return &ListItemsResult{Total: total, Items: items}, nil
}

func (r *Repository) modelToEntity(m db.TodoItem) *Item {
	var desc *string
	if m.Description.Valid {
		desc = &m.Description.String
	}
	return &Item{
		ID:          m.ID,
		StepID:      m.StepID,
		Name:        m.Name,
		Description: desc,
		Priority:    ItemPriority(m.Priority),
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

package comment

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
	return &Repository{queries: queries}
}

type ListCommentsFilter struct {
	Limit  *int32
	Offset *int32
}

type ListCommentsResult struct {
	Total    int64
	Comments []*Comment
}

func (r *Repository) CreateComment(ctx context.Context, c Comment) error {
	return r.queries.CreateComment(ctx, db.CreateCommentParams{
		ID:        c.ID,
		ItemID:    c.ItemID,
		AuthorID:  c.AuthorID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	})
}

func (r *Repository) GetCommentByID(ctx context.Context, id uuid.UUID) (*Comment, error) {
	m, err := r.queries.GetCommentByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return modelToEntity(m), nil
}

func (r *Repository) UpdateComment(ctx context.Context, c Comment) error {
	return r.queries.UpdateCommentByID(ctx, db.UpdateCommentByIDParams{
		ID:        c.ID,
		Content:   c.Content,
		UpdatedAt: c.UpdatedAt,
	})
}

func (r *Repository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCommentByID(ctx, id)
}

func (r *Repository) ListCommentsByItemID(ctx context.Context, itemID uuid.UUID, filter ListCommentsFilter) (*ListCommentsResult, error) {
	total, err := r.queries.CountCommentsByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	models, err := r.queries.ListCommentsByItemID(ctx, db.ListCommentsByItemIDParams{
		ItemID: itemID,
		Offset: sql.NullInt32{Int32: ptrVal(filter.Offset), Valid: filter.Offset != nil},
		Limit:  sql.NullInt32{Int32: ptrVal(filter.Limit), Valid: filter.Limit != nil},
	})
	if err != nil {
		return nil, err
	}

	comments := make([]*Comment, len(models))
	for i, m := range models {
		comments[i] = modelToEntity(m)
	}

	return &ListCommentsResult{Total: total, Comments: comments}, nil
}

func (r *Repository) IsItemInOwnedProject(ctx context.Context, itemID, ownerID uuid.UUID) (bool, error) {
	return r.queries.IsItemInOwnedProject(ctx, db.IsItemInOwnedProjectParams{
		ID:      itemID,
		OwnerID: ownerID,
	})
}

func modelToEntity(m db.TodoItemComment) *Comment {
	return &Comment{
		ID:        m.ID,
		ItemID:    m.ItemID,
		AuthorID:  m.AuthorID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ptrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// ensure Repository implements the repository interface at compile time
var _ repository = (*Repository)(nil)

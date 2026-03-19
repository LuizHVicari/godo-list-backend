package comment

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	ItemID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComment(id, itemID, authorID uuid.UUID, content string, createdAt, updatedAt time.Time) *Comment {
	return &Comment{
		ID:        id,
		ItemID:    itemID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

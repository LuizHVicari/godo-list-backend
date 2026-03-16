package step

import (
	"time"

	"github.com/google/uuid"
)

type Step struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Name        string
	Position    int32
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewStep(id, projectID uuid.UUID, name string, position int32, createdAt, updatedAt time.Time) *Step {
	return &Step{
		ID:          id,
		ProjectID:   projectID,
		Name:        name,
		Position:    position,
		IsCompleted: false,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

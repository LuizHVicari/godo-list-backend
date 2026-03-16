package item

import (
	"time"

	"github.com/google/uuid"
)

type ItemPriority string

const (
	ItemPriorityNone     ItemPriority = "none"
	ItemPriorityLow      ItemPriority = "low"
	ItemPriorityMedium   ItemPriority = "medium"
	ItemPriorityHigh     ItemPriority = "high"
	ItemPriorityCritical ItemPriority = "critical"
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Priority    ItemPriority
	Position    int32
	IsCompleted bool
	StepID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewItem(
	id uuid.UUID,
	name string,
	description *string,
	priority ItemPriority,
	position int32,
	stepID uuid.UUID,
	createdAt, updatedAt time.Time,
) *Item {
	return &Item{
		ID:          id,
		Name:        name,
		Description: description,
		Priority:    priority,
		Position:    position,
		StepID:      stepID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

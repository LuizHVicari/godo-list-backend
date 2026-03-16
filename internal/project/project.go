package project

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Description *string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(id uuid.UUID, name string, description *string, ownerID uuid.UUID, createdAt, updatedAt time.Time) *Project {
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

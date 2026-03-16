package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSession(id, userId uuid.UUID) Session {
	return Session{
		ID:        id,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

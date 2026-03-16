package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id uuid.UUID, email, passwordHash string) User {
	return User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
}

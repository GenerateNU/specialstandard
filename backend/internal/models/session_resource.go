package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionResource struct {
	SessionID  uuid.UUID `json:"session_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateSessionResource struct {
	SessionID  uuid.UUID `json:"session_id"`
	ResourceID uuid.UUID `json:"resource_id"`
}

type DeleteSessionResource = CreateSessionResource

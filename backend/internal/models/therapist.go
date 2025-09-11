package models

import (
	"time"

	"github.com/google/uuid"
)

type Therapist struct {
	ID              uuid.UUID  `json:"id"`
	First_name      string     `json:"first_name"`
	Last_name       int        `json:"last_name"`
	Email           int        `json:"email"`
	Active          bool       `json:"active"`
	Created_at      time.Time  `json:"created_at"`
	Updated_at      time.Time  `json:"updated_at"`
}

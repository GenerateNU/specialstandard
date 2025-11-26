package models

import (
	"time"

	"github.com/google/uuid"
)

type Newsletter struct {
	ID        uuid.UUID `db:"id" json:"id"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
	S3URL     string    `db:"s3_url" json:"s3_url"`
}

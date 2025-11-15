package models

import "time"

type School struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	DistrictID int       `db:"district_id" json:"district_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
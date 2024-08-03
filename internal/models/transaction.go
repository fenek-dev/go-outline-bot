package models

import (
	"time"
)

type Transaction struct {
	ID         string    `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Amount     int       `json:"amount" db:"amount"`
	Type       string    `json:"type" db:"type"`
	Status     string    `json:"status" db:"status"`
	ExternalID string    `json:"external_id" db:"external_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

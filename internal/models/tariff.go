package models

import "time"

type Tariff struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     uint32    `json:"price" db:"price"`
	Bandwidth uint64    `json:"bandwidth" db:"bandwidth"`
	ServerID  uint64    `json:"server_id" db:"server_id"`
	Active    bool      `json:"active" db:"active"`
	Duration  uint32    `json:"duration" db:"duration"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

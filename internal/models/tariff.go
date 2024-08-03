package models

import "time"

type Tariff struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"`
	Bandwidth int       `json:"bandwidth" db:"bandwidth"`
	ServerID  int64     `json:"server_id" db:"server_id"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

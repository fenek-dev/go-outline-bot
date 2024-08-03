package models

import "time"

type Subscription struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	ServerID     int64     `json:"server_id" db:"server_id"`
	TariffID     int64     `json:"tariff_id" db:"tariff_id"`
	InitialPrice int       `json:"initial_price" db:"initial_price"`
	Key          string    `json:"-" db:"key"`
	Status       string    `json:"status" db:"status"`
	ExpiredAt    time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

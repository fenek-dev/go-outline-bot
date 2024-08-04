package models

import "time"

type Server struct {
	ID             uint64    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	CountryCode    string    `json:"country_code" db:"country_code"`
	Ping           uint32    `json:"ping" db:"ping"`
	IP             string    `json:"-" db:"ip"`
	Port           string    `json:"-" db:"port"`
	APIKey         string    `json:"-" db:"api_key"`
	Certificate    string    `json:"-" db:"certificate"`
	MaxConnections uint32    `json:"max_connections" db:"max_connections"`
	TotalBandwidth uint32    `json:"total_bandwidth" db:"total_bandwidth"`
	Active         bool      `json:"active" db:"active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

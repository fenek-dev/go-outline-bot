package models

import "time"

type Server struct {
	ID             int64     `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	CountryCode    string    `json:"country_code" db:"country_code"`
	Ping           int       `json:"ping" db:"ping"`
	ServerIP       string    `json:"-" db:"server_ip"`
	ServerPort     string    `json:"-" db:"server_port"`
	APIKey         string    `json:"-" db:"api_key"`
	Certificate    string    `json:"-" db:"certificate"`
	MaxConnections int       `json:"max_connections" db:"max_connections"`
	TotalBandwidth int       `json:"total_bandwidth" db:"total_bandwidth"`
	Active         bool      `json:"active" db:"active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

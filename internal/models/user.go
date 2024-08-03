package models

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	Username  *string   `json:"username" db:"username"`
	Balance   int       `json:"balance" db:"balance"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

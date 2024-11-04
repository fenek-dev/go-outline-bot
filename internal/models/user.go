package models

import "time"

type User struct {
	ID        uint64    `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	Username  *string   `json:"username" db:"username"`
	Balance   uint32    `json:"balance" db:"balance"`
	PartnerID *uint64   `json:"partner_id" db:"partner_id"`
	BonusUsed bool      `json:"bonus_used" db:"bonus_used"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

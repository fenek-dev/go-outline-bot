package models

import "time"

type Coupon struct {
	ID            uint64    `json:"id" db:"id"`
	OwnerID       uint64    `json:"owner_id" db:"owner_id"`
	ReceiverID    *uint64   `json:"receiver_id" db:"receiver_id"`
	TransactionID *string   `json:"transaction_id" db:"transaction_id"`
	Amount        uint32    `json:"amount" db:"amount"`
	Code          string    `json:"code" db:"code"`
	Status        string    `json:"status" db:"status"`
	ExpiredAt     time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

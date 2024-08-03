package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type TransactionMeta struct {
	IsCoupon       *bool   `json:"is_coupon"`
	IsProlongation *bool   `json:"is_prolongation"`
	CouponID       *uint64 `json:"coupon_id"`
	SubscriptionID *uint64 `json:"subscription_id"`
}

type Transaction struct {
	ID string `json:"id" db:"id"`

	UserID     uint64            `json:"user_id" db:"user_id"`
	Amount     uint32            `json:"amount" db:"amount"`
	Meta       string            `json:"meta" db:"meta"`
	ExternalID *string           `json:"external_id" db:"external_id"`
	Type       TransactionType   `json:"type" db:"type"`
	Status     TransactionStatus `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

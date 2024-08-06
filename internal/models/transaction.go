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
	IsCoupon        *bool   `json:"is_coupon"`
	IsDiscounted    *bool   `json:"is_discounted"`
	IsProlongation  *bool   `json:"is_prolongation"`
	IsCommission    *bool   `json:"is_commission"`
	ReferalID       *uint64 `json:"referal_id"`
	CouponID        *uint64 `json:"coupon_id"`
	SubscriptionID  *uint64 `json:"subscription_id"`
	DiscountPercent *uint8  `json:"discount_percent"`
}

type Transaction struct {
	ID uint64 `json:"id" db:"id"`

	UserID     uint64            `json:"user_id" db:"user_id"`
	Amount     uint32            `json:"amount" db:"amount"`
	Meta       string            `json:"meta" db:"meta"`
	ExternalID *string           `json:"external_id" db:"external_id"`
	Type       TransactionType   `json:"type" db:"type"`
	Status     TransactionStatus `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

package storage

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNegativeBalance   = errors.New("balance cannot be negative")
)

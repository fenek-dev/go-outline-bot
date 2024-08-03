package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type BalanceErrors error

var (
	ErrPhoneNumberEmpty BalanceErrors = errors.New("phone number is empty")
)

func (s *Service) RefreshBalance(ctx context.Context, user models.User) (balance int64, err error) {
	// TODO: calc balance by all user transactions
	return 0, nil
}

func (s *Service) RequestDeposit(ctx context.Context, user models.User, amount uint32) (redirectUri string, err error) {
	// TODO: Create transaction
	// TODO: Request to payment service

	if user.Phone == nil {
		return "", ErrPhoneNumberEmpty
	}

	transaction := &models.Transaction{
		UserID:     user.ID,
		Amount:     amount,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeDeposit,
		ExternalID: lo.ToPtr(uuid.New().String()),
	}

	err = s.storage.CreateTransaction(ctx, transaction)
	if err != nil {
		return "", err
	}

	response, err := s.paymentClient.CreateTransaction(ctx, user.ID, payment_service.CreateTransactionRequest{
		TxUUID:       transaction.ID,
		Amount:       transaction.Amount,
		UserID:       string(user.ID),
		CurrencyCode: "RUB",
		Description:  "Пополнение баланса",

		MethodID:   1,
		MethodType: payment_service.TransactionMethodTypePayment,

		PostbackURL: "",
		SuccessURL:  "",
		FailURL:     "",

		Customer: payment_service.TransactionCustomer{
			ID:    string(user.ID),
			Email: "notify@myshelf.shop",
			Phone: lo.FromPtr(user.Phone),
		},

		Items: []payment_service.TransactionItem{
			{
				ID:       string(user.ID),
				Name:     fmt.Sprintf("Пополнение баланса #%d", user.ID),
				Price:    transaction.Amount,
				Quantity: 1,
				SKU:      lo.ToPtr(fmt.Sprintf("balance-%d", user.ID)),
			},
		},
	})

	if err != nil {
		// TODO: Update transaction status
	}

	if response.Action.Type != payment_service.ResultTypeRedirect || response.Action.RedirectURL == nil {
		return "", fmt.Errorf("invalid payment service response: %v", response)
	}

	return *response.Action.RedirectURL, nil
}

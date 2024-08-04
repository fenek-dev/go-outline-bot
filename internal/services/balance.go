package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
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

		PostbackURL: "", // @TODO: Get from config PAYMENT_SERVICE_POSTBACK_URL
		SuccessURL:  "", // @TODOD: Get from config PAYMENT_SERVICE_SUCCESS_URL
		FailURL:     "", // @TODOD: Get from config PAYMENT_SERVICE_FAIL_URL

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

func (s *Service) ConfirmDeposit(ctx context.Context, user models.User, transactionID uint64) (err error) {
	transaction, err := s.storage.GetTransaction(ctx, transactionID)
	if err != nil {
		return err
	}

	if transaction.UserID != user.ID {
		return fmt.Errorf("transaction %s not found", transactionID)
	}

	if transaction.Status != models.TransactionStatusPending {
		return fmt.Errorf("transaction %s already confirmed", transactionID)
	}

	err = s.storage.WithTx(ctx, "ConfirmDeposit", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.IncBalanceTx(ctx, tx, user.ID, transaction.Amount)
		if err != nil {
			return err
		}

		err = s.storage.UpdateTransactionStatusTx(ctx, tx, transactionID, models.TransactionStatusSuccess)
		if err != nil {
			return err
		}

		// @TODO: Send notification to user

		// Partner commission
		// @TODO: Outbox for partner commission
		if user.PartnerID != nil {
			commission := s.calcPartnerComission(transaction.Amount)

			err = s.storage.IncBalanceTx(ctx, tx, *user.PartnerID, commission)
			if err != nil {
				return err
			}

			meta, err := json.Marshal(models.TransactionMeta{
				IsCommission: lo.ToPtr(true),
				ReferalID:    lo.ToPtr(user.ID),
			})

			if err != nil {
				return err
			}

			err = s.storage.CreateTransactionTx(ctx, tx, &models.Transaction{
				UserID: *user.PartnerID,
				Amount: commission,
				Type:   models.TransactionTypeDeposit,
				Status: models.TransactionStatusSuccess,
				Meta:   string(meta),
			})

			// @TODO: Send notification to partner
		}

		return nil
	})

	return err
}

func (s *Service) calcPartnerComission(amount uint32) uint32 {
	commission := amount / 10
	return commission - (commission % 10)
}

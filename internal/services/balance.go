package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

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

func (s *Service) GetBalance(ctx context.Context, userId uint64) (balance uint32, err error) {
	return s.storage.GetBalance(ctx, userId)
}

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
		TxUUID:       string(transaction.ID),
		Amount:       transaction.Amount,
		UserID:       string(user.ID),
		CurrencyCode: "RUB",
		Description:  "Пополнение баланса",

		MethodID:   1,
		MethodType: payment_service.TransactionMethodTypePayment,

		PostbackURL: s.config.Payment.PostbackUrl,
		SuccessURL:  s.config.Payment.SuccessUrl,
		FailURL:     s.config.Payment.FailUrl,

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
		return "", fmt.Errorf("failed to request deposit: %w", err)
	}

	if response.Action.Type != payment_service.ResultTypeRedirect || response.Action.RedirectURL == nil {
		return "", fmt.Errorf("invalid payment service response: %v", response)
	}

	return *response.Action.RedirectURL, nil
}

func (s *Service) ConfirmDeposit(ctx context.Context, transactionExternalID string) (err error) {
	var partnerTransaction *models.Transaction
	transaction, err := s.storage.GetTransactionByExternalID(ctx, transactionExternalID)
	if err != nil {
		return fmt.Errorf("failed to get transaction by external id: %w", err)
	}

	if transaction.Status != models.TransactionStatusPending {
		return fmt.Errorf("transaction %d already confirmed", transaction.ID)
	}

	user, err := s.storage.GetUser(ctx, transaction.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.storage.WithTx(ctx, "ConfirmDeposit", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.IncBalanceTx(ctx, tx, user.ID, transaction.Amount)
		if err != nil {
			return fmt.Errorf("failed to increase balance: %w", err)
		}

		err = s.storage.UpdateTransactionStatusTx(ctx, tx, transaction.ID, models.TransactionStatusSuccess)
		if err != nil {
			return fmt.Errorf("failed to update transaction status: %w", err)
		}

		// @TODO: Send notification to user

		// Partner commission
		// @TODO: Outbox for partner commission
		if user.PartnerID != nil {
			commission := s.CalcPartnerComission(transaction.Amount)

			err = s.storage.IncBalanceTx(ctx, tx, *user.PartnerID, commission)
			if err != nil {
				return fmt.Errorf("failed to increase partner balance: %w", err)
			}

			meta, err := json.Marshal(models.TransactionMeta{
				IsCommission: lo.ToPtr(true),
				ReferalID:    lo.ToPtr(user.ID),
			})

			if err != nil {
				return fmt.Errorf("failed to marshal meta: %w", err)
			}

			partnerTransaction := &models.Transaction{
				UserID: *user.PartnerID,
				Amount: commission,
				Type:   models.TransactionTypeDeposit,
				Status: models.TransactionStatusSuccess,
				Meta:   string(meta),
			}

			err = s.storage.CreateTransactionTx(ctx, tx, partnerTransaction)

			if err != nil {
				return fmt.Errorf("failed to create partner transaction: %w", err)
			}
		}

		return nil
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to confirm deposit: %w", err)
	}

	if user.PartnerID != nil && partnerTransaction != nil {
		// @TODO: To outbox
		s.NotifyPartnerAboutDeposit(ctx, *partnerTransaction)
	}

	return err
}

func (s *Service) CancelDeposit(ctx context.Context, transactionExternalID string) (err error) {
	transaction, err := s.storage.GetTransactionByExternalID(ctx, transactionExternalID)
	if err != nil {
		return fmt.Errorf("failed to get transaction by external id: %w", err)
	}

	if transaction.Status != models.TransactionStatusPending {
		return fmt.Errorf("transaction %d already confirmed", transaction.ID)
	}

	err = s.storage.WithTx(ctx, "CancelDeposit", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.UpdateTransactionStatusTx(ctx, tx, transaction.ID, models.TransactionStatusFailed)
		if err != nil {
			return err
		}

		return nil
	}, nil)

	// @TODO: Send notification to user

	return err
}

func (s *Service) CalcPartnerComission(amount uint32) uint32 {
	commissionPercent := s.config.Partner.CommissionPercent
	if commissionPercent == 0 {
		return 0
	}

	commission := math.Floor(float64(amount) * float64(commissionPercent) / 100)

	return uint32(commission)
}

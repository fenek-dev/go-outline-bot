package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/pkg/utils"
	"github.com/samber/lo"
)

type SubscriptionError error

var (
	ErrAlreadyHasTrial  SubscriptionError = errors.New("already has trial")
	ErrNotEnoughBalance SubscriptionError = errors.New("not enough balance")
)

func (s *Service) CreateSubscription(ctx context.Context, user models.User, tariffID uint64) (subscription *models.Subscription, err error) {
	tariff, err := s.storage.GetTariff(ctx, tariffID)
	if err != nil {
		return subscription, fmt.Errorf("get tariff: %w", err)
	}

	if tariff.IsTrial {
		hasTrial, err := s.storage.TrialSubscriptionExists(ctx, user.ID)
		if err != nil {
			return subscription, fmt.Errorf("has trial subscription: %w", err)
		}

		if hasTrial {
			return subscription, ErrAlreadyHasTrial
		}
	}

	price, discountPercent := s.CalcTariffPrice(ctx, tariff, user)

	if (user.Balance - price) < 0 {
		return nil, ErrNotEnoughBalance
	}

	s.balanceMu.Lock()
	defer s.balanceMu.Unlock()

	key, err := s.CreateKey(ctx, tariff)
	if err != nil {
		return nil, fmt.Errorf("create key: %w", err)
	}

	subscription = &models.Subscription{
		UserID:       user.ID,
		ServerID:     tariff.ServerID,
		TariffID:     tariff.ID,
		InitialPrice: price,
		KeyUUID:      key.ID,
		AccessUrl:    key.AccessURL,
		ServerIP:     key.Method,
		ServerPort:   key.Port,
		Password:     key.Password,
		Method:       key.Method,
		ExpiredAt:    utils.CalcExpiredAt(tariff.Duration),
		Status:       "pending",
	}

	txErr := s.storage.WithTx(ctx, "CreateSubscription", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.CreateSubscriptionTx(ctx, tx, subscription)

		meta, err := json.Marshal(models.TransactionMeta{
			IsDiscounted:    lo.ToPtr(discountPercent > 0),
			SubscriptionID:  &subscription.ID,
			DiscountPercent: &discountPercent,
		})

		if err != nil {
			return err
		}

		err = s.storage.CreateTransactionTx(ctx, tx, &models.Transaction{
			UserID: user.ID,
			Amount: price,
			Type:   models.TransactionTypeWithdrawal,
			Status: models.TransactionStatusSuccess,
			Meta:   string(meta),
		})

		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		err = s.storage.DecBalanceTx(ctx, tx, user.ID, price)
		if err != nil {
			return fmt.Errorf("dec balance: %w", err)
		}

		err = s.storage.SetUserBonusUsedTx(ctx, tx, user.ID)

		if err != nil {
			return fmt.Errorf("set user bonus used: %w", err)
		}

		return nil
	}, nil)

	if txErr != nil {
		return nil, fmt.Errorf("create subscription tx: %w", txErr)
	}

	return subscription, nil
}

func (s *Service) GetSubscriptions(ctx context.Context, user models.User) (subscriptions []models.Subscription, err error) {
	return s.storage.GetSubscriptionsByUser(ctx, user.ID)
}

// TODO: EnableAutoProlongation
// TODO: DisableAutoProlongation
func (s *Service) ExpireSubscription(ctx context.Context, subscription models.Subscription) (err error) {

	if subscription.AutoProlong {
		err = s.ProlongSubscription(ctx, subscription)
		return fmt.Errorf("prolong subscription: %w", err)
	}

	// TODO: to outbox in db transaction
	err = s.DeactivateKey(ctx, subscription)
	if err != nil {
		return fmt.Errorf("deactivate key: %w", err)
	}

	txErr := s.storage.WithTx(ctx, "ExpireSubscription", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.UpdateSubscriptionStatusTx(ctx, tx, subscription.ID, "expired")
		if err != nil {
			return fmt.Errorf("update subscription status: %w", err)
		}

		return nil
	}, nil)

	if txErr != nil {
		return fmt.Errorf("expire subscription tx: %w", txErr)
	}

	go s.NotifySubscriptionExpired(ctx, subscription)

	return nil
}

func (s *Service) ProlongSubscription(ctx context.Context, subscription models.Subscription) (err error) {
	tariff, err := s.storage.GetTariff(ctx, subscription.TariffID)
	if err != nil {
		return fmt.Errorf("get tariff: %w", err)
	}

	user, err := s.storage.GetUser(ctx, subscription.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if (user.Balance - tariff.Price) < 0 {
		return ErrNotEnoughBalance
	}

	txErr := s.storage.WithTx(ctx, "ProlongSubscription", func(ctx context.Context, tx pg.Executor) error {
		meta, err := json.Marshal(models.TransactionMeta{
			IsProlongation: lo.ToPtr(true),
		})

		if err != nil {
			return fmt.Errorf("marshal transaction meta: %w", err)
		}

		err = s.storage.CreateTransactionTx(ctx, tx, &models.Transaction{
			UserID: subscription.UserID,
			Amount: subscription.InitialPrice,
			Type:   models.TransactionTypeWithdrawal,
			Status: models.TransactionStatusSuccess,
			Meta:   string(meta),
		})

		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		err = s.storage.IncBalanceTx(ctx, tx, subscription.UserID, subscription.InitialPrice)
		if err != nil {
			return fmt.Errorf("inc balance: %w", err)
		}

		err = s.storage.ProlongSubscriptionTx(ctx, tx, subscription.ID, utils.CalcExpiredAt(tariff.Duration))
		if err != nil {
			return fmt.Errorf("prolong subscription: %w", err)
		}

		return nil
	}, nil)

	if txErr != nil {
		return fmt.Errorf("prolong subscription tx: %w", txErr)
	}

	go s.NotifySubscriptionProlongation(ctx, subscription)

	return nil
}

func (s *Service) CalcTariffPrice(ctx context.Context, tariff models.Tariff, user models.User) (price uint32, discountPercent uint8) {
	if user.PartnerID != nil && !user.BonusUsed {
		discountPercent := s.config.Partner.DiscountPercent // @TODO: To config PARTNER_DISCOUNT_PERCENT
		price := math.Floor(float64(tariff.Price) * (1 - float64(discountPercent/100)))

		return uint32(price), discountPercent
	}

	return tariff.Price, 0
}

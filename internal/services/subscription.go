package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/pkg/utils"
	"github.com/samber/lo"
)

type SubscriptionError error

var (
	ErrNotEnoughBalance SubscriptionError = errors.New("not enough balance")
)

func (s *Service) CreateSubscription(ctx context.Context, user models.User, tariffID uint64) (subscription *models.Subscription, err error) {
	tariff, err := s.storage.GetTariff(ctx, tariffID)
	if err != nil {
		return subscription, err
	}

	price, discountPercent := s.GetTariffPrice(ctx, tariff, user)

	if (user.Balance - price) < 0 {
		return nil, ErrNotEnoughBalance
	}

	s.balanceMu.Lock()
	defer s.balanceMu.Unlock()

	key, err := s.CreateKey(ctx, tariff)
	if err != nil {
		return nil, err
	}

	subscription = &models.Subscription{
		UserID:       user.ID,
		ServerID:     tariff.ServerID,
		TariffID:     tariff.ID,
		InitialPrice: price,
		KeyUUID:      key.ID,
		AccessUrl:    "",
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
			Type:   models.TransactionTypeDeposit,
			Status: models.TransactionStatusSuccess,
			Meta:   string(meta),
		})

		if err != nil {
			return err
		}

		err = s.storage.DecBalanceTx(ctx, tx, user.ID, price)
		if err != nil {
			return err
		}

		err = s.storage.UserBonusUsedTx(ctx, tx, user.ID)

		return err
	})

	return subscription, txErr
}

func (s *Service) GetSubscriptions(ctx context.Context, user models.User) (subscriptions []models.Subscription, err error) {
	return s.storage.GetSubscriptions(ctx, user.ID)
}

// TODO: EnableAutoProlongation
// TODO: DisableAutoProlongation

func (s *Service) ExpireSubscription(ctx context.Context, subscription models.Subscription) (err error) {

	if subscription.AutoProlong {
		err = s.ProlongSubscription(ctx, subscription)
		return err
	}

	// TODO: to outbox in db transaction
	err = s.DeactivateKey(ctx, subscription)
	if err != nil {
		return err
	}

	txErr := s.storage.WithTx(ctx, "ExpireSubscription", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.UpdateSubscriptionStatusTx(ctx, tx, subscription.ID, "expired")
		if err != nil {
			return err
		}

		return nil
	})

	// TODO: Send notification to user

	return txErr
}

func (s *Service) ProlongSubscription(ctx context.Context, subscription models.Subscription) (err error) {
	tariff, err := s.storage.GetTariff(ctx, subscription.TariffID)
	if err != nil {
		return err
	}

	user, err := s.storage.GetUser(ctx, subscription.UserID)
	if err != nil {
		return err
	}

	if (user.Balance - tariff.Price) < 0 {
		return ErrNotEnoughBalance
	}

	txErr := s.storage.WithTx(ctx, "ProlongSubscription", func(ctx context.Context, tx pg.Executor) error {
		meta, err := json.Marshal(models.TransactionMeta{
			IsProlongation: lo.ToPtr(true),
		})

		err = s.storage.CreateTransactionTx(ctx, tx, &models.Transaction{
			UserID: subscription.UserID,
			Amount: subscription.InitialPrice,
			Type:   models.TransactionTypeWithdrawal,
			Status: models.TransactionStatusSuccess,
			Meta:   string(meta),
		})

		err = s.storage.IncBalanceTx(ctx, tx, subscription.UserID, subscription.InitialPrice)
		if err != nil {
			return err
		}

		err = s.storage.ProlongSubscriptionTx(ctx, tx, subscription.ID, utils.CalcExpiredAt(tariff.Duration))
		if err != nil {
			return err
		}

		return nil
	})

	// TODO: Send notification to user

	return txErr
}

func (s *Service) NotifySubscriptionProlongation(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

func (s *Service) NotifySubscriptionBandwidthLimit(ctx context.Context, subscription models.Subscription, totalBytes uint64) (err error) {
	// TODO: Send notification to user
	return nil
}

func (s *Service) GetTariffPrice(ctx context.Context, tariff models.Tariff, user models.User) (price uint32, discountPercent uint8) {
	if user.PartnerID != nil && !user.BonusUsed {
		discountPercent := uint32(10) // @TODO: To config PARTNER_DISCOUNT_PERCENT
		price := tariff.Price - (tariff.Price * discountPercent / 100)
		price = tariff.Price - (tariff.Price % 10) // round to 10

		return price, uint8(discountPercent)
	}

	return tariff.Price, 0
}

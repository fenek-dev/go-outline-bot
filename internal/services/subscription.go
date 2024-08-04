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

	if (user.Balance - tariff.Price) < 0 {
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
		InitialPrice: tariff.Price,
		KeyUUID:      key.ID,
		AccessUrl:    "",
		ExpiredAt:    utils.CalcExpiredAt(tariff.Duration),
		Status:       "pending",
	}

	txErr := s.storage.WithTx(ctx, "CreateSubscription", func(ctx context.Context, tx pg.Executor) error {
		err = s.storage.CreateSubscriptionTx(ctx, tx, subscription)

		meta, err := json.Marshal(models.TransactionMeta{
			SubscriptionID: &subscription.ID,
		})

		if err != nil {
			return err
		}

		err = s.storage.CreateTransactionTx(ctx, tx, &models.Transaction{
			UserID: user.ID,
			Amount: tariff.Price,
			Type:   models.TransactionTypeDeposit,
			Status: models.TransactionStatusSuccess,
			Meta:   string(meta),
		})

		err = s.storage.DecBalanceTx(ctx, tx, user.ID, uint64(tariff.Price))
		if err != nil {
			return err
		}

		return err
	}, nil)

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
	}, nil)

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

		err = s.storage.IncBalanceTx(ctx, tx, subscription.UserID, uint64(subscription.InitialPrice))
		if err != nil {
			return err
		}

		err = s.storage.ProlongSubscriptionTx(ctx, tx, subscription.ID, utils.CalcExpiredAt(tariff.Duration))
		if err != nil {
			return err
		}

		return nil
	}, nil)

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

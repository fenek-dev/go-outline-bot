package services

import (
	"context"

	"github.com/fenek-dev/go-outline-bot/internal/models"
)

// Уведомление о том, что подписка скоро закончится (за 3 дня)
func (s *Service) NotifySubscriptionProlongationComing(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

// Уведомление о том, что подписка продлена
func (s *Service) NotifySubscriptionProlongation(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

// Уведомление о том, что трафик подписки скоро закончится (использовано 80%)
func (s *Service) NotifySubscriptionBandwidthLimitComing(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

// Уведомление о том, что исчерпан лимит трафика
func (s *Service) NotifySubscriptionBandwidthLimitReached(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

// Уведомление о том, что подписка закончилась
func (s *Service) NotifySubscriptionExpired(ctx context.Context, subscription models.Subscription) (err error) {
	// TODO: Send notification to user
	return nil
}

func (s *Service) NotifyPartnerAboutDeposit(ctx context.Context, transaction models.Transaction) (err error) {
	// TODO: Send notification to partner
	return nil
}

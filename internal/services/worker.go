package services

import (
	"context"
	"fmt"

	"github.com/fenek-dev/go-outline-bot/internal/models"
)

func (s *Service) CheckExpireSubscriptions(ctx context.Context) (err error) {
	// Get expired subscriptions without auto prolongation or without enough balance
	subscriptions, err := s.storage.GetExpiredSubscriptions(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("GetExpiredSubscriptions: %v", err))
		return fmt.Errorf("GetExpiredSubscriptions: %w", err)
	}

	// Diactivate subscriptions
	subscriptionIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		subscriptionIDs = append(subscriptionIDs, subscription.ID)
	}

	s.storage.UpdateSubscriptionsStatus(ctx, subscriptionIDs, models.SubscriptionStatusExpired)

	// Notify users about expired subscription
	for _, subscription := range subscriptions {
		go s.NotifySubscriptionExpired(ctx, subscription)
	}

	return nil
}

func (s *Service) CheckBandwidthLimits(ctx context.Context) (err error) {
	// Get subscriptions with bandwidth reached
	subscriptions, err := s.storage.GetSubscriptionsByBandwidthReached(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("GetSubscriptionsByBandwidthReached: %v", err))
		return err
	}

	// Diactivate subscriptions
	subscriptionIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		subscriptionIDs = append(subscriptionIDs, subscription.ID)
	}

	s.storage.UpdateSubscriptionsStatus(ctx, subscriptionIDs, models.SubscriptionStatusBandwidthReached)

	// Notify users about bandwidth reached
	for _, subscription := range subscriptions {
		go s.NotifySubscriptionBandwidthLimitReached(ctx, subscription)
	}

	return nil
}

func (s *Service) UpdateBandwidths(ctx context.Context) (err error) {
	// Get all servers
	servers, err := s.storage.GetAllServers(ctx)
	if err != nil {
		s.log.Error(fmt.Sprintf("GetAllServers: %v", err))
		return err
	}

	for _, server := range servers {
		metrics, err := s.GetBandwidthMetrics(ctx, server)
		if err != nil {
			s.log.Error(fmt.Sprintf("GetBandwidthMetrics: %v", err))
			continue
		}

		s.storage.UpdateSubscriptionsBandwidthByKeyID(ctx, server.ID, metrics)
	}

	return nil
}
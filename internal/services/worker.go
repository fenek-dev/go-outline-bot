package services

import "context"

func (s *Service) CheckProlongations(ctx context.Context) (err error) {
	// @TODO: Get expired subscriptions with auto prolongation and enough balance
	// @TODO: Dispatch prolongation for each subscription
	return nil
}

func (s *Service) CheckExpireSubscriptions(ctx context.Context) (err error) {
	// @TODO: Get expired subscriptions without auto prolongation or without enough balance
	// @TODO: Diactivate subscriptions
	// @TODD: Notify users about expired subscription
	return nil
}

func (s *Service) UpdateBandwidths(ctx context.Context) (err error) {
	// @TODO: Get all active subscriptions
	// @TODO: Get bandwidths for each subscription
	// @TODO: Batch update bandwidths
	return nil
}

func (s *Service) CheckBandwidthLimits(ctx context.Context) (err error) {
	// @TODO: Get all active subscriptions where bandwidth limit is reaching
	// @TODO: Notify users about bandwidth limit
	return nil
}

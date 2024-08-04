package services

import (
	"context"
	"github.com/fenek-dev/go-outline-bot/internal/models"
)

func (s *Service) GetTariffsByServer(ctx context.Context, serverId uint64) (tariffs []models.Tariff, err error) {
	return s.storage.GetTariffsByServer(ctx, serverId)
}

func (s *Service) GetTariff(ctx context.Context, tariffId uint64) (tariff models.Tariff, err error) {
	return s.storage.GetTariff(ctx, tariffId)
}

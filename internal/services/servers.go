package services

import (
	"context"
	"github.com/fenek-dev/go-outline-bot/internal/models"
)

func (s *Service) GetAllServers(ctx context.Context) (servers []models.Server, err error) {
	return s.storage.GetAllServers(ctx)
}

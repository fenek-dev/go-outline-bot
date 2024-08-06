package services

import (
	"context"
	"fmt"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/pkg/outline_client"
	"github.com/fenek-dev/go-outline-bot/pkg/utils"
	"github.com/google/uuid"
)

func (s *Service) CreateKey(ctx context.Context, tariff models.Tariff) (key *outline_client.OutlineKey, err error) {
	server, err := s.storage.GetServer(ctx, tariff.ServerID)
	if err != nil {
		return nil, fmt.Errorf("GetServer: %w", err)
	}

	client, err := outline_client.NewOutlineVPN(server.IP, server.APIKey)
	if err != nil {
		return nil, fmt.Errorf("CreateKey: %w", err)
	}

	key, err = client.AddKey(&outline_client.OutlineKey{
		ID:   uuid.New().String(),
		Name: tariff.Name,
	})

	if err != nil {
		return nil, fmt.Errorf("AddKey: %w", err)
	}

	err = client.SetKeyLimit(key.ID, utils.ConvertGBtoBytes(float32(tariff.Bandwidth)))
	if err != nil {
		return nil, fmt.Errorf("SetKeyLimit: %w", err)
	}

	return key, nil
}

func (s *Service) DeactivateKey(ctx context.Context, subscription models.Subscription) (err error) {
	server, err := s.storage.GetServer(ctx, subscription.ServerID)
	if err != nil {
		return fmt.Errorf("GetServer: %w", err)
	}

	client, err := outline_client.NewOutlineVPN(server.IP, server.APIKey)
	if err != nil {
		return fmt.Errorf("NewOutlineVPN: %w", err)
	}

	// TODO: Retries

	// Set limit to 1 byte
	err = client.SetKeyLimit(subscription.KeyUUID, 1)
	if err != nil {
		return fmt.Errorf("SetKeyLimit: %w", err)
	}

	return nil
}

func (s *Service) GetBandwidthMetrics(ctx context.Context, server models.Server) (metrics map[string]uint64, err error) {
	client, err := outline_client.NewOutlineVPN(server.IP, server.APIKey)
	if err != nil {
		return nil, fmt.Errorf("NewOutlineVPN: %w", err)
	}

	response, err := client.GetTransferMetrics()
	if err != nil {
		return nil, fmt.Errorf("GetTransferMetrics: %w", err)
	}

	return response.BytesTransferredByUserId, nil
}

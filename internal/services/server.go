package services

import (
	"context"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/pkg/outline_client"
	"github.com/fenek-dev/go-outline-bot/pkg/utils"
	"github.com/google/uuid"
)

func (s *Service) CreateKey(ctx context.Context, tariff models.Tariff) (key *outline_client.OutlineKey, err error) {
	server, err := s.storage.GetServer(ctx, tariff.ServerID)
	if err != nil {
		return nil, err
	}

	client, err := outline_client.NewOutlineVPN(server.IP, server.APIKey)
	if err != nil {
		return nil, err
	}

	key, err = client.AddKey(&outline_client.OutlineKey{
		ID:   uuid.New().String(),
		Name: tariff.Name,
	})

	if err != nil {
		return nil, err
	}

	err = client.SetKeyLimit(key.ID, utils.ConvertGBtoBytes(float32(tariff.Bandwidth)))
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *Service) DeactivateKey(ctx context.Context, subscription models.Subscription) (err error) {
	server, err := s.storage.GetServer(ctx, subscription.ServerID)
	if err != nil {
		return err
	}

	client, err := outline_client.NewOutlineVPN(server.IP, server.APIKey)
	if err != nil {
		return err
	}

	// TODO: Retries

	// Set limit to 1 byte
	err = client.SetKeyLimit(subscription.KeyUUID, 1)
	if err != nil {
		return err
	}

	return nil
}

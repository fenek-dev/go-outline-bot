package services

import (
	"context"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func (s *Service) CreateUser(ctx context.Context, user *telebot.User) (err error) {
	return s.storage.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, userID uint64) (user models.User, err error) {
	return s.storage.GetUser(ctx, userID)
}

func (s *Service) SetUserPhone(ctx context.Context, userID uint64, phone string) (err error) {
	return s.storage.SetUserPhone(ctx, userID, phone)
}

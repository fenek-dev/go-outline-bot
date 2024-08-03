package services

import (
	"context"
	"gopkg.in/telebot.v3"
)

func (s *Service) CreateUser(ctx context.Context, user *telebot.User) (err error) {
	return s.storage.CreateUser(ctx, user)
}

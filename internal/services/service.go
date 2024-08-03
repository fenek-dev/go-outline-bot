package services

import (
	"context"
	"gopkg.in/telebot.v3"
	"log/slog"
)

type Option func(*Service)

type Storage interface {
	CreateUser(ctx context.Context, user *telebot.User) (err error)
}

type Service struct {
	storage Storage

	log *slog.Logger
}

func New(storage Storage, opts ...Option) *Service {
	s := &Service{
		storage: storage,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Service) {
		s.log = logger
	}
}

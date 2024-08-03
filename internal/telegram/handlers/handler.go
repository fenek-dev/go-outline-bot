package handlers

import (
	"context"
	"gopkg.in/telebot.v3"
	"log/slog"
	"time"
)

type Option func(*Handlers)

type Service interface {
	CreateUser(ctx context.Context, user *telebot.User) (err error)
}

type Handlers struct {
	service Service

	log     *slog.Logger
	timeout time.Duration
}

func New(service Service, opts ...Option) *Handlers {
	defaultHandlers := &Handlers{
		service: service,
		timeout: time.Second * 10,
	}

	for _, opt := range opts {
		opt(defaultHandlers)
	}

	return defaultHandlers
}

func WithTimeout(duration time.Duration) Option {
	return func(h *Handlers) {
		h.timeout = duration
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(h *Handlers) {
		h.log = logger
	}
}

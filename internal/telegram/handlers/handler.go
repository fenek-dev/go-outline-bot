package handlers

import (
	"context"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"gopkg.in/telebot.v3"
	"log/slog"
	"time"
)

type Option func(*Handlers)

type Service interface {
	CreateUser(ctx context.Context, user *telebot.User) (err error)

	GetBalance(ctx context.Context, userId uint64) (balance uint32, err error)

	GetAllServers(ctx context.Context) (servers []models.Server, err error)

	GetTariff(ctx context.Context, tariffId uint64) (tariff models.Tariff, err error)
	GetTariffsByServer(ctx context.Context, serverId uint64) (tariffs []models.Tariff, err error)
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

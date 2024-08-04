package services

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage/pg"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
	"gopkg.in/telebot.v3"
)

type Option func(*Service)

type Storage interface {
	CreateUser(ctx context.Context, user *telebot.User) (err error)
	GetUser(ctx context.Context, userID uint64) (user models.User, err error)
	SetUserBonusUsedTx(ctx context.Context, tx pg.Executor, userID uint64) (err error)

	GetBalance(ctx context.Context, userID uint64) (balance uint32, err error)
	IncBalanceTx(ctx context.Context, tx pg.Executor, userID uint64, amount uint32) (err error)
	DecBalanceTx(ctx context.Context, tx pg.Executor, userID uint64, amount uint32) (err error)

	GetTariff(ctx context.Context, tariffID uint64) (tariff models.Tariff, err error)
	GetTariffsByServer(ctx context.Context, serverId uint64) (tariffs []models.Tariff, err error)

	GetTransaction(ctx context.Context, transactionID uint64) (transaction models.Transaction, err error)
	GetTransactionsByUser(ctx context.Context, userID uint64) (transactions []models.Transaction, err error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (err error)
	CreateTransactionTx(ctx context.Context, tx pg.Executor, transaction *models.Transaction) (terr error)
	UpdateTransactionStatusTx(ctx context.Context, tx pg.Executor, transactionID uint64, status models.TransactionStatus) (err error)

	GetSubscriptionsByUser(ctx context.Context, userID uint64) (subscriptions []models.Subscription, err error)
	GetProlongableSubscriptions(ctx context.Context) (subscriptions []models.Subscription, err error)
	GetExpiredSubscriptions(ctx context.Context) (subscriptions []models.Subscription, err error)
	GetSubscriptionsByBandwidthReached(ctx context.Context) (subscriptions []models.Subscription, err error)
	CreateSubscription(ctx context.Context, subscription *models.Subscription) (err error)
	CreateSubscriptionTx(ctx context.Context, tx pg.Executor, subscription *models.Subscription) (err error)
	UpdateSubscriptionsStatus(ctx context.Context, subscriptionIDs []uint64, status models.SubscriptionStatus) (err error)
	UpdateSubscriptionStatusTx(ctx context.Context, tx pg.Executor, subscriptionID uint64, status models.SubscriptionStatus) (err error)
	UpdateSubscriptionsBandwidthByKeyID(ctx context.Context, serverID uint64, metrics map[string]uint64) (err error)
	ProlongSubscriptionTx(ctx context.Context, tx pg.Executor, subscriptionID uint64, expiredAt time.Time) (err error)

	GetServer(ctx context.Context, serverID uint64) (server models.Server, err error)
	GetAllServers(ctx context.Context) (servers []models.Server, err error)

	WithTx(ctx context.Context, label string, fn func(ctx context.Context, tx pg.Executor) error, options *pgx.TxOptions, opts ...pg.TxOption) (err error)
}

type Service struct {
	storage       Storage
	paymentClient *payment_service.Client

	log *slog.Logger

	balanceMu sync.Mutex
}

func New(storage Storage, paymentClient *payment_service.Client, opts ...Option) *Service {
	s := &Service{
		storage:       storage,
		paymentClient: paymentClient,
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

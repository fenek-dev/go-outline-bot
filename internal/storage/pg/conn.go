package pg

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*pgxpool.Config, *Postgres)

type Postgres struct {
	conn *pgxpool.Pool

	log *slog.Logger
}

func (p *Postgres) GetUser(ctx context.Context, userID uint64) (user models.User, err error) {
	//TODO implement me
	return user, err
}

func (p *Postgres) IncBalanceTx(ctx context.Context, tx Executor, userID uint64, amount uint64) (err error) {
	//TODO implement me
	return err
}

func (p *Postgres) DecBalanceTx(ctx context.Context, tx Executor, userID uint64, amount uint64) (err error) {
	//TODO implement me
	return err
}

func (p *Postgres) CreateTransaction(ctx context.Context, transaction *models.Transaction) (err error) {
	//TODO implement me
	return err
}

func (p *Postgres) CreateTransactionTx(ctx context.Context, tx Executor, transaction *models.Transaction) (err error) {
	//TODO implement me
	return err
}

func (p *Postgres) GetSubscriptions(ctx context.Context, userID uint64) (subscriptions []models.Subscription, err error) {
	//TODO implement me
	return subscriptions, err
}

func (p *Postgres) UpdateSubscriptionStatusTx(ctx context.Context, tx Executor, subscriptionID uint64, status models.SubscriptionStatus) (err error) {
	//TODO implement me
	return err
}

func (p *Postgres) ProlongSubscriptionTx(ctx context.Context, tx Executor, subscriptionID uint64, expiredAt time.Time) (err error) {
	//TODO implement me
	return err
}

func New(ctx context.Context, DBUrl string, opts ...Option) *Postgres {

	cfg, err := pgxpool.ParseConfig(DBUrl)
	if err != nil {
		panic(fmt.Sprintf("can not parse db config: %s", err.Error()))
	}

	p := &Postgres{conn: nil}

	for _, opt := range opts {
		opt(cfg, p)
	}

	conn, err := pgxpool.New(ctx, DBUrl)

	if err != nil {
		panic(fmt.Sprintf("can not connect to db: %s", err.Error()))
	}

	if err := conn.Ping(ctx); err != nil {
		panic(fmt.Sprintf("db ping didn't work: %s", err.Error()))
	}

	return &Postgres{conn: conn}
}

func WithMaxConnections(maxConn int32) Option {
	return func(config *pgxpool.Config, _ *Postgres) {
		config.MaxConns = maxConn
	}
}

func WithMinConnections(minConn int32) Option {
	return func(config *pgxpool.Config, _ *Postgres) {
		config.MinConns = minConn
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(_ *pgxpool.Config, p *Postgres) {
		p.log = logger
	}
}

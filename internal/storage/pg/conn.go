package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Option func(*pgxpool.Config, *Postgres)

type Postgres struct {
	conn *pgxpool.Pool

	log *slog.Logger
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

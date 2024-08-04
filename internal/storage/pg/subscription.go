package pg

import (
	"context"
	"errors"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) CreateSubscription(ctx context.Context, subscription *models.Subscription) (err error) {
	return p.CreateSubscriptionTx(ctx, p.conn, subscription)
}

func (p *Postgres) CreateSubscriptionTx(ctx context.Context, tx Executor, subscription *models.Subscription) (err error) {
	_, err = p.conn.Exec(ctx, "INSERT INTO subscriptions (id, first_name, username) VALUES ($1, $2, $3)")
	var pgErr *pgconn.PgError
	// Error when user already exists in db
	if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode {
		return storage.ErrUserAlreadyExists
	}

	return err
}

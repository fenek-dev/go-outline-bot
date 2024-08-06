package pg

import (
	"context"
	"errors"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"time"
)

func (p *Postgres) CreateSubscription(ctx context.Context, subscription *models.Subscription) (err error) {
	return p.CreateSubscriptionTx(ctx, p.conn, subscription)
}
func (p *Postgres) CreateSubscriptionTx(ctx context.Context, tx Executor, subscription *models.Subscription) (err error) {
	err = tx.QueryRow(
		ctx,
		"INSERT INTO subscriptions (user_id, server_id, tariff_id, initial_price, key_uuid, access_url, status, expired_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		subscription.UserID,
		subscription.ServerID,
		subscription.TariffID,
		subscription.InitialPrice,
		subscription.KeyUUID,
		subscription.AccessUrl,
		subscription.Status,
		subscription.ExpiredAt,
	).Scan(&subscription.ID)

	return err
}

func (p *Postgres) GetSubscriptionsByUser(ctx context.Context, userID uint64) (subscriptions []models.Subscription, err error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&subscriptions, rows)
	if err != nil {
		return nil, err
	}

	return subscriptions, err
}

func (p *Postgres) UpdateSubscriptionStatusTx(ctx context.Context, tx Executor, subscriptionID uint64, status models.SubscriptionStatus) (err error) {
	_, err = tx.Exec(ctx, "UPDATE subscriptions SET status = $1 WHERE id = $2", status, subscriptionID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrSubscriptionNotFound
	}

	return err
}

func (p *Postgres) ProlongSubscriptionTx(ctx context.Context, tx Executor, subscriptionID uint64, expiredAt time.Time) (err error) {
	_, err = tx.Exec(ctx, "UPDATE subscriptions SET expired_at = $1 WHERE id = $2", expiredAt, subscriptionID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrSubscriptionNotFound
	}

	return err
}

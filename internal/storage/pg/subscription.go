package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) CreateSubscription(ctx context.Context, subscription *models.Subscription) (err error) {
	return p.CreateSubscriptionTx(ctx, p.conn, subscription)
}

func (p *Postgres) CreateSubscriptionTx(ctx context.Context, tx Executor, subscription *models.Subscription) (err error) {
	err = tx.QueryRow(
		ctx,
		"INSERT INTO subscriptions (user_id, server_id, tariff_id, initial_price, key_uuid, accessurl, status, expired_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
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

func (p *Postgres) GetExpiredSubscriptions(ctx context.Context) (subscriptions []models.Subscription, err error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM subscriptions WHERE expired_at < $1", time.Now())
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&subscriptions, rows)
	if err != nil {
		return nil, err
	}

	return subscriptions, err
}

func (p *Postgres) GetSubscriptionsByBandwidthReached(ctx context.Context) (subscriptions []models.Subscription, err error) {
	rows, err := p.conn.Query(
		ctx,
		`
		SELECT s.* FROM subscriptions s 
    	JOIN public.tariffs t ON s.tariff_id = t.id 
		    WHERE t.bandwidth <= s.bandwidth_spent
		`,
	)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&subscriptions, rows)
	if err != nil {
		return nil, err
	}

	return subscriptions, err
}

func (p *Postgres) GetProlongableSubscriptions(ctx context.Context) (subscriptions []models.Subscription, err error) {
	rows, err := p.conn.Query(
		ctx,
		"SELECT * FROM subscriptions WHERE expired_at < $1 AND auto_prolong = TRUE",
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&subscriptions, rows)
	if err != nil {
		return nil, err
	}

	return subscriptions, err
}

func (p *Postgres) UpdateSubscriptionsStatus(ctx context.Context, subscriptionIDs []uint64, status models.SubscriptionStatus) (err error) {
	_, err = p.conn.Exec(ctx, "UPDATE subscriptions SET status = $1 WHERE id = ANY($2)", status, subscriptionIDs)

	return err
}

func (p *Postgres) UpdateSubscriptionsBandwidthByKeyID(ctx context.Context, serverID uint64, metrics map[string]uint64) (err error) {
	values := make([]string, 0, len(metrics))
	for keyID, bandwidth := range metrics {
		values = append(values, fmt.Sprintf("(%d, %s)", bandwidth, keyID))
	}
	result := strings.Join(values, ",")

	_, err = p.conn.Exec(
		ctx,
		`UPDATE subscriptions s
		SET bandwidth_spent=tmp.spent 
		FROM (VALUES $1) AS tmp (spent, key) WHERE s.key_uuid=tmp.key AND s.server_id=$2`,
		result,
		serverID,
	)

	return nil
}

func (p *Postgres) TrialSubscriptionExists(ctx context.Context, userID uint64) (has bool, err error) {
	var count int
	err = p.conn.QueryRow(ctx, "SELECT COUNT(id) FROM subscriptions WHERE is_trial = true AND user_id = $1", userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

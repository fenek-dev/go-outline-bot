package pg

import (
	"context"
	"errors"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) GetBalance(ctx context.Context, userId uint64) (balance uint32, err error) {
	err = p.conn.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userId).Scan(&balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrUserNotFound
		}
		return 0, err
	}

	return balance, nil
}

func (p *Postgres) IncBalanceTx(ctx context.Context, tx Executor, userID uint64, amount uint32) (err error) {
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	return err
}

func (p *Postgres) DecBalanceTx(ctx context.Context, tx Executor, userID uint64, amount uint32) (err error) {
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == CheckViolationCode {
		return storage.ErrNegativeBalance
	}

	return err
}

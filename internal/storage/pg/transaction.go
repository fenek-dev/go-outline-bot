package pg

import (
	"context"
	"errors"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) CreateTransactionTx(ctx context.Context, tx Executor, transaction *models.Transaction) (err error) {
	err = tx.QueryRow(
		ctx,
		"INSERT INTO transactions (user_id, amount, type, status, external_id, meta) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		transaction.UserID,
		transaction.Amount,
		transaction.Type,
		transaction.Status,
		transaction.ExternalID,
		transaction.Meta,
	).Scan(&transaction.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrTransactionNotFound
	}

	return err
}

func (p *Postgres) CreateTransaction(ctx context.Context, transaction *models.Transaction) (err error) {
	return p.CreateTransactionTx(ctx, p.conn, transaction)
}

func (p *Postgres) GetTransaction(ctx context.Context, transactionID uint64) (transaction models.Transaction, err error) {
	err = pgxscan.Get(ctx, p.conn, &transaction, "SELECT * FROM transactions WHERE id = $1", transactionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return transaction, storage.ErrTransactionNotFound
		}
		return
	}

	return transaction, err
}

func (p *Postgres) GetTransactionsByUser(ctx context.Context, userID uint64) (transactions []models.Transaction, err error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM transactions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&transactions, rows)
	if err != nil {
		return nil, err
	}

	return transactions, err
}

func (p *Postgres) UpdateTransactionStatusTx(ctx context.Context, tx Executor, transactionID uint64, status models.TransactionStatus) (err error) {
	_, err = tx.Exec(ctx, "UPDATE transactions SET status = $1 WHERE id = $2", status, transactionID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrTransactionNotFound
	}

	return err
}

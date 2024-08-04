package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type TxOption func(ctx context.Context, tx pgx.Tx) error

func (s *Postgres) WithTx(
	ctx context.Context,
	label string,
	fn func(ctx context.Context, tx Executor) error,
	opts *pgx.TxOptions,
	txOpts ...TxOption,
) (err error) {
	var tx pgx.Tx
	if opts != nil {
		tx, err = s.conn.BeginTx(ctx, *opts)
	} else {
		tx, err = s.conn.Begin(ctx)
	}
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			if rollBackErr := tx.Rollback(ctx); err != nil {
				panic(fmt.Errorf("rollback failed: %v, during panic: %v", rollBackErr, p))
			}
			panic(p)
		}

		if err != nil {
			if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
				err = fmt.Errorf("rollback failed: %v, during failed transaction: %v", rollBackErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %v, during failed transaction: %v", commitErr, err)
			}
		}
	}()

	for _, txOpt := range txOpts {
		if err = txOpt(ctx, tx); err != nil {
			return err
		}
	}

	return fn(ctx, tx)
}

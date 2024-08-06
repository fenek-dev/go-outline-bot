package pg

import (
	"context"
	"errors"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gopkg.in/telebot.v3"
)

func (p *Postgres) CreateUser(ctx context.Context, user *telebot.User) (err error) {

	_, err = p.conn.Exec(ctx, "INSERT INTO users (id, first_name, username) VALUES ($1, $2, $3)", user.ID, user.FirstName, user.Username)
	var pgErr *pgconn.PgError
	// Error when user already exists in db
	if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode {
		return storage.ErrUserAlreadyExists
	}

	return err

}

func (p *Postgres) GetUser(ctx context.Context, userID uint64) (user models.User, err error) {
	err = pgxscan.Get(ctx, p.conn, &user, "SELECT * FROM users WHERE id = $1", userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, storage.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (p *Postgres) SetUserBonusUsedTx(ctx context.Context, tx Executor, userID uint64) (err error) {
	_, err = tx.Exec(ctx, "UPDATE users SET bonus_used = true WHERE id = $1", userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	return err
}

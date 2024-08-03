package pg

import (
	"context"
	"errors"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"gopkg.in/telebot.v3"
)

func (p *Postgres) CreateUser(ctx context.Context, user *telebot.User) (err error) {

	_, err = p.conn.Exec(ctx, "INSERT INTO users (id, first_name) VALUES ($1, $2)", user.ID, user.FirstName)

	var pgErr *pgconn.PgError
	// Error when user already exists in db
	if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode {
		return storage.ErrUserAlreadyExists
	}

	return err

}

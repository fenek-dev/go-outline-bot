package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/markup"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"gopkg.in/telebot.v3"
	"time"
)

func (h *Handlers) Start(c telebot.Context) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		user        = c.Sender()
	)
	defer cancel()

	err = h.service.CreateUser(ctx, user)

	if err != nil && !errors.Is(err, storage.ErrUserAlreadyExists) {
		c.Send("Something gone wrong, try again")
		return fmt.Errorf("unexpected error: %e", err)
	}

	return c.Send("Welcome to the club, buddy", markup.Menu)
}

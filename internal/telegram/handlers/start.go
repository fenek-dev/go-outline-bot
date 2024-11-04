package handlers

import (
	"context"
	"errors"
	"fmt"
	markup2 "github.com/fenek-dev/go-outline-bot/internal/markup"
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
		h.log.Error("error on start", "error", err)
		_ = c.Send("Ой, что-то пошло не так. Попробуйте еще раз", markup2.OnlyClose)
		return fmt.Errorf("unexpected error: %w", err)
	}

	return c.Send("Привет, я бот который бла бла бла", markup2.Menu)
}

func (h *Handlers) Close(c telebot.Context) error {
	return c.Delete()
}

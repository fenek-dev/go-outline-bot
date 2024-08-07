package handlers

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	t "gopkg.in/telebot.v3"
)

func (h *Handlers) OpenBalance(c t.Context) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	ID := uint64(c.Sender().ID)

	balance, err := h.service.GetBalance(ctx, ID)

	if err != nil {
		h.log.Error("can not get balance", "error", err)
		return err
	}

	return c.Send(fmt.Sprintf("Ваш баланс: %d RUB", balance), markup.Balance)
}

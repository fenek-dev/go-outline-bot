package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	"gopkg.in/telebot.v3"
	"math"
	"strconv"
	"time"
)

func (h *Handlers) OpenKeysMenu(c telebot.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	userId := c.Sender().ID

	subs, err := h.service.GetSubscriptionsByUser(ctx, uint64(userId))

	if err != nil {
		h.log.Error("Can not get subs", "error", err)
		return c.Send("Произошла ошибка. Попробуйте еще раз", markup.OnlyClose)
	}

	rows := make([]telebot.Row, 0, len(subs)+2)
	rows = append(rows, markup.KeysMenu.Row(markup.KeysGetNewBtn))

	for _, sub := range subs {
		if time.Now().After(sub.ExpiredAt) {
			continue
		}

		btn := telebot.Btn{
			Text:   fmt.Sprintf("%d, (%vд)", sub.ID, math.Trunc(sub.ExpiredAt.Sub(time.Now()).Hours()/24)),
			Data:   strconv.FormatUint(sub.ID, 10),
			Unique: markup.KeyItem.Unique,
		}

		rows = append(rows, markup.KeysMenu.Row(btn))
	}

	rows = append(rows, markup.KeysMenu.Row(markup.CloseBtn))

	markup.KeysMenu.Inline(
		rows...,
	)

	return c.Send("Ключи:", markup.KeysMenu)
}

func (h *Handlers) CloseKeysMenu(c telebot.Context) error {
	return c.Delete()
}

func (h *Handlers) OpenKeyInfo(c telebot.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	keyID := c.Data()

	id, err := strconv.ParseUint(keyID, 10, 64)
	if err != nil {
		return err
	}

	sub, err := h.service.GetSubscription(ctx, id)
	if err != nil {
		h.log.Error("Can not get sub", "error", err)
		if errors.Is(err, storage.ErrSubscriptionNotFound) {
			return c.Send("Ключ не найдена", markup.OnlyClose)
		}
		return c.Send("Произошла ошибка. Попробуйте еще раз", markup.OnlyClose)
	}

	return c.Send(fmt.Sprintf("Ключ %d, действителен до %v", sub.ID, sub.ExpiredAt), markup.OnlyClose)
}

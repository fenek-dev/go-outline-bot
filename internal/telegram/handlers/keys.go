package handlers

import (
	"github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	"gopkg.in/telebot.v3"
)

func (h *Handlers) OpenKeysMenu(c telebot.Context) error {
	return c.Send("Ключи:", markup.KeysMenu)
}

func (h *Handlers) CloseKeysMenu(c telebot.Context) error {
	return c.Delete()
}

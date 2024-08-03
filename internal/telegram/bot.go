package telegram

import (
	"fmt"
	"github.com/fenek-dev/go-outline-bot/configs"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/handlers"
	"gopkg.in/telebot.v3"
)

func InitBot(cfg *configs.TelegramConfig, h *handlers.Handlers) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.Webhook{},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("can not connect to telegram bot: %w", err)
	}

	b.Handle("/start", h.Start)

	return b, nil
}

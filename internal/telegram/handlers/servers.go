package handlers

import (
	"github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	"golang.org/x/net/context"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (h *Handlers) OpenServersMenu(c telebot.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	servers, err := h.service.GetAllServers(ctx)

	if err != nil {
		return err
	}

	rows := make([]telebot.Row, 0, len(servers))

	for _, server := range servers {
		btn := telebot.Btn{
			Text:   server.Name,
			Data:   strconv.FormatUint(server.ID, 10),
			Unique: markup.ServerItem.Unique,
		}

		rows = append(rows, markup.ServersMenu.Row(btn))
	}

	rows = append(rows, markup.ServersMenu.Row(markup.ServersBackBtn))

	markup.ServersMenu.Inline(
		rows...,
	)

	return c.Edit("Сервера:", markup.ServersMenu)
}

func (h *Handlers) BackServersMenu(c telebot.Context) error {
	return c.Edit("Ключи:", markup.KeysMenu)
}

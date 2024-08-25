package handlers

import (
	markup2 "github.com/fenek-dev/go-outline-bot/internal/markup"
	"golang.org/x/net/context"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (h *Handlers) OpenServersMenu(c telebot.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	servers, err := h.service.GetAllServers(ctx)

	if err != nil {
		h.log.Error("can not get servers", "error", err)
		return err
	}

	rows := make([]telebot.Row, 0, len(servers))

	for _, server := range servers {
		btn := telebot.Btn{
			Text:   server.Name,
			Data:   strconv.FormatUint(server.ID, 10),
			Unique: markup2.ServerItem.Unique,
		}

		rows = append(rows, markup2.ServersMenu.Row(btn))
	}

	rows = append(rows, markup2.ServersMenu.Row(markup2.ServersBackBtn))

	markup2.ServersMenu.Inline(
		rows...,
	)

	return c.Edit("Сервера:", markup2.ServersMenu)
}

func (h *Handlers) BackServersMenu(c telebot.Context) error {
	return c.Edit("Ключи:", markup2.KeysMenu)
}

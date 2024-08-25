package handlers

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/markup"
	t "gopkg.in/telebot.v3"
)

func (h *Handlers) HandlePhone(c t.Context) (err error) {
	fmt.Println("HandlePhone")
	ctx := context.Background()
	contact := c.Message().Contact

	if contact == nil {
		return c.Send("Пожалуйста, отправьте номер телефона")
	}

	id := uint64(c.Sender().ID)
	err = h.service.SetUserPhone(ctx, id, contact.PhoneNumber)
	if err != nil {
		h.log.Error("can not set phone", "error", err)
		return c.Send("Произошла ошибка при добавлении номера телефона")
	}

	return c.Send("Номер телефона успешно добавлен ✅", markup.Menu)
}

func (h *Handlers) ClosePhone(c t.Context) (err error) {
	return c.Send("Добавление номера телефона отменено", markup.Menu)
}

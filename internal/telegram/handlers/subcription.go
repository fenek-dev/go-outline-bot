package handlers

import (
	"gopkg.in/telebot.v3"
)

func (h *Handlers) BuySubscription(c telebot.Context) error {
	//ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	//defer cancel()

	//tariffID := c.Data()

	return c.Delete()
}

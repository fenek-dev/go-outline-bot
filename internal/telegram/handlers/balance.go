package handlers

import (
	"context"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/markup"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/state"
	"github.com/fenek-dev/go-outline-bot/pkg/utils"
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

func (h *Handlers) TopUpBalance(c t.Context) (err error) {
	user, err := h.service.GetUser(context.Background(), uint64(c.Sender().ID))
	if err != nil {
		h.log.Error("can not get user", "error", err)
		return c.Send("Странно, вас не нашлось в нашей базе данных. Попробуйте начать все заново /start", markup.OnlyClose)
	}

	if user.Phone == nil || *user.Phone == "" {
		return c.Send("Для пополнения балансе необходимо указать номер телефона на который придут уведомления об операции:", markup.Phone)
	}

	state.SetUserCallback(c.Sender().ID, h.HandleTopUpBalanceAmount)
	return c.Send("Введи сумму пополнения (мин. 50 руб.). Комиcсия не взимается.")
}

func (h *Handlers) TopUpClose(c t.Context) (err error) {
	state.DeleteUserCallback(c.Sender().ID)
	return c.Send("Пополнение баланса отменено", markup.OnlyClose)
}

func (h *Handlers) HandleTopUpBalanceAmount(c t.Context) error {
	ctx := context.Background()
	amount, err := utils.ParseAmount(c.Text())
	if err != nil {
		return c.Send("Некорректная сумма. Попробуй еще раз")
	}

	if amount < 50 {
		return c.Send("Минимальная сумма пополнения 50 руб. Попробуй еще раз")
	}

	user, err := h.service.GetUser(context.Background(), uint64(c.Sender().ID))
	if err != nil {
		h.log.Error("can not get user", "error", err)
		return c.Send("Произошла ошибка. Попробуй еще раз", markup.OnlyClose)
	}

	uri, err := h.service.RequestDeposit(ctx, user, uint32(amount))
	if err != nil {
		h.log.Error("can not request deposit", "error", err)
		return c.Send("Произошла ошибка. Попробуй еще раз", markup.OnlyClose)
	}

	state.DeleteUserCallback(c.Sender().ID)

	markup.Confirm.Inline(
		markup.Confirm.Row(markup.Confirm.URL("💸 Оплатить", uri)),
	)

	return c.Send(fmt.Sprintf("Пополнить баланс на %d руб.?", amount), markup.Confirm)

}

package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/storage"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (h *Handlers) BuySubscription(c telebot.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	userId := c.Sender().ID

	user, err := h.service.GetUser(ctx, uint64(userId))

	if err != nil {
		h.log.Error("can not get user", "error", err)
		if errors.Is(err, storage.ErrUserNotFound) {
			return c.Send("Странно, в нашей базе данных вас нет. Для регистрации введите команду /start", markup.OnlyClose)
		}

		_ = c.Send("Произошла ошибка. Попробуйте еще раз", markup.OnlyClose)
		return err
	}

	tariffID := c.Data()

	id, err := strconv.ParseUint(tariffID, 10, 64)
	if err != nil {
		h.log.Error("can not parse tariff id", "error", err)
		_ = c.Send("Произошла ошибка. Попробуйте еще раз", markup.OnlyClose)
		return err
	}

	tariff, err := h.service.GetTariff(ctx, id)

	if err != nil {
		h.log.Error("can not get tariff", "error", err)
		if errors.Is(err, storage.ErrTariffNotFound) {
			_ = c.Send("Такого тарифа не существует", markup.OnlyClose)
			return err
		}
		_ = c.Send("Произошла ошибка. Попробуйте еще раз", markup.OnlyClose)
		return err
	}

	if user.Balance < tariff.Price {
		return c.Send(fmt.Sprintf("На вашем счету недостаточно средств. Стоимость тарифа %d RUB. Ваш баланс: %d RUB", tariff.Price, user.Balance), markup.Balance)
	}

	sub, err := h.service.CreateSubscription(ctx, user, tariff)

	if err != nil {
		h.log.Error("can not create subscription", "error", err)
		_ = c.Send("Не удалось создать подписку", markup.OnlyClose)
		return err
	}

	_ = c.Edit(fmt.Sprintf("Поздравляем! Вы успешно приобрели тариф %s на %d дней. Ваш ключ: %s", tariff.Name, tariff.Duration, sub.AccessUrl), markup.OnlyClose)
	return c.Send(fmt.Sprintf("Ваш ключ: %s", sub.AccessUrl), markup.KeysMenu)
}

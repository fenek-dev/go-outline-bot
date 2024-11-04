package handlers

import (
	"context"
	"fmt"
	markup2 "github.com/fenek-dev/go-outline-bot/internal/markup"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (h *Handlers) OpenTariffsMenu(c telebot.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	serverID := c.Data()

	id, err := strconv.ParseUint(serverID, 10, 64)
	if err != nil {
		return err
	}

	tariffs, err := h.service.GetTariffsByServer(ctx, id)

	rows := make([]telebot.Row, 0, len(tariffs))

	for _, tariff := range tariffs {
		if !tariff.Active {
			continue
		}
		btn := telebot.Btn{
			Text:   fmt.Sprintf("%s, %vGB, %v₽/%vд", tariff.Name, tariff.Bandwidth, tariff.Price, tariff.Duration),
			Data:   strconv.FormatUint(tariff.ID, 10),
			Unique: markup2.TariffItem.Unique,
		}

		rows = append(rows, markup2.TariffsMenu.Row(btn))
	}

	rows = append(rows, markup2.TariffsMenu.Row(markup2.TariffsBackBtn))

	markup2.TariffsMenu.Inline(
		rows...,
	)

	return c.Edit("Тарифы:", markup2.TariffsMenu)
}

func (h *Handlers) BackTariffsMenu(c telebot.Context) error {
	return c.Edit("Сервера:", markup2.ServersMenu)
}

func (h *Handlers) OpenTariff(c telebot.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	tariffID := c.Data()

	id, err := strconv.ParseUint(tariffID, 10, 64)
	if err != nil {
		return err
	}

	tariff, err := h.service.GetTariff(ctx, id)
	if err != nil {
		return err
	}

	err = c.Delete()
	if err != nil {
		return err
	}

	markup2.TariffInfo.Inline(
		markup2.TariffInfo.Row(markup2.WithData(tariffID, markup2.TariffBuyBtn)),
		markup2.TariffInfo.Row(markup2.CloseBtn),
	)

	return c.Send(fmt.Sprintf("Тариф %s, %vGB, %v₽/%vд", tariff.Name, tariff.Bandwidth, tariff.Price, tariff.Duration), markup2.TariffInfo)
}

func (h *Handlers) BackTariff(c telebot.Context) error {
	return c.Edit("Тарифы:", markup2.TariffsMenu)
}

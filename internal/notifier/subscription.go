package notifier

import (
	"context"
	"fmt"
	m "github.com/fenek-dev/go-outline-bot/internal/markup"
	t "gopkg.in/telebot.v3"
	"strconv"
)

func (n *Notifier) NotifySubscriptionProlongationComing(ctx context.Context, data NotifyDTO) (err error) {
	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, "Подписка на тариф скоро будет продлена")
	return err
}

func (n *Notifier) NotifySubscriptionProlongation(ctx context.Context, data NotifyDTO) (err error) {
	id := strconv.FormatUint(data.SubscriptionID, 10)
	m.ProlongSuccess.Inline(
		m.ProlongSuccess.Row(m.WithDataAndText(id, "Перейти к ключу", m.KeyItem)),
		m.ProlongSuccess.Row(m.CloseBtn),
	)

	text := fmt.Sprintf("Подписка на тариф %s успешно продлена", data.TariffName)

	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, text, m.ProlongSuccess)

	return err
}

func (n *Notifier) NotifySubscriptionBandwidthLimitComing(ctx context.Context, data NotifyBandwidthDTO) (err error) {
	id := strconv.FormatUint(data.SubscriptionID, 10)
	m.BandwidthLimit.Inline(
		m.BandwidthLimit.Row(m.WithDataAndText(id, "Перейти к ключу", m.KeyItem)),
		m.BandwidthLimit.Row(m.CloseBtn),
	)

	text := fmt.Sprintf(
		"Подписка на тариф %s скоро достигнет лимита.  Потрачено - %d, всего доступно - %d",
		data.TariffName,
		data.BandwidthSpent,
		data.BandwidthTotal,
	)
	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, text, m.BandwidthLimit)
	return err
}

func (n *Notifier) NotifySubscriptionBandwidthLimitReached(ctx context.Context, data NotifyBandwidthDTO) (err error) {
	id := strconv.FormatUint(data.SubscriptionID, 10)
	m.BandwidthLimit.Inline(
		m.BandwidthLimit.Row(m.WithDataAndText(id, "Перейти к ключу", m.KeyItem)),
		m.BandwidthLimit.Row(m.CloseBtn),
	)

	text := fmt.Sprintf(
		"Подписка на тариф %s достигла лимита.  Потрачено - %d, всего доступно - %d",
		data.TariffName,
		data.BandwidthSpent,
		data.BandwidthTotal,
	)
	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, text, m.BandwidthLimit)
	return err
}

func (n *Notifier) NotifySubscriptionExpired(ctx context.Context, data NotifyDTO) (err error) {
	m.ProlongSuccess.Inline(
		m.ProlongSuccess.Row(m.WithText("Перейти к доступным тарифам", m.KeysGetNewBtn)),
		m.ProlongSuccess.Row(m.CloseBtn),
	)

	text := fmt.Sprintf("Подписка на тариф %s истекла", data.TariffName)

	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, text, m.ProlongSuccess)

	return err
}

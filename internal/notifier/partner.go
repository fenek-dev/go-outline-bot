package notifier

import (
	"context"
	"fmt"
	m "github.com/fenek-dev/go-outline-bot/internal/markup"
	t "gopkg.in/telebot.v3"
)

func (n *Notifier) NotifyPartnerAboutDeposit(ctx context.Context, data NotifyPartnerDTO) (err error) {
	m.PartnerDeposit.Inline(
		m.PartnerDeposit.Row(m.WithText("Открыть баланс", m.PartnerDepositOpenBalanceBtn)),
		m.PartnerDeposit.Row(m.CloseBtn),
	)

	text := fmt.Sprintf(
		"Пользователь %s пополнил баланс по вашей ссылке на сумму %d RUB. Ваше вознаграждение в сумме %d RUB было начислено на ваш счет",
		data.RecipientUsername,
		data.Amount,
		data.RewardAmount,
	)

	_, err = n.bot.Send(&t.User{ID: int64(data.UserID)}, text, m.PartnerDeposit)
	return err
}

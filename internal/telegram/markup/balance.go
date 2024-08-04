package markup

import t "gopkg.in/telebot.v3"

var (
	Balance         = &t.ReplyMarkup{}
	BalanceRecharge = Balance.Data("💳 Пополнить", "BalanceRecharge")
	BalanceHistory  = Balance.Data("📜 История", "BalanceHistory")
)

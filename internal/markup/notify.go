package markup

import t "gopkg.in/telebot.v3"

var (
	ProlongSuccess               = &t.ReplyMarkup{}
	BandwidthLimit               = &t.ReplyMarkup{}
	PartnerDeposit               = &t.ReplyMarkup{}
	PartnerDepositOpenBalanceBtn = PartnerDeposit.Data("💳 Открыть баланс", "PartnerDepositOpenBalance")
)

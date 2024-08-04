package markup

import t "gopkg.in/telebot.v3"

var (
	TariffsMenu    = &t.ReplyMarkup{}
	TariffsBackBtn = TariffsMenu.Data("⬅ Назад", "TariffsClose")
	TariffItem     = TariffsMenu.Data("Тариф", "TariffItem")

	TariffInfo   = &t.ReplyMarkup{}
	TariffBuyBtn = TariffInfo.Data("💳 Купить", "TariffBuy")
)

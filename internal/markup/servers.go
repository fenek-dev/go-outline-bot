package markup

import t "gopkg.in/telebot.v3"

var (
	ServersMenu    = &t.ReplyMarkup{}
	ServersBackBtn = ServersMenu.Data("⬅ Назад", "ServersClose")
	ServerItem     = ServersMenu.Data("REPLACE", "Tariff")
)

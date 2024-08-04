package markup

import (
	t "gopkg.in/telebot.v3"
)

var (
	KeysMenu      = &t.ReplyMarkup{}
	KeysGetNewBtn = KeysMenu.Data("🔑Новый ключ", "KeysGetNew")
	KeysCloseBtn  = KeysMenu.Data("❌Закрыть", "KeysClose")
)

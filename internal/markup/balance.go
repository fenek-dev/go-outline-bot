package markup

import t "gopkg.in/telebot.v3"

var (
	Balance        = &t.ReplyMarkup{}
	BalanceTopUp   = Balance.Data("💳 Пополнить", "BalanceTopUp")
	BalanceHistory = Balance.Data("📜 История", "BalanceHistory")

	TopUp      = &t.ReplyMarkup{}
	TopUpClose = TopUp.Data("❌ Отменить операцию", "TopUpClose")

	Confirm = &t.ReplyMarkup{}

	Phone      = &t.ReplyMarkup{OneTimeKeyboard: true, ResizeKeyboard: true}
	PhoneSend  = Phone.Contact("📱 Отправить номер телефона")
	PhoneClose = Phone.Text("❌ Отмена")
)

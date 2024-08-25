package markup

import (
	t "gopkg.in/telebot.v3"
)

var (
	Menu              = &t.ReplyMarkup{ResizeKeyboard: true}
	ClientListOpenBtn = Menu.Text("ℹ️Инструкции")
	KeysOpenBtn       = Menu.Text("🔑Ключи")
	BalanceOpenBtn    = Menu.Text("💳Баланс")
	InviteOpenBtn     = Menu.Text("📩Пригласить друга")

	ClientList        = &t.ReplyMarkup{}
	ClientListAndroid = ClientList.Data("🤖Android", "Android")
	ClientListIOS     = ClientList.Data("📱IOS", "IOS")
	ClientListWindows = ClientList.Data("🖥️Windows", "Windows")
	ClientListMacOS   = ClientList.Data("💻MacOS", "MacOS")

	AndroidList        = &t.ReplyMarkup{}
	AndroidListBackBtn = AndroidList.Data("⬅ Назад", "AndroidListBack")

	IOSList        = &t.ReplyMarkup{}
	IOSListBackBtn = IOSList.Data("⬅ Назад", "IOSListBack")

	WindowsList        = &t.ReplyMarkup{}
	WindowsListBackBtn = WindowsList.Data("⬅ Назад", "WindowsListBack")

	MacOSList        = &t.ReplyMarkup{}
	MacOSListBackBtn = MacOSList.Data("⬅ Назад", "MacOSListBack")
)

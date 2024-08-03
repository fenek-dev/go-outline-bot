package markup

import (
	t "gopkg.in/telebot.v3"
)

var (
	Menu        = &t.ReplyMarkup{ResizeKeyboard: true}
	InfoOpenBtn = Menu.Text("ℹ Информация")

	Info              = &t.ReplyMarkup{}
	InfoClose         = ClientList.Data("❌Закрыть", "InfoClose")
	ClientListOpenBtn = Info.Data("📱Список клиентов", "ClientList")

	ClientList        = &t.ReplyMarkup{}
	ClientListAndroid = ClientList.Data("🤖Android", "Android")
	ClientListIOS     = ClientList.Data("📱IOS", "IOS")
	ClientListWindows = ClientList.Data("🖥️Windows", "Windows")
	ClientListMacOS   = ClientList.Data("💻MacOS", "MacOS")
	ClientListBack    = ClientList.Data("⬅ Назад", "ClientListBack")

	AndroidList        = &t.ReplyMarkup{}
	AndroidListBackBtn = AndroidList.Data("⬅ Назад", "AndroidListBack")

	IOSList        = &t.ReplyMarkup{}
	IOSListBackBtn = IOSList.Data("⬅ Назад", "IOSListBack")

	WindowsList        = &t.ReplyMarkup{}
	WindowsListBackBtn = WindowsList.Data("⬅ Назад", "WindowsListBack")

	MacOSList        = &t.ReplyMarkup{}
	MacOSListBackBtn = MacOSList.Data("⬅ Назад", "MacOSListBack")
)

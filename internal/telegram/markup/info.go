package markup

import (
	t "gopkg.in/telebot.v3"
)

var (
	Menu              = &t.ReplyMarkup{ResizeKeyboard: true}
	ClientListOpenBtn = Menu.Text("📱Список клиентов")
	KeysOpenBtn       = Menu.Text("🔑Ключи")

	ClientList        = &t.ReplyMarkup{}
	ClientListAndroid = ClientList.Data("🤖Android", "Android")
	ClientListIOS     = ClientList.Data("📱IOS", "IOS")
	ClientListWindows = ClientList.Data("🖥️Windows", "Windows")
	ClientListMacOS   = ClientList.Data("💻MacOS", "MacOS")
	ClientListClose   = ClientList.Data("❌Закрыть", "ClientListClose")

	AndroidList        = &t.ReplyMarkup{}
	AndroidListBackBtn = AndroidList.Data("⬅ Назад", "AndroidListBack")

	IOSList        = &t.ReplyMarkup{}
	IOSListBackBtn = IOSList.Data("⬅ Назад", "IOSListBack")

	WindowsList        = &t.ReplyMarkup{}
	WindowsListBackBtn = WindowsList.Data("⬅ Назад", "WindowsListBack")

	MacOSList        = &t.ReplyMarkup{}
	MacOSListBackBtn = MacOSList.Data("⬅ Назад", "MacOSListBack")
)

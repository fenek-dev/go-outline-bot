package handlers

import (
	"github.com/fenek-dev/go-outline-bot/internal/markup"
	"gopkg.in/telebot.v3"
)

var (
	SelectAPlatform = "Выберите платформу:"
)

func AvailableClients(platform string) string {
	return "Список доступных клиентов под " + platform + ":"
}

func (h *Handlers) CloseInfo(c telebot.Context) error {
	return c.Delete()
}

func (h *Handlers) OpenClientsList(c telebot.Context) error {
	return c.Send(SelectAPlatform, markup.ClientList)
}

func (h *Handlers) CloseClientsList(c telebot.Context) error {
	return c.Delete()
}

func (h *Handlers) ClientsListIOS(c telebot.Context) error {
	return c.Edit(AvailableClients("IOS"), markup.IOSList)
}

func (h *Handlers) BackIosList(c telebot.Context) error {
	return c.Edit(SelectAPlatform, markup.ClientList)
}

func (h *Handlers) ClientsListAndroid(c telebot.Context) error {
	return c.Edit(AvailableClients("Android"), markup.AndroidList)
}

func (h *Handlers) BackAndroidList(c telebot.Context) error {
	return c.Edit(SelectAPlatform, markup.ClientList)
}

func (h *Handlers) ClientsListWindows(c telebot.Context) error {
	return c.Edit(AvailableClients("Windows"), markup.WindowsList)
}

func (h *Handlers) BackWindowsList(c telebot.Context) error {
	return c.Edit(SelectAPlatform, markup.ClientList)
}

func (h *Handlers) ClientsListMacOS(c telebot.Context) error {
	return c.Edit(AvailableClients("MacOS"), markup.MacOSList)
}

func (h *Handlers) BackMacOSList(c telebot.Context) error {
	return c.Edit(SelectAPlatform, markup.ClientList)
}

package telegram

import (
	"fmt"
	"time"

	"github.com/fenek-dev/go-outline-bot/configs"
	"github.com/fenek-dev/go-outline-bot/internal/markup"
	"github.com/fenek-dev/go-outline-bot/internal/telegram/handlers"
	"gopkg.in/telebot.v3"
)

func InitBot(cfg *configs.TelegramConfig, h *handlers.Handlers) (*telebot.Bot, error) {

	var poller telebot.Poller = &telebot.Webhook{}
	if cfg.Debug {
		poller = &telebot.LongPoller{Timeout: time.Second * 1}
	}

	pref := telebot.Settings{
		Token:  cfg.Token,
		Poller: poller,
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("can not connect to telegram bot: %w", err)
	}

	markup.Init()

	b.Handle("/start", h.Start)

	b.Handle(&markup.InfoOpenBtn, h.OpenInfo)
	b.Handle(&markup.InfoClose, h.CloseInfo)

	b.Handle(&markup.ClientListOpenBtn, h.OpenClientsList)

	b.Handle(&markup.ClientListBack, h.BackClientsList)

	b.Handle(&markup.ClientListIOS, h.ClientsListIOS)
	b.Handle(&markup.IOSListBackBtn, h.BackMacOSList)

	b.Handle(&markup.ClientListAndroid, h.ClientsListAndroid)
	b.Handle(&markup.AndroidListBackBtn, h.BackAndroidList)

	b.Handle(&markup.ClientListWindows, h.ClientsListWindows)
	b.Handle(&markup.WindowsListBackBtn, h.BackWindowsList)

	b.Handle(&markup.ClientListMacOS, h.ClientsListMacOS)
	b.Handle(&markup.MacOSListBackBtn, h.BackMacOSList)

	return b, nil
}

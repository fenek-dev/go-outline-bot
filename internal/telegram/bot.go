package telegram

import (
	"fmt"
	m "github.com/fenek-dev/go-outline-bot/internal/telegram/markup"
	"time"

	"github.com/fenek-dev/go-outline-bot/configs"
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

	m.Init()

	b.Handle("/start", h.Start)

	b.Handle(&m.ClientListOpenBtn, h.OpenClientsList)

	b.Handle(&m.CloseBtn, h.Close)

	b.Handle(&m.ClientListIOS, h.ClientsListIOS)
	b.Handle(&m.IOSListBackBtn, h.BackMacOSList)

	b.Handle(&m.ClientListAndroid, h.ClientsListAndroid)
	b.Handle(&m.AndroidListBackBtn, h.BackAndroidList)

	b.Handle(&m.ClientListWindows, h.ClientsListWindows)
	b.Handle(&m.WindowsListBackBtn, h.BackWindowsList)

	b.Handle(&m.ClientListMacOS, h.ClientsListMacOS)
	b.Handle(&m.MacOSListBackBtn, h.BackMacOSList)

	// Keys
	b.Handle(&m.KeysOpenBtn, h.OpenKeysMenu)

	//Servers
	b.Handle(&m.KeysGetNewBtn, h.OpenServersMenu)
	b.Handle(&m.ServersBackBtn, h.BackServersMenu)

	b.Handle(&m.ServerItem, h.OpenTariffsMenu)
	b.Handle(&m.TariffsBackBtn, h.BackTariffsMenu)

	// Tariffs
	b.Handle(&m.TariffItem, h.OpenTariff)
	b.Handle(&m.TariffBuyBtn, h.BuySubscription)

	// Balance
	b.Handle(&m.BalanceOpenBtn, h.OpenBalance)
	return b, nil
}

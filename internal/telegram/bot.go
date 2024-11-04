package telegram

import (
	"fmt"
	"github.com/fenek-dev/go-outline-bot/internal/markup"
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

	markup.Init()

	b.Handle(telebot.OnText, h.TextHandler)

	b.Handle("/start", h.Start)

	b.Handle(&markup.ClientListOpenBtn, h.OpenClientsList)

	b.Handle(&markup.CloseBtn, h.Close)

	b.Handle(&markup.ClientListIOS, h.ClientsListIOS)
	b.Handle(&markup.IOSListBackBtn, h.BackMacOSList)

	b.Handle(&markup.ClientListAndroid, h.ClientsListAndroid)
	b.Handle(&markup.AndroidListBackBtn, h.BackAndroidList)

	b.Handle(&markup.ClientListWindows, h.ClientsListWindows)
	b.Handle(&markup.WindowsListBackBtn, h.BackWindowsList)

	b.Handle(&markup.ClientListMacOS, h.ClientsListMacOS)
	b.Handle(&markup.MacOSListBackBtn, h.BackMacOSList)

	// Keys
	b.Handle(&markup.KeysOpenBtn, h.OpenKeysMenu)
	b.Handle(&markup.KeyItem, h.OpenKeyInfo)
	b.Handle(&markup.KeyAutoProlongBtn, h.ToggleAutoProlong)

	//Servers
	b.Handle(&markup.KeysGetNewBtn, h.OpenServersMenu)
	b.Handle(&markup.ServersBackBtn, h.BackServersMenu)

	b.Handle(&markup.ServerItem, h.OpenTariffsMenu)
	b.Handle(&markup.TariffsBackBtn, h.BackTariffsMenu)

	// Tariffs
	b.Handle(&markup.TariffItem, h.OpenTariff)
	b.Handle(&markup.TariffBuyBtn, h.BuySubscription)

	// Balance
	b.Handle(&markup.BalanceOpenBtn, h.OpenBalance)
	b.Handle(&markup.PartnerDepositOpenBalanceBtn, h.OpenBalance)
	b.Handle(&markup.BalanceTopUp, h.TopUpBalance)
	b.Handle(&markup.TopUpClose, h.TopUpClose)

	// Phone
	b.Handle(telebot.OnContact, h.HandlePhone)
	b.Handle(&markup.PhoneClose, h.ClosePhone)

	return b, nil
}

package notifier

import "gopkg.in/telebot.v3"

type Notifier struct {
	bot *telebot.Bot
}

func New(bot *telebot.Bot) *Notifier {
	return &Notifier{bot: bot}
}

func (n *Notifier) SetBot(bot *telebot.Bot) {
	n.bot = bot
}

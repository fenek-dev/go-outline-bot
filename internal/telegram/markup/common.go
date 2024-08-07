package markup

import t "gopkg.in/telebot.v3"

var (
	CloseBtn = t.Btn{
		Unique: "Close",
		Text:   "❌ Закрыть",
	}

	OnlyClose = &t.ReplyMarkup{}
)

func WithData(data string, btn t.Btn) t.Btn {
	btn.Data = data
	return btn
}

func WithText(text string, btn t.Btn) t.Btn {
	btn.Text = text
	return btn
}

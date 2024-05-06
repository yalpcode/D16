package bot

import tele "gopkg.in/telebot.v3"

var (
	deleter_markup = &tele.ReplyMarkup{}
	btnDelete      = deleter_markup.Data("❌ Удалить", "dlt", "deletemsg")
)

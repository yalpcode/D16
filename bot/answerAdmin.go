package bot

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var queueMsg sync.Map

func answerAdmin(c tele.Context) error {
	bot := c.Bot()

	channel, err := bot.ChatByID(channel_id)
	if err != nil {
		log.Fatal(err)
	}

	admins, err := bot.AdminsOf(channel)
	if err != nil {
		log.Fatal(err)
	}

	selector := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0)

	for _, admin := range admins {
		if !admin.User.IsBot {
			admin_id := strconv.FormatInt(admin.User.ID, 10)
			rows = append(rows, selector.Row(selector.Data(admin.User.FirstName, "ad", admin_id)))
		}
	}

	selector.Inline(
		rows...,
	)

	return c.Send("Выберите, кому хотите задать вопрос:", selector)
}

func selectAdmin(c tele.Context, state fsm.Context) error {
	go state.Update("ID_ADMIN", strings.Split(c.Data(), "|")[1])
	go state.Set(ID_ADMINSG)

	return c.Send("Введите сообщение для админа (Анонимно):\n<b>ОТВЕТ БУДЕТ В КАНАЛЕ</b>")
}

func inputAnswerAdmin(c tele.Context, state fsm.Context) error {
	v, ok := queueMsg.Load(c.Chat().ID)
	if ok && v != nil {
		go state.Finish(true)
		return c.Send("Вы уже отправляли сообщение! Подождите ответа")
	}

	var id_admin_s string
	state.MustGet("ID_ADMIN", &id_admin_s)
	go state.Finish(true)

	id_admin, err := strconv.ParseInt(id_admin_s, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	bot := c.Bot()

	chat_admin, err := bot.ChatByID(id_admin)
	if err != nil {
		log.Fatal(err)
	}

	bot.Send(
		chat_admin,
		"Анонимный вопрос: "+c.Text(),
		deleter_markup,
	)

	go queueMsg.Store(c.Chat().ID, true)

	return c.Send("Сообщение отправлено")
}

func deleteAnsMsg(c tele.Context) error {
	go queueMsg.Delete(c.Chat().ID)
	return c.Delete()
}

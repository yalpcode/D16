package bot

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var QueueAns *TTLMap = NewTTLMap(30 * time.Second)

type TTLMap struct {
	data       sync.Map
	expiration time.Duration
}

func NewTTLMap(expiration time.Duration) *TTLMap {
	return &TTLMap{
		data:       sync.Map{},
		expiration: expiration,
	}
}

func (m *TTLMap) Set(key, value interface{}) {
	m.data.Store(key, value)
	time.AfterFunc(m.expiration, func() { m.data.Delete(key) })
}

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
	v, ok := QueueAns.data.Load(c.Chat().ID)
	if ok && v != nil {
		go state.Finish(true)
		return c.Send("Вы уже отправляли сообщение! Подождите ответа или 30 секунд")
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

	if c.Message().Photo != nil {
		photo := c.Message().Photo
		photo.Caption = "Анонимный вопрос: " + c.Text()
		photo.Send(bot, chat_admin, &tele.SendOptions{ReplyMarkup: deleter_markup})
	} else {
		bot.Send(
			chat_admin,
			"Анонимный вопрос: "+c.Text(),
			deleter_markup,
		)
	}

	QueueAns.Set(c.Chat().ID, true)

	return c.Send("Сообщение отправлено")
}

func deleteAnsMsg(c tele.Context) error {
	QueueAns.data.Delete(c.Chat().ID)
	return c.Delete()
}

func getAnsAdmin(c tele.Context, _ fsm.Context) error {
	QueueAns.data.Delete(c.Chat().ID)

	bot := c.Bot()

	channel, err := bot.ChatByID(channel_id)
	if err != nil {
		log.Fatal(err)
	}

	text := "\n\n```Ответ: " + c.Text() + "```\n\n" + c.Sender().FirstName

	if c.Message().ReplyTo.Photo != nil && c.Message().Photo != nil {
		photo := c.Message().Photo
		photo.Caption = c.Message().ReplyTo.Caption + text
		album := tele.Album{
			c.Message().ReplyTo.Photo,
			photo,
		}

		bot.SendAlbum(
			channel,
			album,
		)
	} else if c.Message().ReplyTo.Photo != nil {
		photo := c.Message().ReplyTo.Photo
		photo.Caption = c.Message().ReplyTo.Caption + text
		photo.Send(bot, channel, &tele.SendOptions{})
	} else if c.Message().Photo != nil {
		photo := c.Message().Photo
		photo.Caption = c.Message().ReplyTo.Text + text
		photo.Send(bot, channel, &tele.SendOptions{})
	} else {
		_, err = bot.Send(channel, c.Message().ReplyTo.Text+text)
	}

	bot.Delete(c.Message().ReplyTo)

	return err
}

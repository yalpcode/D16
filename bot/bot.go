package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var channel_id int64

func start(c tele.Context) error {
	return c.Send("Hello!")
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
		admin_id := strconv.FormatInt(admin.User.ID, 10)
		rows = append(rows, selector.Row(selector.Data(admin.User.FirstName, "ad", admin_id)))
	}

	selector.Inline(
		rows...,
	)

	return c.Send("Выберите кому хотите задать вопрос:", selector)
}

// func sendAnswerAdmin(c tele.Context) error {
// 	c.Data()
// }

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	channel_id, err = strconv.ParseInt(os.Getenv("CHANNEL"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", start)
	b.Handle("/answer_admin", answerAdmin)
	// b.Handle(tele.OnQuery)

	b.Start()
}

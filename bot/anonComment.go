package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

func sendAnonComment(c tele.Context, state fsm.Context) error {
	parts := strings.Split(c.Text(), "/")
	lastMessageNumber := parts[len(parts)-1]
	message_id, _ := strconv.Atoi(lastMessageNumber)

	go state.Update("messageID", message_id)
	go state.Set(IS_CHANNELSG)
	return c.Send("Введите текст комментария: ")
}

func anonComment(c tele.Context, state fsm.Context) error {
	bot := c.Bot()

	var message_id int
	state.MustGet("messageID", &message_id)
	go state.Finish(true)

	chat, err := bot.ChatByID(chat_id)
	if err != nil {
		log.Fatal(err)
	}

	messageN := &tele.Message{ID: message_id, Chat: chat, Text: c.Text()}

	_, err = bot.Send(chat, c.Text(), &tele.SendOptions{ReplyTo: messageN})
	if err != nil {
		log.Fatal(err)
	}

	return err
}

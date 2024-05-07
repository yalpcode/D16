package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func sendAnonComment(c tele.Context) error {
	bot := c.Bot()
	_, err := bot.Send(c.Message().ReplyTo.Chat, c.Text())
	fmt.Println(c.Text())
	return err
}

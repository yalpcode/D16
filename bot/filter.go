package bot

import (
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

func CheckSubscription() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			bot := c.Bot()
			channel, err := bot.ChatByID(channel_id)
			if err != nil {
				return err
			}
			member, err := bot.ChatMemberOf(channel, c.Sender())

			if err != nil {
				return err
			}
			if member.Role == "kicked" || member.Role == "left" || c.Sender().IsBot {
				return nil
			}
			return next(c)
		}
	}
}

func FilterCallBack(c tele.Context, state fsm.Context) error {
	data := c.Data()
	if strings.Contains(data, "ad") {
		return selectAdmin(c, state)
	}

	return nil
}

func FilterText(c tele.Context, state fsm.Context) error {
	if c.Message().IsReply() && c.Message().ReplyTo.Sender.IsBot && strings.HasPrefix(c.Message().ReplyTo.Text, "Анонимный вопрос:") {
		return getAnsAdmin(c, state)
	}
	if strings.HasPrefix(c.Text(), "https://t.me/c/"+chat_link) {
		return sendAnonComment(c, state)
	}

	return nil
}

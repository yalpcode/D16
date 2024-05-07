package bot

import (
	"fmt"
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
			if member.Role == "kicked" || member.Role == "left" {
				return nil
			}
			return next(c)
		}
	}
}

func CheckPost() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			fmt.Println(c.Message().IsReply())
			if c.Message().IsReply() && c.Message().ReplyTo.Chat.ID == channel_id {
				return next(c)
			}
			return nil
		}
	}
}

func OnCallbackF(c tele.Context, state fsm.Context) error {
	data := c.Data()

	if strings.Contains(data, "ad") {
		return selectAdmin(c, state)
	}

	return nil
}

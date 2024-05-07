package bot

import (
	"strings"

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
			if !c.Sender().IsBot && strings.HasPrefix(c.Text(), "https://t.me/c/"+chat_link) {
				return next(c)
			}
			return nil
		}
	}
}

func CheckAnsAdmin() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().IsReply() && c.Message().ReplyTo.Sender.IsBot && strings.HasPrefix(c.Message().ReplyTo.Text, "Анонимный вопрос:") {
				return next(c)
			}
			return nil
		}
	}
}

func CheckCallBack(pref string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if strings.Contains(c.Data(), pref) {
				return next(c)
			}
			return nil
		}
	}
}

package bot

import (
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

func OnCallbackF(c tele.Context, state fsm.Context) error {
	data := c.Data()

	if strings.Contains(data, "ad") {
		return selectAdmin(c, state)
	}

	return nil
}

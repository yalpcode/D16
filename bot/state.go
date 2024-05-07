package bot

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
)

var (
	ansAdminSG = fsm.NewStateGroup("ansAd")
	ID_ADMINSG = ansAdminSG.New("ID_ADMIN")
)

var (
	anonCommentSG = fsm.NewStateGroup("anonCom")
	IS_CHANNELSG  = anonCommentSG.New("messageID")
)

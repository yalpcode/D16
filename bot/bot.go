package bot

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	tele "gopkg.in/telebot.v3"
)

var channel_id, chat_id int64
var chat_link string

func start(c tele.Context) error {
	text := "Здравствуйте, <b>" + c.Chat().FirstName + "</b>! Я бот-помощник по каналу!\nВы можете увидеть мой функционал используя '/' или во вкладке 'меню'."
	return c.Send(text)
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	channel_id, err = strconv.ParseInt(os.Getenv("CHANNEL"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	chat_link = os.Getenv("CHAT_LINK")
	chat_id, err = strconv.ParseInt(os.Getenv("CHAT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	pref := tele.Settings{
		Token:     os.Getenv("TOKEN"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	storage := memory.NewStorage()
	defer storage.Close()

	manager := fsm.NewManager(b, nil, storage, nil)

	deleter_markup.Inline(
		deleter_markup.Row(btnDelete),
	)

	b.Use(CheckSubscription())

	b.Handle("/start", start)

	b.Handle("/answer_admin", answerAdmin)
	manager.Bind(tele.OnCallback, fsm.AnyState, selectAdmin, CheckCallBack("ad"))
	manager.Bind(tele.OnText, ID_ADMINSG, inputAnswerAdmin)
	b.Handle(&btnDelete, deleteAnsMsg)
	manager.Bind(tele.OnText, fsm.DefaultState, getAnsAdmin, CheckAnsAdmin())

	manager.Bind(tele.OnText, fsm.AnyState, sendAnonComment, CheckPost())
	manager.Bind(tele.OnText, IS_CHANNELSG, anonComment)

	b.Start()
}

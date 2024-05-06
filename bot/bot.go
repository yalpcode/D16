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

var channel_id int64

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

	b.Handle("/start", start)
	b.Handle("/answer_admin", answerAdmin)
	manager.Bind(tele.OnCallback, fsm.DefaultState, selectAdmin)
	manager.Bind(tele.OnText, ID_ADMINSG, inputAnswerAdmin)

	b.Start()
}

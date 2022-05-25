package mainpac

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
	"math/rand"
	"streamtg/go-log"
	"streamtg/minidb"
	"time"
)

type Service struct {
	Bot    *Bot
	MiniDB *minidb.Pudge
	Loc    *time.Location
	Rand   *rand.Rand

	YouTube *YouTube
}

type YouTube struct {
	ChannelID     string
	CycleDuration int // minutes

	Text       Text
	Locs       []string
	TimeFormat string
	TimeCity   bool
}

func TimeFormatCity(withCity bool) string {
	if withCity {
		return "2 Jan 15:04 MST"
	}
	return "2 Jan 15:04"
}

type Text struct {
	Live, Upcoming, Start, End string
}

type Bot struct {
	*tb.Bot
	*layout.Layout

	UserList   []int64
	AdminList  []int64
	NotifyList []int64
	ErrorList  []int64

	Username string
	Uptime   time.Time

	CallbackQuery map[int64]string //контекстный ввод
}

func (s Service) Start() {
	log.Info("tgbot init")
	s.InitBot()
	log.Info("tgbot launch...")
	fmt.Println("tgbot @" + s.Bot.Me.Username)
	go s.GoCheckErrs()
	s.Bot.Start()
}

func (s Service) GoCheckErrs() {
	time.Sleep(time.Second * 30)
	nErr := log.GetErrN()
	if nErr > 0 {
		s.Bot.sendToSlice(s.Bot.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
	}

	for range time.Tick(time.Minute * 5) {
		nErr = log.GetErrN()
		if nErr > 0 {
			s.Bot.sendToSlice(s.Bot.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
		}
	}
}

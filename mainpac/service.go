package mainpac

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"streamtg/go-log"
	"streamtg/minidb"
	"streamtg/twitch"

	"github.com/mmcdole/gofeed"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Service struct {
	Bot    *Bot
	MiniDB *minidb.MiniDB
	Loc    *time.Location
	Rand   *rand.Rand

	YouTubeTwitch *YouTubeTwitch
}

type YouTubeTwitch struct {
	Parser   *gofeed.Parser
	LogLevel string

	Twitch *twitch.Twitch

	LastRSS_YT       *gofeed.Feed
	LastTime         time.Time
	NumberIterations int

	PauseMutex       sync.Mutex
	Pause            int // 0 false 1 true 2 true & wait
	PauseWaitChannel chan struct{}

	ChannelID            string
	TwitchNick           string
	CycleDurationMinutes int // minutes

	Locs       []string
	TimeFormat string
	TimeCity   bool
}

func TimeFormatCity(withCity bool) string {
	if withCity {
		return "🕒 2 Jan 15:04 MST"
	}
	return "🕒 2 Jan 15:04"
}

type Bot struct {
	*tb.Bot
	*layout.Layout

	UserList   []int64
	AdminList  []int64
	NotifyList []minidb.Channel
	ErrorList  []int64

	Username string
	Uptime   time.Time

	CallbackQuery map[int64]string // контекстный ввод
}

func (s Service) Start() {
	log.Info("tgbot init")
	s.InitBot()
	log.Info("tgbot launch...")
	fmt.Println("tgbot @" + s.Bot.Me.Username)
	go s.GoCheckErrs()
	go s.GoYouTubeTwitch()
	s.Bot.Start()
}

func (s Service) GoCheckErrs() {
	checkErrs := func() {
		nErr := log.GetErrN()
		if nErr > 0 {
			s.Bot.sendToSlice(s.Bot.ErrorList, fmt.Sprintf("Новых ошибок: %v.\nЗагляните в логи - /logs.", nErr))
		}
	}

	time.Sleep(time.Second * 30)
	checkErrs()

	for range time.Tick(time.Minute * 5) {
		checkErrs()
	}
}

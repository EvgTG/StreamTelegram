package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"streamtg/go-log"
	"streamtg/mainpac"
	"streamtg/minidb"
	"streamtg/twitch"
	"streamtg/util"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mmcdole/gofeed"
	"github.com/rotisserie/eris"
	"go.uber.org/fx"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

func New() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			NewDB,
			NewService,
		),

		fx.Invoke(
			ReadConfig,
			Logger,
			PingServe,
			Start,
		),
	)
	return
}

func ReadConfig() {
	err := cleanenv.ReadConfig("files/cfg.env", &CFG)
	if err != nil {
		panic(eris.Wrap(err, "ReadConfig"))
	}
}

func Logger() {
	log.SetLogger(log.New(CFG.LogLevel, true))
	log.Info("Go!")
}

func NewDB() *minidb.MiniDB {
	db, err := minidb.NewDB()
	util.ErrCheckFatal(err, "minidb.NewDB()", "NewDB", "init")
	return db
}

func NewService(lc fx.Lifecycle, db *minidb.MiniDB) *mainpac.Service {
	// Telegram
	lt, err := layout.New("bot.yml")
	util.ErrCheckFatal(err, "layout.New()", "NewService", "init")

	bot, err := tb.NewBot(tb.Settings{
		Token: CFG.TgApiToken,
		Poller: &tb.Webhook{
			Listen:      ":" + CFG.LocalPort,
			IP:          CFG.IP,
			SecretToken: CFG.SecretToken,

			Endpoint: &tb.WebhookEndpoint{
				PublicURL: "https://" + CFG.IP + ":" + CFG.Port + CFG.Path,
				Cert:      "cert/pem.pem",
			},

			// +.tls если без обратного прокси
		},
		ParseMode: tb.ModeHTML,
	})
	// long-polling &tb.LongPoller{Timeout: 30 * time.Second},
	util.ErrCheckFatal(err, "tb.NewBot()", "NewService", "init")
	bot.Use(lt.Middleware("ru"))

	// YouTube Twitch
	twitch, err := twitch.NewTwitch(db)
	util.ErrCheckFatal(err, "twitch.NewTwitch()", "NewService", "init")

	channelID, err := db.GetChannelID()
	util.ErrCheckFatal(err, "db.GetChannelID()", "NewService", "init")

	twitchNick, err := db.GetTwitchNick()
	util.ErrCheckFatal(err, "db.GetTwitchNick()", "NewService", "init")

	cycleDuration, err := db.GetCycleDuration()
	util.ErrCheckFatal(err, "db.GetCycleDuration()", "NewService", "init")
	if cycleDuration == 0 {
		cycleDuration = 5
	}

	locs, err := db.GetLocs()
	util.ErrCheckFatal(err, "db.GetLocs()", "NewService", "init")

	timeWithCity, err := db.GetTimeWithCity()
	util.ErrCheckFatal(err, "db.GetTimeWithCity()", "NewService", "init")

	notifyList, err := db.GetNotifyList()
	util.ErrCheckFatal(err, "db.GetNotifyList()", "NewService", "init")

	service := &mainpac.Service{
		Bot: &mainpac.Bot{
			Bot:           bot,
			Layout:        lt,
			Username:      bot.Me.Username,
			UserList:      CFG.UserList,
			AdminList:     CFG.AdminList,
			NotifyList:    notifyList,
			ErrorList:     CFG.ErrorList,
			Uptime:        time.Now(),
			CallbackQuery: make(map[int64]string, 0),
		},
		MiniDB: db,
		Loc:    CFG.TimeLocation.Get(),
		Rand:   rand.New(rand.NewSource(time.Now().UnixNano())),

		YouTubeTwitch: &mainpac.YouTubeTwitch{
			Parser:               gofeed.NewParser(),
			LogLevel:             CFG.LogLevel,
			Twitch:               twitch,
			LastTime:             time.Unix(0, 0),
			NumberIterations:     0,
			PauseMutex:           sync.Mutex{},
			Pause:                0,
			PauseWaitChannel:     make(chan struct{}),
			ChannelID:            channelID,
			TwitchNick:           twitchNick,
			CycleDurationMinutes: cycleDuration,
			Locs:                 locs,
			TimeFormat:           mainpac.TimeFormatCity(timeWithCity),
			TimeCity:             timeWithCity,
		},
	}

	err = service.Bot.RemoveWebhook()
	util.ErrCheckFatal(err, "Bot.RemoveWebhook()", "NewService", "init")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// service.Bot.Stop()
			service.Bot.RemoveWebhook()
			return nil
		},
	})

	return service
}

func Start(s *mainpac.Service) {
	s.Start()
}

func PingServe() {
	if !CFG.PingOn {
		log.Info("PingServer off")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pingstreamtg", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})
	log.Info("PingServer on")
	go http.ListenAndServe(":"+CFG.PingPort, mux)
}

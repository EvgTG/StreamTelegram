package mainpac

import (
	"StreamTelegram/model"
	"context"
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	tg      *tg
	yt      *yt
	db      *model.Model
	loc     *time.Location
	envVars envVars
}

type tg struct {
	tgBot            *tgbotapi.BotAPI
	updateConfig     tgbotapi.UpdateConfig
	toID             []int64
	errorToID        int64
	userList         []int64
	numberIterations int
	uptime           time.Time
	callbackQuery    string
}

type yt struct {
	yts       *youtube.Service
	channelID string
	texts     map[string]string
	CycleTime time.Duration
	stop      int8
	stopch    chan bool
	lastRSS   gofeed.Feed
	lastTime  time.Time
}

type InitConfig struct {
	Proxy, TgApiToken   string
	TOID                []int64
	ErrorToID           int64
	UserList            []int64
	ChannelID, YTApiKey string
	Loc                 *time.Location
	LanguageOFText      string
	CycleTime           time.Duration
}

type envVars struct {
	toID      []int64
	cycleTime time.Duration
}

func New(cfg InitConfig, db *model.Model) (*Service, error) {
	var err error

	//Telegram
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	if cfg.Proxy != "" {
		proxyURL, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, fmt.Errorf("mainpac.New - proxy url.Parse(): %s", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}
	}

	tgBot, err := tgbotapi.NewBotAPIWithClient(cfg.TgApiToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		return nil, fmt.Errorf("mainpac.New - tgbotapi.NewBotAPIWithClient(): %s", err)
	}

	tgBot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//YouTube
	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YTApiKey))
	if err != nil {
		return nil, fmt.Errorf("mainpac.New - youtube.NewService(): %s", err)
	}

	service := &Service{
		tg: &tg{
			tgBot:            tgBot,
			updateConfig:     u,
			toID:             cfg.TOID,
			errorToID:        cfg.ErrorToID,
			userList:         cfg.UserList,
			numberIterations: 0,
			uptime:           time.Now(),
			callbackQuery:    "",
		},
		yt: &yt{
			yts:       yts,
			channelID: cfg.ChannelID,
			texts:     GetTexts(cfg.LanguageOFText),
			CycleTime: cfg.CycleTime,
			stop:      0,
			stopch:    make(chan bool),
			lastRSS:   gofeed.Feed{},
			lastTime:  time.Time{},
		},
		db:  db,
		loc: cfg.Loc,
		envVars: envVars{
			toID:      cfg.TOID,
			cycleTime: cfg.CycleTime,
		},
	}

	st := db.GetLs()
	if st.DBPriority.ToIDBL {
		service.tg.toID = st.DBPriority.ToID
	}
	if st.DBPriority.CycleTimeBL {
		service.yt.CycleTime = st.DBPriority.CycleTime
	}

	return service, nil
}

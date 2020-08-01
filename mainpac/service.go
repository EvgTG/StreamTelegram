package mainpac

import (
	"StreamTelegram/model"
	"context"
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	tg *tg
	yt *yt
	db *model.Model
}

type tg struct {
	tgBot            *tgbotapi.BotAPI
	updateConfig     tgbotapi.UpdateConfig
	toID             int64
	errorToID        int64
	userList         []int64
	numberIterations int
	uptime           time.Time
}

type yt struct {
	yts       *youtube.Service
	channelID string
}

type InitConfig struct {
	Proxy, TgApiToken   string //tgBot
	TOID, ErrorToID     int64
	UserList            []int64
	ChannelID, YTApiKey string //yt
}

func New(cfg InitConfig, db *model.Model) (*Service, error) {
	service := Service{
		tg: &tg{},
		yt: &yt{},
		db: nil,
	}
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

	service.tg.tgBot, err = tgbotapi.NewBotAPIWithClient(cfg.TgApiToken, client)
	if err != nil {
		return nil, fmt.Errorf("mainpac.New - tgbotapi.NewBotAPIWithClient(): %s", err)
	}

	service.tg.tgBot.Debug = false

	service.tg.updateConfig = tgbotapi.NewUpdate(0)
	service.tg.updateConfig.Timeout = 60

	service.tg.toID = cfg.TOID
	service.tg.errorToID = cfg.ErrorToID
	service.tg.userList = cfg.UserList
	service.tg.numberIterations = 0
	service.tg.uptime = time.Now()

	//YouTube
	ctx := context.Background()
	service.yt.yts, err = youtube.NewService(ctx, option.WithAPIKey(cfg.YTApiKey))
	if err != nil {
		return nil, fmt.Errorf("mainpac.New - youtube.NewService(): %s", err)
	}

	service.yt.channelID = cfg.ChannelID

	//DB
	service.db = db

	return &service, nil
}

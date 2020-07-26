package main

import (
	"StreamTelegram/go-config"
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"StreamTelegram/mongodb"
	"StreamTelegram/tgbot"
	u "StreamTelegram/utility"
	"context"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"strings"
)

func New() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			Config,
			NewYT,
			NewDB,
			NewTGBot,
		),

		fx.Invoke(
			Logger,
			Start,
		),
	)
	return
}

func Config() *config.Config {
	conf := config.New()
	return conf
}

func Logger(conf *config.Config) {
	logLevel := conf.GetString("LOGLVL")
	log.SetLogger(log.New(logLevel))
}

func NewDB(conf *config.Config) *model.Model {
	return model.New(mongodb.NewDB(conf.GetString("NAMEDB")))
}

func NewYT(conf *config.Config) *youtube.Service {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(conf.GetString("YTAPIKEY")))
	u.Fatal("main.NewYT - youtube.NewService", err)
	return youtubeService
}

func NewTGBot(conf *config.Config, youtubeService *youtube.Service) *tgbot.TGBot {
	toIDs := conf.GetString("TOID")
	if toIDs == "" {
		log.Fatal("ERR main.NewTGBot - empty TOID")
	}
	toID, err := strconv.ParseInt(toIDs, 10, 64)
	u.Fatal("main.NewTGBot - TOID strconv.ParseInt", err)

	var errorToID int64
	errorToIDs := conf.GetString("ERRORTOID")
	if toIDs != "" {
		errorToID, err = strconv.ParseInt(errorToIDs, 10, 64)
		u.Fatal("main.NewTGBot - ERRORTOID strconv.ParseInt", err)
	}

	var uList []int64
	uListL := strings.Split(conf.GetString("USERLIST"), ",")
	for _, v := range uListL {
		id, err := strconv.ParseInt(v, 10, 64)
		u.Fatal("main.NewTGBot - strconv.ParseInt", err)
		uList = append(uList, id)
	}

	tgBot, err := tgbot.New(conf.GetString("PROXY"), conf.GetString("TOKEN"), toID, errorToID, youtubeService, uList)
	u.Fatal("main.NewTGBot - tgbot.New", err)

	return tgBot
}

func Start(db *model.Model, tg *tgbot.TGBot, conf *config.Config, youtubeService *youtube.Service) {
	go tg.Start(youtubeService)
	start(db, tg, conf.GetString("CHANNELID"), youtubeService)
}

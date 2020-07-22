package main

import (
	"StreamTelegram/go-config"
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"StreamTelegram/mongodb"
	"StreamTelegram/tgbot"
	u "StreamTelegram/utility"
	"go.uber.org/fx"
	"strconv"
)

func New() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			Config,
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

func NewTGBot(conf *config.Config) *tgbot.TGBot {
	toIDs := conf.GetString("TOID")
	if toIDs == "" {
		log.Fatal("empty TOID")
	}
	toID, err := strconv.ParseInt(toIDs, 10, 64)
	u.Fatal("main.NewTGBot - TOID strconv.ParseInt", err)

	var errorToID int64
	errorToIDs := conf.GetString("ERRORTOID")
	if toIDs != "" {
		errorToID, err = strconv.ParseInt(errorToIDs, 10, 64)
		u.Fatal("main.NewTGBot - ERRORTOID strconv.ParseInt", err)
	}

	tgBot, err := tgbot.New(conf.GetString("PROXY"), conf.GetString("TOKEN"), toID, errorToID)
	u.Fatal("main.NewTGBot - tgbot.New", err)

	return tgBot
}

func Start(db *model.Model, tg *tgbot.TGBot, conf *config.Config) {
	go tg.Start()
	start(db, tg, conf.GetString("CHANNELID"), conf.GetString("YTAPIKEY"))
}

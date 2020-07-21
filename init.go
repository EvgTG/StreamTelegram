package main

import (
	"StreamTelegram/go-config"
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"StreamTelegram/mongodb"
	"StreamTelegram/tgbot"
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
	if err != nil {
		log.Fatal("strconv.ParseInt: ", err)
	}

	tgBot, err := tgbot.New(conf.GetString("PROXY"), conf.GetString("TOKEN"), toID)
	if err != nil {
		log.Fatal("tgbot.New: ", err)
	}

	return tgBot
}

func Start(db *model.Model, tg *tgbot.TGBot, conf *config.Config) {
	go tg.Start()
	start(db, tg, conf.GetString("CHANNELID"), conf.GetString("YTAPIKEY"))
}

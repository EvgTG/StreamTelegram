package main

import (
	"StreamTelegram/go-config"
	"StreamTelegram/go-log"
	"StreamTelegram/mainpac"
	"StreamTelegram/model"
	"StreamTelegram/mongodb"
	"go.uber.org/fx"
	"time"
)

func New() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			Config,
			NewDB,
			NewService,
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
	return model.New(mongodb.NewDB(conf.GetString("NAMEDB"), conf.GetString("MONGOURL")))
}

func NewService(conf *config.Config, db *model.Model) *mainpac.Service {
	loc, err := conf.GetTimeLocation("LOC")
	if err != nil {
		mainpac.Fatal("main.NewService - time zone time.LoadLocation()", err)
	}

	var cycleTime time.Duration
	cycleTimeInt := conf.GetInt64("CYCLETIME")
	if cycleTimeInt == 0 {
		cycleTime = time.Duration(3) * time.Minute
	} else {
		cycleTime = time.Duration(cycleTimeInt) * time.Minute
	}

	cfg := mainpac.InitConfig{
		Proxy:              conf.GetString("PROXY"),
		TgApiToken:         conf.GetString("TOKEN"),
		TOID:               conf.GetIntSlice64("TOID"),
		ErrorToID:          conf.GetInt64("ERRORTOID"),
		UserList:           conf.GetIntSlice64("USERLIST"),
		ChannelID:          conf.GetString("CHANNELID"),
		YTApiKey:           conf.GetString("YTAPIKEY"),
		Loc:                loc,
		LanguageOFText:     conf.GetString("LANGUAGETEXT"),
		CycleTime:          cycleTime,
		TimeFormatWithCity: conf.GetBool("TIMECITY"),
	}

	service, err := mainpac.New(cfg, db)
	if err != nil {
		mainpac.Fatal("main.NewService - mainpac.New()", err)
	}
	return service
}

func Start(service *mainpac.Service) {
	go service.StartTG()
	service.StartYT()
}

package main

import (
	"StreamTelegram/go-config"
	"StreamTelegram/go-log"
	"StreamTelegram/mainpac"
	"StreamTelegram/model"
	"StreamTelegram/mongodb"
	"go.uber.org/fx"
	"strconv"
	"strings"
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
	return model.New(mongodb.NewDB(conf.GetString("NAMEDB")))
}

func NewService(conf *config.Config, db *model.Model) *mainpac.Service {
	//TODO сделать лист, добавить уведомление и управление через тг
	toIDs := conf.GetString("TOID")
	if toIDs == "" {
		log.Fatal("ERR main.NewTGBot - empty TOID")
	}
	toID, err := strconv.ParseInt(toIDs, 10, 64)
	mainpac.Fatal("main.NewService - TOID strconv.ParseInt()", err)

	var errorToID int64
	errorToIDs := conf.GetString("ERRORTOID")
	if errorToIDs != "" {
		errorToID, err = strconv.ParseInt(errorToIDs, 10, 64)
		mainpac.Fatal("main.NewService - ERRORTOID strconv.ParseInt()", err)
	}

	var userList []int64
	userListL := strings.Split(conf.GetString("USERLIST"), ",")
	for _, v := range userListL {
		id, err := strconv.ParseInt(v, 10, 64)
		mainpac.Fatal("main.NewService - strconv.ParseInt()", err)
		userList = append(userList, id)
	}

	cfg := mainpac.InitConfig{
		Proxy:      conf.GetString("PROXY"),
		TgApiToken: conf.GetString("TOKEN"),
		TOID:       toID,
		ErrorToID:  errorToID,
		UserList:   userList,
		ChannelID:  conf.GetString("CHANNELID"),
		YTApiKey:   conf.GetString("YTAPIKEY"),
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

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
	var err error
	var toID []int64
	toIDs := strings.Split(conf.GetString("TOID"), ",")
	for _, v := range toIDs {
		id, err := strconv.ParseInt(v, 10, 64)
		mainpac.Fatal("main.NewService - USERLIST strconv.ParseInt()", err)
		toID = append(toID, id)
	}
	if len(toID) == 0 {
		log.Fatal("ERR main.NewTGBot - empty TOID")
	}

	var errorToID int64
	errorToIDs := conf.GetString("ERRORTOID")
	if errorToIDs != "" {
		errorToID, err = strconv.ParseInt(errorToIDs, 10, 64)
		mainpac.Fatal("main.NewService - ERRORTOID strconv.ParseInt()", err)
	} else {
		log.Info("ERR main.NewTGBot - empty ERRORTOID")
	}

	var userList []int64
	userListL := strings.Split(conf.GetString("USERLIST"), ",")
	for _, v := range userListL {
		id, err := strconv.ParseInt(v, 10, 64)
		mainpac.Fatal("main.NewService - USERLIST strconv.ParseInt()", err)
		userList = append(userList, id)
	}
	if len(userList) == 0 {
		log.Fatal("main.NewService - empty USERLIST")
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

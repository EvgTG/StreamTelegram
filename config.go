package main

import (
	"time"

	"github.com/rotisserie/eris"
)

var CFG InitConfig

type InitConfig struct {
	LogLevel string `env:"LOG_LVL" env-default:"INFO"`

	ProxyTG      string         `env:"TG_PROXY"`
	TgApiToken   string         `env:"TG_TOKEN"`
	UserList     []int64        `env:"LIST_USER"`
	AdminList    []int64        `env:"LIST_ADMIN"`
	NotifyList   []int64        `env:"LIST_NOTIFY"`
	ErrorList    []int64        `env:"LIST_ERROR"`
	TimeLocation MyTimeLocation `env:"LOC" env-default:"UTC"`

	// TG Webhook
	IP          string `env:"WH_IP"`
	Path        string `env:"WH_PATH"`
	Port        string `env:"WH_PORT" env-default:"8443"`
	LocalPort   string `env:"WH_LOCAL_PORT"`
	SecretToken string `env:"WH_SEC"`

	// Ping
	PingPort string `env:"PING_PORT" env-default:"6970"`
	PingOn   bool   `env:"PING_ON" env-default:"false"`
}

type MyTimeLocation string

func (l *MyTimeLocation) SetValue(s string) error {
	*l = MyTimeLocation(s)
	return nil
}

func (l MyTimeLocation) Get() *time.Location {
	loc, err := time.LoadLocation(string(l))
	if err != nil {
		panic(eris.Wrap(err, "cfg.GetTimeLocation()"))
	}
	return loc
}

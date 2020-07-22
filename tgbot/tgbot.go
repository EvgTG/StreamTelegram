package tgbot

import (
	"StreamTelegram/go-log"
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"net/url"
	"time"
)

type TGBot struct {
	tgBot        *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	toID         int64
	errorToID    int64
}

func New(proxy, token string, toID, errorToID int64) (*TGBot, error) {
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("url.Parse(): %s", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}
	}

	tgBot, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPIWithClient(): %s", err)
	}

	tgBot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &TGBot{tgBot, u, toID, errorToID}, nil
}

func (bt *TGBot) SendNotification(text string) {
	msg := tgbotapi.NewMessage(bt.toID, text)
	bt.tgBot.Send(msg)
}

func (bt *TGBot) SendLog(text string) {
	if bt.errorToID == 0 {
		return
	}
	msg := tgbotapi.NewMessage(bt.errorToID, text)
	bt.tgBot.Send(msg)
}

func (tb *TGBot) Start() {
	uptime := time.Now()
	updates, err := tb.tgBot.GetUpdatesChan(tb.updateConfig)
	if err != nil {
		log.Fatal("tgBot.GetUpdatesChan: ", err)
	}

	log.Info("Start tg bot!")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("Uptime: %s", time.Since(uptime).Round(time.Second))
				tb.tgBot.Send(msg)
			}
			continue
		}
	}
}

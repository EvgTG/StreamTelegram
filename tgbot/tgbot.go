package tgbot

import (
	"StreamTelegram/go-log"
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
	"time"
)

type TGBot struct {
	tgBot            *tgbotapi.BotAPI
	updateConfig     tgbotapi.UpdateConfig
	toID             int64
	errorToID        int64
	uList            []int64
	NumberIterations int
	youtubeService   *youtube.Service
}

func New(proxy, token string, toID, errorToID int64, youtubeService *youtube.Service, uList []int64) (*TGBot, error) {
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

	return &TGBot{tgBot, u, toID, errorToID, uList, 0, youtubeService}, nil
}

func (tb *TGBot) SendNotification(text string) {
	msg := tgbotapi.NewMessage(tb.toID, text)
	tb.tgBot.Send(msg)
}

func (tb *TGBot) SendLog(text string) {
	if tb.errorToID == 0 {
		return
	}
	msg := tgbotapi.NewMessage(tb.errorToID, text)
	tb.tgBot.Send(msg)
}

func (tb *TGBot) Start() {
	uptime := time.Now()
	updates, err := tb.tgBot.GetUpdatesChan(tb.updateConfig)
	if err != nil {
		tb.SendLog("ERR tgbot.Start - tgBot.GetUpdatesChan: " + err.Error())
		log.Fatal("ERR tgbot.Start - tgBot.GetUpdatesChan: ", err.Error())
	}

	log.Info("Start tg bot!")
	for update := range updates {
		if update.Message != nil {
			if !userInList(tb.uList, update.Message.Chat.ID) {
				continue
			}
		}

		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
			if !userInList(tb.uList, update.CallbackQuery.Message.Chat.ID) {
				continue
			}
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("Hi!")
				tb.tgBot.Send(msg)
			case "status":
				msg.Text = fmt.Sprintf("Uptime: %s\nNumber of iterations: %v", time.Since(uptime).Round(time.Second), tb.NumberIterations)
				tb.tgBot.Send(msg)
			}
			continue
		}
	}
}

func userInList(list []int64, a int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

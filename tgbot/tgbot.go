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
	uptime           time.Time
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

	return &TGBot{
		tgBot:            tgBot,
		updateConfig:     u,
		toID:             toID,
		errorToID:        errorToID,
		uList:            uList,
		NumberIterations: 0,
		youtubeService:   youtubeService,
		uptime:           time.Now(),
	}, nil
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

			if update.CallbackQuery.Data == "update_status" {
				text, inlineKeyboard := tb.textStatus()
				tb.tgBot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text))
				tb.tgBot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, inlineKeyboard))
			}

			tb.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
			continue
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
				msg.Text, msg.ReplyMarkup = tb.textStatus()
				tb.tgBot.Send(msg)
			}
			continue
		}
	}
}

func (tb *TGBot) textStatus() (string, tgbotapi.InlineKeyboardMarkup) {
	text := fmt.Sprintf("Uptime: %s\nNumber of iterations: %v", time.Since(tb.uptime).Round(time.Second), tb.NumberIterations)
	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("🔄Update", "update_status"))
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	return text, inlineKeyboard
}

func userInList(list []int64, a int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

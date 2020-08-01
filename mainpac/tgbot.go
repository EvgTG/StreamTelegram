package mainpac

import (
	"StreamTelegram/go-log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"time"
)

func (s *Service) StartTG() {
	updates, err := s.tg.tgBot.GetUpdatesChan(s.tg.updateConfig)
	s.FatalTG("StartTG - tgBot.GetUpdatesChan()", err)

	log.Info("Start tg bot!")
	for update := range updates {
		if update.Message != nil {
			if !userInList(s.tg.userList, update.Message.Chat.ID) {
				continue
			}
		}

		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
			if !userInList(s.tg.userList, update.CallbackQuery.Message.Chat.ID) {
				continue
			}

			if update.CallbackQuery.Data == "update_status" {
				text, inlineKeyboard := s.tg.textStatus()
				s.tg.tgBot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text))
				s.tg.tgBot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, inlineKeyboard))
			}

			s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
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
				s.tg.tgBot.Send(msg)
			case "status":
				msg.Text, msg.ReplyMarkup = s.tg.textStatus()
				s.tg.tgBot.Send(msg)
			case "search":
				if update.Message.ReplyToMessage != nil {
					update.Message.Text = "/search " + update.Message.ReplyToMessage.Text
				}

				urls := strings.Replace(update.Message.Text, "/search ", "", 1)
				id := strings.Replace(urls, "https://youtu.be/", "", 1)
				id = strings.Replace(id, "https://www.youtube.com/watch?v=", "", 1)

				if strings.Contains(id, "&") {
					id = id[0:strings.Index(id, "&")]
				}

				video := s.yt.yts.Videos.List([]string{"snippet"})
				video.Id(id)
				videoRes, err := video.Do()
				if err != nil {
					msg.Text = fmt.Sprintf("Error youtube request - %v", err.Error())
					s.tg.tgBot.Send(msg)
					continue
				} else if len(videoRes.Items) == 0 {
					msg.Text = "Not found"
					s.tg.tgBot.Send(msg)
					continue
				}

				msg.ParseMode = "markdown"
				msg.Text = fmt.Sprintf("%v\nID: `%v`\n[URL](https://www.youtube.com/channel/%v),  [RSS](https://www.youtube.com/feeds/videos.xml?channel_id=%v)",
					videoRes.Items[0].Snippet.ChannelTitle, videoRes.Items[0].Snippet.ChannelId, videoRes.Items[0].Snippet.ChannelId, videoRes.Items[0].Snippet.ChannelId)
				s.tg.tgBot.Send(msg)
			}
			continue
		}
	}
}

func (tg *tg) SendNotification(text string) {
	msg := tgbotapi.NewMessage(tg.toID, text)
	tg.tgBot.Send(msg)
}

func (tg *tg) SendLog(text string) {
	if tg.errorToID == 0 {
		return
	}
	msg := tgbotapi.NewMessage(tg.errorToID, text)
	tg.tgBot.Send(msg)
}

func (tg *tg) textStatus() (string, tgbotapi.InlineKeyboardMarkup) {
	tm := time.Since(tg.uptime).Round(time.Second)
	var hours int
	var hoursStr string
	for tm.Hours() > 24 {
		tm -= time.Hour * 24
		hours++
		hoursStr = fmt.Sprintf("%vd", hours)
	}
	text := fmt.Sprintf("Uptime: %s\nNumber of iterations: %v", hoursStr+tm.String(), tg.numberIterations)
	buttons := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🔄Update", "update_status")}
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

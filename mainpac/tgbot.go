package mainpac

import (
	"StreamTelegram/go-log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
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

			switch update.CallbackQuery.Data {
			case "update_status":
				text, inlineKeyboard := s.tg.textStatus(s.yt.stop > 0)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text))
				s.tg.tgBot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
			case "start":
				if s.yt.stop == 2 {
					s.yt.stopch <- true
				}
				s.yt.stop = 0
				text, inlineKeyboard := s.tg.textStatus(s.yt.stop > 0)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text))
				s.tg.tgBot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Ok"))
			case "stop":
				s.yt.stop = 1
				text, inlineKeyboard := s.tg.textStatus(s.yt.stop > 0)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text))
				s.tg.tgBot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Ok"))
			case "delete":
				s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			}

			if strings.Contains(update.CallbackQuery.Data, "get_rss") {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
				if update.CallbackQuery.Data[0:6] == "nodel_" {
					update.CallbackQuery.Data = update.CallbackQuery.Data[6:len(update.CallbackQuery.Data)]
				} else {
					s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				}

				feed, err := s.tg.GetRSSFeed(strings.Replace(update.CallbackQuery.Data, "get_rss", "", 1))
				if err != nil {
					msg.Text = "Failed to get"
					s.tg.tgBot.Send(msg)
					s.tg.SendLog(err.Error())
					s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					continue
				}

				msg.ParseMode = "markdown"
				msg.Text = s.tg.textRSS(feed, s.loc)
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🗞Update", update.CallbackQuery.Data), tgbotapi.NewInlineKeyboardButtonData("❌Delete", "delete")}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
				s.tg.tgBot.Send(msg)
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			}

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
				msg.Text, msg.ReplyMarkup = s.tg.textStatus(s.yt.stop > 0)
				s.tg.tgBot.Send(msg)
			case "lastrss":
				msg.ParseMode = "markdown"
				msg.Text = s.tg.textRSS(&s.yt.lastRSS, s.loc)
				s.tg.tgBot.Send(msg)
			case "getrss":
				id := update.Message.Text[8:len(update.Message.Text)]
				feed, err := s.tg.GetRSSFeed(id)
				if err != nil {
					msg.Text = "Failed to get"
					s.tg.tgBot.Send(msg)
					s.tg.SendLog(err.Error())
				}
				msg.ParseMode = "markdown"
				msg.Text = s.tg.textRSS(feed, s.loc)
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🗞Update", "get_rss"+id), tgbotapi.NewInlineKeyboardButtonData("❌Delete", "delete")}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
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

				//TODO приделать проверку id

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
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🗞Get RSS", "nodel_get_rss"+videoRes.Items[0].Snippet.ChannelId)}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
				s.tg.tgBot.Send(msg)
			}
			continue
		}
	}
}

func (tg *tg) SendNotification(text string) {
	for _, id := range tg.toID {
		msg := tgbotapi.NewMessage(id, text)
		tg.tgBot.Send(msg)
	}
}

func (tg *tg) SendLog(text string) {
	if tg.errorToID == 0 {
		return
	}
	msg := tgbotapi.NewMessage(tg.errorToID, text)
	tg.tgBot.Send(msg)
}

func (tg *tg) textStatus(stop bool) (string, tgbotapi.InlineKeyboardMarkup) {
	tm := time.Since(tg.uptime).Round(time.Second)
	var hours int
	var hoursStr string
	for tm.Hours() > 24 {
		tm -= time.Hour * 24
		hours++
		hoursStr = fmt.Sprintf("%vd", hours)
	}
	text := fmt.Sprintf("Uptime: %s\nPause: %v\nNumber of iterations: %v", hoursStr+tm.String(), stop, tg.numberIterations)
	buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🔄Update", "update_status")}
	buttons2 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("▶️", "start"), tgbotapi.NewInlineKeyboardButtonData("⏸", "stop")}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(buttons1, buttons2)
	return text, inlineKeyboard
}

func (tg *tg) textRSS(feed *gofeed.Feed, loc *time.Location) string {
	str := fmt.Sprintf("[%v](%v)\n", feed.Title, feed.Link)
	for n, item := range feed.Items {
		str += fmt.Sprintf("%v. [%v](%v)\n%v\n", n+1, item.Title, item.Link, item.UpdatedParsed.In(loc).Format("2006 01.02 15:04"))
	}
	return str
}

func (tg *tg) GetRSSFeed(channelID string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID)
	if err != nil {
		return nil, fmt.Errorf("mainpac.GetRSSFeed - fp.ParseURL(): %s", err)
	}
	return feed, nil
}

func userInList(list []int64, a int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

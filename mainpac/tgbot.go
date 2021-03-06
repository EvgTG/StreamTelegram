package mainpac

import (
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (s *Service) StartTG() {
	rgxVID := regexp.MustCompile(`[A-Za-z0-9_-]{11}`)
	rgxChID := regexp.MustCompile(`UC[A-Za-z0-9_-]{22}`)

	updates, err := s.tg.tgBot.GetUpdatesChan(s.tg.updateConfig)
	s.FatalTG("StartTG - tgBot.GetUpdatesChan()", err)

	log.Info("Start tg bot!")
	for update := range updates {
		if update.Message != nil {
			if !userInList(s.tg.userList, update.Message.Chat.ID) {
				continue
			}
		}

		//context input (settings)
		if s.tg.callbackQuery != "" && update.Message != nil && update.Message.Chat != nil { //st_edit_chiddelid_***
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if strings.Contains(s.tg.callbackQuery, "st_") {
				if strings.Contains(s.tg.callbackQuery, "edit_") {
					switch {
					case strings.Contains(s.tg.callbackQuery, "chid"):
						var ints []int64
						var errbl bool

						strs := strings.Split(update.Message.Text, ",")
						deleteSpaces(&strs)
						for _, v := range strs {
							a, err := strconv.ParseInt(v, 10, 64)
							if err != nil {
								errbl = true
							}
							ints = append(ints, a)
						}

						if errbl {
							msg.Text = "Invalid formatting. Try again"
							s.tg.tgBot.Send(msg)
							continue
						}

						s.tg.toID = ints
						st := s.db.GetLs()
						st.DBPriority.ToID = ints
						st.DBPriority.ToIDBL = true
						err := s.db.SetLs(&st)
						s.FatalTG("StartTG - s.db.SetLs()", err)

						s.tg.callbackQuery = s.tg.callbackQuery[12:len(s.tg.callbackQuery)]
						msg.Text = "Changed"
						mes, err := s.tg.tgBot.Send(msg)
						if err == nil && mes.MessageID != 0 {
							go func(chatid int64, id int) {
								time.Sleep(time.Second * 3)
								s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(chatid, id))
							}(update.Message.Chat.ID, mes.MessageID)
						}
					case strings.Contains(s.tg.callbackQuery, "cycletime"):
						cycleTimeInt, err := strconv.ParseInt(update.Message.Text, 10, 64)
						if err != nil {
							msg.Text = "Invalid formatting. Try again"
							s.tg.tgBot.Send(msg)
							continue
						}

						cycleTime := time.Duration(cycleTimeInt) * time.Minute

						s.yt.CycleTime = cycleTime
						st := s.db.GetLs()
						st.DBPriority.CycleTime = cycleTime
						st.DBPriority.CycleTimeBL = true
						err = s.db.SetLs(&st)
						s.FatalTG("StartTG - s.db.SetLs()", err)

						s.tg.callbackQuery = s.tg.callbackQuery[17:len(s.tg.callbackQuery)]
						msg.Text = "Changed"
						mes, err := s.tg.tgBot.Send(msg)
						if err == nil && mes.MessageID != 0 {
							go func(chatid int64, id int) {
								time.Sleep(time.Second * 3)
								s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(chatid, id))
							}(update.Message.Chat.ID, mes.MessageID)
						}
					}
				}
			}

			if strings.Contains(s.tg.callbackQuery, "delid_") { //delid_***
				id, _ := strconv.Atoi(s.tg.callbackQuery[6:len(s.tg.callbackQuery)])
				s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, id))
			}

			s.tg.callbackQuery = ""
			continue
		}

		//buttons
		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat != nil {
			if !userInList(s.tg.userList, update.CallbackQuery.Message.Chat.ID) {
				continue
			}

			switch {
			case update.CallbackQuery.Data == "update_status":
				text, inlineKeyboard := s.tg.statusMes(s.yt.stop > 0, s.yt.lastTime, s.loc, s.yt.CycleTime)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
			case update.CallbackQuery.Data == "start":
				if s.yt.stop == 2 {
					s.yt.stopch <- true
				}
				s.yt.stop = 0
				text, inlineKeyboard := s.tg.statusMes(s.yt.stop > 0, s.yt.lastTime, s.loc, s.yt.CycleTime)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Ok"))
			case update.CallbackQuery.Data == "stop":
				s.yt.stop = 1
				text, inlineKeyboard := s.tg.statusMes(s.yt.stop > 0, s.yt.lastTime, s.loc, s.yt.CycleTime)
				s.tg.tgBot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text, inlineKeyboard))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Ok"))
			case update.CallbackQuery.Data == "delete":
				s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			case update.CallbackQuery.Data == "cancel":
				s.tg.callbackQuery = ""
				s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Cancelled"))
			case strings.Contains(update.CallbackQuery.Data, "get_rss"):
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
				s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
			case strings.Contains(update.CallbackQuery.Data, "st_"):
				if strings.Contains(update.CallbackQuery.Data, "st_update_") { //st_update_
					switch {
					case strings.Contains(update.CallbackQuery.Data, "chid"):
						s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
						s.tg.toIDMes(s.db, update.CallbackQuery.Message.Chat.ID)
						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
						continue
					case strings.Contains(update.CallbackQuery.Data, "cycletime"):
						s.tg.tgBot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
						s.cycleTimeMes(update.CallbackQuery.Message.Chat.ID)
						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
						continue
					}
				}

				msgCancel := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
				buttons := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("❌Cancel", "cancel")}
				msgCancel.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)

				switch {
				case strings.Contains(update.CallbackQuery.Data, "edit_"):
					switch {
					case strings.Contains(update.CallbackQuery.Data, "chid"):
						s.tg.callbackQuery = update.CallbackQuery.Data
						msgCancel.Text = "Enter the IDs separated by commas."
						mes, err := s.tg.tgBot.Send(msgCancel)
						if err == nil && mes.MessageID != 0 {
							s.tg.callbackQuery += "delid_" + strconv.Itoa(mes.MessageID)
						}
						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					case strings.Contains(update.CallbackQuery.Data, "cycletime"):
						s.tg.callbackQuery = update.CallbackQuery.Data
						msgCancel.Text = "Enter the channel check time."
						mes, err := s.tg.tgBot.Send(msgCancel)
						if err == nil && mes.MessageID != 0 {
							s.tg.callbackQuery += "delid_" + strconv.Itoa(mes.MessageID)
						}
						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
					}
				case strings.Contains(update.CallbackQuery.Data, "bl_"):
					switch {
					case strings.Contains(update.CallbackQuery.Data, "chid"):
						settings := s.db.GetLs()

						if settings.DBPriority.ToIDBL {
							s.tg.toID = s.envVars.toID
							settings.DBPriority.ToIDBL = false
						} else {
							s.tg.toID = settings.DBPriority.ToID
							settings.DBPriority.ToIDBL = true
						}

						err = s.db.SetLs(&settings)
						s.FatalTG("StartTG - s.db.SetLs()", err)

						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
					case strings.Contains(update.CallbackQuery.Data, "cycletime"):
						settings := s.db.GetLs()

						if settings.DBPriority.CycleTimeBL {
							s.yt.CycleTime = s.envVars.cycleTime
							settings.DBPriority.CycleTimeBL = false
						} else {
							s.yt.CycleTime = settings.DBPriority.CycleTime
							settings.DBPriority.CycleTimeBL = true
						}

						err = s.db.SetLs(&settings)
						s.FatalTG("StartTG - s.db.SetLs()", err)

						s.tg.tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Updated"))
					}
				}
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		//commands
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("Hi!")
				s.tg.tgBot.Send(msg)
			case "status":
				msg.Text, msg.ReplyMarkup = s.tg.statusMes(s.yt.stop > 0, s.yt.lastTime, s.loc, s.yt.CycleTime)
				s.tg.tgBot.Send(msg)
			case "lastrss":
				msg.ParseMode = "markdown"
				msg.Text = s.tg.textRSS(&s.yt.lastRSS, s.loc)
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("❌Delete", "delete")}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
				s.tg.tgBot.Send(msg)
			case "getrss":
				id := strings.ReplaceAll(update.Message.Text[7:len(update.Message.Text)], " ", "")
				if !rgxChID.MatchString(id) {
					msg.Text = "Invalid id"
					s.tg.tgBot.Send(msg)
					continue
				}
				feed, err := s.tg.GetRSSFeed(id)
				if err != nil {
					msg.Text = "Failed to get"
					s.tg.tgBot.Send(msg)
					s.tg.SendLog(err.Error())
					continue
				}
				msg.ParseMode = "markdown"
				msg.Text = s.tg.textRSS(feed, s.loc)
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🗞Update", "get_rss"+id), tgbotapi.NewInlineKeyboardButtonData("❌Delete", "delete")}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
				s.tg.tgBot.Send(msg)
			case "settings":
				msg.Text = "/toid - modify targets for notification\n/cycletime - change the channel check time"
				s.tg.tgBot.Send(msg)
			case "toid":
				s.tg.toIDMes(s.db, update.Message.Chat.ID)
			case "cycletime":
				s.cycleTimeMes(update.Message.Chat.ID)
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

				if !rgxVID.MatchString(id) {
					msg.Text = "Invalid video id"
					s.tg.tgBot.Send(msg)
					continue
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
				buttons1 := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("🗞Get RSS", "nodel_get_rss"+videoRes.Items[0].Snippet.ChannelId)}
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1)
				s.tg.tgBot.Send(msg)
			}
			continue
		}
	}
}

func (tg *tg) SendNotification(text string) {
	log.Debug("notif.1 ", time.Now(), " text= ", text)
	for _, id := range tg.toID {
		msg := tgbotapi.NewMessage(id, text)
		tg.tgBot.Send(msg)
		log.Debug("notif.2 ", time.Now(), " text= ", id)
	}
}

func (tg *tg) toIDMes(db *model.Model, id int64) {
	msg := tgbotapi.NewMessage(id, "")
	var change, status string
	var toIDstrs []string

	settings := db.GetLs()

	if settings.DBPriority.ToIDBL {
		status = "DataBase"
		change = "Select Environment"
	} else {
		status = "Environment"
		change = "Select DataBase"
	}

	for _, v := range tg.toID {
		toIDstrs = append(toIDstrs, strconv.FormatInt(v, 10))
	}
	msg.ParseMode = "markdown"
	msg.Text = fmt.Sprintf("To ID: `%v`\nSource: %v", strings.Join(toIDstrs, ","), status)
	buttons1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("✏️Edit", "st_edit_chid"),
		tgbotapi.NewInlineKeyboardButtonData("🔄Update", "st_update_chid"),
	}
	buttons2 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(change, "st_bl_chid"),
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1, buttons2)
	tg.tgBot.Send(msg)
}

func (s *Service) cycleTimeMes(id int64) {
	msg := tgbotapi.NewMessage(id, "")
	var change, status string

	settings := s.db.GetLs()

	if settings.DBPriority.CycleTimeBL {
		status = "DataBase"
		change = "Select Environment"
	} else {
		status = "Environment"
		change = "Select DataBase"
	}

	msg.Text = fmt.Sprintf("Cycle time: %vm\nSource: %v", s.yt.CycleTime.Minutes(), status)
	buttons1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("✏️Edit", "st_edit_cycletime"),
		tgbotapi.NewInlineKeyboardButtonData("🔄Update", "st_update_cycletime"),
	}
	buttons2 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(change, "st_bl_cycletime"),
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons1, buttons2)
	s.tg.tgBot.Send(msg)
}

func (tg *tg) SendLog(text string) {
	if tg.errorToID == 0 {
		return
	}
	msg := tgbotapi.NewMessage(tg.errorToID, text)
	tg.tgBot.Send(msg)
}

func (tg *tg) statusMes(stop bool, lastTm time.Time, loc *time.Location, cycleTime time.Duration) (string, tgbotapi.InlineKeyboardMarkup) {
	tm := time.Since(tg.uptime).Round(time.Second)
	var hours int
	var hoursStr string
	for tm.Hours() > 24 {
		tm -= time.Hour * 24
		hours++
		hoursStr = fmt.Sprintf("%vd", hours)
	}
	text := fmt.Sprintf("Uptime: %s\nPause: %v\nNumber of iterations: %v\nCycle time: %vm\nTime of last check RSS:\n%v", hoursStr+tm.String(), stop, tg.numberIterations, cycleTime.Minutes(), lastTm.In(loc).Format("01.02 15:04 -07:00 MST"))
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

func deleteSpaces(list *[]string) {
	for i := 0; i < len(*list); i++ {
		(*list)[i] = strings.ReplaceAll((*list)[i], " ", "")
	}
}

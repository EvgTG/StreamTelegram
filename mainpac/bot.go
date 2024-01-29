package mainpac

import (
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (s *Service) InitBot() {
	adminOnly := s.Bot.Group()
	adminOnly.Use(middleware.Whitelist(s.Bot.AdminList...))

	// Команды

	adminOnly.Handle("/start", s.TgStart)
	adminOnly.Handle("/help", s.TgStart)
	adminOnly.Handle("/test", s.TgTest)
	adminOnly.Handle("/test_notify", s.TgTestNotify)
	adminOnly.Handle("/adm", s.TgAdm)
	adminOnly.Handle("/status", s.TgStatus)
	adminOnly.Handle("/logs", s.TgLogs)
	adminOnly.Handle("/set_commands", s.TgSetCommands)
	adminOnly.Handle("/set_channel", s.TgSetChannelID)
	adminOnly.Handle("/set_twitch", s.TgSetTwitchNick)
	adminOnly.Handle("/set_twitch_client", s.TgTwitchClient)
	adminOnly.Handle("/twitch_auth_url", s.TgTwitchAuthURL)
	adminOnly.Handle("/twitch_auth", s.TgTwitchAuth)
	adminOnly.Handle("/get_channel", s.TgGetChannelID)
	adminOnly.Handle("/set_dur", s.TgSetCycleDuration)
	adminOnly.Handle("/locs", s.TgLocs)
	adminOnly.Handle("/set_loc", s.TgSetLoc)
	adminOnly.Handle("/notify", s.TgNotify)
	adminOnly.Handle("/notify_add", s.TgNotifyAdd)
	adminOnly.Handle("/notify_del", s.TgNotifyDel)
	adminOnly.Handle("/last_rss", s.TgLastRSS)
	adminOnly.Handle("/type_of_vid", s.TgTypeOfVideo)
	adminOnly.Handle(tb.OnText, s.TgCallbackQuery)

	// Кнопки

	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "test"), s.TgTestBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "delete"), s.TgDeleteBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "cancel"), s.TgCancelReplyMarkup)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "status_update"), s.TgStatusUpdate)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "pause"), s.TgPause)

	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "get_logs"), s.TgGetLogsBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "clear_logs"), s.TgClearLogsBtn)

	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "set_loc"), s.TgSetLoc)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "locs_update"), s.TgLocsUpdateBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "locs_clear"), s.TgLocsClearBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "time_city"), s.TgLocsCity)

	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "notify_up"), s.TgNotifyUpdateBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "notify_add"), s.TgNotifyAdd)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "notify_del"), s.TgNotifyDel)
}

/*

s.Bot.Handle("/", s.Tg)
s.Bot.addBtn(rm.Data("", ""), "", s.Tg)

*/

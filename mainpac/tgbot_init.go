package mainpac

import tb "gopkg.in/tucnak/telebot.v3"

func (s *Service) InitBot() {
	// Команды

	s.Bot.Handle("/start", s.TgStart)
	s.Bot.Handle("/help", s.TgStart)

	// Админские команды

	s.Bot.Handle("/test", s.TgTest)
	s.Bot.Handle("/adm", s.TgAdm)
	s.Bot.Handle("/status", s.TgStatus)
	s.Bot.Handle("/logs", s.TgLogs)
	s.Bot.Handle("/setCommands", s.TgSetCommands)
	s.Bot.Handle(tb.OnText, s.TgCallbackQuery)

	// Кнопки

	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "test"), s.TgTestBtn)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "delete"), s.TgDeleteBtn)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "cancel"), s.TgCancelReplyMarkup)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "status_update"), s.TgStatusUpdate)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "get_logs"), s.TgGetLogsBtn)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "clear_logs"), s.TgClearLogsBtn)
}

/*

s.Bot.Handle("/", s.Tg)
s.Bot.addBtn(rm.Data("", ""), "", s.Tg)

*/

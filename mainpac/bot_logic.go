package mainpac

import (
	"fmt"
	"github.com/rotisserie/eris"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
	"streamtg/go-log"
)

/*
TgSome            - команда
TgSome+Update/Btn - кнопка обновления/обычная
TgSome+Func       - логика работы
Но они обязательны только все вместе
*/

func (s *Service) TgStart(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "start"), s.Bot.Markup(x, "remove_keyboard"))
	return
}

func (s *Service) TgTest(x tb.Context) (errReturn error) {
	x.Send("Test", s.Bot.Markup(x, "test"), tb.NoPreview)
	return
}

func (s *Service) TgTestBtn(x tb.Context) (errReturn error) {
	x.Send("Test", &tb.SendOptions{ReplyTo: x.Message()}, s.Bot.Markup(x, "test"), tb.NoPreview)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "test"})
	return
}

func (s *Service) TgAdm(x tb.Context) (errReturn error) {
	text := fmt.Sprintf("" +
		"\n/start - приветствие" +
		"\n/status - статус работы" +
		"\n/logs - действия над логами" +
		"\n/set_commands - установить меню бота" +
		"\n/set_dur - время проверки стримов" +
		"\n/notify - каналы для вывода уведомлений" +
		"\n/test_notify - проверка уведомлений" +
		"\n/locs - настройка часовых поясов в уведомлениях" +
		"\n\nYouTube" +
		"\n/set_channel - установить канал" +
		"\n/get_channel - получить id канала" +
		"\n/type_of_vid - получить тип видео" +
		"\n/last_rss - последние видео" +
		"\n\nTwitch" +
		"\n/set_twitch_client - настройка клиента" +
		"\n/twitch_auth_url - ссылка для аутентификации" +
		"\n/twitch_auth - аутентификация" +
		"\n/set_twitch - установить twitch ник",
	)

	x.Send(text)
	return
}

func (s *Service) TgStatus(x tb.Context) (errReturn error) {
	text, rm := s.TgStatusFunc(x)
	mes, err := s.Bot.Send(x.Sender(), text, rm, tb.NoPreview)
	if err == nil && mes != nil {
		s.Bot.Pin(mes)
	}
	return
}

func (s *Service) TgStatusUpdate(x tb.Context) (errReturn error) {
	text, rm := s.TgStatusFunc(x)
	_, err := s.Bot.Edit(x.Message(), text, rm, tb.NoPreview)
	if err != nil {
		s.Bot.Delete(x.Message())
		mes, err := s.Bot.Send(x.Sender(), text, rm)
		if err == nil && mes != nil {
			s.Bot.Pin(mes)
		}
	}

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Обновлено"})
	return
}

func (s *Service) TgStatusFunc(x tb.Context) (string, *tb.ReplyMarkup) {
	pause := false
	s.YouTubeTwitch.PauseMutex.Lock()
	if s.YouTubeTwitch.Pause > 0 {
		pause = true
	}
	s.YouTubeTwitch.PauseMutex.Unlock()

	text := fmt.Sprintf(""+
		"Launch time: %s"+
		"\nUptime: %s"+
		"\nPause: %v"+
		"\n"+
		"\nChannel ID: <a href=\"youtube.com/channel/%s\">%s</a>"+
		"\nTwitch nick: <a href=\"twitch.tv/%s\">%s</a>"+
		"\nTwitch client-%v auth-%v"+
		"\nCycle duration: %vmin"+
		"\nN iterations: %v"+
		"\nLast check: %s"+
		"\n<a href=\"%s\">RSS url</a>",

		s.Bot.Uptime.In(s.Loc).Format("2006.01.02 15:04:05 MST"), s.Bot.uptimeString(s.Bot.Uptime), pause,
		s.YouTubeTwitch.ChannelID, s.YouTubeTwitch.ChannelID, s.YouTubeTwitch.TwitchNick, s.YouTubeTwitch.TwitchNick,
		s.YouTubeTwitch.Twitch.ClientOK(), s.YouTubeTwitch.Twitch.AuthOK(),
		s.YouTubeTwitch.CycleDurationMinutes, s.YouTubeTwitch.NumberIterations,
		s.YouTubeTwitch.LastTime.In(s.Loc).Format("2006.01.02 15:04:05 MST"),
		"https://www.youtube.com/feeds/videos.xml?channel_id="+s.YouTubeTwitch.ChannelID,
	)

	rm := s.Bot.Markup(x, "status")

	return text, rm
}

func (s *Service) TgPause(x tb.Context) (errReturn error) {
	s.YouTubeTwitch.SetPause()
	x.Respond()
	s.TgStatusUpdate(x)
	return
}

func (s *Service) TgLogs(x tb.Context) (errReturn error) {
	text := "1. Получить файл логов\n2. Очистить файл логов"
	x.Send(text, s.Bot.Markup(x, "logs"))
	return
}

func (s *Service) TgGetLogsBtn(x tb.Context) (errReturn error) {
	_, err := s.Bot.Send(x.Sender(), &tb.Document{File: tb.FromDisk("files/logrus.log"), FileName: "logrus.log"})
	if err != nil {
		s.Bot.Send(x.Sender(), eris.Wrap(err, "Ошибка отправки файла.").Error())
	}
	x.Respond()
	return
}

func (s *Service) TgClearLogsBtn(x tb.Context) (errReturn error) {
	os.Truncate("files/logrus.log", 0)
	log.Info("Очищено")

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Очищено", ShowAlert: true})
	return
}

func (s *Service) TgCallbackQuery(x tb.Context) (errReturn error) {
	if x.Message().IsForwarded() {
		x.Send(fmt.Sprintf("id <code>%v</code>", x.Message().OriginalChat.ID))
	}

	switch s.Bot.CallbackQuery[x.Chat().ID] {
	case "": //Нет в CallbackQuery - игнор
	case "test":

	}
	return
}

func (s *Service) TgSetCommands(x tb.Context) (errReturn error) {
	err := x.Bot().SetCommands(s.Bot.Layout.Commands())
	if err != nil {
		x.Send(eris.Wrap(err, "x.Bot().SetCommands()").Error())
		return
	}

	x.Send("Готово.")
	return
}

func (s *Service) TgDeleteBtn(x tb.Context) (errReturn error) {
	x.Respond()
	x.Delete()
	return
}

func (s *Service) TgCancelReplyMarkup(x tb.Context) (errReturn error) {
	delete(s.Bot.CallbackQuery, x.Chat().ID)
	x.Send("Отменено.", &tb.ReplyMarkup{RemoveKeyboard: true})
	return
}

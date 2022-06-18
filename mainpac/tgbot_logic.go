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
	if s.Bot.isNotAdmin(x) {
		return
	}

	x.Send(s.Bot.Text(x, "start"), s.Bot.Markup(x, "remove_keyboard"))
	return
}

func (s *Service) TgTest(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	x.Send("Test", s.Bot.Markup(x, "test"), tb.ModeHTML, tb.NoPreview)
	return
}

func (s *Service) TgTestBtn(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	x.Send("Test", &tb.SendOptions{ReplyTo: x.Message()}, s.Bot.Markup(x, "test"), tb.ModeHTML, tb.NoPreview)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "test"})
	return
}

func (s *Service) TgAdm(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	text := fmt.Sprintf("" +
		"\n/start - приветствие" +
		"\n/status - статус работы" +
		"\n/logs - действия над логами" +
		"\n/set_commands - установить меню бота" +
		"\n/set_channel - установить канал" +
		"\n/get_channel - получить id канала" +
		"\n/set_dur - время обновления информации" +
		"\n/locs - настройка часовых поясов в уведомлениях" +
		"\n/set_loc - добавить часовой пояс",
	)

	x.Send(text, tb.ModeHTML)
	return
}

func (s *Service) TgStatus(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	text, rm := s.TgStatusFunc(x)
	mes, err := s.Bot.Send(x.Sender(), text, rm)
	if err == nil && mes != nil {
		s.Bot.Pin(mes)
	}
	return
}

func (s *Service) TgStatusUpdate(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	text, rm := s.TgStatusFunc(x)
	_, err := s.Bot.Edit(x.Message(), text, rm)
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
	text := fmt.Sprintf("Запущен: %s\nUptime: %s\n\nChannel ID: <a href=\"youtube.com/channel/%s\">%s</a>\nCycle duration: %vmin",
		s.Bot.Uptime.In(s.Loc).Format("2006.01.02 15:04:05 MST"), s.Bot.uptimeString(s.Bot.Uptime),
		s.YouTube.ChannelID, s.YouTube.ChannelID, s.YouTube.CycleDuration,
	)

	rm := s.Bot.Markup(x, "status")

	return text, rm
}

func (s *Service) TgLogs(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	text := "1. Получить файл логов\n2. Очистить файл логов"
	x.Send(text, s.Bot.Markup(x, "logs"), tb.ModeHTML)
	return
}

func (s *Service) TgGetLogsBtn(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	_, err := s.Bot.Send(x.Sender(), &tb.Document{File: tb.FromDisk("files/logrus.log"), FileName: "logrus.log"})
	if err != nil {
		s.Bot.Send(x.Sender(), eris.Wrap(err, "Ошибка отправки файла.").Error())
	}
	x.Respond()
	return
}

func (s *Service) TgClearLogsBtn(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	os.Truncate("files/logrus.log", 0)
	log.Info("Очищено")

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Очищено", ShowAlert: true})
	return
}

func (s *Service) TgCallbackQuery(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	switch s.Bot.CallbackQuery[x.Chat().ID] {
	case "": //Нет в CallbackQuery - игнор
	case "test":

	}
	return
}

func (s *Service) TgSetCommands(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

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
	if s.Bot.isNotAdmin(x) {
		return
	}

	delete(s.Bot.CallbackQuery, x.Chat().ID)
	x.Send("Отменено.", &tb.ReplyMarkup{RemoveKeyboard: true})
	return
}

package mainpac

import (
	"fmt"
	"github.com/rotisserie/eris"
	tb "gopkg.in/tucnak/telebot.v3"
	"strconv"
	"streamtg/util"
	"strings"
	"time"
)

func (s *Service) TgSetChannelID(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/set_channel" {
		x.Send(s.Bot.Text(x, "set_channel_empty"))
		return
	}

	url := strings.Replace(x.Text(), "/set_channel ", "", 1)
	url = strings.Replace(url, " ", "", -1)

	id, err := util.GetChannelIDByUrl(url)
	if err != nil {
		x.Send(eris.Wrap(err, "GetChannelIDByUrl()").Error())
		return
	}

	err = s.MiniDB.SetChannelID(id)
	if err != nil {
		x.Send(eris.Wrap(err, "MiniDB.SetChannelID()").Error())
		return
	}
	s.YouTube.ChannelID = id

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgGetChannelID(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/get_channel" {
		x.Send(s.Bot.Text(x, "get_channel_empty"))
		return
	}

	url := strings.Replace(x.Text(), "/get_channel ", "", 1)
	url = strings.Replace(url, " ", "", -1)

	id, err := util.GetChannelIDByUrl(url)
	if err != nil {
		x.Send(eris.Wrap(err, "GetChannelIDByUrl()").Error())
		return
	}

	x.Send("id - " + id)
	return
}

func (s *Service) TgSetCycleDuration(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/set_dur" {
		x.Send(s.Bot.Text(x, "set_dur"))
		return
	}

	text := strings.Replace(x.Text(), "/set_dur ", "", 1)
	text = strings.Replace(text, " ", "", -1)

	dur, err := strconv.Atoi(text)
	if err != nil {
		x.Send(s.Bot.Text(x, "err_format"))
		return
	}

	err = s.MiniDB.SetCycleDuration(dur)
	if err != nil {
		x.Send(eris.Wrap(err, "MiniDB.SetCycleDuration()").Error())
		return
	}
	s.YouTube.CycleDuration = dur

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgLocsFunc(x tb.Context) (string, *tb.ReplyMarkup) {
	text := fmt.Sprintf("%v\n\n", s.YouTube.Locs)
	tm := time.Now()

	for _, locStr := range s.YouTube.Locs {
		loc, _ := time.LoadLocation(locStr)
		text += tm.In(loc).Format(s.YouTube.TimeFormat) + "\n"
	}

	rm := s.Bot.Layout.Markup(x, "locs")

	return text, rm
}

func (s *Service) TgLocs(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	x.Send(s.TgLocsFunc(x))
	return
}

func (s *Service) TgLocsUpdateBtn(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}
	x.Respond()
	x.Edit(s.TgLocsFunc(x))
	return
}

func (s *Service) TgLocsClearBtn(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	s.YouTube.Locs = []string{}
	err := s.MiniDB.SetLocs(s.YouTube.Locs)
	if err != nil {
		x.Send(err.Error())
		x.Respond()
		return
	}

	x.Respond()
	x.Edit(s.TgLocsFunc(x))
	return
}

func (s *Service) TgSetLoc(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/set_loc" || x.Callback() != nil {
		x.Send(s.Bot.Text(x, "set_loc"), tb.NoPreview)
		x.Respond()
		return
	}

	text := strings.Replace(x.Text(), "/set_loc ", "", 1)
	text = strings.Replace(text, " ", "", -1)

	loc, err := time.LoadLocation(text)
	if err != nil {
		x.Send(err.Error())
		return
	}

	s.YouTube.Locs = append(s.YouTube.Locs, loc.String())
	err = s.MiniDB.SetLocs(s.YouTube.Locs)
	if err != nil {
		x.Send(err.Error())
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgLocsCity(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	s.YouTube.TimeCity = !s.YouTube.TimeCity
	s.YouTube.TimeFormat = TimeFormatCity(s.YouTube.TimeCity)

	err := s.MiniDB.SetTimeWithCity(s.YouTube.TimeCity)
	if err != nil {
		x.Send(err.Error())
		return
	}

	x.Respond()
	x.Edit(s.TgLocsFunc(x))
	return
}

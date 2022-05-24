package mainpac

import (
	"github.com/rotisserie/eris"
	tb "gopkg.in/tucnak/telebot.v3"
	"strconv"
	"streamtg/util"
	"strings"
)

func (s *Service) TgSetChannelID(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/set_channel" {
		x.Send(s.Bot.Text(x, "set_channel_empty"), tb.ModeHTML)
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

	if x.Text() == "/set_channel" {
		x.Send(s.Bot.Text(x, "get_channel_empty"), tb.ModeHTML)
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
		x.Send(s.Bot.Text(x, "set_dur"), tb.ModeHTML)
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

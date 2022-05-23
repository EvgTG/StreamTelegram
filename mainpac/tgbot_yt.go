package mainpac

import (
	"github.com/rotisserie/eris"
	tb "gopkg.in/tucnak/telebot.v3"
	"streamtg/util"
	"strings"
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

	x.Send("Сделано. " + id)
	return
}

func (s *Service) TgGetChannelID(x tb.Context) (errReturn error) {
	if s.Bot.isNotAdmin(x) {
		return
	}

	if x.Text() == "/set_channel" {
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

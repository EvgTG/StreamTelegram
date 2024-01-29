package mainpac

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"streamtg/minidb"
	"streamtg/util"

	"github.com/rotisserie/eris"
	tb "gopkg.in/telebot.v3"
)

func (s *Service) TgSetChannelID(x tb.Context) (errReturn error) {
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
	s.YouTubeTwitch.ChannelID = id

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgSetTwitchNick(x tb.Context) (errReturn error) {
	if x.Text() == "/set_twitch" {
		x.Send(s.Bot.Text(x, "set_twitch_empty"))
		return
	}

	nick := strings.Replace(x.Text(), "/set_twitch ", "", 1)
	nick = strings.Replace(nick, " ", "", -1)

	err := s.MiniDB.SetTwitchNick(nick)
	if err != nil {
		x.Send(eris.Wrap(err, "MiniDB.SetTwitchNick()").Error())
		return
	}
	s.YouTubeTwitch.TwitchNick = nick

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgGetChannelID(x tb.Context) (errReturn error) {
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
	s.YouTubeTwitch.CycleDurationMinutes = dur

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgLocsFunc(x tb.Context) (string, *tb.ReplyMarkup) {
	text := fmt.Sprintf("%v\n\n", s.YouTubeTwitch.Locs)
	tm := time.Now()

	for _, locStr := range s.YouTubeTwitch.Locs {
		loc, _ := time.LoadLocation(locStr)
		text += tm.In(loc).Format(s.YouTubeTwitch.TimeFormat) + "\n"
	}

	rm := s.Bot.Layout.Markup(x, "locs")
	return text, rm
}

func (s *Service) TgLocs(x tb.Context) (errReturn error) {
	x.Send(s.TgLocsFunc(x))
	return
}

func (s *Service) TgLocsUpdateBtn(x tb.Context) (errReturn error) {
	x.Respond()
	x.Edit(s.TgLocsFunc(x))
	return
}

func (s *Service) TgLocsClearBtn(x tb.Context) (errReturn error) {
	s.YouTubeTwitch.Locs = []string{}
	err := s.MiniDB.SetLocs(s.YouTubeTwitch.Locs)
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

	s.YouTubeTwitch.Locs = append(s.YouTubeTwitch.Locs, loc.String())
	err = s.MiniDB.SetLocs(s.YouTubeTwitch.Locs)
	if err != nil {
		x.Send(err.Error())
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgLocsCity(x tb.Context) (errReturn error) {
	s.YouTubeTwitch.TimeCity = !s.YouTubeTwitch.TimeCity
	s.YouTubeTwitch.TimeFormat = TimeFormatCity(s.YouTubeTwitch.TimeCity)

	err := s.MiniDB.SetTimeWithCity(s.YouTubeTwitch.TimeCity)
	if err != nil {
		x.Send(err.Error())
		return
	}

	x.Respond()
	x.Edit(s.TgLocsFunc(x))
	return
}

func (s *Service) TgNotify(x tb.Context) (errReturn error) {
	x.Send(s.TgNotifyFunc(x))
	return
}

func (s *Service) TgNotifyUpdateBtn(x tb.Context) (errReturn error) {
	x.Respond()
	x.Edit(s.TgNotifyFunc(x))
	return
}

func (s *Service) TgNotifyFunc(x tb.Context) (string, *tb.ReplyMarkup) {
	text := ""

	for i, channel := range s.Bot.NotifyList {
		text += fmt.Sprintf("%v. id: <code>%v</code> start", i+1, channel.ID)
		if channel.EndOfStream {
			text += " end"
		}
		text += "\n"
	}

	if len(s.Bot.NotifyList) == 0 {
		text = s.Bot.TextLocale("ru", "notify_nil")
	}

	rm := s.Bot.Layout.Markup(x, "notify")
	return text, rm
}

func (s *Service) TgNotifyAdd(x tb.Context) (errReturn error) {
	if x.Text() == "/notify_add" || x.Callback() != nil {
		x.Send(s.Bot.Text(x, "notify_add"), tb.NoPreview)
		x.Respond()
		return
	}

	text := strings.Replace(x.Text(), "/notify_add ", "", 1)
	text = strings.Replace(text, " ", "", -1)

	end := strings.Contains(text, "end")
	if end {
		text = strings.Replace(text, "end", "", 1)
	}

	id, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		x.Send(eris.Wrap(err, "parse id").Error())
		return
	}

	ch := minidb.Channel{
		ID:          id,
		EndOfStream: end,
	}

	s.Bot.NotifyList = append(s.Bot.NotifyList, ch)
	err = s.MiniDB.SetNotifyList(s.Bot.NotifyList)
	if err != nil {
		x.Send(eris.Wrap(err, "SetNotifyList").Error())
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgNotifyDel(x tb.Context) (errReturn error) {
	if x.Text() == "/notify_del" || x.Callback() != nil {
		x.Send(s.Bot.Text(x, "notify_del"), tb.NoPreview)
		x.Respond()
		return
	}

	text := strings.Replace(x.Text(), "/notify_del ", "", 1)
	text = strings.Replace(text, " ", "", -1)

	id, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		x.Send(eris.Wrap(err, "parse id").Error())
		return
	}

	iForDel := -1
	for i, channel := range s.Bot.NotifyList {
		if channel.ID == id {
			iForDel = i
		}
	}
	if iForDel >= 0 {
		s.Bot.NotifyList = append(s.Bot.NotifyList[:iForDel], s.Bot.NotifyList[iForDel+1:]...)
	} else {
		x.Send(s.Bot.Text(x, "notify_del_nil"))
		return
	}

	err = s.MiniDB.SetNotifyList(s.Bot.NotifyList)
	if err != nil {
		x.Send(eris.Wrap(err, "SetNotifyList").Error())
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgLastRSS(x tb.Context) (errReturn error) {
	if s.YouTubeTwitch.LastRSS_YT != nil {
		feed_yt := s.YouTubeTwitch.LastRSS_YT

		str1 := fmt.Sprintf("[%v](%v)\n", feed_yt.Title, feed_yt.Link)
		for n, item := range feed_yt.Items {
			str1 += fmt.Sprintf("%v. [%v](%v)\n%v\n", n+1, item.Title, item.Link, item.UpdatedParsed.In(s.Loc).Format("2006 01.02 15*:*04"))
		}

		x.Send(str1, &tb.SendOptions{ParseMode: tb.ModeMarkdown})

	}

	return
}

func (s *Service) TgTestNotify(x tb.Context) (errReturn error) {
	content := &NotifyContent{
		Type:    "",
		Title:   "Название стрима",
		VideoID: "dQw4w9WgXcQ",
		Time:    "",
		TimePub: nil,
	}

	// Live
	content.Type = util.Live
	s.SendNotify(content)
	time.Sleep(time.Millisecond * 500)

	// Live Twitch
	content.Type = util.LiveTwitch
	s.SendNotify(content)
	time.Sleep(time.Millisecond * 500)

	// Upcoming
	content.Type = util.Upcoming
	tm := time.Now().Add(time.Minute * 5)
	content.TimePub = &tm
	s.SendNotify(content)
	content.Time = ""
	time.Sleep(time.Millisecond * 500)

	// Upcoming 2
	content.Type = util.Upcoming
	tm = time.Now().Add(time.Hour * 24).Add(time.Hour).Add(time.Minute * 5)
	content.TimePub = &tm
	s.SendNotify(content)
	content.Time = ""
	time.Sleep(time.Millisecond * 500)

	// LiveGo
	content.Type = util.LiveGo
	s.SendNotify(content)
	time.Sleep(time.Millisecond * 500)

	// End
	content.Type = util.End
	s.SendNotify(content)
	time.Sleep(time.Millisecond * 500)

	// End404
	content.Type = util.End404
	s.SendNotify(content)

	x.Send("Готово.")
	return
}

func (s *Service) TgTypeOfVideo(x tb.Context) (errReturn error) {
	if x.Text() == "/type_of_vid" {
		x.Send(s.Bot.Text(x, "type_of_vid_empty"))
		return
	}

	url := strings.Replace(x.Text(), "/type_of_vid ", "", 1)
	url = strings.Replace(url, " ", "", -1)

	ok, err := regexp.MatchString(`^https:\/\/www\.youtube\.com\/watch\?v=`, url)
	if err != nil || !ok {
		x.Send(s.Bot.Text(x, "type_of_vid_404"))
		return
	}

	url = strings.Replace(url, "https://www.youtube.com/watch?v=", "", -1)
	typeVid, tmStart, err := util.TypeVideo(url, false)
	if err != nil {
		x.Send(eris.Wrap(err, "util.TypeVideo()"))
		return
	}

	textTime := "nil"
	if tmStart != nil {
		textTime = tmStart.In(s.Loc).String()
	}
	text := fmt.Sprintf("type: %v\ntime: %v", typeVid, textTime)
	x.Send(text)
	return
}

func (s *Service) TgTwitchClient(x tb.Context) (errReturn error) {
	if x.Text() == "/set_twitch_client" {
		x.Send(s.Bot.Text(x, "set_twitch_client_empty"))
		return
	}

	strs := strings.Split(strings.Replace(x.Text(), "/set_twitch_client ", "", 1), " ")
	if len(strs) != 2 {
		x.Send(s.Bot.Text(x, "set_twitch_client_err_len"))
		return
	}

	err := s.YouTubeTwitch.Twitch.SetClient(strs[0], strs[1])
	if err != nil {
		x.Send(eris.Wrap(err, "Ошибка"))
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

func (s *Service) TgTwitchAuthURL(x tb.Context) (errReturn error) {
	if !s.YouTubeTwitch.Twitch.ClientOK() {
		x.Send(s.Bot.Text(x, "twitch_auth_url_err"))
		return
	}

	url, err := s.YouTubeTwitch.Twitch.GetAuthURL()
	if err != nil {
		x.Send(eris.Wrap(err, "Ошибка"))
		return
	}

	x.Send(url, tb.NoPreview)
	x.Send(s.Bot.Text(x, "twitch_auth_url"))
	return
}

func (s *Service) TgTwitchAuth(x tb.Context) (errReturn error) {
	if !s.YouTubeTwitch.Twitch.ClientOK() {
		x.Send(s.Bot.Text(x, "twitch_auth_url_err"))
		return
	}

	if x.Text() == "/twitch_auth" {
		x.Send(s.Bot.Text(x, "twitch_auth_empty"))
		return
	}

	code := strings.NewReplacer("/twitch_auth", "", " ", "", "http://localhost/?code=", "", "&scope=user%3Aread%3Aemail&state=some-state", "").Replace(x.Text())
	err := s.YouTubeTwitch.Twitch.SetCode(code)
	if err != nil {
		x.Send(eris.Wrap(err, "Ошибка"))
		return
	}

	x.Send(s.Bot.Text(x, "done"))
	return
}

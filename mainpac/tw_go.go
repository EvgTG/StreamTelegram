package mainpac

import (
	"github.com/rotisserie/eris"
	"streamtg/go-log"
	"streamtg/util"
	"time"
)

func (s *Service) TwitchCheck() {
	if s.YouTubeTwitch.TwitchNick == "" || s.YouTubeTwitch.TwitchNick == "0" {
		return
	}

	// получение стрима
	stream, err := s.YouTubeTwitch.Twitch.GetStream(s.YouTubeTwitch.TwitchNick)
	if err != nil {
		log.Error(eris.Wrap(err, "TwitchCheck - s.YouTubeTwitch.Twitch.GetStream(nick)"))
		return
	}

	if stream == nil {
		return
	}

	type item struct {
		id, title string
		isLive    bool
		timePub   *time.Time
	}

	// проверка на новинку
	repetition, err := s.MiniDB.CheckTwitchVideo(stream.ID)
	if err != nil {
		log.Error(eris.Wrap(err, "TwitchCheck - MiniDB.CheckTwitchVideo(videoID)"))
		return
	}
	if repetition {
		return
	}

	itm := item{
		id:      stream.ID,
		title:   stream.Title,
		isLive:  true,
		timePub: &stream.StartedAt,
	}

	// обработка нового стрима
	log.Infof("TwitchCheck new item %v %v", itm.id, itm.title)

	content := &NotifyContent{
		Type:       util.LiveTwitch,
		TwitchNick: s.YouTubeTwitch.TwitchNick,
		Title:      itm.title,
		VideoID:    itm.id,
		Time:       "",
		TimePub:    itm.timePub,
	}

	s.SendNotify(content)
	go s.GoEndWaitTwitch(content)
}

func (s *Service) GoEndWaitTwitch(content *NotifyContent) {
	log.Debug("GoEndWait", content.VideoID)

	{
		ok := false
		for _, channel := range s.Bot.NotifyList {
			if channel.EndOfStream {
				ok = true
				break
			}
		}
		if !ok {
			return
		}
	}

	for {
		time.Sleep(time.Minute * 7)

		stream, err := s.YouTubeTwitch.Twitch.GetStream(s.YouTubeTwitch.TwitchNick)
		if err != nil {
			log.Error(eris.Wrap(err, "TwitchCheck - s.YouTubeTwitch.Twitch.GetStream(nick)"))
			return
		}

		if stream == nil || stream.ID != content.VideoID {
			content.Type = util.EndTwitch
			s.SendNotify(content)
		}
	}

}

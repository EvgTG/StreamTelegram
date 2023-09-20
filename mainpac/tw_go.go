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

	// получение rss
	feed, err := s.YouTubeTwitch.Parser.ParseURL("https://twitchrss.appspot.com/vod/" + s.YouTubeTwitch.TwitchNick)
	if err != nil {
		log.Error(eris.Wrap(err, "TwitchCheck - ParseURL()"))
		return
	}
	util.ClearFeed(feed)
	s.YouTubeTwitch.LastRSS_TW = feed

	type item struct {
		id, title string
		isLive    bool
		timePub   *time.Time
	}

	// проверка на новые видео
	items := make([]item, 0, 1)
	for _, value := range feed.Items {
		repetition, err := s.MiniDB.CheckTwitchVideo(value.GUID)
		if err != nil {
			log.Error(eris.Wrap(err, "TwitchCheck - MiniDB.CheckTwitchVideo(videoID)"))
			return
		}

		if repetition {
			continue
		}
		items = append(items, item{
			id:      value.GUID,
			title:   value.Title,
			isLive:  util.IsTwitchLiveItem(value),
			timePub: value.PublishedParsed,
		})
	}

	if len(items) == 0 {
		return
	}

	// обработка новых видео
	var typeVideo string
	for _, itm := range items {
		log.Infof("TwitchCheck new item %v %v", itm.id, itm.title)

		if itm.isLive {
			typeVideo = util.LiveTwitch
		} else {
			typeVideo = util.ArchiveTwitch
		}

		content := &NotifyContent{
			Type:       typeVideo,
			TwitchNick: s.YouTubeTwitch.TwitchNick,
			Title:      itm.title,
			VideoID:    itm.id,
			Time:       "",
			TimePub:    itm.timePub,
		}

		switch typeVideo {
		case util.LiveTwitch:
			s.SendNotify(content)
			go s.GoEndWaitTwitch(content)
		}

		time.Sleep(time.Second * 1)
	}

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

		feed, err := s.YouTubeTwitch.Parser.ParseURL("https://twitchrss.appspot.com/vod/" + s.YouTubeTwitch.TwitchNick)
		if err != nil {
			log.Error(eris.Wrap(err, "GoEndWaitTwitch - ParseURL()"))
			return
		}
		util.ClearFeed(feed)

		br := false
		ok := false
		for _, item := range feed.Items {
			if item.GUID != content.VideoID {
				continue
			}

			ok = true
			if !util.IsTwitchLiveItem(item) {
				content.Type = util.EndTwitch
				s.SendNotify(content)
				br = true
			}
		}
		if br {
			break
		}
		if !ok {
			content.Type = util.EndTwitch404
			s.SendNotify(content)
			break
		}
	}

}

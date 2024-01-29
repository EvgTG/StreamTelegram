package mainpac

import (
	"strings"
	"time"

	"streamtg/go-log"
	"streamtg/util"

	"github.com/rotisserie/eris"
	tb "gopkg.in/telebot.v3"
)

func (s *Service) GoYouTubeTwitch() {
	for {
		s.YouTubeCheck()
		s.TwitchCheck()

		s.YouTubeTwitch.NumberIterations++
		s.YouTubeTwitch.LastTime = time.Now()
		time.Sleep(time.Minute * time.Duration(s.YouTubeTwitch.CycleDurationMinutes))

		// Pause
		s.YouTubeTwitch.PauseMutex.Lock()
		if s.YouTubeTwitch.Pause == 1 {
			s.YouTubeTwitch.Pause = 2
			s.YouTubeTwitch.PauseMutex.Unlock()
			<-s.YouTubeTwitch.PauseWaitChannel
		} else {
			s.YouTubeTwitch.PauseMutex.Unlock()
		}
	}
}

func (s *Service) YouTubeCheck() {
	if s.YouTubeTwitch.ChannelID == "" || s.YouTubeTwitch.ChannelID == "UCBR8-60-B28hp2BmDPdntcQ" {
		return
	}

	// получение rss
	feed, err := s.YouTubeTwitch.Parser.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + s.YouTubeTwitch.ChannelID)
	if err != nil {
		log.Error(eris.Wrap(err, "YouTubeCheck - ParseURL()"))
		return
	}
	s.YouTubeTwitch.LastRSS_YT = feed

	type item struct {
		id, title string
	}

	// проверка на новые видео
	items := make([]item, 0, 1)
	for _, value := range feed.Items {
		if !strings.Contains(value.GUID, "yt:video:") {
			continue
		}

		videoID := strings.Replace(value.GUID, "yt:video:", "", 1)
		repetition, err := s.MiniDB.CheckVideo(videoID)
		if err != nil {
			log.Error(eris.Wrap(err, "YouTubeCheck - MiniDB.CheckVideo(videoID)"))
			return
		}

		if repetition {
			continue
		}
		items = append(items, item{id: videoID, title: value.Title})
	}

	if len(items) == 0 {
		return
	}

	// обработка новых видео
	for _, itm := range items {
		log.Infof("YouTubeCheck new item %v %v", itm.id, itm.title)

		typeVideo, timePub, err := util.TypeVideo(itm.id, s.YouTubeTwitch.DebugLevel())
		if err != nil {
			log.Error(eris.Wrap(err, "YouTubeCheck util.TypeVideo()"))
			time.Sleep(time.Second * 10)
			continue
		}
		log.Infof("YouTubeCheck %v %v", itm.id, typeVideo)

		content := &NotifyContent{
			Type:    typeVideo,
			Title:   itm.title,
			VideoID: itm.id,
			Time:    "",
			TimePub: timePub,
		}

		switch typeVideo {
		case util.Live:
			fallthrough
		case util.Wait:
			s.SendNotify(content)
			go s.GoEndWait(content)
		case util.Upcoming:
			s.SendNotify(content)
			go s.GoStartWait(content)
		}

		time.Sleep(time.Second * 10)
	}

}

func (s *Service) GoStartWait(content *NotifyContent) {
	log.Debug("GoStartWait", content.VideoID)

	time.Sleep(time.Second*time.Duration(content.TimePub.Unix()-time.Now().Unix()) + 35)
	for i := 0; i < 30; i++ {
		typeVideo, _, err := util.TypeVideo(content.VideoID, s.YouTubeTwitch.DebugLevel())
		log.Debug("GoStartWait for", content.VideoID, typeVideo)
		if err != nil {
			log.Error(eris.Wrap(err, "GoStartWait util.TypeVideo()"))
			time.Sleep(time.Minute * 2)
			continue
		}

		if typeVideo == util.Live {
			content.Type = util.LiveGo
			s.SendNotify(content)
			go s.GoEndWait(content)
			break
		}

		if !(typeVideo == util.Upcoming || typeVideo == util.Wait) {
			break
		}

		time.Sleep(time.Minute * 2)
	}
}

func (s *Service) GoEndWait(content *NotifyContent) {
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

		typeVideo, _, err := util.TypeVideo(content.VideoID, s.YouTubeTwitch.DebugLevel())
		if err != nil {
			log.Error(eris.Wrap(err, "GoEndWait util.TypeVideo()"))
			continue
		}

		if typeVideo == util.End {
			content.Type = util.End
			s.SendNotify(content)
			break
		}

		if typeVideo == util.Err404 {
			content.Type = util.End404
			s.SendNotify(content)
			break
		}

		if typeVideo != util.Live {
			break
		}
	}
}

type NotifyContent struct {
	Type       string
	TwitchNick string
	Title      string
	VideoID    string
	Time       string
	TimePub    *time.Time
}

func (s *Service) SendNotify(content *NotifyContent) {
	log.Debug("SendNotify", content.VideoID, content.Type)

	for _, channel := range s.Bot.NotifyList {
		switch content.Type {
		case util.Live:
			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "live", content))
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify Live"))
			}
		case util.Upcoming:
			now := time.Now()
			content.Time = ""

			bl := false
			for _, locStr := range s.YouTubeTwitch.Locs {
				loc, err := time.LoadLocation(locStr)
				if err != nil {
					log.Error(eris.Wrap(err, "SendNotify time.LoadLocation(locStr)"))
					continue
				}

				tmFormat := s.YouTubeTwitch.TimeFormat
				if content.TimePub.In(loc).YearDay() == now.In(loc).YearDay() && content.TimePub.In(loc).Year() == now.In(loc).Year() {
					tmFormat = strings.Replace(tmFormat, "2 Jan ", "", 1)
				}
				tm := content.TimePub.In(loc).Format(tmFormat)
				tm = util.MonthReplacer.Replace(tm)
				tm = util.CityReplacer.Replace(tm)

				if bl {
					content.Time += "\n" + tm
				} else {
					content.Time += tm
					bl = true
				}
			}

			content.Time += "\nЧерез " + timeToStream((time.Second * time.Duration(content.TimePub.Unix()-time.Now().Unix())))

			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "upcoming", content))
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify Upcoming"))
			}
		case util.LiveGo:
			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "live_go", content))
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify LiveGo"))
			}
		case util.End, util.EndTwitch:
			if !channel.EndOfStream {
				continue
			}

			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "end", content))
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify End"))
			}
		case util.End404, util.EndTwitch404:
			if !channel.EndOfStream {
				continue
			}

			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "end404", content))
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify End404"))
			}
		case util.LiveTwitch:
			_, err := s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "live_twitch", content), tb.NoPreview)
			if err != nil {
				log.Error(eris.Wrap(err, "SendNotify LiveTwitch"))
			}
		}
	}
}

func (y *YouTubeTwitch) SetPause() {
	y.PauseMutex.Lock()
	defer y.PauseMutex.Unlock()

	switch y.Pause {
	case 0:
		y.Pause = 1
	case 1:
		y.Pause = 0
	case 2:
		y.Pause = 0
		y.PauseWaitChannel <- struct{}{}
	}
}

func (y *YouTubeTwitch) DebugLevel() bool {
	return strings.ToLower(y.LogLevel) == "debug"
}

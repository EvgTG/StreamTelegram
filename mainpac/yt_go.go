package mainpac

import (
	"github.com/rotisserie/eris"
	tb "gopkg.in/tucnak/telebot.v3"
	"streamtg/go-log"
	"streamtg/util"
	"strings"
	"time"
)

func (s *Service) GoYouTube() {
	for {
		s.YouTubeCheck()

		s.YouTube.NumberIterations++
		s.YouTube.LastTime = time.Now()
		time.Sleep(time.Minute * time.Duration(s.YouTube.CycleDurationMinutes))

		// Pause
		s.YouTube.PauseMutex.Lock()
		if s.YouTube.Pause == 1 {
			s.YouTube.Pause = 2
			s.YouTube.PauseMutex.Unlock()
			<-s.YouTube.PauseWaitChannel
		} else {
			s.YouTube.PauseMutex.Unlock()
		}
	}
}

func (s *Service) YouTubeCheck() {
	if s.YouTube.ChannelID == "" {
		return
	}

	// получение rss
	feed, err := s.YouTube.Parser.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + s.YouTube.ChannelID)
	if err != nil {
		log.Error(eris.Wrap(err, "YouTubeCheck - ParseURL()"))
		return
	}
	s.YouTube.LastRSS = *feed

	type item struct {
		id, title string
	}

	// проверка на новые видео
	items := []item{}
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
		typeVideo, timePub, err := util.TypeVideo("https://www.youtube.com/watch?v=" + itm.id)
		if err != nil {
			log.Error(eris.Wrap(err, "YouTubeCheck util.TypeVideo()"))
			break
		}

		content := &NotifyContent{
			Type:    typeVideo,
			Title:   itm.title,
			VideoID: itm.id,
			Time:    "",
			TimePub: timePub,
		}

		switch typeVideo {
		case util.Err404:
			log.Error(eris.New("YouTubeCheck util.Err404 util.TypeVideo()"))
		case util.Live:
			s.SendNotify(content)
			go s.GoEndWait(content)
		case util.Upcoming:
			s.SendNotify(content)
			go s.GoStartWait(content)
		}
	}

}

func (s *Service) GoStartWait(content *NotifyContent) {
	time.Sleep(time.Second*time.Duration(content.TimePub.Unix()-time.Now().Unix()) + 15)
	for i := 0; i < 30; i++ {
		typeVideo, _, err := util.TypeVideo("https://www.youtube.com/watch?v=" + content.VideoID)
		if err != nil {
			log.Error(eris.Wrap(err, "GoStartWait util.TypeVideo()"))
			break
		}

		if typeVideo == util.Live {
			content.Type = util.LiveGo
			s.SendNotify(content)
			go s.GoEndWait(content)
			break
		}

		if typeVideo != util.Upcoming {
			break
		}

		time.Sleep(time.Minute * 2)
	}
}

func (s *Service) GoEndWait(content *NotifyContent) {
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

		typeVideo, _, err := util.TypeVideo("https://www.youtube.com/watch?v=" + content.VideoID)
		if err != nil {
			log.Error(eris.Wrap(err, "GoEndWait util.TypeVideo()"))
			break
		}

		if typeVideo == util.End {
			content.Type = util.End
			s.SendNotify(content)
			break
		}

		if typeVideo != util.Live {
			break
		}
	}

}

type NotifyContent struct {
	Type    string
	Title   string
	VideoID string
	Time    string
	TimePub *time.Time
}

func (s *Service) SendNotify(content *NotifyContent) {
	for _, channel := range s.Bot.NotifyList {
		switch content.Type {
		case util.Live:
			s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "live", content))
		case util.Upcoming:
			if len(s.YouTube.Locs) > 1 {
				content.Time += ":\n"
			}

			bl := false
			for _, locStr := range s.YouTube.Locs {
				loc, err := time.LoadLocation(locStr)
				if err != nil {
					log.Error(eris.Wrap(err, "SendNotify time.LoadLocation(locStr)"))
					continue
				}

				tm := content.TimePub.In(loc).Format(s.YouTube.TimeFormat)
				if locStr == "Europe/Moscow" {
					tm = util.MonthReplacer.Replace(tm)
				}

				if bl {
					content.Time += "\n" + tm
					bl = true
				} else {
					content.Time += tm
				}
			}

			s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "upcoming", content))
		case util.LiveGo:
			s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "live_go", content))
		case util.End:
			s.Bot.Send(&tb.User{ID: channel.ID}, s.Bot.TextLocale("ru", "end", content))
		}
	}
}

func (s *Service) YouTubePause() {
	s.YouTube.PauseMutex.Lock()
	defer s.YouTube.PauseMutex.Unlock()

	switch s.YouTube.Pause {
	case 0:
		s.YouTube.Pause = 1
	case 1:
		s.YouTube.Pause = 0
	case 2:
		s.YouTube.Pause = 0
		s.YouTube.PauseWaitChannel <- struct{}{}
	}
}

package mainpac

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func (s *Service) StartYT() {
	fp := gofeed.NewParser()
	loc, err := time.LoadLocation("Europe/Moscow")
	s.FatalTG("StartYT - time.LoadLocation()", err)

	for {
		idsForCheck := []string{}
		feed, _ := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + s.yt.channelID)
		for _, value := range feed.Items {
			if !strings.Contains(value.GUID, "yt:video:") {
				continue
			}
			vID := strings.Replace(value.GUID, "yt:video:", "", 1)
			bl, err := s.db.Check(vID)
			s.FatalTG("StartYT - db.Check()", err)
			if bl {
				continue
			}
			idsForCheck = append(idsForCheck, vID)
		}

		if len(idsForCheck) != 0 {
			video := s.yt.yts.Videos.List([]string{"snippet", "liveStreamingDetails"})
			video.Id(strings.Join(idsForCheck, ","))
			videoRes, err := video.Do()
			s.FatalTG("StartYT - youtubeService.Videos.List.Do()", err)

			for _, value := range videoRes.Items {
				if value.Snippet.LiveBroadcastContent != "live" && value.Snippet.LiveBroadcastContent != "upcoming" {
					continue
				}

				switch value.Snippet.LiveBroadcastContent {
				case "live":
					text := fmt.Sprintf("%v\n\nyoutube.com/watch?v=%v", value.Snippet.Title, value.Id)
					s.tg.SendNotification(text)
				case "upcoming":
					layout := "2006-01-02T15:04:05Z"
					t, err := time.Parse(layout, value.LiveStreamingDetails.ScheduledStartTime)
					s.FatalTG("StartYT - time.Parse()", err)
					t = t.In(loc)
					text := fmt.Sprintf("%v\n\nЗапланировано на %v по Мск\nyoutube.com/watch?v=%v", value.Snippet.Title, t.Format("2 Jan 15:04"), value.Id)
					s.tg.SendNotification(text)
				}
			}
		}

		s.tg.numberIterations++
		time.Sleep(time.Minute * 5)
		if s.yt.stop == 1 {
			s.yt.stop = 2
			<-s.yt.stopch
		}
	}
}

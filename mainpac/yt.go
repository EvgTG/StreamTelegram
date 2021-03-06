package mainpac

import (
	"StreamTelegram/go-log"
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func (s *Service) StartYT() {
	fp := gofeed.NewParser()

	for {
		idsForCheck := []string{}
		feed, err := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + s.yt.channelID)
		if err != nil {
			s.tg.SendLog(fmt.Sprintf("ERR StartYT - fp.ParseURL(): %v", err.Error()))
			s.cycleActions()
			continue
		}
		s.yt.lastRSS = *feed
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
			if err != nil {
				s.tg.SendLog(fmt.Sprintf("ERR StartYT - youtubeService.Videos.List.Do(): %v", err.Error()))
				s.cycleActions()
				continue
			}

			for _, value := range videoRes.Items {
				if value.Snippet.LiveBroadcastContent != "live" && value.Snippet.LiveBroadcastContent != "upcoming" {
					continue
				}

				log.Debug("notif.go ", time.Now(), " Title= ", value.Snippet.Title)
				switch value.Snippet.LiveBroadcastContent {
				case "live":
					text := fmt.Sprintf(s.yt.texts["live"], value.Snippet.Title, value.Id)
					s.tg.SendNotification(text)
				case "upcoming":
					t, err := time.Parse("2006-01-02T15:04:05Z", value.LiveStreamingDetails.ScheduledStartTime)
					s.FatalTG("StartYT - time.Parse()", err)
					t = t.In(s.loc)
					text := fmt.Sprintf(s.yt.texts["upcoming"], value.Snippet.Title, t.Format("2 Jan 15:04 MST"), value.Id)
					s.tg.SendNotification(text)

					//уведомление о начале планируемого стрима
					go func(tmSec int64, id string) {
						time.Sleep(time.Second*time.Duration(tmSec) + 15)
						for exit := 0; exit < 60; exit++ {
							vid := s.yt.yts.Videos.List([]string{"snippet", "liveStreamingDetails"})
							vid.Id(id)
							vidRes, err := vid.Do()
							if err != nil {
								s.tg.SendLog("StartYT - waiting for stream, youtubeService.Videos.List.Do()")
								time.Sleep(time.Second * 30)
								continue
							}
							if len(vidRes.Items) == 1 {
								if vidRes.Items[0].Snippet.LiveBroadcastContent == "live" {
									s.tg.SendNotification(s.yt.texts["upcoming_go"])
									break
								}
							}
							time.Sleep(time.Second * 30)
						}
					}(t.Unix()-time.Now().Unix(), value.Id)
				}
			}
		}

		s.cycleActions()
	}
}

func (s *Service) cycleActions() {
	s.tg.numberIterations++
	s.yt.lastTime = time.Now()
	time.Sleep(s.yt.CycleTime)
	if s.yt.stop == 1 {
		s.yt.stop = 2
		<-s.yt.stopch
	}
}

package main

import (
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"StreamTelegram/tgbot"
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"strings"
	"time"
)

//TODO заменить фаталы на одну функцию
//TODO закинуть в цикл
func start(db *model.Model, tg *tgbot.TGBot, chID, ytAPIkey string) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(ytAPIkey))
	if err != nil {
		log.Fatal(err)
	}
	fp := gofeed.NewParser()
	//TODO перевести в конфиг
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}

	idsForCheck := []string{}
	feed, _ := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=" + chID)
	for _, value := range feed.Items {
		if !strings.Contains(value.GUID, "yt:video:") {
			continue
		}
		vID := strings.Replace(value.GUID, "yt:video:", "", 1)
		bl, err := db.Check(vID)
		if err != nil {
			log.Fatal(err)
		}
		if bl {
			continue
		}
		idsForCheck = append(idsForCheck, vID)
	}

	video := youtubeService.Videos.List([]string{"snippet", "liveStreamingDetails"})
	video.Id(strings.Join(idsForCheck, ","))
	videoRes, err := video.Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range videoRes.Items {
		if value.Snippet.LiveBroadcastContent != "live" && value.Snippet.LiveBroadcastContent != "upcoming" {
			continue
		}

		switch value.Snippet.LiveBroadcastContent {
		case "live":
			text := fmt.Sprintf("%v\n\nyoutube.com/watch?v=%v", value.Snippet.Title, value.Id)
			tg.SendNotification(text)
		case "upcoming":
			layout := "2006-01-02T15:04:05Z"
			t, err := time.Parse(layout, value.LiveStreamingDetails.ScheduledStartTime)
			if err != nil {
				log.Fatal(err)
			}
			t = t.In(loc)
			text := fmt.Sprintf("%v\n\nЗапланировано на %v по Мск\nyoutube.com/watch?v=%v", value.Snippet.Title, t.Format("2 Jan 15:04"), value.Id)
			tg.SendNotification(text)
		}
	}

	//fmt.Println(feed.Title)

	/*ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(ytAPIkey))
	if err != nil {
		log.Fatal(err)
	}*/

	/*ts := youtubeService.Channels.List([]string{"topicDetails"})
	ts.Id(chID)
	tss, err := ts.Do()
	for i, v := range tss.Items {
		log.Infof("%v - %s", i, v.Snippet.Title)
	}*/

	/*for {
		liveStReg := youtubeService.Search.List([]string{})
		liveStReg.ChannelId(chID)
		liveStReg.EventType("live")
		liveStReg.Type("video")
		liveSt, err := liveStReg.Do()
		if err != nil {
			log.Fatal(err)
		}

		if liveSt.Items != nil {
			bl, err := db.Check(liveSt.Items[0].Id.VideoId)
			if err != nil {
				log.Fatal(err)
			}

			if !bl {
				liveStSnippetReg := youtubeService.Search.List([]string{"snippet"})
				liveStSnippetReg.ChannelId(chID)
				liveStSnippetReg.EventType("live")
				liveStSnippetReg.Type("video")
				liveStSnippet, err := liveStSnippetReg.Do()
				if err != nil {
					log.Fatal(err)
				} else if liveStSnippet.Items == nil {
					log.Fatal("empty liveStSnippet.Items")
				}

				title := liveStSnippet.Items[0].Snippet.Title
				text := title + "\n\nyoutube.com/watch?v=" + liveStSnippet.Items[0].Id.VideoId
				tg.SendNotification(text)
			}
		}

		time.Sleep(time.Second * 15)
	}*/
}

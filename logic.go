package main

import (
	"StreamTelegram/model"
	"StreamTelegram/tgbot"
)

func start(db *model.Model, tg *tgbot.TGBot, chID, ytAPIkey string) {

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

package util

import (
	"github.com/mmcdole/gofeed"
	"regexp"
)

const (
	LiveTwitch    = "livetw"
	EndTwitch     = "endtw"
	EndTwitch404  = "end404tw"
	ArchiveTwitch = "archivetw"
)

var titleLIVE = regexp.MustCompile(" - LIVE$")
var guidLIVE = regexp.MustCompile("_live$")

func ClearFeed(feed *gofeed.Feed) {
	for _, item := range feed.Items {
		if titleLIVE.MatchString(item.Title) {
			item.Title = item.Title[:len(item.Title)-7]
		}
		if guidLIVE.MatchString(item.GUID) {
			item.GUID = item.GUID[:len(item.GUID)-5]
		}
	}
}

func IsTwitchLive(feed *gofeed.Feed) bool {
	for _, item := range feed.Items {
		if len(item.Categories) < 1 {
			continue
		}

		for _, category := range item.Categories {
			if category == "live" {
				return true
			}
		}
	}
	return false
}

func IsTwitchLiveItem(item *gofeed.Item) bool {
	for _, category := range item.Categories {
		if category == "live" {
			return true
		}
	}
	return false
}

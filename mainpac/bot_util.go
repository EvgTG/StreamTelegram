package mainpac

import (
	"fmt"
	"strings"
	"time"

	"streamtg/util"

	tb "gopkg.in/telebot.v3"
)

func (bot *Bot) isNotAdmin(x tb.Context) bool {
	if x.Chat().ID >= 0 && util.IntInSlice(bot.AdminList, x.Sender().ID) {
		return false
	}
	return true
}

func (bot *Bot) sendToSlice(slice []int64, mesText string) {
	for _, chatID := range slice {
		bot.Send(&tb.User{ID: chatID}, mesText)
	}
}

// 4d7h6m34s
func (bot *Bot) uptimeString(timestamp time.Time) string {
	uptime := time.Since(timestamp).Round(time.Second)
	hours, hoursStr := 0, ""
	for uptime.Hours() >= 24 {
		uptime -= time.Hour * 24
		hours++
	}
	if hours > 0 {
		hoursStr = fmt.Sprintf("%vd", hours)
	}
	return hoursStr + uptime.String()
}

var durReplacer = strings.NewReplacer("d", "д ", "h", "ч ", "m", "м", "0s", "")

func timeToStream(tm time.Duration) (str string) {
	tm = tm.Round(time.Minute)

	hours, hoursStr := 0, ""
	for tm.Hours() >= 24 {
		tm -= time.Hour * 24
		hours++
	}
	if hours > 0 {
		hoursStr = fmt.Sprintf("%vd", hours)
	}

	return durReplacer.Replace(hoursStr + tm.String())
}

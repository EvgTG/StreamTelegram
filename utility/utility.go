package utility

import (
	"StreamTelegram/go-log"
	"StreamTelegram/tgbot"
	"fmt"
)

func Fatal(text string, err error) {
	if err != nil {
		log.Fatalf("ERR %v: %v", text, err)
	}
}

func FatalTG(text string, tg *tgbot.TGBot, err error) {
	if err != nil {
		text := fmt.Sprintf("ERR %v: %v", text, err)

		tg.SendLog(text)
		log.Fatal(text)
	}
}

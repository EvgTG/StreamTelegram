package mainpac

import (
	"StreamTelegram/go-log"
	"fmt"
)

func Fatal(text string, err error) {
	if err != nil {
		log.Fatalf("ERR %v: %v", text, err)
	}
}

func (s Service) FatalTG(text string, err error) {
	if err != nil {
		text := fmt.Sprintf("ERR %v: %v", text, err)

		s.tg.SendLog(text)
		log.Fatal(text)
	}
}

func GetTexts(language string) map[string]string {
	mp := make(map[string]string, 3)

	switch language {
	case "rus":
		mp["live"] = "%v\n\nyoutube.com/watch?v=%v"
		mp["upcoming"] = "%v\n\nЗапланировано на %v\nyoutube.com/watch?v=%v"
		mp["upcoming_go"] = "Стрим начался!"
	case "eng":
		fallthrough
	default:
		mp["live"] = "%v\n\nyoutube.com/watch?v=%v"
		mp["upcoming"] = "%v\n\nScheduled for %v\nyoutube.com/watch?v=%v"
		mp["upcoming_go"] = "Stream started!"
	}

	return mp
}

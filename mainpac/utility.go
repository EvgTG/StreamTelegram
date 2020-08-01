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

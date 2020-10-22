package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type QuestsService struct {
}

func (qs *QuestsService) start() {

	for {
		// log.Println("tick")

		time.Sleep(time.Second * 5)
	}
}

func (qs *QuestsService) isFbLinkValid(link string) bool {
	if !strings.Contains(link, "https://www.facebook.com/") {
		logTelegram(fmt.Sprintf("[quests.go - 21] link: %s", link))
		return false
	}

	if r, err := http.Get(link); err != nil {
		logTelegram("[quests.go - 29]" + err.Error())
		return false
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		if !strings.Contains(body, "Anote") ||
			!strings.Contains(body, "anonutopia.com") ||
			!strings.Contains(body, "follow") ||
			!strings.Contains(body, "nickname") {

			return false
		}
	}

	return true
}

func initQuestsService() *QuestsService {
	qs := &QuestsService{}
	go qs.start()
	return qs
}

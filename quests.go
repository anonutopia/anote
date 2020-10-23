package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
)

type QuestsService struct {
}

func (qs *QuestsService) start() {
	for {
		qs.checkQuests()
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

func (qs *QuestsService) createQuest(u *User, link string) {
	u.FbPostLink = link
	now := time.Now()
	u.LastFbQuest = &now
	if !u.SentAint {
		u.MinedAnotes += 10 * int(satInBtc)
		messageTelegram(tr(u.TelegramID, "fbQuestAnotesAdded"), int64(u.TelegramID))
		u.SentFbAnotes = true
	} else {
		u.SentFbAnotes = false
	}
	db.Save(u)
}

func (qs *QuestsService) isQuestAvailable(u *User) bool {
	if u.LastFbQuest != nil &&
		time.Since(*u.LastFbQuest) < time.Duration(time.Hour*24*7) {
		return false
	}

	return true
}

func (qs *QuestsService) checkQuests() {
	var users []*User

	first := time.Now().Add(time.Hour * -15)
	second := time.Now().Add(time.Hour * -12)

	db.Where("sent_aint = false AND last_fb_quest BETWEEN ? AND ?", first, second).Find(&users)

	for _, u := range users {
		if !u.SentAint {
			qs.sendAint(u)
			u.SentAint = true
			db.Save(u)
			messageTelegram(tr(u.TelegramID, "fbQuestAintSent"), int64(u.TelegramID))
		}
	}

	db.Where("sent_fb_anotes = false AND last_fb_quest BETWEEN ? AND ?", first, second).Find(&users)

	for _, u := range users {
		if u.SentAint {
			u.MinedAnotes += 10 * int(satInBtc)
			u.SentFbAnotes = true
			db.Save(u)
			messageTelegram(tr(u.TelegramID, "fbQuestAnotesAdded2"), int64(u.TelegramID))
		}
	}
}

func (qs *QuestsService) sendAint(u *User) {
	atr := &gowaves.AssetsTransferRequest{
		Amount:    1,
		AssetID:   conf.AintID,
		Fee:       100000,
		Recipient: u.Address,
		Sender:    conf.NodeAddress,
	}

	if _, err := wnc.AssetsTransfer(atr); err != nil {
		logTelegram("[quests.go - 107]" + err.Error())
	}
}

func initQuestsService() *QuestsService {
	qs := &QuestsService{}
	go qs.start()
	return qs
}

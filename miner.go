package main

import (
	"time"

	ui18n "github.com/unknwon/i18n"
)

// MinerMonitor represents Anote mining monitoring object
type MinerMonitor struct {
}

func (mm *MinerMonitor) checkMiners() {
	var users []*User
	db.Find(&users)
	for _, u := range users {
		if u.MiningActivated != nil && time.Since(*u.MiningActivated).Hours() >= float64(24) && !u.SentWarning {
			messageTelegram(ui18n.Tr(lang, "miningWarning"), int64(u.TelegramID))
			u.SentWarning = true
			u.Mining = false
			db.Save(&u)
		}
	}
}

func (mm *MinerMonitor) start() {
	for {
		mm.checkMiners()

		time.Sleep(time.Second * 5)
	}
}

func initMinerMonitor() {
	mm := &MinerMonitor{}
	go mm.start()
}

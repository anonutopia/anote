package main

import (
	"time"
)

// MinerMonitor represents Anote mining monitoring object
type MinerMonitor struct {
}

func (mm *MinerMonitor) checkMiners() {
	var users []*User
	db.Find(&users)
	for _, u := range users {
		if time.Since(*u.MintingActivated).Hours() >= float64(24) && !u.SentWarning {
			messageTelegram("Mining warning!", int64(u.TelegramID))
			u.SentWarning = true
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

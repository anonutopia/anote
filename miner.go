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
		if u.MiningActivated != nil && time.Since(*u.MiningActivated).Hours() >= float64(24) && !u.SentWarning {
			msg := tr(u.TelegramID, "miningWarning")
			msg += "\n\n"
			msg += tr(u.TelegramID, "purchaseHowto")
			messageTelegram(msg, int64(u.TelegramID))
			u.SentWarning = true
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				logTelegram(err.Error())
			}
		} else if time.Since(u.CreatedAt).Hours() >= float64(24) && !u.SentWarning {
			msg := tr(u.TelegramID, "miningWarningFirst")
			msg += "\n\n"
			msg += tr(u.TelegramID, "purchaseHowto")
			messageTelegram(msg, int64(u.TelegramID))
			u.SentWarning = true
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				logTelegram(err.Error())
			}
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

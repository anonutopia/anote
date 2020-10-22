package main

import (
	"strings"
	"time"
)

// MinerMonitor represents Anote mining monitoring object
type MinerMonitor struct {
}

func (mm *MinerMonitor) checkMiners() {
	var users []*User
	db.Find(&users)
	for _, u := range users {
		now := time.Now()
		if u.MiningActivated != nil &&
			u.MiningWarning != nil &&
			time.Since(*u.MiningActivated).Hours() >= float64(24) &&
			time.Since(*u.MiningWarning).Hours() >= float64(24) {

			msg := tr(u.TelegramID, "miningWarning")
			msg += "\n\n"
			msg += tr(u.TelegramID, "purchaseHowto")

			u.MiningWarning = &now
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				logTelegram("[mm.checkMiners - db.Save - 60] " + err.Error())
			}

			if u.team() >= 3 {
				messageTelegram(msg, int64(u.TelegramID))
			} else {
				minerMsg := strings.Replace(tr(u.TelegramID, "minerWarning"), "\\n", "\n", -1)
				messageTelegram(minerMsg, int64(u.TelegramID))

				go func(u *User) {
					time.Sleep(time.Minute * 5)
					messageTelegram(msg, int64(u.TelegramID))
				}(u)
			}
		} else if u.MiningActivated == nil &&
			(u.MiningWarning == nil || time.Since(*u.MiningWarning).Hours() >= float64(24)) &&
			time.Since(u.CreatedAt).Hours() >= float64(24) {

			u.MiningWarning = &now
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				logTelegram("[mm.checkMiners - db.Save - 67] " + err.Error())
			}

			msg := tr(u.TelegramID, "miningWarningFirst")
			msg += "\n\n"
			msg += tr(u.TelegramID, "purchaseHowto")
			messageTelegram(msg, int64(u.TelegramID))
		}
	}
}

func (mm *MinerMonitor) start() {
	for {
		mm.checkMiners()

		time.Sleep(time.Minute)
	}
}

func initMinerMonitor() {
	mm := &MinerMonitor{}
	go mm.start()
}

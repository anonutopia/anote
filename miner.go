package main

import (
	"log"
	"strings"
	"time"
)

// MinerMonitor represents Anote mining monitoring object
type MinerMonitor struct {
}

func (mm *MinerMonitor) checkMiners() {
	var users []*User
	counter := 1
	db.Find(&users)
	log.Println(len(users))
	for _, u := range users {
		var err error
		log.Println(counter)
		counter++
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
				logTelegram("[mm.checkMiners - db.Save - 30] " + err.Error())
			}

			if u.team() >= 3 {
				err = messageTelegram(msg, int64(u.TelegramID))
				if err != nil {
					db.Delete(u)
				}
			} else {
				minerMsg := strings.Replace(tr(u.TelegramID, "minerWarning"), "\\n", "\n", -1)
				err = messageTelegram(minerMsg, int64(u.TelegramID))
				if err != nil {
					db.Delete(u)
				}
				go func(u *User) {
					time.Sleep(time.Minute * 5)
					err = messageTelegram(msg, int64(u.TelegramID))
				}(u)
			}
		} else if u.MiningActivated == nil &&
			(u.MiningWarning == nil || time.Since(*u.MiningWarning).Hours() >= float64(24)) &&
			time.Since(u.CreatedAt).Hours() >= float64(24) {

			msg := tr(u.TelegramID, "miningWarningFirst")
			msg += "\n\n"
			msg += tr(u.TelegramID, "purchaseHowto")
			err = messageTelegram(msg, int64(u.TelegramID))
			if err != nil {
				db.Delete(u)
			} else {
				u.MiningWarning = &now
				u.Mining = false
				if err := db.Save(&u).Error; err != nil {
					logTelegram("[mm.checkMiners - db.Save - 30] " + err.Error())
				}
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

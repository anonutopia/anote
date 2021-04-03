package main

import (
	"log"
	"strings"
	"time"

	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UserManager struct {
}

func (um *UserManager) createUser(m *tb.Message) {
	code := randString(10)
	tNick := m.Sender.Username

	if len(tNick) == 0 {
		tNick = code
	}

	u := &User{
		TelegramID:        m.Sender.ID,
		Address:           code,
		Nickname:          tNick,
		Code:              code,
		AnoteRobotStarted: true,
	}

	r := &User{}

	db.FirstOrCreate(u, u)

	if err := db.Where("code = ?", m.Payload).First(r).Error; err != nil {
		db.FirstOrCreate(u, u)
		db.Where("nickname = ?", m.Payload).First(r)
	}

	if r.ID != 0 && r.ID != u.ID {
		u.Referral = r
		db.Save(u)
	}
}

func (um *UserManager) getUser(m *tb.Message) *User {
	u := &User{TelegramID: m.Sender.ID}
	db.First(u, u)

	u.checkMining()
	u.addMined()
	db.Save(u)

	return u
}

func (um *UserManager) checkNick(m *tb.Message) *User {
	user := um.getUser(m)
	if user.Nickname != m.Sender.Username {
		user.Nickname = m.Sender.Username
		if err := db.Save(user).Error; err != nil {
			log.Println(err)
		}
	}
	return user
}

func (um *UserManager) checkMiners() {
	var users []*User
	db.Find(&users)
	for _, u := range users {
		now := time.Now()
		if u.MiningActivated != nil &&
			u.MiningWarning != nil &&
			time.Since(*u.MiningActivated).Hours() >= float64(24) &&
			time.Since(*u.MiningWarning).Hours() >= float64(24) {

			msg := gotrans.T("miningWarning")
			msg += "\n\n"
			msg += gotrans.T("purchaseHowto")

			u.MiningWarning = &now
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				log.Println(err)
			}

			if u.AnoteRobotStarted {
				if u.team() >= 3 {
					messageTelegram(msg, u.TelegramID)
				} else {
					minerMsg := strings.Replace(gotrans.T("minerWarning"), "\\n", "\n", -1)
					messageTelegram(minerMsg, u.TelegramID)

					go func(u *User) {
						time.Sleep(time.Minute * 5)
						messageTelegram(msg, u.TelegramID)
					}(u)
				}
			}
		} else if u.MiningActivated == nil &&
			(u.MiningWarning == nil || time.Since(*u.MiningWarning).Hours() >= float64(24)) &&
			time.Since(u.CreatedAt).Hours() >= float64(24) {

			u.MiningWarning = &now
			u.Mining = false
			if err := db.Save(&u).Error; err != nil {
				log.Println(err)
			}

			if u.AnoteRobotStarted {
				msg := gotrans.T("miningWarningFirst")
				msg += "\n\n"
				msg += gotrans.T("purchaseHowto")
				messageTelegram(msg, u.TelegramID)
			}
		}
	}
}

func (um *UserManager) start() {
	for {
		um.checkMiners()

		time.Sleep(time.Minute)
	}
}

func initUserManager() *UserManager {
	um := &UserManager{}
	go um.start()
	return um
}

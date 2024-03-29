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
	u := um.getUser(m)
	code := randString(10)

	u.Code = &code
	u.AnoteRobotStarted = true

	if u.ID != 0 {
		db.Save(u)
		return
	}

	u.Nickname = &m.Sender.Username
	u.TelegramID = &m.Sender.ID
	u.MinedAnotes = int(SatInBTC)

	db.Create(u)

	r := &User{}

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
	u := &User{TelegramID: &m.Sender.ID}
	db.First(u, u)

	if u.ID == 0 {
		return u
	}

	u.checkMining()
	u.addMined()
	db.Save(u)

	return u
}

func (um *UserManager) checkNick(m *tb.Message) *User {
	user := um.getUser(m)
	if user.Nickname != nil && *user.Nickname != m.Sender.Username {
		user.Nickname = &m.Sender.Username
		if err := db.Save(user).Error; err != nil {
			log.Println(err)
			// logTelegram(err.Error())
		}
	}
	return user
}

func (um *UserManager) isPayloadRef(m *tb.Message) (bool, *User) {
	r := &User{}

	if err := db.Where("code = ?", m.Payload).First(r).Error; err != nil {
		if err := db.Where("nickname = ?", m.Payload).First(r).Error; err != nil {
			return false, nil
		}
	}

	if r.ID != 0 {
		return true, r
	}

	return false, nil
}

func (um *UserManager) isPayloadMe(m *tb.Message) (bool, *User) {
	u := &User{}

	if err := db.Where("temp_code = ?", m.Payload).First(u).Error; err != nil {
		return false, nil
	}

	if u.ID != 0 {
		return true, u
	}

	return false, nil
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
			if err := db.Save(&u).Error; err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			if u.AnoteRobotStarted {
				if u.team() >= 3 {
					messageTelegram(msg, *u.TelegramID)
				} else {
					minerMsg := strings.Replace(gotrans.T("minerWarning"), "\\n", "\n", -1)
					messageTelegram(minerMsg, *u.TelegramID)

					go func(u *User) {
						time.Sleep(time.Minute * 5)
						messageTelegram(msg, *u.TelegramID)
					}(u)
				}
			}
		} else if u.MiningActivated == nil &&
			(u.MiningWarning == nil ||
				time.Since(*u.MiningWarning).Hours() >= float64(24)) &&
			time.Since(u.CreatedAt).Hours() >= float64(24) {

			u.MiningWarning = &now

			if err := db.Save(&u).Error; err != nil {
				log.Println(err)
				logTelegram(err.Error())
			}

			msg := gotrans.T("miningWarningFirst")
			msg += "\n\n"
			msg += gotrans.T("purchaseHowto")

			if u.AnoteRobotStarted {
				messageTelegram(msg, *u.TelegramID)
			}
		}

		if u.Mining {
			if time.Since(*u.MiningActivated).Hours() >= float64(24) {
				u.addMined()
				u.Mining = false
				if err := db.Save(&u).Error; err != nil {
					log.Println(err)
					logTelegram(err.Error())
				}
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
	// go um.start()
	return um
}

package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type UserManager struct {
	Running bool
}

func (um *UserManager) createUser(m *tb.Message) {
	code := randString(10)
	tNick := m.Sender.Username

	if len(tNick) == 0 {
		tNick = code
	}

	u := &User{
		TelegramID: m.Sender.ID,
		Address:    code,
		Nickname:   tNick,
		Code:       code,
	}
	r := &User{}

	db.FirstOrCreate(u, u)
	db.First(r, m.Payload)

	if r.ID != 0 && r.ID != u.ID {
		u.Referral = r
		db.Save(u)
	}
}

func (um *UserManager) getUser(m *tb.Message) *User {
	u := &User{TelegramID: m.Sender.ID}
	db.First(u, u)
	return u
}

func (um *UserManager) checkNick(m *tb.Message) {
	user := um.getUser(m)
	if user.Nickname != m.Sender.Username {
		user.Nickname = m.Sender.Username
		if err := db.Save(user).Error; err != nil {
			log.Println(err)
		}
	}
}

func (um *UserManager) saveState() {
	log.Println("Saving state.")
	um.Running = false
}

func (um *UserManager) start() {
	um.Running = true
	go func() {
		for um.Running {
			time.Sleep(time.Second * 10)
		}
	}()
}

func initUserManager() *UserManager {
	um := &UserManager{}
	um.start()
	return um
}

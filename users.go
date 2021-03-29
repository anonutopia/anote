package main

import (
	"log"

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
		TelegramID: m.Sender.ID,
		Address:    code,
		Nickname:   tNick,
		Code:       code,
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

func initUserManager() *UserManager {
	um := &UserManager{}
	return um
}

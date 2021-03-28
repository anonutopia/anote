package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type UserManager struct {
	Running bool
	Users   map[uint]*User
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

	um.Users[u.ID] = u

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
	for _, u := range um.Users {
		db.Save(u)
	}
	um.Running = false
}

func (um *UserManager) loadState() {
	um.Users = make(map[uint]*User)
	var users []*User
	db.Find(&users)
	for _, u := range users {
		um.Users[u.ID] = u
		if len(u.Code) == 0 {
			u.Code = randString(10)
		}
		if len(u.TempCode) == 0 {
			u.TempCode = randString(10)
		}
		db.Save(u)
	}
}

func (um *UserManager) checkMining() {
	for _, u := range um.Users {
		u.checkMining()
		u.addMined()
	}
}

func (um *UserManager) start() {
	um.loadState()
	um.Running = true
	go func() {
		for um.Running {
			um.checkMining()
			time.Sleep(time.Second * 10)
		}
	}()
}

func initUserManager() *UserManager {
	um := &UserManager{}
	um.start()
	return um
}

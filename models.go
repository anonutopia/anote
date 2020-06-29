package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64 `sql:"type:int"`
	ValueStr string `sql:"type:string"`
}

// Transaction represents node's transaction
type Transaction struct {
	gorm.Model
	TxID      string `sql:"size:255"`
	Processed bool   `sql:"DEFAULT:false"`
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address          string `sql:"size:255;unique_index"`
	TelegramUsername string `sql:"size:255"`
	TelegramID       int    `sql:"unique_index"`
	ReferralID       uint
	Referral         *User
	MiningActivated  *time.Time
	MinedAnotes      uint
	SentWarning      bool `sql:"DEFAULT:false"`
	Mining           bool `sql:"DEFAULT:false"`
}

func (u *User) status() string {
	if len(u.Address) == 0 {
		return "Not Registered"
	} else if u.team() >= 5 {
		return "Miner"
	} else if u.Mining {
		return "Pioneer"
	} else {
		return "Dissident"
	}
}

func (u *User) isMiningStr() string {
	if u.Mining {
		return "yes"
	} else {
		return "no"
	}
}

func (u *User) miningPower() float64 {
	power := float64(0)

	if u.Mining {
		power += 0.2
	}

	if u.teamActive() > 0 {
		power += float64(u.teamActive()) * 0.05
	}

	return power
}

func (u *User) miningPowerStr() string {
	return fmt.Sprintf("%.2f A/h", u.miningPower())
}

func (u *User) team() int {
	count := 0
	db.Where(&User{ReferralID: u.ID}).Find(&User{}).Count(&count)
	return count
}

func (u *User) teamStr() string {
	return fmt.Sprintf("%d", u.team())
}

func (u *User) teamInactive() int {
	team := u.team()
	active := 0
	db.Where(&User{ReferralID: u.ID, Mining: true}).Find(&User{}).Count(&active)
	return team - active
}

func (u *User) teamActive() int {
	active := 0
	db.Where(&User{ReferralID: u.ID, Mining: true}).Find(&User{}).Count(&active)
	return active
}

func (u *User) teamInactiveStr() string {
	return fmt.Sprintf("%d", u.teamInactive())
}

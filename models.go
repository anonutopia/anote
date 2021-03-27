package main

import (
	"time"

	"gorm.io/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64 `sql:"type:int"`
	ValueStr string `sql:"type:string"`
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address         string `sql:"size:255;unique_index"`
	TelegramID      int    `sql:"unique_index"`
	ReferralID      uint
	Referral        *User
	MiningActivated *time.Time
	MinedAnotes     int
	Mining          bool `sql:"DEFAULT:false"`
	LastWithdraw    *time.Time
	Language        string `sql:"size:255;"`
	MiningWarning   *time.Time
	Nickname        string `sql:"size:255;unique_index"`
}

func (u *User) getAddress() string {
	if len(u.Address) > 0 {
		return u.Address
	}

	return "no address"
}

package main

import (
	"time"

	"github.com/anonutopia/gowaves"
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

func (u *User) miningPower() float64 {
	power := float64(0)

	power += 0.02

	if u.teamActive() > 0 {
		power += float64(u.teamActive()) * 0.005
	}

	if u.teamActive() >= 3 {
		power *= 10
	}

	if len(u.Address) > 0 {
		avr, err := gowaves.WNC.AddressValidate(u.Address)
		if err == nil && avr.Valid {
			abr, err := gowaves.WNC.AssetsBalance(u.Address, AINTId)
			if err == nil {
				power += float64(abr.Balance) / float64(SatInBTC)
			}
		}
	}

	return power
}

func (u *User) team() int64 {
	var users []*User
	count := int64(0)
	db.Where(&User{ReferralID: u.ID}).Find(&users).Count(&count)
	return count
}

func (u *User) teamInactive() int64 {
	return u.team() - u.teamActive()
}

func (u *User) teamActive() int64 {
	var users []*User
	active := int64(0)
	db.Where("referral_id = ? AND mining_activated >= ?", u.ID, time.Now().Add(-24*time.Hour).Format("2006-01-02")).Find(&users).Count(&active)
	return active
}

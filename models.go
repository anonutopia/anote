package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/bykovme/gotrans"
	"gorm.io/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;uniqueIndex"`
	ValueInt uint64 `sql:"type:int"`
	ValueStr string `sql:"type:string"`
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address         string `gorm:"size:255;uniqueIndex"`
	TelegramID      int    `gorm:"uniqueIndex"`
	ReferralID      uint
	Referral        *User
	MiningActivated *time.Time
	MinedAnotes     int
	Mining          bool `sql:"DEFAULT:false"`
	LastWithdraw    *time.Time
	Language        string `sql:"size:255;"`
	MiningWarning   *time.Time
	Nickname        string `gorm:"size:255;uniqueIndex"`
	Code            string `gorm:"size:255;uniqueIndex"`
	UpdatedAddress  bool   `sql:"DEFAULT:false"`
	TempCode        string `gorm:"size:255;uniqueIndex"`
}

func (u *User) getAddress() string {
	if len(u.Address) > 0 && u.Address != u.Code {
		return u.Address
	}

	return "no wallet address"
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

func (u *User) status() string {
	if len(u.Address) == 0 {
		return "Guest"
	} else if u.team() >= 3 {
		return "Miner"
	} else if u.Mining {
		return "Pioneer"
	} else {
		return "Dissident"
	}
}

func (u *User) isMiningStr() string {
	if u.Mining {
		return gotrans.T("yes")
	}

	return gotrans.T("no")
}

func (u *User) miningPowerStr() string {
	return fmt.Sprintf("%.5f A/h", u.miningPower())
}

func (u *User) withdraw() {
	log.Println("withdraw")
}

func (u *User) mine() {
	log.Println("mine")
}

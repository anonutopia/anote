package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64
	ValueStr string
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
	Address         string `sql:"size:255"`
	TelegramID      int    `sql:"unique_index"`
	ReferralID      uint
	Referral        *User
	MiningActivated *time.Time
	LastStatus      *time.Time
	MinedAnotes     int
	Mining          bool `sql:"DEFAULT:false"`
	LastWithdraw    *time.Time
	Language        string `sql:"size:255;"`
	MiningWarning   *time.Time
	Nickname        string `sql:"size:255;unique_index"`
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
		return tr(u.TelegramID, "yes")
	}

	return tr(u.TelegramID, "no")
}

func (u *User) miningPower() float64 {
	multipliedByTen := false
	power := float64(0)
	limit := uint64(0)

	power += 0.02

	if u.teamActive() > 0 {
		power += float64(u.teamActive()) * 0.005
	}

	if u.teamActive() >= 3 {
		power *= 10
		multipliedByTen = true
	}

	if len(u.Address) > 0 {
		abr, _ := wnc.AssetsBalance(u.Address, conf.TokenID)
		hma, _ := pc.HashMultiplierAmount()
		limit = uint64(abr.Balance) / hma
	}

	if limit > 4 {
		limit = 4
	}

	for i := 1; uint64(i) < limit; i++ {
		if multipliedByTen {
			power *= 2
		} else {
			power *= 10
			multipliedByTen = true
		}
	}

	return power
}

func (u *User) miningPowerStr() string {
	return fmt.Sprintf("%.3f A/h", u.miningPower())
}

func (u *User) team() int {
	var users []*User
	count := 0
	db.Where(&User{ReferralID: u.ID}).Find(&users).Count(&count)
	return count
}

func (u *User) teamInactive() int {
	return u.team() - u.teamActive()
}

func (u *User) teamActive() int {
	var users []*User
	active := 0
	db.Where("referral_id = ? AND mining_activated IS NOT NULL", u.ID).Find(&users).Count(&active)
	return active
}

func hashingPower() {
	var hp float64
	var users []*User
	db.Find(&users)
	for _, u := range users {
		if u.Mining {
			hp += u.miningPower()
		}
	}
	log.Printf("power: %.2f", hp)
}

// Shout models is used for storing shouts and auctions for ads
type Shout struct {
	gorm.Model
	Message   string
	Link      string `sql:"size:255"`
	Price     int
	OwnerID   uint
	Owner     *User
	ChatID    int
	Finished  bool `sql:"DEFAULT:false"`
	Published bool `sql:"DEFAULT:false"`
}

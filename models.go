package main

import (
	"fmt"
	"log"
	"math"
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
	MiningActivated *time.Time `gorm:"index"`
	MinedAnotes     int
	Mining          bool `sql:"DEFAULT:false"`
	LastWithdraw    *time.Time
	Language        string `sql:"size:255;"`
	MiningWarning   *time.Time
	Nickname        string `gorm:"size:255;uniqueIndex"`
	Code            string `gorm:"size:255;uniqueIndex"`
	UpdatedAddress  bool   `sql:"DEFAULT:false"`
	TempCode        string `gorm:"size:255;uniqueIndex"`
	LastAdd         *time.Time
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
				power += u.miningPowerAint(float64(abr.Balance) / float64(SatInBTC))
			}
		}
	}

	return power
}

func (u *User) miningPowerAint(amount float64) float64 {
	power := float64(0)
	factor := float64(1)

	for amount > 1.0 {
		power += factor
		if factor > 0.01 {
			factor = factor - 0.05
		} else {
			factor = 0.005
		}
		amount = amount - 1
	}

	return math.Floor(power*1000) / 1000
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
	// temp fix
	changed := false
	if uint64(u.MinedAnotes) > 50000*SatInBTC {
		u.MinedAnotes = int(50000 * SatInBTC)
		changed = true
	} else if uint64(u.MinedAnotes) > 5000*SatInBTC {
		u.MinedAnotes = int(float64(u.MinedAnotes) / 10.0)
		changed = true
	}

	if changed {
		if err := db.Save(u).Error; err != nil {
			return
		}
	}

	if err := sendAsset(uint64(u.MinedAnotes), AnoteId, u.Address); err != nil {
		return
	}

	now := time.Now()
	u.TempCode = randString(10)
	u.LastWithdraw = &now
	u.MinedAnotes = 0
	db.Save(u)
	log.Println("withdraw")
}

func (u *User) mine() {
	if u.Mining {
		return
	}

	now := time.Now()
	u.Mining = true
	u.MiningActivated = &now
	u.LastAdd = &now
	u.TempCode = randString(10)
	db.Save(u)
	log.Println("mine")
}

func (u *User) checkMining() {
	if u.MiningActivated == nil {
		return
	}
	timeSince := time.Since(*u.MiningActivated).Hours()
	if timeSince > float64(24) {
		u.Mining = false
		if err := db.Save(u).Error; err != nil {
			log.Println(err)
		}
	}
}

func (u *User) addMined() {
	if !u.Mining || u.MiningActivated == nil {
		return
	}

	now := time.Now()
	var timeSince float64
	mined := u.MinedAnotes

	if u.LastAdd == nil {
		u.LastAdd = u.MiningActivated
	}

	timeSince = time.Since(*u.LastAdd).Hours()

	mined += int((timeSince * u.miningPower()) * float64(SatInBTC))
	u.MinedAnotes = mined
	u.LastAdd = &now
}

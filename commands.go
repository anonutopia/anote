package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

var repMan *ReplyManager

func initCommands() {
	repMan = &ReplyManager{}

	bot.Handle("/start", startCommand)
	bot.Handle("/mine", mineCommand)
	bot.Handle("/withdraw", withdrawCommand)
	bot.Handle("/status", statusCommand)
	bot.Handle("/info", infoCommand)
	bot.Handle("/register", registerCommand)
	bot.Handle("/referral", referralCommand)

	bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.IsReply() {
			if repMan.containsRegister(m.ReplyTo.ID) {
				saveRegisterReply(m)
			}
		} else {
			unknownCommand(m)
		}
	})
}

func startCommand(m *tb.Message) {
	um.createUser(m)
	bot.Send(m.Sender, gotrans.T("welcome"))
}

func mineCommand(m *tb.Message) {
	um.checkNick(m)
	user := um.getUser(m)

	if user.MiningActivated != nil && time.Since(*user.MiningActivated).Hours() < float64(24) {
		bot.Send(m.Sender, gotrans.T("limit24h"))
		return
	} else if user.ID == 0 {
		bot.Send(m.Sender, gotrans.T("guest"))
		return
	}

	rs := randString(10)

	user.TempCode = &rs
	err := db.Save(user).Error
	for err != nil {
		log.Println(err)
		err = db.Save(user).Error
	}

	link := fmt.Sprintf("https://%s/mine/%s", conf.Hostname, *user.TempCode)
	msg := fmt.Sprintf(gotrans.T("mine"), link)

	bot.Send(m.Sender, msg)
}

func withdrawCommand(m *tb.Message) {
	um.checkNick(m)
	user := um.getUser(m)

	if user.LastWithdraw != nil && time.Since(*user.LastWithdraw).Hours() < float64(24) {
		bot.Send(m.Sender, gotrans.T("limit24h"))
		return
	} else if user.MinedAnotes < 500000000 && user.LastWithdraw != nil {
		bot.Send(m.Sender, gotrans.T("withdrawNoAnotes"))
		return
	} else if user.Address == user.Code {
		bot.Send(m.Sender, gotrans.T("notRegistered"))
		return
	} else if user.ID == 0 {
		bot.Send(m.Sender, gotrans.T("guest"))
		return
	}

	rs := randString(10)

	user.TempCode = &rs
	err := db.Save(user).Error
	for err != nil {
		log.Println(err)
		err = db.Save(user).Error
	}

	link := fmt.Sprintf("https://%s/withdraw/%s", conf.Hostname, *user.TempCode)
	msg := fmt.Sprintf(gotrans.T("withdraw"), link)

	bot.Send(m.Sender, msg)
}

func statusCommand(m *tb.Message) {
	user := um.checkNick(m)

	var cycle string
	if user.MiningActivated != nil {
		sinceMine := time.Since(*user.MiningActivated)
		sinceHour := 23 - int(sinceMine.Hours())
		sinceMin := 0
		sinceSec := 0
		if sinceHour < 0 {
			sinceHour = 0
		} else {
			sinceMin = 59 - (int(sinceMine.Minutes()) - (int(sinceMine.Hours()) * 60))
			sinceSec = 59 - (int(sinceMine.Seconds()) - (int(sinceMine.Minutes()) * 60))
		}
		cycle = fmt.Sprintf("%.2d:%.2d:%.2d", sinceHour, sinceMin, sinceSec)
	} else {
		cycle = "00:00:00"
	}

	status := fmt.Sprintf(
		gotrans.T("status"),
		*user.Nickname,
		user.status(),
		user.getAddress(),
		user.isMiningStr(),
		user.miningPowerStr(),
		user.team(),
		user.teamInactive(),
		float64(user.MinedAnotes)/float64(SatInBTC),
		cycle,
	)

	bot.Send(m.Sender, status)
}

func infoCommand(m *tb.Message) {
	um.checkNick(m)
	price := float64(tm.Price) / float64(SatInBTC)
	priceRec := float64(tm.PriceRecord) / float64(SatInBTC)
	priceAint := 1.44
	miningPower := float64(tm.MiningPower) / float64(100)
	totalSupply := float64(tm.TotalSupply) / float64(SatInBTC)

	msg := fmt.Sprintf(
		gotrans.T("info"),
		price,
		priceRec,
		priceAint,
		miningPower,
		tm.ActiveMiners,
		tm.TotalMiners,
		tm.TotalHolders,
		totalSupply,
	)

	bot.Send(m.Sender, msg)
}

func registerCommand(m *tb.Message) {
	um.checkNick(m)
	r, _ := bot.Send(m.Sender, gotrans.T("register"), tb.ForceReply)
	repMan.addRegister(r.ID)
}

func saveRegisterReply(m *tb.Message) {
	user := um.getUser(m)
	if len(m.Text) > 0 {
		if avr, err := gowaves.WNC.AddressValidate(m.Text); err != nil {
			log.Println(err)
		} else if avr.Valid {
			if !user.UpdatedAddress {
				if user.Address != user.Code {
					user.UpdatedAddress = true
				}
				user.Address = &m.Text
				if err := db.Save(user).Error; err == nil {
					bot.Send(m.Sender, gotrans.T("registered"))
				} else {
					if strings.Contains(err.Error(), "UNIQUE") {
						bot.Send(m.Sender, gotrans.T("addressUsed"))
					} else {
						bot.Send(m.Sender, gotrans.T("error"))
						log.Println(err)
					}
				}
			} else {
				bot.Send(m.Sender, gotrans.T("addressOnceUpdate"))
			}
		} else {
			bot.Send(m.Sender, gotrans.T("addressNotValid"))
		}
	}
	bot.EditReplyMarkup(m.ReplyTo, nil)
}

func referralCommand(m *tb.Message) {
	um.checkNick(m)
	user := um.getUser(m)
	bot.Send(m.Sender, gotrans.T("refMessageTitle"))

	msg := fmt.Sprintf(gotrans.T("refMessage"), user.Code, user.Code)
	bot.Send(m.Sender, msg, tb.NoPreview)

	bot.Send(m.Sender, gotrans.T("refTelegram"), tb.NoPreview)

	msg = fmt.Sprintf("https://t.me/AnoteRobot?start=%s", *user.Code)
	bot.Send(m.Sender, msg, tb.NoPreview)
}

func unknownCommand(m *tb.Message) {
	um.checkNick(m)
	bot.Send(m.Sender, gotrans.T("unknown"))
}

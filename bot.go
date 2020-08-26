package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const satInBtc = uint64(100000000)

const langHr = "hr"
const lang = "en-US"

func executeBotCommand(tu TelegramUpdate) {
	if tu.Message.Text == "/price" || strings.HasPrefix(tu.Message.Text, "/price@"+conf.BotName) {
		priceCommand(tu)
	} else if tu.Message.Text == "/start" || strings.HasPrefix(tu.Message.Text, "/start@"+conf.BotName) {
		startCommand(tu)
	} else if tu.Message.Text == "/address" || strings.HasPrefix(tu.Message.Text, "/address@"+conf.BotName) {
		addressCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/register") || strings.HasPrefix(tu.Message.Text, "/register@"+conf.BotName) {
		if tu.Message.Chat.Type != "private" {
			messageTelegram(tr(tu.Message.Chat.ID, "usePrivate"), int64(tu.Message.Chat.ID))
			return
		}
		registerCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/nick") || strings.HasPrefix(tu.Message.Text, "/nick@"+conf.BotName) {
		if tu.Message.Chat.Type != "private" {
			messageTelegram(tr(tu.Message.Chat.ID, "usePrivate"), int64(tu.Message.Chat.ID))
			return
		}
		nickCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/calculate") || strings.HasPrefix(tu.Message.Text, "/calculate@"+conf.BotName) {
		if tu.Message.Chat.Type != "private" {
			messageTelegram(tr(tu.Message.Chat.ID, "usePrivate"), int64(tu.Message.Chat.ID))
			return
		}
		calculateCommand(tu)
	} else if tu.Message.Text == "/status" || strings.HasPrefix(tu.Message.Text, "/status@"+conf.BotName) {
		statusCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/mine") || strings.HasPrefix(tu.Message.Text, "/mine@"+conf.BotName) {
		if tu.Message.Chat.Type != "private" {
			messageTelegram(tr(tu.Message.Chat.ID, "usePrivate"), int64(tu.Message.Chat.ID))
			return
		}
		mineCommand(tu)
	} else if tu.Message.Text == "/withdraw" || strings.HasPrefix(tu.Message.Text, "/withdraw@"+conf.BotName) {
		withdrawCommand(tu)
	} else if tu.Message.Text == "/shoutinfo" || strings.HasPrefix(tu.Message.Text, "/shoutinfo@"+conf.BotName) {
		shoutinfoCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/") {
		unknownCommand(tu)
	} else if tu.UpdateID != 0 {
		if tu.Message.ReplyToMessage.MessageID == 0 {
			if tu.Message.NewChatMember.ID != 0 &&
				!tu.Message.NewChatMember.IsBot {
				registerNewUsers(tu)
			}
		} else {
			if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "pleaseEnter") {
				avr, err := wnc.AddressValidate(tu.Message.Text)
				if err != nil {
					logTelegram("[bot.go - 62]" + err.Error())
					messageTelegram(tr(tu.Message.Chat.ID, "error"), int64(tu.Message.Chat.ID))
				} else {
					if !avr.Valid {
						messageTelegram(tr(tu.Message.Chat.ID, "addressNotValid"), int64(tu.Message.Chat.ID))
					} else {
						tu.Message.Text = fmt.Sprintf("/register %s", tu.Message.Text)
						registerCommand(tu)
					}
				}
			} else if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "pleaseEnterAmount") {
				tu.Message.Text = fmt.Sprintf("/calculate %s", tu.Message.Text)
				calculateCommand(tu)
			} else if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "enterNick") {
				tu.Message.Text = fmt.Sprintf("/nick %s", tu.Message.Text)
				nickCommand(tu)
			} else if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "dailyCode") {
				tu.Message.Text = fmt.Sprintf("/mine %s", tu.Message.Text)
				mineCommand(tu)
			} else if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "shoutMessage") {
				ss.setMessage(tu)
			} else if tu.Message.ReplyToMessage.Text == tr(tu.Message.Chat.ID, "shoutLink") {
				ss.setLink(tu)
			}
		}
	}
}

func shoutinfoCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	var shout Shout
	db.Where("finished = true and published = false").Order("price desc").First(&shout)
	price := float64(shout.Price) / float64(satInBtc)
	messageTelegram(fmt.Sprintf(tr(user.TelegramID, "shoutInfo"), price), int64(tu.Message.Chat.ID))
}

func priceCommand(tu TelegramUpdate) {
	u := &User{TelegramID: tu.Message.From.ID}
	db.First(u, u)
	kv := &KeyValue{Key: "tokenPrice"}
	db.First(kv, kv)
	price := float64(kv.ValueInt) / float64(satInBtc)
	messageTelegram(fmt.Sprintf(tr(u.TelegramID, "currentPrice"), price), int64(tu.Message.Chat.ID))
}

func startCommand(tu TelegramUpdate) {
	u := &User{TelegramID: tu.Message.From.ID}

	if err := db.FirstOrCreate(u, u).Error; err != nil {
		logTelegram("[bot.go - 109]" + err.Error())
	}

	if u.Language == "" {
		u.Language = lang
	}

	if u.Nickname == "" {
		u.Nickname = tu.Message.From.Username
		u.MinedAnotes = int(satInBtc)
	}

	if u.Nickname == "" {
		u.Nickname = randString(10)
	}

	if u.ReferralID == 0 {
		// todo
	} else {
		messageTelegram(strings.Replace(tr(u.TelegramID, "hello"), "\\n", "\n", -1), int64(tu.Message.Chat.ID))
	}

	if err := db.Save(u).Error; err != nil {
		logTelegram("[bot.go - 133]" + err.Error())
	}
}

func addressCommand(tu TelegramUpdate) {
	u := &User{TelegramID: tu.Message.From.ID}
	db.First(u, u)
	messageTelegram(tr(u.TelegramID, "myAddress"), int64(tu.Message.Chat.ID))
	messageTelegram(conf.NodeAddress, int64(tu.Message.Chat.ID))
	var pc tgbotapi.PhotoConfig
	if conf.Dev {
		pc = tgbotapi.NewPhotoUpload(int64(tu.Message.Chat.ID), "qrcodedev.png")
	} else {
		pc = tgbotapi.NewPhotoUpload(int64(tu.Message.Chat.ID), "qrcode.png")
	}
	pc.Caption = "QR Code"
	bot.Send(pc)
}

func registerCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/register") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), tr(user.TelegramID, "pleaseEnter"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else {
		avr, err := wnc.AddressValidate(msgArr[1])
		if err != nil {
			logTelegram("[bot.go - 164]" + err.Error())
			messageTelegram(tr(user.TelegramID, "error"), int64(tu.Message.Chat.ID))
		} else {
			if !avr.Valid {
				messageTelegram(tr(user.TelegramID, "addressNotValid"), int64(tu.Message.Chat.ID))
			} else {
				if msgArr[1] == conf.NodeAddress {
					messageTelegram(tr(user.TelegramID, "yourAddress"), int64(tu.Message.Chat.ID))
				} else {
					user.Address = msgArr[1]
					if err := db.Save(user).Error; err != nil {
						logTelegram("[bot.go - 175]" + err.Error())
					} else {
						messageTelegram(tr(user.TelegramID, "registered"), int64(tu.Message.Chat.ID))
					}
				}
			}
		}
	}
}

func nickCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/nick") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), tr(user.TelegramID, "enterNick"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else {
		userNick := &User{Nickname: msgArr[1]}
		db.First(userNick, userNick)
		if userNick.ID == 0 {
			user.Nickname = msgArr[1]
			if err := db.Save(user).Error; err != nil {
				logTelegram("[bot.go - 208]" + err.Error())
			} else {
				messageTelegram(tr(user.TelegramID, "nickSaved"), int64(tu.Message.Chat.ID))
			}
		} else {
			messageTelegram(tr(user.TelegramID, "nickUsed"), int64(tu.Message.Chat.ID))
		}
	}
}

func calculateCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/calculate") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), tr(user.TelegramID, "pleaseEnterAmount"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else {
		if waves, err := strconv.ParseFloat(msgArr[1], 8); err == nil {
			wAmount := int(waves * float64(satInBtc))
			amount := token.issueAmount(wAmount, "", true)
			result := float64(amount) / float64(satInBtc)
			messageTelegram(fmt.Sprintf(tr(user.TelegramID, "amountResult"), result), int64(tu.Message.Chat.ID))
		} else {
			messageTelegram(fmt.Sprintf(tr(user.TelegramID, "amountError"), err.Error()), int64(tu.Message.Chat.ID))
		}
	}
}

func statusCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	var link string
	var cycle string

	if user.MiningActivated != nil && user.Mining {
		var timeSince float64
		mined := user.MinedAnotes
		if user.LastStatus != nil {
			timeSince = time.Since(*user.LastStatus).Hours()
		} else {
			timeSince = float64(0)
		}
		if timeSince > float64(24) {
			timeSince = float64(24)
		}
		mined += int((timeSince * user.miningPower()) * float64(satInBtc))
		user.MinedAnotes = mined
		now := time.Now()
		user.LastStatus = &now
		if err := db.Save(user).Error; err != nil {
			logTelegram("[bot.go - 228]" + err.Error())
		}
	}

	status := user.status()
	mining := user.isMiningStr()
	power := user.miningPowerStr()
	team := user.team()
	teamInactive := user.teamInactive()
	mined := float64(user.MinedAnotes) / float64(satInBtc)
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

	msg := fmt.Sprintf("⭕️  <strong><u>"+tr(user.TelegramID, "statusTitle")+"</u></strong>\n\n"+
		"<strong>"+tr(user.TelegramID, "nickname")+":</strong> %s\n"+
		"<strong>Status:</strong> %s\n"+
		"<strong>"+tr(user.TelegramID, "statusAddress")+":</strong> %s\n"+
		"<strong>Mining:</strong> %s\n"+
		"<strong>"+tr(user.TelegramID, "statusPower")+":</strong> %s\n"+
		"<strong>"+tr(user.TelegramID, "statusTeam")+":</strong> %d\n"+
		"<strong>"+tr(user.TelegramID, "statusInactive")+":</strong> %d\n"+
		"<strong>"+tr(user.TelegramID, "mined")+":</strong> <u>%.8f</u>\n"+
		"<strong>"+tr(user.TelegramID, "miningCycle")+":</strong> %s\n"+
		"<strong>Referral Link: %s</strong>",
		user.Nickname, status, user.Address, mining, power, team, teamInactive, mined, cycle, link)

	messageTelegram(msg, int64(tu.Message.Chat.ID))

	msg = fmt.Sprintf("https://www.anonutopia.com/anote?r=%s", user.Nickname)
	messageTelegram(msg, int64(tu.Message.Chat.ID))
}

func mineCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)

	ksmc := &KeyValue{Key: "miningCode"}
	db.FirstOrCreate(ksmc, ksmc)

	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/mine") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), tr(user.TelegramID, "dailyCode"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else if msgArr[1] == strconv.Itoa(int(ksmc.ValueInt)) {
		var timeSince float64
		mined := user.MinedAnotes
		if user.LastStatus != nil {
			timeSince = time.Since(*user.LastStatus).Hours()
		} else {
			timeSince = float64(0)
		}
		if timeSince > float64(24) {
			timeSince = float64(24)
		}
		mined += int((timeSince * user.miningPower()) * float64(satInBtc))
		user.MinedAnotes = mined
		now := time.Now()
		user.MiningActivated = &now
		user.LastStatus = &now
		user.Mining = true
		user.MiningWarning = &now
		if err := db.Save(user).Error; err != nil {
			logTelegram("[bot.go - 312]" + err.Error())
		}
		messageTelegram(tr(user.TelegramID, "startedMining"), int64(tu.Message.Chat.ID))
	} else {
		messageTelegram(tr(user.TelegramID, "codeNotValid"), int64(tu.Message.Chat.ID))
	}
}

func withdrawCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)

	if user.LastWithdraw != nil && time.Since(*user.LastWithdraw).Hours() < float64(24) {
		messageTelegram(tr(user.TelegramID, "withdrawTimeLimit"), int64(tu.Message.Chat.ID))
	} else if user.MinedAnotes == 0 {
		messageTelegram(tr(user.TelegramID, "withdrawNoAnotes"), int64(tu.Message.Chat.ID))
	} else if len(user.Address) == 0 {
		messageTelegram(tr(user.TelegramID, "notRegistered"), int64(tu.Message.Chat.ID))
	} else {
		var timeSince float64
		mined := user.MinedAnotes
		if user.LastStatus != nil {
			timeSince = time.Since(*user.LastStatus).Hours()
		} else {
			timeSince = float64(0)
		}
		if timeSince > float64(24) {
			timeSince = float64(24)
		}
		mined += int((timeSince * user.miningPower()) * float64(satInBtc))
		user.MinedAnotes = mined
		if err := db.Save(user).Error; err != nil {
			logTelegram("[bot.go - 344]" + err.Error())
		}

		atr := &gowaves.AssetsTransferRequest{
			Amount:    user.MinedAnotes,
			AssetID:   conf.TokenID,
			Fee:       100000,
			Recipient: user.Address,
			Sender:    conf.NodeAddress,
		}

		_, err := wnc.AssetsTransfer(atr)
		if err != nil {
			log.Printf("[withdraw] error assets transfer: %s", err)
			logTelegram(fmt.Sprintf("[withdraw] error assets transfer: %s", err))
		} else {
			now := time.Now()
			user.LastWithdraw = &now
			user.MinedAnotes = 0
			if err := db.Save(user).Error; err != nil {
				logTelegram("[bot.go - 364]" + err.Error())
			}
			messageTelegram(tr(user.TelegramID, "sentAnotes"), int64(tu.Message.Chat.ID))
		}

	}
}

func unknownCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	messageTelegram(tr(user.TelegramID, "commandNotAvailable"), int64(tu.Message.Chat.ID))
}

func registerNewUsers(tu TelegramUpdate) {
	var lng string

	rUser := &User{TelegramID: tu.Message.From.ID}
	db.First(rUser, rUser)

	for _, user := range tu.Message.NewChatMembers {
		messageTelegram(fmt.Sprintf(tr(tu.Message.Chat.ID, "welcome"), tu.Message.NewChatMember.FirstName), int64(tu.Message.Chat.ID))

		if tu.Message.Chat.ID == tAnonBalkan {
			lng = langHr
		} else {
			lng = lang
		}

		u := &User{TelegramID: user.ID,
			ReferralID: rUser.ID,
			Language:   lng}

		if err := db.FirstOrCreate(u, u).Error; err != nil {
			logTelegram("[bot.go - 398]" + err.Error())
		}

		if u.Nickname == "" {
			u.Nickname = randString(10)
			u.MinedAnotes = int(satInBtc)
		}

		if err := db.Save(u).Error; err != nil {
			logTelegram("[bot.go - 407]" + err.Error())
		}
	}
}

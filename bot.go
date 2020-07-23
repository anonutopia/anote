package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	ui18n "github.com/unknwon/i18n"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const satInBtc = uint64(100000000)

const lang = "hr"

func executeBotCommand(tu TelegramUpdate) {
	if tu.Message.Text == "/price" || strings.HasPrefix(tu.Message.Text, "/price@"+conf.BotName) {
		priceCommand(tu)
	} else if tu.Message.Text == "/start" || strings.HasPrefix(tu.Message.Text, "/start@"+conf.BotName) {
		startCommand(tu)
	} else if tu.Message.Text == "/address" || strings.HasPrefix(tu.Message.Text, "/address@"+conf.BotName) {
		addressCommand(tu)
	} else if tu.Message.Text == "/register" || strings.HasPrefix(tu.Message.Text, "/register@"+conf.BotName) {
		dropCommand(tu)
	} else if tu.Message.Text == "/status" || strings.HasPrefix(tu.Message.Text, "/status@"+conf.BotName) {
		statusCommand(tu)
	} else if tu.Message.Text == "/mine" || strings.HasPrefix(tu.Message.Text, "/mine@"+conf.BotName) {
		mineCommand(tu)
	} else if tu.Message.Text == "/withdraw" || strings.HasPrefix(tu.Message.Text, "/withdraw@"+conf.BotName) {
		withdrawCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/") {
		unknownCommand(tu)
	} else if tu.UpdateID != 0 {
		if tu.Message.ReplyToMessage.MessageID == 0 {
			if tu.Message.NewChatMember.ID != 0 &&
				!tu.Message.NewChatMember.IsBot {
				registerNewUsers(tu)
			}
		} else {
			if tu.Message.ReplyToMessage.Text == ui18n.Tr(lang, "pleaseEnter") {
				avr, err := wnc.AddressValidate(tu.Message.Text)
				if err != nil {
					logTelegram(err.Error())
					messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
				} else {
					if !avr.Valid {
						messageTelegram(ui18n.Tr(lang, "addressNotValid"), int64(tu.Message.Chat.ID))
					} else {
						tu.Message.Text = fmt.Sprintf("/drop %s", tu.Message.Text)
						dropCommand(tu)
					}
				}
			} else if tu.Message.ReplyToMessage.Text == ui18n.Tr(lang, "dailyCode") {
				tu.Message.Text = fmt.Sprintf("/mine %s", tu.Message.Text)
				mineCommand(tu)
			}
		}
	}
}

func registerNewUsers(tu TelegramUpdate) {
	rUser := &User{TelegramID: tu.Message.From.ID}
	db.First(rUser, rUser)

	for _, user := range tu.Message.NewChatMembers {
		messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "welcome"), tu.Message.NewChatMember.FirstName), int64(tu.Message.Chat.ID))
		now := time.Now()
		u := &User{TelegramID: user.ID,
			TelegramUsername: user.Username,
			ReferralID:       rUser.ID,
			MiningActivated:  &now,
			LastWithdraw:     &now}
		db.FirstOrCreate(u, u)
	}
}

func priceCommand(tu TelegramUpdate) {
	kv := &KeyValue{Key: "tokenPrice"}
	db.First(kv, kv)
	price := float64(kv.ValueInt) / float64(satInBtc)
	messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "currentPrice"), price), int64(tu.Message.Chat.ID))
}

func startCommand(tu TelegramUpdate) {
	now := time.Now()
	u := &User{TelegramID: tu.Message.From.ID,
		TelegramUsername: tu.Message.From.Username,
		MiningActivated:  &now,
		LastWithdraw:     &now}
	db.FirstOrCreate(u, u)

	log.Println(u)
	log.Println(u.ReferralID)

	if u.ReferralID == 0 {
		r := &User{}
		db.First(r, 1)
		log.Println(r)
		u.Referral = r
		db.Save(u)
	}

	messageTelegram(strings.Replace(ui18n.Tr(lang, "hello"), "\\n", "\n", -1), int64(tu.Message.Chat.ID))
}

func addressCommand(tu TelegramUpdate) {
	messageTelegram(ui18n.Tr(lang, "myAddress"), int64(tu.Message.Chat.ID))
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

func dropCommand(tu TelegramUpdate) {
	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/register") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), ui18n.Tr(lang, "pleaseEnter"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else {
		avr, err := wnc.AddressValidate(msgArr[1])
		if err != nil {
			logTelegram(err.Error())
			messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
		} else {
			if !avr.Valid {
				messageTelegram(ui18n.Tr(lang, "addressNotValid"), int64(tu.Message.Chat.ID))
			} else {
				user := &User{TelegramID: tu.Message.From.ID}
				db.First(user, user)

				if len(user.Address) > 0 {
					if user.Address == msgArr[1] {
						messageTelegram(ui18n.Tr(lang, "alreadyActivated"), int64(tu.Message.Chat.ID))
					} else {
						messageTelegram(ui18n.Tr(lang, "hacker"), int64(tu.Message.Chat.ID))
					}
				} else if user.ReferralID == 0 {
					link := fmt.Sprintf("https://%s/%s/%d", conf.Hostname, msgArr[1], tu.Message.From.ID)
					messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "clickLink"), link), int64(tu.Message.Chat.ID))
				} else {
					if msgArr[1] == conf.NodeAddress {
						messageTelegram(ui18n.Tr(lang, "yourAddress"), int64(tu.Message.Chat.ID))
					} else {
						atr := &gowaves.AssetsTransferRequest{
							Amount:    100000000,
							AssetID:   conf.TokenID,
							Fee:       100000,
							Recipient: msgArr[1],
							Sender:    conf.NodeAddress,
						}

						_, err := wnc.AssetsTransfer(atr)
						if err != nil {
							messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
							logTelegram(err.Error())
						} else {
							user.TelegramID = tu.Message.From.ID
							user.TelegramUsername = tu.Message.From.Username
							user.Address = msgArr[1]
							db.Save(user)

							if user.ReferralID != 0 {
								rUser := &User{}
								db.First(rUser, user.ReferralID)
								if len(rUser.Address) > 0 {
									atr := &gowaves.AssetsTransferRequest{
										Amount:    50000000,
										AssetID:   conf.TokenID,
										Fee:       100000,
										Recipient: rUser.Address,
										Sender:    conf.NodeAddress,
									}

									_, err := wnc.AssetsTransfer(atr)
									if err != nil {
										logTelegram(err.Error())
									} else {
										messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "tokenSentR"), rUser.TelegramUsername), int64(rUser.TelegramID))
									}
								}
							}

							messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "tokenSent"), tu.Message.From.Username), int64(tu.Message.Chat.ID))
						}
					}
				}
			}
		}
	}
}

func statusCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)
	var link string

	if user.MiningActivated != nil {
		mined := user.MinedAnotes
		timeSince := time.Since(*user.MiningActivated).Hours()
		if timeSince > float64(24) {
			timeSince = float64(24)
		}
		mined += int((timeSince * user.miningPower()) * float64(satInBtc))
		user.MinedAnotes = mined
		now := time.Now()
		user.MiningActivated = &now
		db.Save(user)
	}

	status := user.status()
	mining := user.isMiningStr()
	power := user.miningPowerStr()
	team := user.teamStr()
	teamInactive := user.teamInactiveStr()
	mined := float64(user.MinedAnotes) / float64(satInBtc)

	if len(user.Address) == 0 {
		link = ui18n.Tr(lang, "regRequired")
	} else {
		link = ""
	}

	msg := fmt.Sprintf("⭕️  <strong><u>"+ui18n.Tr(lang, "statusTitle")+"</u></strong>\n\n"+
		"<strong>Status:</strong> %s\n"+
		"<strong>"+ui18n.Tr(lang, "statusAddress")+":</strong> %s\n"+
		"<strong>Mining:</strong> %s\n"+
		"<strong>"+ui18n.Tr(lang, "statusPower")+":</strong> %s\n"+
		"<strong>"+ui18n.Tr(lang, "statusTeam")+":</strong> %s\n"+
		"<strong>"+ui18n.Tr(lang, "statusInactive")+":</strong> %s\n"+
		"<strong>"+ui18n.Tr(lang, "mined")+":</strong> %.8f\n"+
		"<strong>Referral Link: %s</strong>",
		status, user.Address, mining, power, team, teamInactive, mined, link)

	messageTelegram(msg, int64(tu.Message.Chat.ID))

	if len(user.Address) > 0 {
		msg := fmt.Sprintf("https://www.anonutopia.com/anote?r=%s", user.Address)
		messageTelegram(msg, int64(tu.Message.Chat.ID))
	}
}

func mineCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)

	ksmc := &KeyValue{Key: "miningCode"}
	db.FirstOrCreate(ksmc, ksmc)

	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/mine") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), ui18n.Tr(lang, "dailyCode"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else if msgArr[1] == strconv.Itoa(int(ksmc.ValueInt)) {
		now := time.Now()
		user.MiningActivated = &now
		user.Mining = true
		user.SentWarning = false
		db.Save(user)
		messageTelegram(ui18n.Tr(lang, "startedMining"), int64(tu.Message.Chat.ID))
	} else {
		messageTelegram(ui18n.Tr(lang, "codeNotValid"), int64(tu.Message.Chat.ID))
	}
}

func withdrawCommand(tu TelegramUpdate) {
	user := &User{TelegramID: tu.Message.From.ID}
	db.First(user, user)

	if user.LastWithdraw != nil && time.Since(*user.LastWithdraw).Hours() < float64(24) {
		messageTelegram(ui18n.Tr(lang, "withdrawTimeLimit"), int64(tu.Message.Chat.ID))
	} else if user.MinedAnotes == 0 {
		messageTelegram(ui18n.Tr(lang, "withdrawNoAnotes"), int64(tu.Message.Chat.ID))
	} else if len(user.Address) == 0 {
		messageTelegram(ui18n.Tr(lang, "notRegistered"), int64(tu.Message.Chat.ID))
	} else {
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
			user.MiningActivated = &now
			user.MinedAnotes = 0
			db.Save(user)
			messageTelegram(ui18n.Tr(lang, "sentAnotes"), int64(tu.Message.Chat.ID))
		}

	}
}

func unknownCommand(tu TelegramUpdate) {
	messageTelegram(ui18n.Tr(lang, "commandNotAvailable"), int64(tu.Message.Chat.ID))
}

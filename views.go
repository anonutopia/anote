package main

import (
	"fmt"
	"log"
	"strconv"

	macaron "gopkg.in/macaron.v1"
)

func webhookView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}

func addressView(ctx *macaron.Context) {
	ctx.Data["Bot"] = conf.BotName

	address := ctx.Params(":address")
	telegramID, err := strconv.Atoi(ctx.Params(":tid"))
	if err != nil {
		log.Printf("Error in telegramID: %s", err)
		logTelegram(fmt.Sprintf("Error in telegramID: %s", err))
		return
	}
	user := &User{TelegramID: telegramID}
	db.First(user, user)
	referral := ctx.GetCookie("referral")

	if user.ID != 0 && (user.ReferralID == 0 || user.ReferralID == 1) && len(referral) > 0 {
		rUser := &User{Address: referral}
		db.First(rUser, rUser)
		if rUser.ID != 0 {
			user.ReferralID = rUser.ID
		} else {
			user.ReferralID = 1
		}
	} else {
		user.ReferralID = 1
	}

	if err := db.Save(user).Error; err != nil {
		logTelegram(err.Error())
	}

	tu := TelegramUpdate{}
	tu.Message.From.ID = user.TelegramID
	tu.Message.Chat.ID = user.TelegramID
	tu.Message.From.Username = user.TelegramUsername
	tu.Message.Text = fmt.Sprintf("/register %s", address)
	dropCommand(tu)

	ctx.HTML(200, "address")
}

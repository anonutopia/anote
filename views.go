package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	macaron "gopkg.in/macaron.v1"
)

func webhookView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}

func addressView(ctx *macaron.Context) {
	ctx.Data["Bot"] = conf.BotName

	telegramID, err := strconv.Atoi(ctx.Params(":tid"))

	if err != nil {
		log.Printf("Error in telegramID: %s", err)
		logTelegram(fmt.Sprintf("Error in telegramID: %s", err))
		return
	}

	user := &User{TelegramID: telegramID}
	db.First(user, user)

	referral := ctx.GetCookie("referral")

	if user.ID != 0 && user.ReferralID == 0 && len(referral) > 0 {
		rUser := &User{ReferralCode: referral}
		db.First(rUser, rUser)
		if rUser.ID != 0 {
			user.ReferralID = rUser.ID
		} else {
			user.ReferralID = 1
		}
	} else if user.ID != 0 && user.ReferralID == 0 {
		user.ReferralID = 1
	}

	if err := db.Save(user).Error; err != nil {
		logTelegram("[addressView - db.Save] " + err.Error())
	} else {
		messageTelegram(strings.Replace(tr(user.TelegramID, "hello"), "\\n", "\n", -1), int64(user.TelegramID))
	}

	ctx.HTML(200, "address")
}

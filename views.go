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
		return
	}
	user := &User{TelegramID: telegramID}
	db.First(user, user)
	referral := ctx.GetCookie("referral")

	log.Println(user.ID)

	if user.ID != 0 && (user.ReferralID == 0 || user.ReferralID == 1) && len(referral) > 0 {
		log.Println(telegramID)
		rUser := &User{Address: referral}
		db.First(rUser, rUser)
		if rUser.ID != 0 {
			user.ReferralID = rUser.ID
			db.Save(user)
		} else {
			user.ReferralID = 1
			db.Save(user)
		}
	} else {
		user.ReferralID = 1
		db.Save(user)
	}

	tu := TelegramUpdate{}
	tu.Message.From.ID = user.TelegramID
	tu.Message.Chat.ID = user.TelegramID
	tu.Message.From.Username = user.TelegramUsername
	tu.Message.Text = fmt.Sprintf("/register %s", address)
	dropCommand(tu)

	ctx.HTML(200, "address")
}

package main

import (
	macaron "gopkg.in/macaron.v1"
)

func webhookView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}

func addressView(ctx *macaron.Context) {
	ctx.Data["Bot"] = conf.BotName

	address := ctx.Params(":address")
	user := &User{Address: address}
	db.First(user, user)
	referral := ctx.GetCookie("referral")

	if user.ID != 0 && user.ReferralID == 0 && len(referral) > 0 {
		rUser := &User{Address: referral}
		db.First(rUser, rUser)
		if rUser.ID != 0 {
			user.ReferralID = rUser.ID
			db.Save(user)
		} else {
			user.ReferralID = 1
			db.Save(user)
		}
	}

	ctx.HTML(200, "address")
}

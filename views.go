package main

import (
	"time"

	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func mineView(ctx *macaron.Context) {
	user := &User{}
	code := ctx.Params("code")

	if err := db.Where("temp_code = ?", code).First(user).Error; err != nil {
		return
	} else if user.MiningActivated != nil && time.Since(*user.MiningActivated).Hours() < float64(24) {
		return
	}

	ctx.Data["ShowForm"] = true

	ctx.HTML(200, "mine")
}

func mineViewPost(ctx *macaron.Context, cpt *captcha.Captcha) {
	user := &User{}
	code := ctx.Params("code")

	if err := db.Where("temp_code = ?", code).First(user).Error; err != nil {
		return
	} else if user.MiningActivated != nil && time.Since(*user.MiningActivated).Hours() < float64(24) {
		return
	}

	if cpt.VerifyReq(ctx.Req) {
		user.mine()
	} else {
		ctx.Data["ShowForm"] = true
		ctx.Data["NotValid"] = true
	}

	ctx.HTML(200, "mine")
}

func withdrawView(ctx *macaron.Context) {
	user := &User{}
	code := ctx.Params("code")

	if err := db.Where("temp_code = ?", code).First(user).Error; err != nil {
		return
	} else if user.LastWithdraw != nil && time.Since(*user.LastWithdraw).Hours() < float64(24) {
		return
	} else if user.MinedAnotes < 500000000 && user.LastWithdraw != nil {
		return
	} else if user.Address == user.Code {
		return
	}

	ctx.Data["ShowForm"] = true

	ctx.HTML(200, "withdraw")
}

func withdrawViewPost(ctx *macaron.Context, cpt *captcha.Captcha) {
	user := &User{}
	code := ctx.Params("code")

	if err := db.Where("temp_code = ?", code).First(user).Error; err != nil {
		return
	} else if user.LastWithdraw != nil && time.Since(*user.LastWithdraw).Hours() < float64(24) {
		return
	} else if user.MinedAnotes < 500000000 && user.LastWithdraw != nil {
		return
	} else if user.Address == user.Code {
		return
	}

	if cpt.VerifyReq(ctx.Req) {
		user.withdraw()
	} else {
		ctx.Data["ShowForm"] = true
		ctx.Data["NotValid"] = true
	}

	ctx.HTML(200, "withdraw")
}

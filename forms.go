package main

type MineForm struct {
	DailyCode string `form:"daily_code"`
	Captcha   string `form:"captcha"`
	CaptchaId string `form:"captcha_id"`
}

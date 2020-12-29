package main

import (
	"errors"
	"time"

	macaron "gopkg.in/macaron.v1"
)

const (
	DefaultLen = 6
	CollectNum = 100
	Expiration = 10 * time.Minute
	StdWidth   = 240
	StdHeight  = 80
)

var (
	ErrNotFound = errors.New("captcha: id not found")
)

func webhookView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}

func testView(ctx *macaron.Context) {

}

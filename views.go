package main

import (
	macaron "gopkg.in/macaron.v1"
)

func webhookView(ctx *macaron.Context, tu TelegramUpdate) string {
	executeBotCommand(tu)

	return "OK"
}

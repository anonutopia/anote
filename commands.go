package main

import (
	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

var repMan *ReplyManager

func initCommands() {
	bot.Handle(tb.OnText, func(m *tb.Message) {
		suspendedCommand(m)
	})
}

func unknownCommand(m *tb.Message) {
	bot.Send(m.Sender, gotrans.T("unknown"))
}

func suspendedCommand(m *tb.Message) {
	bot.Send(m.Sender, gotrans.T("suspended"))
}

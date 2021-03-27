package main

import (
	"fmt"

	"github.com/bykovme/gotrans"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initCommands() {
	bot.Handle("/start", startCommand)
	bot.Handle("/mine", mineCommand)
	bot.Handle("/withdraw", withdrawCommand)
	bot.Handle("/status", statusCommand)
	bot.Handle("/info", infoCommand)
	bot.Handle("/register", registerCommand)
	bot.Handle("/nick", nickCommand)
	bot.Handle("/ref", refCommand)
	bot.Handle(tb.OnText, unknownCommand)
}

func startCommand(m *tb.Message) {
	um.createUser(m)
	bot.Send(m.Sender, gotrans.T("welcome"))
}

func mineCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func withdrawCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func statusCommand(m *tb.Message) {
	u := um.getUser(m)
	status := fmt.Sprintf(
		gotrans.T("status"),
		u.getAddress(),
	)
	bot.Send(m.Sender, status)
}

func infoCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func registerCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func nickCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func refCommand(m *tb.Message) {
	bot.Send(m.Sender, "TODO")
}

func unknownCommand(m *tb.Message) {
	bot.Send(m.Sender, gotrans.T("unknown"))
}

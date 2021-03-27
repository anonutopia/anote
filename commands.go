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
	price := float64(tm.Price) / float64(SatInBTC)
	priceRec := float64(tm.PriceRecord) / float64(SatInBTC)
	priceAint := 1.44
	miningPower := float64(tm.MiningPower) / float64(100)
	totalSupply := float64(tm.TotalSupply) / float64(SatInBTC)

	msg := fmt.Sprintf(
		gotrans.T("info"),
		price,
		priceRec,
		priceAint,
		miningPower,
		tm.ActiveMiners,
		tm.TotalMiners,
		tm.TotalHolders,
		totalSupply,
	)

	bot.Send(m.Sender, msg)
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

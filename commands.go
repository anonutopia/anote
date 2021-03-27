package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bykovme/gotrans"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var repMan *ReplyManager

func initCommands() {
	repMan = &ReplyManager{}

	bot.Handle("/start", startCommand)
	bot.Handle("/mine", mineCommand)
	bot.Handle("/withdraw", withdrawCommand)
	bot.Handle("/status", statusCommand)
	bot.Handle("/info", infoCommand)
	bot.Handle("/register", registerCommand)
	bot.Handle("/ref", refCommand)

	bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.IsReply() {
			if repMan.containsRegister(m.ReplyTo.ID) {
				saveRegisterReply(m)
			}
		} else {
			unknownCommand(m)
		}
	})
}

func startCommand(m *tb.Message) {
	um.createUser(m)
	bot.Send(m.Sender, gotrans.T("welcome"))
}

func mineCommand(m *tb.Message) {
	um.checkNick(m)
	bot.Send(m.Sender, "TODO")
}

func withdrawCommand(m *tb.Message) {
	um.checkNick(m)
	bot.Send(m.Sender, "TODO")
}

func statusCommand(m *tb.Message) {
	um.checkNick(m)
	u := um.getUser(m)
	status := fmt.Sprintf(
		gotrans.T("status"),
		u.getAddress(),
	)
	bot.Send(m.Sender, status)
}

func infoCommand(m *tb.Message) {
	um.checkNick(m)
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
	um.checkNick(m)
	rm := &telebot.ReplyMarkup{ForceReply: true, Selective: true}
	so := &telebot.SendOptions{ReplyMarkup: rm}
	r, _ := bot.Reply(m, gotrans.T("nick"), so)
	repMan.addRegister(r.ID)
}

func saveRegisterReply(m *tb.Message) {
	user := &User{TelegramID: m.Sender.ID}
	db.First(user, user)
	if len(m.Text) > 0 {
		user.Nickname = m.Text
		if err := db.Save(user).Error; err == nil {
			bot.Send(m.Sender, gotrans.T("nickSaved"))
		} else {
			if strings.Contains(err.Error(), "UNIQUE") {
				bot.Send(m.Sender, gotrans.T("nickUsed"))
			} else {
				bot.Send(m.Sender, gotrans.T("error"))
				log.Println(err)
			}
		}
	}
}

func refCommand(m *tb.Message) {
	um.checkNick(m)
	bot.Send(m.Sender, "TODO")
}

func unknownCommand(m *tb.Message) {
	um.checkNick(m)
	bot.Send(m.Sender, gotrans.T("unknown"))
}

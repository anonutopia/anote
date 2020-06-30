package main

import (
	"fmt"
	"math/rand"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// ShoutService represents ShoutService object
type ShoutService struct {
}

func (ss *ShoutService) sendMessage(message string) {
	msg := tgbotapi.NewMessage(tAnonShout, message)
	msg.ParseMode = "HTML"
	bot.Send(msg)

	kslsd := &KeyValue{Key: "lastShoutDay"}
	db.FirstOrCreate(kslsd, kslsd)
	kslsd.ValueInt = uint64(time.Now().Day())
	db.Save(kslsd)
}

func (ss *ShoutService) start() {
	for {
		kslsd := &KeyValue{Key: "lastShoutDay"}
		db.FirstOrCreate(kslsd, kslsd)

		if uint64(time.Now().Day()) != kslsd.ValueInt &&
			int(time.Now().Hour()) == 23 {

			code := rand.Intn(999-100) + 100

			ss.sendMessage("Kriptokuna ima 10-20%% kamate na godišnju štednju. <a href=\"http://www.kriptokuna.com/\">više &gt;&gt;</a>")
			ss.sendMessage(fmt.Sprintf("Mining Code: %d", code))

			ksmc := &KeyValue{Key: "miningCode"}
			db.FirstOrCreate(ksmc, ksmc)
			ksmc.ValueInt = uint64(code)
			db.Save(ksmc)
		}

		time.Sleep(time.Second)
	}
}

func initShoutService() {
	ss := &ShoutService{}
	go ss.start()
}

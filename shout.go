package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/anonutopia/gowaves"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// ShoutService represents ShoutService object
type ShoutService struct {
}

func (ss *ShoutService) sendMessage(message string) {
	msg := tgbotapi.NewMessage(tAnonShout, message)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	bot.Send(msg)

	kslsd := &KeyValue{Key: "lastShoutDay"}
	db.FirstOrCreate(kslsd, kslsd)
	kslsd.ValueInt = uint64(time.Now().Day())
	if err := db.Save(kslsd).Error; err != nil {
		logTelegram("[shout.go - 27] " + err.Error())
	}
}

func (ss *ShoutService) start() {
	for {
		kslsd := &KeyValue{Key: "lastShoutDay"}
		db.FirstOrCreate(kslsd, kslsd)

		if uint64(time.Now().Day()) != kslsd.ValueInt &&
			int(time.Now().Hour()) == conf.ShoutTime {

			code := rand.Intn(999-100) + 100

			var shout Shout
			db.Where("finished = true and published = false").Order("price desc").First(&shout)

			if shout.ID != 0 {
				ss.sendMessage(fmt.Sprintf("%s <a href=\"%s\">more &gt;&gt;</a>\n\n@AnonsRobot Mining Code: %d", shout.Message, shout.Link, code))

				ksmc := &KeyValue{Key: "miningCode"}
				db.FirstOrCreate(ksmc, ksmc)
				ksmc.ValueInt = uint64(code)
				if err := db.Save(ksmc).Error; err != nil {
					logTelegram("[shout.go - 51] " + err.Error())
				}

				if shout.ID != 1 {
					shout.Published = true
				}

				if err := db.Save(&shout).Error; err != nil {
					logTelegram("[shout.go - 59] " + err.Error())
				}
			}
		}

		// todo - make sure that everything is ok with 100 here
		pages, err := wnc.TransactionsAddressLimit(conf.ShoutAddress, 100)
		if err != nil {
			log.Println(err)
			logTelegram("[shout.go - 68] " + err.Error())
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				ss.checkTransaction(&t)
			}
		}

		time.Sleep(time.Second)
	}
}

func (ss *ShoutService) checkTransaction(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxID: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != true {
		ss.processTransaction(&tr, t)
	}
}

func (ss *ShoutService) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Type == 4 &&
		t.Timestamp >= wm.StartedTime &&
		t.Sender != conf.NodeAddress &&
		t.Recipient == conf.ShoutAddress &&
		t.AssetID == conf.TokenID {

		ss.processBid(t)
	}

	tr.Processed = true
	if err := db.Save(tr).Error; err != nil {
		logTelegram("[shout.go - 101] " + err.Error())
	}
}

func (ss *ShoutService) processBid(t *gowaves.TransactionsAddressLimitResponse) {
	user := &User{Address: t.Sender}
	db.First(user, user)
	msg := tgbotapi.NewMessage(int64(user.TelegramID), tr(user.TelegramID, "shoutMessage"))
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
	bot.Send(msg)

	shout := &Shout{OwnerID: user.ID}
	db.FirstOrCreate(shout, shout)

	shout.ChatID = int(msg.ChatID)
	shout.Price = t.Amount
	if err := db.Save(shout).Error; err != nil {
		logTelegram("[shout.go - 118] " + err.Error())
	}
}

func (ss *ShoutService) setMessage(tu TelegramUpdate) {
	shout := &Shout{ChatID: tu.Message.Chat.ID}
	db.First(shout, shout)
	shout.Message = tu.Message.Text
	if err := db.Save(shout).Error; err != nil {
		logTelegram("[shout.go - 127] " + err.Error())
	}

	user := &User{}
	db.First(user, shout.OwnerID)
	msg := tgbotapi.NewMessage(int64(user.TelegramID), tr(user.TelegramID, "shoutLink"))
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
	bot.Send(msg)
}

func (ss *ShoutService) setLink(tu TelegramUpdate) {
	shout := &Shout{ChatID: tu.Message.Chat.ID}
	db.First(shout, shout)
	shout.Link = tu.Message.Text
	shout.Finished = true
	if err := db.Save(shout).Error; err != nil {
		logTelegram("[shout.go - 143] " + err.Error())
	}

	user := &User{}
	db.First(user, shout.OwnerID)
	msg := tgbotapi.NewMessage(int64(user.TelegramID), tr(user.TelegramID, "shoutFinish"))
	bot.Send(msg)
}

func initShoutService() *ShoutService {
	ss := &ShoutService{}
	go ss.start()
	return ss
}

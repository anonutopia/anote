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
	db.Save(kslsd)
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
				db.Save(ksmc)

				if shout.ID != 1 {
					shout.Published = true
				}

				db.Save(&shout)
			}
		}

		// todo - make sure that everything is ok with 100 here
		pages, err := wnc.TransactionsAddressLimit(conf.ShoutAddress, 100)
		if err != nil {
			log.Println(err)
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
	db.Save(tr)
}

func (ss *ShoutService) processBid(t *gowaves.TransactionsAddressLimitResponse) {
	user := &User{Address: t.Sender}
	db.First(user, user)
	msg := tgbotapi.NewMessage(int64(user.TelegramID), tr(user.TelegramID, "shoutMessage"))
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: false}
	bot.Send(msg)

	shout := &Shout{Owner: user}
	db.FirstOrCreate(shout)

	shout.ChatID = int(msg.ChatID)
	shout.Price = t.Amount
	db.Save(shout)
}

func (ss *ShoutService) setMessage(tu TelegramUpdate) {
	shout := &Shout{ChatID: tu.Message.Chat.ID}
	db.First(shout, shout)
	shout.Message = tu.Message.Text
	db.Save(shout)

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
	db.Save(shout)

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

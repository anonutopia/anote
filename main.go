package main

import (
	"fmt"

	"github.com/anonutopia/gowaves"
	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	macaron "gopkg.in/macaron.v1"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var wmc *gowaves.WavesMatcherClient

var db *gorm.DB

// var dbBak *gorm.DB

var bot *tgbotapi.BotAPI

var m *macaron.Macaron

var pc *PriceClient

var tm *TokenMonitor

var ss *ShoutService

var wm *WavesMonitor

func main() {
	conf = initConfig()

	db = initDb()

	// dbBak = initDbBak()

	wnc, wmc = initWaves()

	bot = initBot()

	pc = initPriceClient()

	tm = initTokenMonitor()

	m = initMacaron()
	m.Post(fmt.Sprintf("/%s", conf.TelegramAPIKey), binding.Json(TelegramUpdate{}), webhookView)

	initMinerMonitor()

	ss = initShoutService()

	// send()

	go hashingPower()

	// go clean1()

	// go restoreBackup()

	initMonitor()
}

package main

import (
	"log"
	"os"

	macaron "gopkg.in/macaron.v1"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var db *gorm.DB

var pc *PriceClient

var wm *WavesMonitor

var bot *telebot.Bot

var m *macaron.Macaron

var um *UserManager

var tm *TokenMonitor

var sigs chan os.Signal

var umDone chan struct{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sigs, umDone = initSignalHandler()

	initLangs()

	conf = initConfig()

	bot = initTelegramBot()

	db = initDb()

	pc = initPriceClient()

	tm = initTokenMonitor()

	m = initMacaron()

	um = initUserManager()

	wm = initWavesMonitor()

	initCommands()

	logTelegram("I've started.")

	go bot.Start()

	<-umDone
}

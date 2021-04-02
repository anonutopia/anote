package main

import (
	"log"

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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initSignalHandler()

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

	logTelegram("Anote daemon successfully started. 🚀")

	users := []*User{}
	db.Find(&users)
	counter := 0

	for _, u := range users {
		if u.MinedAnotes > 5000*int(SatInBTC) {
			counter++
			log.Printf("%s - %.8f", u.Address, float64(u.MinedAnotes)/float64(SatInBTC))
		}
	}

	bot.Start()
}

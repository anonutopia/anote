package main

import (
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

var conf *Config

var pc *PriceClient

var bot *telebot.Bot

var tm *TokenMonitor

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initSignalHandler()

	initLangs()

	conf = initConfig()

	bot = initTelegramBot()

	pc = initPriceClient()

	tm = initTokenMonitor()

	initWavesMonitor()

	initCommands()

	logTelegram("Anote daemon successfully started. 🚀")

	bot.Start()
}

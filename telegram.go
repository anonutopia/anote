package main

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initTelegramBot() *telebot.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:     conf.TelegramAPIKey,
		Poller:    &tb.LongPoller{Timeout: TelPollerTimeout * time.Second},
		Verbose:   conf.Debug,
		ParseMode: tb.ModeHTML,
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func logTelegram(message string) {
	group := &telebot.Chat{ID: TelAnonOps}
	if _, err := bot.Send(group, message, tb.NoPreview); err != nil {
		log.Println(err)
	}
}

func messageTelegram(message string, groupId int) {
	var group *telebot.Chat
	if conf.Dev {
		group = &telebot.Chat{ID: TelAnonOps}
	} else {
		group = &telebot.Chat{ID: int64(groupId)}
	}
	if _, err := bot.Send(group, message, tb.NoPreview); err != nil {
		log.Println(err)
	}
}

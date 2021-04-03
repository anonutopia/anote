package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func initTelegramBot() *tb.Bot {
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
	group := &tb.Chat{ID: TelAnonOps}
	if _, err := bot.Send(group, message, tb.NoPreview); err != nil {
		log.Println(err)
	}
}

func messageTelegram(message string, groupId int) {
	var group *tb.Chat
	if conf.Dev {
		group = &tb.Chat{ID: TelAnonOps}
	} else {
		group = &tb.Chat{ID: int64(groupId)}
	}
	if _, err := bot.Send(group, message, tb.NoPreview); err != nil {
		log.Println(err)
	}
}

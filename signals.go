package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func initSignalHandler() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")

		// Stop Telegram bot
		bot.Stop()
	}()
}

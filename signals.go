package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func initSignalHandler() chan struct{} {
	sigs := make(chan os.Signal, 1)
	umDone := make(chan struct{})

	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")

		// Save state of all users
		um.saveState()

		// Stop Telegram bot
		bot.Stop()

		// Finish with background tasks
		for {
			if !um.Running {
				umDone <- struct{}{}
			}
			time.Sleep(time.Second * SigTick)
		}
	}()

	return umDone
}

package main

import (
	"os"
	"os/signal"
	"syscall"
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/channel"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/group"
	"telegram-channel-publisher/publisher"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.WithField("source", "main")
	config.InitLog()
	log.Infoln("Starting server...")
	publisher.InitPublisher()
	bot.InitBot()
	channel.HandleUpdate(bot.ReceiveMessage())
	group.InitMediaHandleCron()
	log.Infoln("Server started successfully!")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("Shutting down server...")
	log.Infoln("Server exit!")
}

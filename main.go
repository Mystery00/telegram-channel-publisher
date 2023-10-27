package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/channel"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/publisher"
)

func main() {
	log := logrus.WithField("source", "main")
	config.InitLog()
	publisher.InitPublisher()
	channel.HandleUpdate(bot.InitBot(), bot.ReceiveMessage())
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("Shutting down server...")
	log.Infoln("Server exit!")
}

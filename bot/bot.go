package bot

import (
	"telegram-channel-publisher/config"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	token       = viper.GetString(config.BotToken)
	apiEndpoint = viper.GetString(config.ApiEndpoint)
)

var bot *tgbot.BotAPI

func InitBot() *tgbot.BotAPI {
	if apiEndpoint == "" {
		apiEndpoint = tgbot.APIEndpoint
	}
	logrus.Debugf("init bot with api endpoint: %s", apiEndpoint)
	newBot, err := tgbot.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	if err != nil {
		logrus.Panic(err)
	}
	newBot.Debug = false
	logrus.Infof("Authorized on account %s", newBot.Self.UserName)
	bot = newBot
	return bot
}

func ReceiveMessage() <-chan tgbot.Update {
	result := make(chan tgbot.Update)
	go func() {
		u := tgbot.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			result <- update
		}
	}()
	return result
}

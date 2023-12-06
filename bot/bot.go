package bot

import (
	"telegram-channel-publisher/config"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	token       = viper.GetString(config.BotToken)
	apiEndpoint = viper.GetString(config.ApiEndpoint)
)

var bot *tgbot.BotAPI

func InitBot() {
	if apiEndpoint == "" {
		apiEndpoint = tgbot.APIEndpoint
	}
	logrus.Debugf("init bot with api endpoint: %s", apiEndpoint)
	newBot, err := tgbot.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}
	newBot.Debug = false
	logrus.Infof("Authorized on account %s", newBot.Self.UserName)
	bot = newBot
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

func DeleteMessage(chatId int64, messageId int, delayTime int) {
	if delayTime > 0 {
		time.Sleep(time.Duration(delayTime) * time.Second)
		if resp, err := bot.Request(tgbot.NewDeleteMessage(chatId, messageId)); err != nil || !resp.Ok {
			logrus.Warnf("Error delete message, resp: %s, err: %v", string(resp.Result), err)
		}
	}
}

func Reply(chatId int64, messageId int, text string) (tgbot.Message, error) {
	msg := tgbot.NewMessage(chatId, text)
	msg.ReplyToMessageID = messageId
	return bot.Send(msg)
}

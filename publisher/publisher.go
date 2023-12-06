package publisher

import (
	"fmt"
	"telegram-channel-publisher/bot"
	"telegram-channel-publisher/config"
	"telegram-channel-publisher/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Publisher interface {
	// Publish 发布消息
	Publish(post model.Post)
}

var publisher Publisher

func InitPublisher() {
	pubType := viper.GetString(config.PublisherType)
	logrus.Debugf("publisher type: %s", pubType)
	switch pubType {
	case "halo":
		publisher = &HaloPublisher{}
		break
	default:
		publisher = &LogPublisher{}
		break
	}
}

func Pub(post model.Post) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("publish error: %v", err)
			replyFail(post, err.(error))
			return
		}
	}()
	publisher.Publish(post)
	replySuccess(post)
}

func replySuccess(post model.Post) {
	var enable bool
	var delayTime int
	if post.IsPrivate {
		enable = viper.GetBool(config.PrivateReplyEnable)
		delayTime = viper.GetInt(config.PrivateReplyDelay)
	} else {
		enable = viper.GetBool(config.ChannelReplyEnable)
		delayTime = viper.GetInt(config.ChannelReplyDelay)
	}
	if !enable {
		return
	}
	logrus.Debugf("reply success message to telegram, delay time: %d", delayTime)
	if delayTime < 0 {
		delayTime = 0
	}
	if delayTime > 60 {
		delayTime = 60
	}
	text := fmt.Sprintf("消息处理成功（本消息 %d秒 后自动删除）", delayTime)
	reply, err := bot.Reply(post.ChatId, post.MessageId, text)
	if err != nil {
		logrus.Errorf("reply message to telegram error: %v", err)
		return
	}
	go bot.DeleteMessage(reply.Chat.ID, reply.MessageID, delayTime)
}

func replyFail(post model.Post, err error) {
	var enable bool
	if post.IsPrivate {
		enable = viper.GetBool(config.PrivateReplyEnable)
	} else {
		enable = viper.GetBool(config.ChannelReplyEnable)
	}
	if !enable {
		return
	}
	logrus.Debugf("reply failed message to telegram")
	text := fmt.Sprintf("消息处理失败，失败原因：%s", err)
	_, err = bot.Reply(post.ChatId, post.MessageId, text)
	if err != nil {
		logrus.Errorf("reply message to telegram error: %v", err)
		return
	}
}

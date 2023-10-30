package publisher

import (
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
			return
		}
	}()
	publisher.Publish(post)
}
